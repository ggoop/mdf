package md

import (
	"fmt"
	"github.com/ggoop/mdf/db"
	"github.com/ggoop/mdf/framework/glog"

	"github.com/ggoop/mdf/utils"
)

func (s *mdSvImpl) UpdateEntity(item MDEntity) error {
	if item.ID == "" {
		return nil
	}
	old := MDEntity{}
	db.Default().Model(old).Where("id=?", item.ID).Take(&old)
	if old.ID != "" {
		updates := utils.Map{}
		if old.Name != item.Name && item.Name != "" {
			updates["Name"] = item.Name
		}
		if old.Code != item.Code && item.Code != "" {
			updates["Code"] = item.Code
		}
		if old.Domain != item.Domain && item.Domain != "" {
			updates["Domain"] = item.Domain
		}
		if old.Tags != item.Tags && item.Tags != "" {
			updates["Tags"] = item.Tags
		}
		if old.Memo != item.Memo && item.Memo != "" {
			updates["Memo"] = item.Memo
		}
		if old.Type != item.Type && item.Type != "" {
			updates["Type"] = item.Type
		}
		if old.TableName != item.TableName && item.TableName != "" {
			updates["TableName"] = item.TableName
		}
		if len(updates) > 0 {
			db.Default().Model(&old).Where("id=?", old.ID).Updates(updates)
		}
	} else {
		db.Default().Create(&item)
	}
	return nil
}
func (s *mdSvImpl) UpdateEnumType(enumType MDEnumType) error {
	if enumType.ID == "" {
		return nil
	}
	entity := MDEntity{}
	db.Default().Model(entity).Where("id=?", enumType.ID).Order("id").Take(&entity)
	if entity.ID == "" {
		entity.ID = enumType.ID
		entity.Code = enumType.ID
		entity.Name = enumType.Name
		entity.Type = utils.TYPE_ENUM
		entity.Domain = enumType.Domain
		db.Default().Create(&entity)
	} else {
		updates := utils.Map{}
		if entity.Name != enumType.Name && enumType.Name != "" {
			updates["Name"] = enumType.Name
		}
		if entity.Domain != entity.Domain && enumType.Domain != "" {
			updates["Domain"] = enumType.Domain
		}
		if len(updates) > 0 {
			db.Default().Model(&entity).Updates(updates)
		}
	}
	if len(enumType.Enums) > 0 {
		for i, enum := range enumType.Enums {
			if enum.Sequence == 0 {
				enumType.Enums[i].Sequence = i
			}
			enumType.Enums[i].EntityID = entity.ID
			if _, err := s.UpdateOrCreateEnum(enumType.Enums[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *mdSvImpl) UpdateOrCreateEnum(enum MDEnum) (*MDEnum, error) {
	entity := MDEntity{}
	if enum.EntityID == "" {
		return nil, nil
	}
	db.Default().Model(entity).Where("id=?", enum.EntityID).Order("id").Take(&entity)
	if entity.ID == "" {
		return nil, glog.Error("找不到枚举类型！")
	}
	old := MDEnum{}
	if db.Default().Where("entity_id=? and id=?", enum.EntityID, enum.ID).Order("id").Take(&old).RecordNotFound() {
		db.Default().Create(&enum)
	} else {
		updates := utils.Map{}
		if old.Name != enum.Name && enum.Name != "" {
			updates["Name"] = enum.Name
		}
		if old.SrcID != enum.SrcID && enum.SrcID != "" {
			updates["SrcID"] = enum.SrcID
		}
		if old.Sequence != enum.Sequence && enum.Sequence >= 0 {
			updates["Sequence"] = enum.Sequence
		}
		if len(updates) > 0 {
			db.Default().Model(&old).Updates(updates)
		}
	}
	return &enum, nil
}

func (s *mdSvImpl) buildColumnNameString4Oracle(item MDField) string {
	fieldStr := s.quote(item.DbName)
	nullable := item.Nullable

	if item.IsPrimaryKey.IsTrue() && item.TypeID == utils.FIELD_TYPE_STRING {
		fieldStr += " VARCHAR2(36)"
		nullable = utils.SBool_False
	} else if item.IsPrimaryKey.IsTrue() && item.TypeID == utils.FIELD_TYPE_INT {
		fieldStr += " NUMBER"
		nullable = utils.SBool_False
	} else if item.TypeID == utils.FIELD_TYPE_STRING {
		if item.Length <= 0 {
			item.Length = 50
		}
		if item.Length >= 4000 {
			fieldStr += " CLOB"
		} else {
			fieldStr += fmt.Sprintf(" VARCHAR2(%d)", item.Length)
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == utils.FIELD_TYPE_BOOL {
		fieldStr += " NUMBER(1,0)"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		nullable = utils.SBool_False
	} else if item.TypeID == utils.FIELD_TYPE_DATE || item.TypeID == utils.FIELD_TYPE_DATETIME {
		fieldStr += " TIMESTAMP"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == utils.FIELD_TYPE_DECIMAL {
		fieldStr += " NUMBER(24,9)"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		nullable = utils.SBool_False
	} else if item.TypeID == utils.FIELD_TYPE_INT {
		fieldStr += " INTEGER"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		nullable = utils.SBool_False
	} else if item.TypeType == utils.TYPE_ENTITY || item.TypeType == utils.TYPE_ENUM {
		fieldStr += " VARCHAR2(36)"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else {
		if item.Length <= 0 {
			item.Length = 255
		}
		fieldStr += fmt.Sprintf(" VARCHAR2(%d)", item.Length)
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	}
	if !nullable.IsTrue() {
		fieldStr += " NOT NULL"
	}
	return fieldStr
}

func (s *mdSvImpl) buildColumnNameString4Mysql(item MDField) string {
	fieldStr := s.quote(item.DbName)
	if item.IsPrimaryKey.IsTrue() && item.TypeID == utils.FIELD_TYPE_STRING {
		fieldStr += " NVARCHAR(36)"
		item.Nullable = utils.SBool_False
	} else if item.IsPrimaryKey.IsTrue() && item.TypeID == utils.FIELD_TYPE_INT {
		fieldStr += " BIGINT"
		item.Nullable = utils.SBool_False
	} else if item.TypeID == utils.FIELD_TYPE_STRING {
		if item.Length <= 0 {
			item.Length = 50
		}
		if item.Length >= 8000 {
			fieldStr += " LONGTEXT"
		} else if item.Length >= 4000 {
			fieldStr += " TEXT"
		} else {
			fieldStr += fmt.Sprintf(" NVARCHAR(%d)", item.Length)
		}

		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == utils.FIELD_TYPE_BOOL {
		fieldStr += " TINYINT"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		item.Nullable = utils.SBool_False
	} else if item.TypeID == utils.FIELD_TYPE_DATE || item.TypeID == utils.FIELD_TYPE_DATETIME {
		fieldStr += " TIMESTAMP"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == utils.FIELD_TYPE_DECIMAL {
		fieldStr += " DECIMAL(24,9)"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		item.Nullable = utils.SBool_False
	} else if item.TypeID == utils.FIELD_TYPE_INT {
		if item.Length >= 8 {
			fieldStr += " BIGINT"
		} else {
			fieldStr += " INT"
		}
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		item.Nullable = utils.SBool_False
	} else if item.TypeType == utils.TYPE_ENTITY || item.TypeType == utils.TYPE_ENUM {
		fieldStr += " nvarchar(36)"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else {
		if item.Length <= 0 {
			item.Length = 255
		}
		fieldStr += fmt.Sprintf(" nvarchar(%d)", item.Length)
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	}
	if !item.Nullable.IsTrue() {
		fieldStr += " NOT NULL"
	}
	fieldStr += fmt.Sprintf(" COMMENT '%s'", item.Name)
	return fieldStr

}
