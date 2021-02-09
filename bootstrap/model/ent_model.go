package model

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

/**
团队/企业
*/
type Ent struct {
	md.Model
	Name        string      `gorm:"size:50" json:"name"`
	Memo        string      `json:"memo"`
	Openid      string      `gorm:"size:50" json:"openid"`
	Token       string      `gorm:"size:50" json:"token"`
	Gateway     string      `gorm:"size:50" json:"gateway"`                //服务网关
	CanInvited  utils.SBool `gorm:"not null;default:0" json:"can_invited"` //是否可邀请
	Developer   utils.SBool `gorm:"not null;default:0;name:开发者" json:"developer"`
	StatusID    string      `gorm:"size:50" json:"status_id"`
	Distributor string      `gorm:"size:50" json:"distributor"` //经销商
	TypeID      string      `gorm:"size:50" json:"type_id"`     //类型,演示demo,测试test,开发dev，其它为正式
	Type        *md.MDEnum  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;limit:sys.ent.type"`
}

func (t Ent) TableName() string {
	return "sys_ents"
}
func (s *Ent) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".ent", Domain: MD_DOMAIN, Name: "团队/企业"}
}

/**
成员
*/
type EntUser struct {
	md.Model
	EntID     string      `gorm:"size:50" json:"ent_id"`
	Ent       *Ent        `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	UserID    string      `gorm:"size:50" json:"user_id"`
	User      *User       `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	TypeID    string      `gorm:"size:50;limit:sys.ent.user.type" json:"type_id"`
	UserName  string      `gorm:"size:50" json:"user_name"`
	Enabled   utils.SBool `gorm:"not null;default:1;name:启用" json:"enabled"`
	IsDefault utils.SBool `gorm:"not null;default:1;name:默认" json:"is_default"`
}

func (t EntUser) TableName() string {
	return "sys_ent_users"
}
func (s *EntUser) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".ent.user", Domain: MD_DOMAIN, Name: "企业用户"}
}
