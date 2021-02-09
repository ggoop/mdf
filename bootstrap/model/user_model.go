package model

import (
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

/**
用户
*/
type User struct {
	md.Model
	Openid    string      `gorm:"size:50" json:"openid"`
	Mobile    string      `gorm:"size:50" json:"mobile"`
	Email     string      `gorm:"size:50" json:"email"`
	Account   string      `gorm:"size:50" json:"account"`
	Password  string      `json:"password"`
	Name      string      `gorm:"size:50" json:"name"`
	AvatarUrl string      `json:"avatar_url"`
	Memo      string      `json:"memo"`
	Token     string      `gorm:"size:50" json:"token"`
	IsSystem  utils.SBool `gorm:"not null;default:0;name:系统的" json:"is_system"`
	Enabled   utils.SBool `gorm:"not null;default:1;name:启用" json:"enabled"`
}

func (t User) TableName() string {
	return "sys_users"
}

func (s *User) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".user", Domain: MD_DOMAIN, Name: "用户"}
}

/**
用户收藏内容
*/
type UserFavorite struct {
	md.Model
	EntID    string `gorm:"size:50" json:"ent_id"`
	UserID   string `gorm:"size:50" json:"user_id"`
	DataID   string `gorm:"size:50" json:"data_id"`
	DataCode string `gorm:"size:50" json:"data_code"`
	DataType string `gorm:"size:50" json:"data_type"`
	Name     string `gorm:"size:50" json:"name"`
}

func (t UserFavorite) TableName() string {
	return "sys_user_favorites"
}
func (s *UserFavorite) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".user.favorite", Domain: MD_DOMAIN, Name: "用户收藏"}
}
