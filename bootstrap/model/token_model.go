package model

import (
	"encoding/json"

	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/framework/md"
	"github.com/ggoop/mdf/utils"
)

type Token struct {
	md.Model
	ClientID  string `gorm:"size:50"`
	EntID     string `gorm:"size:50;index:idx_type"`
	EntName   string `gorm:"size:50"`
	UserID    string `gorm:"size:50"`
	UserName  string `gorm:"size:50"`
	ProductID string `gorm:"size:50"`
	ServiceID string `gorm:"size:50"`
	Token     string `gorm:"size:100;index:idx_token"`
	Name      string
	Type      string `gorm:"size:50;index:idx_type"`
	Content   string
	Scope     string
	ExpireAt  utils.Time
}

func (t Token) TableName() string {
	return "sys_tokens"
}
func (s *Token) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".token", Domain: MD_DOMAIN, Name: "令牌"}
}

func (s *Token) SetContent(value interface{}) error {
	str, err := json.Marshal(value)
	if err != nil {
		glog.Errorf("error:%v", err)
		return err
	}
	s.Content = string(str)
	return nil
}
func (s *Token) GetContent(value interface{}) error {
	if err := json.Unmarshal([]byte(s.Content), value); err != nil {
		glog.Errorf("parse content error:%v", err)
		return err
	}
	return nil
}

type AuthFree struct {
	md.Model
	SrcOpenid       string      `gorm:"size:50;index:idx_type"`
	SrcType         string      `gorm:"size:50;index:idx_type"`
	SrcApp          string      `gorm:"size:50"`
	Account         string      `gorm:"size:50"`
	SrcAccessToken  string      `gorm:"size:500"`
	SrcRefreshToken string      `gorm:"size:500"`
	UserID          string      `gorm:"size:50"`
	UserToken       utils.SJson `gorm:"size:500"`
	EntID           string      `gorm:"size:50"`
	EntToken        utils.SJson `gorm:"size:500"`
	Qty             int         `json:"qty"`
}

func (t AuthFree) TableName() string {
	return "sys_auth_free"
}
func (s *AuthFree) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".auth.free", Domain: MD_DOMAIN, Name: "令牌"}
}
