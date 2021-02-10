package md

import (
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/utils"
)

func initData(db *repositories.MysqlRepo) {
	initEntityData(db)
	initEnumData(db)
}

func initEntityData(db *repositories.MysqlRepo) {
	items := make([]MDEntity, 0)
	//基础数据类型
	items = append(items, MDEntity{ID: "string", Name: "字符", Type: TYPE_SIMPLE, Domain: md_domain, System: utils.SBool_True})
	items = append(items, MDEntity{ID: "int", Name: "整数", Type: TYPE_SIMPLE, Domain: md_domain, System: utils.SBool_True})
	items = append(items, MDEntity{ID: "bool", Name: "布尔", Type: TYPE_SIMPLE, Domain: md_domain, System: utils.SBool_True})
	items = append(items, MDEntity{ID: "decimal", Name: "浮点数", Type: TYPE_SIMPLE, Domain: md_domain, System: utils.SBool_True})
	items = append(items, MDEntity{ID: "text", Name: "文本", Type: TYPE_SIMPLE, Domain: md_domain, System: utils.SBool_True})
	items = append(items, MDEntity{ID: "date", Name: "日期", Type: TYPE_SIMPLE, Domain: md_domain, System: utils.SBool_True})
	items = append(items, MDEntity{ID: "datetime", Name: "时间", Type: TYPE_SIMPLE, Domain: md_domain, System: utils.SBool_True})
	items = append(items, MDEntity{ID: "binary", Name: "二进制", Type: TYPE_SIMPLE, Domain: md_domain, System: utils.SBool_True})
	items = append(items, MDEntity{ID: "xml", Name: "XML", Type: TYPE_SIMPLE, Domain: md_domain, System: utils.SBool_True})
	items = append(items, MDEntity{ID: "json", Name: "JSON", Type: TYPE_SIMPLE, Domain: md_domain, System: utils.SBool_True})

	for i, _ := range items {
		if items[i].Code == "" {
			items[i].Code = items[i].ID
		}
		NewEntitySv(db).UpdateEntity(items[i])
	}
}
func initEnumData(db *repositories.MysqlRepo) {
	items := make([]MDEnumType, 0)
	//md.type.enum
	enumType := MDEnumType{ID: "md.type.enum", Name: "元数据类型", Domain: md_domain}
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_SIMPLE, Name: "简单类型"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_ENUM, Name: "枚举"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_DTO, Name: "对象"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_ENTITY, Name: "实体"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_INTERFACE, Name: "接口"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_VIEW, Name: "视图"})
	items = append(items, enumType)

	//md.type.enum
	enumType = MDEnumType{ID: "md.state.enum", Name: "数据状态", Domain: md_domain}
	enumType.Enums = append(enumType.Enums, MDEnum{ID: STATE_TEMP, Name: "临时"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: STATE_CREATED, Name: "创建的"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: STATE_UPDATED, Name: "更新的"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: STATE_DELETED, Name: "删除的"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: STATE_NORMAL, Name: "正常的"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: STATE_IGNORED, Name: "忽略的"})
	items = append(items, enumType)

	//md.field.type.enum
	enumType = MDEnumType{ID: "md.field.type.enum", Name: "字段数据类型", Domain: md_domain}
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_STRING, Name: "字符"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_INT, Name: "整数"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_BOOL, Name: "布尔"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_DECIMAL, Name: "浮点数"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_TEXT, Name: "文本"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_DATE, Name: "日期"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_DATETIME, Name: "时间"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_XML, Name: "XML"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_JSON, Name: "JSON"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_ENUM, Name: "枚举"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: FIELD_TYPE_ENTITY, Name: "实体"})
	items = append(items, enumType)

	//md.field.type.enum
	enumType = MDEnumType{ID: "md.kind.type.enum", Name: "字段关联关系", Domain: md_domain}
	enumType.Enums = append(enumType.Enums, MDEnum{ID: KIND_TYPE_MANY_TO_MANT, Name: "多对多"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: KIND_TYPE_HAS_MANT, Name: "组合关系"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: KIND_TYPE_HAS_ONE, Name: "一对一"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: KIND_TYPE_BELONGS_TO, Name: "从属关系"})
	items = append(items, enumType)

	//md.rule.type.enum
	enumType = MDEnumType{ID: "md.rule.type.enum", Name: "规则类型", Domain: md_domain}
	enumType.Enums = append(enumType.Enums, MDEnum{ID: RuleType_Widget, Name: "组件"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: RuleType_Entity, Name: "实体"})
	items = append(items, enumType)

	for i, _ := range items {
		NewEntitySv(db).UpdateEnumType(items[i])
	}
}
