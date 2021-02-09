package model

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type Profile struct {
	md.Model
	EntID        string      `gorm:"size:50" json:"ent_id"`
	Code         string      `gorm:"size:50"`
	Name         string      `gorm:"size:50"`
	System       utils.SBool `gorm:"not null;default:0;name:系统的"`
	DefaultValue string
	Value        string
	Memo         string
}

func (t Profile) TableName() string {
	return "sys_profiles"
}
func (s *Profile) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".profile", Domain: MD_DOMAIN, Name: "参数"}
}
