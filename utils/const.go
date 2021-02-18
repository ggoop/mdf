package utils

// 数据状态
const (
	STATE_FIELD = "_state"
	//临时
	STATE_TEMP = "temp"
	//创建的
	STATE_CREATED = "created"
	//更新的
	STATE_UPDATED = "updated"
	//删除的
	STATE_DELETED = "deleted"
	//正常的
	STATE_NORMAL = "normal"
	//忽略的
	STATE_IGNORED = "ignored"
)

//元数据类型
const (
	//简单类型
	TYPE_SIMPLE = "simple"
	//实体
	TYPE_ENTITY = "entity"
	// 枚举
	TYPE_ENUM = "enum"
	// 接口
	TYPE_INTERFACE = "interface"
	// 对象
	TYPE_DTO = "dto"
	// 视图
	TYPE_VIEW = "view"
)

//字段数据类型
const (
	FIELD_TYPE_STRING   = "string"
	FIELD_TYPE_INT      = "int"
	FIELD_TYPE_BOOL     = "bool"
	FIELD_TYPE_DECIMAL  = "decimal"
	FIELD_TYPE_TEXT     = "text"
	FIELD_TYPE_DATE     = "date"
	FIELD_TYPE_DATETIME = "datetime"
	FIELD_TYPE_XML      = "xml"
	FIELD_TYPE_JSON     = "json"
	FIELD_TYPE_ENUM     = "enum"
	FIELD_TYPE_ENTITY   = "entity"
)
