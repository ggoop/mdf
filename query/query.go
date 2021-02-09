package query

import (
	"github.com/ggoop/mdf/md"
	"github.com/ggoop/mdf/utils"
)

const md_domain string = "query"
const (
	//参照
	WHERE_TYPE_REF = "ref"
	// 枚举
	WHERE_TYPE_ENUM = "enum"
	//字符
	WHERE_TYPE_STRING = "string"
	//时间
	WHERE_TYPE_DATE = "date"
	//时间
	WHERE_TYPE_DATETIME = "datetime"
	//布尔
	WHERE_TYPE_BOOL = "bool"
	//数值
	WHERE_TYPE_NUMBER = "number"
	//选择项
	WHERE_TYPE_SELECT = "select"
)

type Query struct {
	md.Model
	ScopeID     string        `gorm:"size:50;name:范围" json:"scope_id"`
	ScopeType   string        `gorm:"size:50;name:范围" json:"scope_type"`
	Code        string        `gorm:"size:50;name:编码" json:"code"`
	Name        string        `gorm:"name:名称"  json:"name"`
	Type        string        `gorm:"size:50;name:查询类型"  json:"type"` //entity:实体,service:服务,sql:SQL脚本
	Entry       string        `gorm:"size:200;name:入口"  json:"entry"`
	Memo        string        `gorm:"name:备注"  json:"memo"`
	PageSize    int           `gorm:"default:30;name:每页显示记录数" json:"page_size"`
	Columns     []QueryColumn `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:栏目集合" json:"columns"`
	Orders      []QueryOrder  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:排序集合" json:"orders"`
	Wheres      []QueryWhere  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:条件集合" json:"wheres"`
	Filters     []QueryFilter `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:QueryID;name:过滤集合" json:"filters"`
	ContextJson string        `gorm:"name:上下文" json:"context_json"`      //上下文参数
	Condition   string        `gorm:"size:200;name:条件" json:"condition"` //条件
}

func (s *Query) MD() *md.Mder {
	return &md.Mder{ID: "query.query", Domain: md_domain, Name: "查询"}
}

type QueryColumn struct {
	md.Model
	QueryID   string      `gorm:"size:50;name:查询" json:"query_id"`
	OwnerID   string      `gorm:"size:50;name:所属ID" json:"owner_id"`
	OwnerType string      `gorm:"size:50;name:所属类型" json:"owner_type"`
	Field     string      `gorm:"size:50;name:字段" json:"field"`
	Expr      string      `gorm:"size:200;name:表达式" json:"expr"` //字段表达式
	Name      string      `gorm:"size:50;name:栏目编码" json:"name"`
	Title     string      `gorm:"name:显示名称" json:"title"`
	Width     string      `gorm:"name:宽度" json:"width"`
	Tags      string      `gorm:"name:标签" json:"tags"`                //is_q
	DataType  string      `gorm:"size:50;name:字段类型" json:"data_type"` //sys.query.data.type，string,enum,entity,bool,datetime
	KeyField  utils.SBool `gorm:"default:false;name:是主键字段" json:"key_field"`
	CodeField utils.SBool `gorm:"default:false;name:是编码字段" json:"code_field"`
	NameField utils.SBool `gorm:"default:false;name:是名称字段" json:"name_field"`
	Sequence  int         `gorm:"name:顺序" json:"sequence"`
	Fixed     utils.SBool `gorm:"default:false;name:固定" json:"fixed"`
	IsDefault utils.SBool `gorm:"default:false;name:默认" json:"is_default"`
	Hidden    utils.SBool `gorm:"default:false;name:隐藏" json:"hidden"`
	Enabled   utils.SBool `gorm:"name:启用;default:true" json:"enabled"`
}

func (s *QueryColumn) MD() *md.Mder {
	return &md.Mder{ID: "query.column", Domain: md_domain, Name: "查询栏目"}
}

