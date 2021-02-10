package md

import "github.com/ggoop/mdf/utils"

type MDActionCommand struct {
	ID        string      `gorm:"primary_key;size:36" json:"id"`
	CreatedAt utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	EntID     string      `gorm:"size:36;name:企业" json:"ent_id"`
	OwnerType string      `gorm:"size:36;name:拥有者类型" json:"owner_type"` //entity,page,ent,domain
	OwnerID   string      `gorm:"size:36;name:拥有者ID" json:"owner_id"`   //common为公共动作
	Code      string      `gorm:"size:36;name:编码" json:"code"`
	Name      string      `gorm:"size:50;name:名称" json:"name"`
	Type      string      `gorm:"size:20;name:类型" json:"type"`
	Action    string      `gorm:"size:50;name:动作" json:"action"`
	Url       string      `gorm:"size:100;name:服务路径" json:"url"`
	Parameter utils.SJson `gorm:"type:text;name:参数" json:"parameter"`
	Method    string      `gorm:"size:20;name:请求方式" json:"method"`
	Target    string      `gorm:"size:36;name:目标" json:"target"`
	Script    string      `gorm:"type:text;name:脚本" json:"script"`
	Enabled   utils.SBool `gorm:"default:true" json:"enabled"`
	System    utils.SBool `gorm:"name:系统的" json:"system"`
}

func (s *MDActionCommand) MD() *Mder {
	return &Mder{ID: "md.action.command", Domain: md_domain, Name: "组件命令"}
}

type MDActionRule struct {
	ID        string      `gorm:"primary_key;size:50" json:"id"` //领域.规则：md.save，ui.save
	CreatedAt utils.Time  `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time  `gorm:"name:更新时间" json:"updated_at"`
	Domain    string      `gorm:"size:36;name:模块" json:"domain"`
	OwnerType string      `gorm:"size:36;name:拥有者类型" json:"owner_type"`
	OwnerID   string      `gorm:"size:36;name:拥有者ID" json:"owner_id"` //common为公共动作
	Code      string      `gorm:"size:50;name:编码" json:"code"`
	Name      string      `gorm:"size:50;name:名称" json:"name"`
	Action    string      `gorm:"size:50;name:动作" json:"action"`
	Url       string      `gorm:"size:100;name:服务路径" json:"url"`
	Sequence  int         `gorm:"size:3;name:顺序" json:"sequence"`
	Replaced  string      `gorm:"size:50;name:被替换的" json:"replaced"`
	Async     utils.SBool `gorm:"name:异步的" json:"async"`
	Enabled   utils.SBool `gorm:"default:true" json:"enabled"`
	System    utils.SBool `gorm:"name:系统的" json:"system"`
}

func (s *MDActionRule) MD() *Mder {
	return &Mder{ID: "md.action.rule", Domain: md_domain, Name: "动作规则"}
}
