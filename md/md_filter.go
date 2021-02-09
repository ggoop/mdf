package md

import "github.com/ggoop/mdf/utils"

type MDFilters struct {
	ID        string      `gorm:"primary_key;size:36" json:"id"`
	CreatedAt utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	EntID     string      `gorm:"size:36;name:企业" json:"ent_id"`
	Code      string      `gorm:"size:36;name:编码" json:"code"`
	Name      string      `gorm:"size:50;name:名称" json:"name"`
	DsType    string      `gorm:"size:36;name:数据源类型" json:"ds_type"`
	DsEntry   string      `gorm:"size:36;name:数据源实体" json:"ds_entry"`
	AutoLoad  utils.SBool `gorm:"name:自动加载" json:"auto_load"`
	PageSize  int         `gorm:"size:3;name:每页大小" json:"page_size"`
	Condition string      `gorm:"type:text;name:条件" json:"condition"`
	Context   utils.SJson `gorm:"type:text;name:上下文" json:"context"`
	System    utils.SBool `gorm:"name:系统的" json:"system"`
}

func (s *MDFilters) MD() *Mder {
	return &Mder{ID: "md.filters", Domain: md_domain, Name: "过滤器"}
}

type MDFilterSolution struct {
	ID        string      `gorm:"primary_key;size:36" json:"id"`
	CreatedAt utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	EntID     string      `gorm:"size:36;name:企业" json:"ent_id"`
	FilterID  string      `gorm:"size:36;name:过滤器ID" json:"filter_id"`
	Code      string      `gorm:"size:36;name:编码" json:"code"`
	Name      string      `gorm:"size:50;name:名称" json:"name"`
	AutoLoad  utils.SBool `gorm:"name:自动加载" json:"auto_load"`
	PageSize  int         `gorm:"size:3;name:每页大小" json:"page_size"`
	Context   utils.SJson `gorm:"type:text;name:上下文" json:"context"`
	System    utils.SBool `gorm:"name:系统的" json:"system"`
}

func (s *MDFilterSolution) MD() *Mder {
	return &Mder{ID: "md.filter.solution", Domain: md_domain, Name: "过滤方案"}
}

type MDFilterItem struct {
	ID          string      `gorm:"primary_key;size:36" json:"id"`
	CreatedAt   utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt   utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	EntID       string      `gorm:"size:36;name:企业" json:"ent_id"`
	FilterID    string      `gorm:"size:36;name:过滤器ID" json:"filter_id"`
	SolutionID  string      `gorm:"size:36;name:过滤方案ID" json:"solution_id"`
	ParentID    string      `gorm:"size:36;name:上级" json:"parent_id"`
	Code        string      `gorm:"size:36;name:编码" json:"code"`
	Name        string      `gorm:"size:50;name:名称" json:"name"`
	Type        string      `gorm:"size:20;name:类型" json:"type"`
	Caption     string      `gorm:"size:50;name:标题" json:"caption"`
	DsType      string      `gorm:"size:36;name:数据源类型" json:"ds_type"`
	DsEntry     string      `gorm:"size:36;name:数据源实体" json:"ds_entry"`
	DsField     string      `gorm:"size:36;name:数据源字段" json:"ds_field"`
	RefType     string      `gorm:"size:36;name:参照类型" json:"ref_type"`
	RefCode     string      `gorm:"size:36;name:参照编码" json:"ref_code"`
	RefReturn   string      `gorm:"size:36;name:参照返回" json:"ref_return"`
	RefFilter   string      `gorm:"type:text;name:参照查询条件" json:"ref_filter"`
	Sequence    int         `gorm:"size:3;name:顺序" json:"sequence"`
	Operator    string      `gorm:"size:36;name:操作符号" json:"operator"`
	Logic       string      `gorm:"size:36;name:逻辑比较" json:"logic"`
	Value1      utils.SJson `gorm:"type:text;name:值1" json:"value1"`
	Value2      utils.SJson `gorm:"type:text;name:值2" json:"value2"`
	Precision   int         `gorm:"size:3;name:精度" json:"precision"`
	Format      string      `gorm:"size:36;name:格式化" json:"format"`
	Placeholder string      `gorm:"size:50;name:占位符" json:"placeholder"`
	Length      int         `gorm:"size:3;name:长度" json:"length"`
	Hidden      utils.SBool `gorm:"name:隐藏" json:"hidden"`
	Multiple    utils.SBool `gorm:"name:多选" json:"multiple"`
	Nullable    utils.SBool `gorm:"name:可空" json:"nullable"`
	Editable    utils.SBool `gorm:"name:可编辑" json:"editable"`
	Fixed       utils.SBool `gorm:"name:固定的" json:"fixed"`
	Extras      utils.SJson `gorm:"type:text;name:扩展属性" json:"extras"` //JSON
	System      utils.SBool `gorm:"name:系统的" json:"system"`
}

func (s *MDFilterItem) MD() *Mder {
	return &Mder{ID: "md.filter.item", Domain: md_domain, Name: "过滤项"}
}
