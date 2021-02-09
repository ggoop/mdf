package model

import "github.com/ggoop/mdf/framework/md"

type Blacklist struct {
	md.Model
	EntID      string `gorm:"size:50"`
	Action     string `gorm:"size:50" json:"action"` //login,applogin
	TargetKey  string `gorm:"size:50" json:"target_key"`
	TargetName string `json:"target_name"`
	Memo       string `json:"memo"`
}

func (s Blacklist) TableName() string {
	return "sys_blacklists"
}
func (s *Blacklist) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".blacklist", Domain: MD_DOMAIN, Name: "黑名单"}
}
