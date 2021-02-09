package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/files"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/framework/mof"
	"github.com/ggoop/mdf/framework/query"
	"github.com/ggoop/mdf/utils"
)

type MdSv struct {
	repo *repositories.MysqlRepo
}

func NewMdSv(repo *repositories.MysqlRepo) *MdSv {
	return &MdSv{repo: repo}
}

func (s *MdSv) InitData(types []string) error {
	f := ""
	f = "./storage/inits/entity.xlsx"
	if utils.PathExists(f) && (len(types) == 0 || utils.StringsContains(types, "entity") >= 0) {
		if data, err := files.NewExcelSv().GetExcelDatas(f); err != nil {
			glog.Error(err)
		} else {
			NewMdSv(s.repo).BatchImport(model.SYS_ENT_ID, data)
		}
	}
	f = "./storage/inits/query.xlsx"
	if utils.PathExists(f) && (len(types) == 0 || utils.StringsContains(types, "query") >= 0) {
		if data, err := files.NewExcelSv().GetExcelDatas(f); err != nil {
			glog.Error(err)
		} else {
			NewQuerySv(s.repo).BatchImport(model.SYS_ENT_ID, data)
		}
	}

	f = "./storage/inits/page.xlsx"
	if utils.PathExists(f) && (len(types) == 0 || utils.StringsContains(types, "page") >= 0) {
		if data, err := files.NewExcelSv().GetExcelDatas(f); err != nil {
			glog.Error(err)
		} else {
			NewMdSv(s.repo).BatchImport(model.SYS_ENT_ID, data)
		}
	}

	f = "./storage/inits/product.xlsx"
	if utils.PathExists(f) && (len(types) == 0 || utils.StringsContains(types, "product") >= 0) {
		sv := NewProductSv(s.repo)
		if data, err := files.NewExcelSv().GetExcelDatas(f); err != nil {
			glog.Error(err)
		} else {
			sv.BatchImport(model.SYS_ENT_ID, data)
		}
	}
	return nil
}
func (s *MdSv) GetPage(page string) (md.MDWidget, error) {
	pageMD := md.MDWidget{}
	if err := s.repo.Model(pageMD).Where("id=?", page).Take(&pageMD).Error; err != nil {
		return pageMD, err
	}
	return pageMD, nil
}
func (s *MdSv) GetPageInfo(req mof.ReqContext) (interface{}, error) {
	return nil, nil
}
func (s *MdSv) TakeDataByQ(req mof.ReqContext) (map[string]interface{}, error) {
	entity := md.GetEntity(req.Entity)
	if entity == nil {
		return nil, nil
	}
	exector := query.NewExector(entity.TableName)
	codeField := &md.MDField{}
	nameField := &md.MDField{}
	entField := &md.MDField{}
	for i, f := range entity.Fields {
		if f.TypeType == md.TYPE_SIMPLE {
			exector.Select("$$" + f.Code + " as \"" + f.DbName + "\"")
			if f.TypeID != "" && f.DbName != "" {
				exector.SetFieldDataType(f.DbName, f.TypeID)
			}
			if strings.Contains(f.Tags, "code") {
				codeField = &entity.Fields[i]
			}
			if strings.Contains(f.Tags, "name") {
				nameField = &entity.Fields[i]
			}
			if strings.Contains(f.Tags, "ent") {
				entField = &entity.Fields[i]
			}
		}
	}
	if entField == nil || entField.Code == "" {
		entField = entity.GetField("EntID")
	}
	if entField != nil && entField.Code != "" {
		exector.Where(fmt.Sprintf("%s=?", entField.DbName), req.EntID)
	}

	if codeField == nil || codeField.Code == "" {
		codeField = entity.GetField("Code")
	}
	if codeField != nil && codeField.Code != "" {
		exec := exector.Clone()
		exec.Where(codeField.Code+" = ?", req.Q)
		if data, err := exec.Take(s.repo); err != nil {
			return nil, err
		} else if len(data) > 0 {
			return data, nil
		}
	}
	if nameField == nil || nameField.Code == "" {
		nameField = entity.GetField("Name")
	}
	if nameField != nil && nameField.Code != "" && nameField.Code != codeField.Code {
		exec := exector.Clone()
		exec.Where(nameField.Code+" = ?", req.Q)
		if data, err := exec.Take(s.repo); err != nil {
			return nil, err
		} else if len(data) > 0 {
			return data, nil
		}
	}
	return nil, nil
}

