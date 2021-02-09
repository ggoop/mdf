package md

import (
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/utils"
)

func initData(db *repositories.MysqlRepo) {
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
		NewEntitySv(db).UpdateEntity(items[i])
	}

	//基础枚举
	enumID := "md.type.enum"
	enumType := MDEnumType{ID: enumID, Name: "元数据类型", Domain: md_domain}
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_SIMPLE, Name: "简单类型"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_ENUM, Name: "枚举"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_DTO, Name: "对象"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_ENTITY, Name: "实体"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_INTERFACE, Name: "接口"})
	enumType.Enums = append(enumType.Enums, MDEnum{ID: TYPE_VIEW, Name: "视图"})

	NewEntitySv(db).UpdateEnumType(enumType)
}