type QueryOrder struct {
	md.Model
	QueryID   string      `gorm:"size:50;name:查询" json:"query_id"`
	OwnerID   string      `gorm:"size:50;name:所属ID" json:"owner_id"`
	OwnerType string      `gorm:"size:50;name:所属类型" json:"owner_type"`
	Field     string      `gorm:"size:50;name:字段"  json:"field"`
	Expr      string      `gorm:"size:200;name:表达式" json:"expr"` //字段表达式
	Title     string      `gorm:"name:显示名称" json:"title"`
	Order     string      `gorm:"name:排序方式"  json:"order"` //desc,asc
	Fixed     utils.SBool `gorm:"default:false;name:固定" json:"fixed"`
	Enabled   utils.SBool `gorm:"default:true;name:启用" json:"enabled"`
	IsDefault utils.SBool `gorm:"default:false;name:默认" json:"is_default"`
	Hidden    utils.SBool `gorm:"default:false;name:隐藏" json:"hidden"`
	Sequence  int         `gorm:"default:0;name:顺序"  json:"sequence"`
}

func (s *QueryOrder) MD() *md.Mder {
	return &md.Mder{ID: "query.order", Domain: md_domain, Name: "查询排序"}
}

type QueryWhere struct {
	md.Model
	QueryID    string        `gorm:"size:50;name:查询" json:"query_id"`
	OwnerID    string        `gorm:"size:50;name:所属ID" json:"owner_id"`
	OwnerType  string        `gorm:"size:50;name:所属类型" json:"owner_type"`
	ParentID   string        `gorm:"size:50" json:"parent_id"`
	Logical    string        `gorm:"size:10;name:逻辑" json:"logical"` //and or
	Field      string        `gorm:"size:50;name:字段" json:"field"`   //如果#开始，则表示标记，需要使用表达式
	Expr       string        `gorm:"size:200;name:表达式" json:"expr"`
	Title      string        `gorm:"name:显示名称" json:"title"`
	DataType   string        `gorm:"size:50;name:字段类型" json:"data_type"` //sys.query.data.type，string,enum,entity,bool,datetime
	DataSource string        `gorm:"name:数据来源" json:"data_source"`
	Operator   string        `gorm:"size:50;name:操作符号" json:"operator"`
	Value      utils.SJson   `gorm:"name:值" json:"value"`
	Sequence   int           `gorm:"name:顺序" json:"sequence"`
	Fixed      utils.SBool   `gorm:"default:false;name:固定" json:"fixed"`
	IsDefault  utils.SBool   `gorm:"default:false;name:默认" json:"is_default"`
	IsBasic    utils.SBool   `gorm:"default:false;name:常用条件" json:"is_basic"`
	Enabled    utils.SBool   `gorm:"default:true;name:启用" json:"enabled"`
	Hidden     utils.SBool   `gorm:"default:false;name:隐藏" json:"hidden"`
	Filters    []QueryFilter `gorm:"-" json:"filters"`
	Children   []QueryWhere  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:ParentID;name:子条件" json:"children"`
}

func (s *QueryWhere) MD() *md.Mder {
	return &md.Mder{ID: "query.where", Domain: md_domain, Name: "查询条件"}
}

type QueryFilter struct {
	md.Model
	QueryID    string      `gorm:"size:50;name:查询ID" json:"query_id"`    //查询
	OwnerField string      `gorm:"size:50;name:所属字段" json:"owner_field"` //字段名称
	Field      string      `gorm:"size:50;name:字段" json:"field"`         //如果#开始，则表示标记，需要使用表达式
	Expr       string      `gorm:"size:200;name:表达式" json:"expr"`
	Title      string      `gorm:"name:显示名称" json:"title"`
	DataType   string      `gorm:"size:50;name:字段类型" json:"data_type"` //string,bool,datetime
	DataSource string      `gorm:"name:数据来源" json:"data_source"`       //条件：where| 常量：const| 变量：var
	Operator   string      `gorm:"size:50;name:操作符号" json:"operator"`
	Value      utils.SJson `gorm:"name:值" json:"value"`
	Sequence   int         `gorm:"name:顺序" json:"sequence"`
	Enabled    utils.SBool `gorm:"default:true;name:启用" json:"enabled"`
}

func (s *QueryFilter) MD() *md.Mder {
	return &md.Mder{ID: "query.filter", Domain: md_domain, Name: "查询过滤"}
}