//执行命令
func (s *MdSv) DoAction(req mof.ReqContext) (interface{}, error) {
	return mof.DoAction(req)
}
func (s *MdSv) BatchImport(entID string, datas []files.ImportData) error {
	if len(datas) > 0 {
		nameList := make(map[string]int)
		nameList["Entity"] = 1
		nameList["Props"] = 2
		nameList["Page"] = 3
		nameList["Widgets"] = 4
		nameList["ActionCommand"] = 5
		nameList["ActionRule"] = 6

		sort.Slice(datas, func(i, j int) bool { return nameList[datas[i].Name] < nameList[datas[j].Name] })

		entities := make([]md.MDEntity, 0)
		fields := make([]md.MDField, 0)
		for _, item := range datas {
			if item.Name == "Entity" {
				if d, err := s.toEntities(entID, item); err != nil {
					return err
				} else if len(d) > 0 {
					entities = append(entities, d...)
				}
			}
			if item.Name == "Props" {
				if d, err := s.toFields(entID, item); err != nil {
					return err
				} else if len(d) > 0 {
					fields = append(fields, d...)
				}
			}
		}

		mofSv := mof.NewMOFSv(s.repo)
		if len(entities) > 0 {
			for i, entity := range entities {
				for _, field := range fields {
					if entity.ID == field.EntityID {
						if entities[i].Fields == nil {
							entities[i].Fields = make([]md.MDField, 0)
						}
						entities[i].Fields = append(entities[i].Fields, field)
					}
				}
			}
			mofSv.AddMDEntities(entities)
		}
	}
	return nil
}
func (s *MdSv) toEntities(entID string, data files.ImportData) ([]md.MDEntity, error) {
	if len(data.Datas) == 0 {
		return nil, nil
	}
	items := make([]md.MDEntity, 0)
	for _, row := range data.Datas {
		item := md.MDEntity{}
		item.ID = files.GetMapStringValue("ID", row)
		item.Name = files.GetMapStringValue("Name", row)
		item.Type = files.GetMapStringValue("Type", row)
		item.TableName = files.GetMapStringValue("TableName", row)
		item.Domain = files.GetMapStringValue("Domain", row)
		item.System = files.GetMapSBoolValue("System", row)
		items = append(items, item)
	}
	return items, nil
}

func (s *MdSv) toFields(entID string, data files.ImportData) ([]md.MDField, error) {
	if len(data.Datas) == 0 {
		return nil, nil
	}
	items := make([]md.MDField, 0)
	for _, row := range data.Datas {
		item := md.MDField{}
		if cValue := files.GetMapStringValue("EntityID", row); cValue != "" {
			item.EntityID = cValue
		}
		if cValue := files.GetMapStringValue("Name", row); cValue != "" {
			item.Name = cValue
		}
		if cValue := files.GetMapStringValue("Code", row); cValue != "" {
			item.Code = cValue
		}
		if cValue := files.GetMapStringValue("TypeID", row); cValue != "" {
			item.TypeID = cValue
		}
		if cValue := files.GetMapStringValue("Kind", row); cValue != "" {
			item.Kind = cValue
		}
		if cValue := files.GetMapStringValue("ForeignKey", row); cValue != "" {
			item.ForeignKey = cValue
		}
		if cValue := files.GetMapStringValue("AssociationKey", row); cValue != "" {
			item.AssociationKey = cValue
		}
		if cValue := files.GetMapStringValue("DbName", row); cValue != "" {
			item.DbName = cValue
		}
		if cValue := files.GetMapStringValue("DbName", row); cValue != "" {
			item.DbName = cValue
		}
		if cValue := files.GetMapIntValue("Length", row); cValue >= 0 {
			item.Length = cValue
		}
		if cValue := files.GetMapIntValue("Precision", row); cValue >= 0 {
			item.Precision = cValue
		}
		if cValue := files.GetMapStringValue("DefaultValue", row); cValue != "" {
			item.DefaultValue = cValue
		}
		if cValue := files.GetMapStringValue("MaxValue", row); cValue != "" {
			item.MaxValue = cValue
		}
		if cValue := files.GetMapStringValue("MinValue", row); cValue != "" {
			item.MinValue = cValue
		}
		if cValue := files.GetMapStringValue("Tags", row); cValue != "" {
			item.Tags = cValue
		}
		if cValue := files.GetMapStringValue("Limit", row); cValue != "" {
			item.Limit = cValue
		}
		item.Nullable = files.GetMapSBoolValue("Nullable", row)
		item.IsPrimaryKey = files.GetMapSBoolValue("IsPrimaryKey", row)
		items = append(items, item)
	}
	return items, nil
}
