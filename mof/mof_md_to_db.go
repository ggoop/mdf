package mof

import (
	"fmt"

	"github.com/ggoop/mdf/md"
	"github.com/ggoop/mdf/utils"
)

func (s *MOFSv) buildColumnNameString4Oracle(item md.MDField) string {
	fieldStr := s.quote(item.DbName)
	nullable := item.Nullable

	if item.IsPrimaryKey.IsTrue() && item.TypeID == md.FIELD_TYPE_STRING {
		fieldStr += " VARCHAR2(36)"
		nullable = utils.SBool_False
	} else if item.IsPrimaryKey.IsTrue() && item.TypeID == md.FIELD_TYPE_INT {
		fieldStr += " NUMBER"
		nullable = utils.SBool_False
	} else if item.TypeID == md.FIELD_TYPE_STRING {
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
	} else if item.TypeID == md.FIELD_TYPE_BOOL {
		fieldStr += " NUMBER(1,0)"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		nullable = utils.SBool_False
	} else if item.TypeID == md.FIELD_TYPE_DATE || item.TypeID == md.FIELD_TYPE_DATETIME {
		fieldStr += " TIMESTAMP"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == md.FIELD_TYPE_DECIMAL {
		fieldStr += " NUMBER(24,9)"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		nullable = utils.SBool_False
	} else if item.TypeID == md.FIELD_TYPE_INT {
		fieldStr += " INTEGER"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		nullable = utils.SBool_False
	} else if item.TypeType == md.TYPE_ENTITY || item.TypeType == md.TYPE_ENUM {
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

func (s *MOFSv) buildColumnNameString4Mysql(item md.MDField) string {
	fieldStr := s.quote(item.DbName)
	if item.IsPrimaryKey.IsTrue() && item.TypeID == md.FIELD_TYPE_STRING {
		fieldStr += " NVARCHAR(36)"
		item.Nullable = utils.SBool_False
	} else if item.IsPrimaryKey.IsTrue() && item.TypeID == md.FIELD_TYPE_INT {
		fieldStr += " BIGINT"
		item.Nullable = utils.SBool_False
	} else if item.TypeID == md.FIELD_TYPE_STRING {
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
	} else if item.TypeID == md.FIELD_TYPE_BOOL {
		fieldStr += " TINYINT"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		item.Nullable = utils.SBool_False
	} else if item.TypeID == md.FIELD_TYPE_DATE || item.TypeID == md.FIELD_TYPE_DATETIME {
		fieldStr += " TIMESTAMP"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		}
	} else if item.TypeID == md.FIELD_TYPE_DECIMAL {
		fieldStr += " DECIMAL(24,9)"
		if item.DefaultValue != "" {
			fieldStr += " DEFAULT " + item.DefaultValue
		} else {
			fieldStr += " DEFAULT 0"
		}
		item.Nullable = utils.SBool_False
	} else if item.TypeID == md.FIELD_TYPE_INT {
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
	} else if item.TypeType == md.TYPE_ENTITY || item.TypeType == md.TYPE_ENUM {
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
