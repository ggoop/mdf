package md

import (
	"fmt"
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/framework/files"
	"sort"
	"strings"
	"sync"

	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

type IMDSv interface {
	Migrate(values ...interface{})
	AddMDEntities(items []MDEntity) error
	BatchImport(datas []files.ImportData) error
	GetEntity(id string) *MDEntity

	GetEnum(typeId string, values ...string) *MDEnum

	TakeDataByQ(token *utils.TokenContext, req *ReqContext) (map[string]interface{}, error)
	UpdateEntity(item MDEntity) error
	UpdateEnumType(enumType MDEnumType) error

	QuotedBy(m MD, ids []string, excludes ...MD) ([]MDEntity, []string)
}

type mdSvImpl struct {
	*sync.Mutex
	mdCache         map[string]*MDEntity
	enumCache       map[string]*MDEnum
	initMDCompleted bool
}

func MDSv() IMDSv {
	return mdSv
}

var mdSv IMDSv = newMDSv()

func newMDSv() *mdSvImpl {
	return &mdSvImpl{
		Mutex:     &sync.Mutex{},
		mdCache:   make(map[string]*MDEntity),
		enumCache: make(map[string]*MDEnum),
	}
}

func (s *mdSvImpl) InitCache() {
	s.enumCache = make(map[string]*MDEnum)
	items, _ := s.GetEnums()
	for i, _ := range items {
		v := items[i]
		s.enumCache[strings.ToLower(v.EntityID+":"+v.ID)] = &v
		s.enumCache[strings.ToLower(v.EntityID+":"+v.Name)] = &v
	}
}

func (s *mdSvImpl) GetEnums() ([]MDEnum, error) {
	items := make([]MDEnum, 0)
	if err := db.Default().Model(&MDEnum{}).Where("entity_id in (?)", db.Default().Model(MDEntity{}).Select("id").Where("type=?", "enum").SubQuery()).Order("entity_id").Order("sequence").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *mdSvImpl) Migrate(values ...interface{}) {
	//先增加模型表
	if !s.initMDCompleted {
		s.initMDCompleted = true
		mds := []interface{}{
			&MDEntity{}, &MDEntityRelation{}, &MDField{}, &MDEnum{},
			&MDActionCommand{}, &MDActionRule{},
			&MDWidget{}, &MDWidgetDatasource{}, &MDWidgetLayout{}, &MDWidgetItem{},
			&MDToolbars{}, &MDToolbarItem{},
			&MDActionCommand{}, &MDActionRule{},
			&MDFilters{}, &MDFilterSolution{}, &MDFilterItem{},
		}
		needDb := make([]interface{}, 0)
		for _, v := range mds {
			m := newMd(v)
			if dd := m.GetMder(); dd == nil || dd.Type == utils.TYPE_ENTITY || dd.Type == utils.TYPE_ENUM || dd.Type == "" {
				needDb = append(needDb, v)
			}
		}
		db.Default().AutoMigrate(needDb...)
		glog.Error("AutoMigrate MD")
		for _, v := range mds {
			m := newMd(v)
			m.Migrate()
		}

		initData()
	}
	if len(values) > 0 {
		needDb := make([]interface{}, 0)
		for _, v := range values {
			m := newMd(v)
			if dd := m.GetMder(); dd == nil || dd.Type == utils.TYPE_ENTITY || dd.Type == utils.TYPE_ENUM || dd.Type == "" {
				needDb = append(needDb, v)
			}
			m.Migrate()
		}
		if err := db.Default().AutoMigrate(needDb...).Error; err != nil {
			glog.Error(err)
		}
	}

}
func (s *mdSvImpl) ImportUIMeta(entityID string) {
	return
}
func (s *mdSvImpl) entityToTables(items []MDEntity, oldItems []MDEntity) error {
	if items == nil || len(items) == 0 {
		return nil
	}
	for i, _ := range items {
		entity := items[i]
		oldItem := MDEntity{}
		for oi, ov := range oldItems {
			if entity.ID == ov.ID {
				oldItem = oldItems[oi]
				break
			}
		}
		if entity.Type != utils.TYPE_ENTITY {
			continue
		}
		if entity.TableName == "" {
			entity.TableName = strings.ReplaceAll(entity.ID, ".", "_")
		}
		if db.Default().Dialect().HasTable(entity.TableName) {
			s.updateTable(entity, oldItem)
		} else {
			s.createTable(entity)
		}
	}
	return nil
}
func (s *mdSvImpl) createTable(item MDEntity) {
	if len(item.Fields) == 0 {
		return
	}
	var primaryKeys []string
	var tags []string
	for i, _ := range item.Fields {
		field := item.Fields[i]
		if field.IsNormal.IsTrue() {
			tags = append(tags, s.buildColumnNameString(field))
			if field.IsPrimaryKey.IsTrue() {
				primaryKeys = append(primaryKeys, s.quote(field.DbName))
			}
		}
	}
	var primaryKeyStr string
	if len(primaryKeys) > 0 {
		primaryKeyStr = fmt.Sprintf(", PRIMARY KEY (%v)", strings.Join(primaryKeys, ","))
	}
	var tableOptions string
	if err := db.Default().Exec(fmt.Sprintf("CREATE TABLE %v (%v %v)%s", s.quote(item.TableName), strings.Join(tags, ","), primaryKeyStr, tableOptions)).Error; err != nil {
		glog.Error(err)
	}
}
func (s *mdSvImpl) quote(str string) string {
	return db.Default().Dialect().Quote(str)
}
func (s *mdSvImpl) updateTable(item MDEntity, old MDEntity) {
	//更新栏目
	for i := range item.Fields {
		field := item.Fields[i]
		if field.DbName == "-" || field.DbName == "" {
			continue
		}
		oldField := MDField{}
		for oi, ov := range old.Fields {
			if strings.ToLower(field.Code) == strings.ToLower(ov.Code) {
				oldField = old.Fields[oi]
				break
			}
		}
		if !field.IsNormal.IsTrue() {
			continue
		}
		newString := s.buildColumnNameString(field)
		if db.Default().Dialect().HasColumn(item.TableName, field.DbName) { //字段已存在
			oldString := s.buildColumnNameString(oldField)
			//修改字段类型、类型长度、默认值、注释
			if oldString != newString && strings.Contains(item.Tags, "update") {
				if err := db.Default().Exec(fmt.Sprintf("ALTER TABLE %v MODIFY %v", s.quote(item.TableName), newString)).Error; err != nil {
					glog.Error(err)
				}
			}
		} else { //新增字段
			if err := db.Default().Exec(fmt.Sprintf("ALTER TABLE %v ADD %v", s.quote(item.TableName), newString)).Error; err != nil {
				glog.Error(err)
			}
		}
	}
}
func (s *mdSvImpl) buildColumnNameString(item MDField) string {
	/*
		column_definition:
		data_type [NOT NULL | NULL] [DEFAULT {literal | (expr)} ]
		[AUTO_INCREMENT] [UNIQUE [KEY]] [[PRIMARY] KEY]
		[COMMENT 'string']

	*/
	dialectName := db.Default().Dialect().GetName()
	if dialectName == "godror" || dialectName == "oracle" {
		return s.buildColumnNameString4Oracle(item)
	} else {
		return s.buildColumnNameString4Mysql(item)
	}
}
func (s *mdSvImpl) AddMDEntities(items []MDEntity) error {
	entityIds := make([]string, 0)
	oldEntities := make([]MDEntity, 0)
	for i, _ := range items {
		entity := items[i]
		if entity.ID == "" {
			continue
		}
		oldEntity := MDEntity{}
		if db.Default().Model(oldEntity).Preload("Fields").Order("id").Where("id=?", entity.ID).Take(&oldEntity); oldEntity.ID != "" {
			oldEntities = append(oldEntities, oldEntity)
			datas := make(map[string]interface{})
			if oldEntity.TableName != entity.TableName {
				datas["TableName"] = entity.TableName
			}
			if oldEntity.Type != entity.Type {
				datas["Type"] = entity.Type
			}
			if oldEntity.Code != entity.Code {
				datas["Code"] = entity.Code
			}
			if oldEntity.Tags != entity.Tags {
				datas["Tags"] = entity.Tags
			}
			if entity.System.Valid() && oldEntity.System.NotEqual(entity.System) {
				datas["System"] = entity.System
			}
			if oldEntity.Domain != entity.Domain {
				datas["Domain"] = entity.Domain
			}
			if oldEntity.Name != entity.Name {
				datas["Name"] = entity.Name
			}
			if oldEntity.Memo != entity.Memo {
				datas["Memo"] = entity.Memo
			}
			if len(datas) > 0 {
				db.Default().Model(MDEntity{}).Where("id=?", oldEntity.ID).Updates(datas)
			}
		} else {
			if entity.Type == utils.TYPE_ENTITY && entity.TableName == "" {
				entity.TableName = strings.ReplaceAll(entity.ID, ".", "_")
			}
			db.Default().Create(&entity)
		}
		entityIds = append(entityIds, entity.ID)
	}
	//属性字段
	for i, _ := range items {
		entity := items[i]
		if entity.ID != "" && entity.Type == utils.TYPE_ENTITY && len(entity.Fields) > 0 {
			itemCodes := make([]string, 0)
			for f, _ := range entity.Fields {
				field := entity.Fields[f]
				itemCodes = append(itemCodes, field.Code)
				field.IsNormal = utils.SBool_True
				if field.DbName == "-" {
					field.IsNormal = utils.SBool_False
				}
				if field.DbName == "" && field.IsNormal.IsTrue() {
					field.DbName = utils.SnakeString(field.Code)
				}
				oldField := MDField{}
				if fieldType := s.GetEntity(field.TypeID); fieldType != nil {
					field.TypeType = fieldType.Type
					field.TypeID = fieldType.ID
				}
				if field.TypeType == utils.TYPE_ENTITY { //实体
					if field.Kind == "" {
						field.Kind = "belongs_to"
					}
					if field.ForeignKey == "" {
						field.ForeignKey = fmt.Sprintf("%sID", field.Code)
					}
					if field.AssociationKey == "" {
						field.AssociationKey = "ID"
					}
					field.IsNormal = utils.SBool_False
				} else if field.TypeType == utils.TYPE_ENUM { //枚举
					if field.Kind == "" {
						field.Kind = "belongs_to"
					}
					if field.ForeignKey == "" {
						field.ForeignKey = fmt.Sprintf("%sID", field.Code)
					}
					if field.AssociationKey == "" {
						field.AssociationKey = "ID"
					}
					field.IsNormal = utils.SBool_False
				}
				if db.Default().Model(MDField{}).Order("id").Where("entity_id=? and code=?", entity.ID, field.Code).Take(&oldField); oldField.ID != "" {
					datas := make(map[string]interface{})
					if oldField.Name != field.Name {
						datas["Name"] = field.Name
					}
					if oldField.Tags != field.Tags {
						datas["Tags"] = field.Tags
					}
					if oldField.DbName != field.DbName {
						datas["DbName"] = field.DbName
					}
					if oldField.IsNormal != field.IsNormal {
						datas["IsNormal"] = field.IsNormal
					}
					if oldField.IsPrimaryKey != field.IsPrimaryKey {
						datas["IsPrimaryKey"] = field.IsPrimaryKey
					}
					if oldField.Length != field.Length {
						datas["Length"] = field.Length
					}
					if oldField.Nullable != field.Nullable {
						datas["Nullable"] = field.Nullable
					}
					if oldField.DefaultValue != field.DefaultValue {
						datas["DefaultValue"] = field.DefaultValue
					}
					if oldField.TypeID != field.TypeID {
						datas["TypeID"] = field.TypeID
					}
					if oldField.TypeType != field.TypeType {
						datas["TypeType"] = field.TypeType
					}
					if oldField.Limit != field.Limit {
						datas["Limit"] = field.Limit
					}
					if oldField.MinValue != field.MinValue {
						datas["MinValue"] = field.MinValue
					}
					if oldField.MaxValue != field.MaxValue {
						datas["MaxValue"] = field.MaxValue
					}
					if oldField.Precision != field.Precision {
						datas["Precision"] = field.Precision
					}
					if oldField.AssociationKey != field.AssociationKey {
						datas["AssociationKey"] = field.AssociationKey
					}
					if oldField.ForeignKey != field.ForeignKey {
						datas["ForeignKey"] = field.ForeignKey
					}
					if oldField.Kind != field.Kind {
						datas["Kind"] = field.Kind
					}
					if oldField.Sequence != field.Sequence {
						datas["Sequence"] = field.Sequence
					}
					if oldField.SrcID != field.SrcID && field.SrcID != "" {
						datas["SrcID"] = field.SrcID
					}
					if len(datas) > 0 {
						db.Default().Model(MDField{}).Where("entity_id=? and code=?", entity.ID, field.Code).Updates(datas)
					}
				} else {
					db.Default().Create(&field)
				}
			}
			db.Default().Delete(MDField{}, "entity_id=? and code not in (?)", entity.ID, itemCodes)
		}
	}
	//枚举
	for _, entity := range items {
		if entity.ID != "" && entity.Type == utils.TYPE_ENUM && len(entity.Fields) > 0 {
			itemCodes := make([]string, 0)
			for f, field := range entity.Fields {
				newEnum := MDEnum{ID: field.Code, EntityID: entity.ID, Sequence: f, Name: field.Name}
				oldEnum := MDEnum{}
				itemCodes = append(itemCodes, newEnum.ID)
				if db.Default().Model(oldEnum).Order("id").Where("entity_id=? and id=?", newEnum.EntityID, newEnum.ID).Take(&oldEnum); oldEnum.ID != "" {
					datas := make(map[string]interface{})
					if oldEnum.Name != field.Name {
						datas["Name"] = field.Name
					}
					if oldEnum.Sequence != newEnum.Sequence {
						datas["Sequence"] = newEnum.Sequence
					}
					if len(datas) > 0 {
						db.Default().Model(MDEnum{}).Where("entity_id=? and id=?", oldEnum.EntityID, oldEnum.ID).Updates(datas)
					}
				} else {
					db.Default().Create(&newEnum)
				}
			}
			db.Default().Delete(MDEnum{}, "entity_id=? and id not in (?)", entity.ID, itemCodes)
		}
	}
	if len(entityIds) > 0 {
		toTables := make([]MDEntity, 0)
		db.Default().Model(MDEntity{}).Preload("Fields").Where("id in (?) and type=?", entityIds, utils.TYPE_ENTITY).Find(&toTables)
		return s.entityToTables(toTables, oldEntities)
	}
	//缓存
	s.InitCache()
	return nil
}

func (s *mdSvImpl) BatchImport(datas []files.ImportData) error {
	if len(datas) > 0 {
		nameList := make(map[string]int)
		nameList["Entity"] = 1
		nameList["Props"] = 2
		nameList["Page"] = 3
		nameList["Widgets"] = 4
		nameList["ActionCommand"] = 5
		nameList["ActionRule"] = 6

		sort.Slice(datas, func(i, j int) bool { return nameList[datas[i].SheetName] < nameList[datas[j].SheetName] })

		entities := make([]MDEntity, 0)
		fields := make([]MDField, 0)
		for _, item := range datas {
			if item.SheetName == "Entity" {
				if d, err := s.toEntities(item); err != nil {
					return err
				} else if len(d) > 0 {
					entities = append(entities, d...)
				}
			}
			if item.SheetName == "Props" {
				if d, err := s.toFields(item); err != nil {
					return err
				} else if len(d) > 0 {
					fields = append(fields, d...)
				}
			}
		}
		if len(entities) > 0 {
			for i, entity := range entities {
				for _, field := range fields {
					if entity.ID == field.EntityID {
						if entities[i].Fields == nil {
							entities[i].Fields = make([]MDField, 0)
						}
						entities[i].Fields = append(entities[i].Fields, field)
					}
				}
			}
			s.AddMDEntities(entities)
		}
	}
	return nil
}
func (s *mdSvImpl) toEntities(data files.ImportData) ([]MDEntity, error) {
	if len(data.Data) == 0 {
		return nil, nil
	}
	items := make([]MDEntity, 0)
	for _, row := range data.Data {
		item := MDEntity{}
		item.ID = files.GetCellValue("ID", row)
		item.Name = files.GetCellValue("Name", row)
		item.Type = files.GetCellValue("Type", row)
		item.TableName = files.GetCellValue("TableName", row)
		item.Domain = files.GetCellValue("Domain", row)
		item.System = utils.ToSBool(files.GetCellValue("System", row))
		items = append(items, item)
	}
	return items, nil
}
func (s *mdSvImpl) toFields(data files.ImportData) ([]MDField, error) {
	if len(data.Data) == 0 {
		return nil, nil
	}
	items := make([]MDField, 0)
	for _, row := range data.Data {
		item := MDField{}
		if cValue := files.GetCellValue("EntityID", row); cValue != "" {
			item.EntityID = cValue
		}
		if cValue := files.GetCellValue("Name", row); cValue != "" {
			item.Name = cValue
		}
		if cValue := files.GetCellValue("Code", row); cValue != "" {
			item.Code = cValue
		}
		if cValue := files.GetCellValue("TypeID", row); cValue != "" {
			item.TypeID = cValue
		}
		if cValue := files.GetCellValue("Kind", row); cValue != "" {
			item.Kind = cValue
		}
		if cValue := files.GetCellValue("ForeignKey", row); cValue != "" {
			item.ForeignKey = cValue
		}
		if cValue := files.GetCellValue("AssociationKey", row); cValue != "" {
			item.AssociationKey = cValue
		}
		if cValue := files.GetCellValue("DbName", row); cValue != "" {
			item.DbName = cValue
		}
		if cValue := files.GetCellValue("DbName", row); cValue != "" {
			item.DbName = cValue
		}
		if cValue := utils.ToInt(files.GetCellValue("Length", row)); cValue >= 0 {
			item.Length = cValue
		}
		if cValue := utils.ToInt(files.GetCellValue("Precision", row)); cValue >= 0 {
			item.Precision = cValue
		}
		if cValue := files.GetCellValue("DefaultValue", row); cValue != "" {
			item.DefaultValue = cValue
		}
		if cValue := files.GetCellValue("MaxValue", row); cValue != "" {
			item.MaxValue = cValue
		}
		if cValue := files.GetCellValue("MinValue", row); cValue != "" {
			item.MinValue = cValue
		}
		if cValue := files.GetCellValue("Tags", row); cValue != "" {
			item.Tags = cValue
		}
		if cValue := files.GetCellValue("Limit", row); cValue != "" {
			item.Limit = cValue
		}
		item.Nullable = utils.ToSBool(files.GetCellValue("Nullable", row))
		item.IsPrimaryKey = utils.ToSBool(files.GetCellValue("IsPrimaryKey", row))
		items = append(items, item)
	}
	return items, nil
}
