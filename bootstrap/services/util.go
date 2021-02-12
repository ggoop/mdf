package services

import (
	"github.com/ggoop/mdf/utils"
)

type synchEntDTO struct {
	ID       string `json:"id"`
	Openid   string `json:"openid"`
	Name     string `json:"name"`
	TypeID   string `json:"type_id"`
	StatusID string `json:"status_id"`
}
type synchPersonDTO struct {
	TypeID string       `json:"type_id"`
	User   synchUserDTO `json:"user"`
}
type synchUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Openid   string `json:"openid"`
	Password string `json:"password"`
}
type synchProductDTO struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
type synchProductHostDTO struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type synchProductModuleDTO struct {
	Code    string           `json:"code"`
	Name    string           `json:"name"`
	Product *synchProductDTO `json:"product"`
}
type synchProductServiceDTO struct {
	Code      string                 `json:"code"`
	Name      string                 `json:"name"`
	Icon      string                 `json:"icon"`
	Uri       string                 `json:"uri"`
	AppUri    string                 `json:"app_uri"`
	InWeb     utils.SBool            `json:"in_web"`
	InApp     utils.SBool            `json:"in_app"`
	Schema    string                 `json:"schema"`
	IsMaster  utils.SBool            `json:"is_master"`
	IsSlave   utils.SBool            `json:"is_slave"`
	IsDefault utils.SBool            `json:"is_default"`
	Sequence  int                    `json:"sequence"`
	BizType   string                 `json:"biz_type"`
	Memo      string                 `json:"memo"`
	Tags      string                 `json:"tags"`
	Product   *synchProductDTO       `json:"product"`
	Host      *synchProductHostDTO   `json:"host"`
	Module    *synchProductModuleDTO `json:"module"`
}

type synchMDPageDTO struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"` //page，ref，app
	Domain     string      `json:"domain"`
	EntID      string      `json:"ent_id"`
	Code       string      `json:"code"`
	Element    string      `json:"element"`
	Name       string      `json:"name"`
	Widgets    utils.SJson `json:"widgets"` //JSON
	MainEntity string      `json:"main_entity"`
	System     utils.SBool `json:"system"`
}

type synchMDRuleDTO struct {
	ID     string      `json:"id"`     //领域.规则：md.save，ui.save
	Domain string      `json:"domain"` //common为公共动作
	Code   string      `json:"code"`
	Name   string      `json:"name"`
	Async  utils.SBool `json:"async"`
	System utils.SBool `json:"system"`
}

type synchMDCommandDTO struct {
	ID      string      `json:"id"`      //save,delete
	PageID  string      `json:"page_id"` //common为公共动作
	Code    string      `json:"code"`
	Name    string      `json:"name"`
	Type    string      `json:"type"` //ui,sv
	Path    string      `json:"path"`
	Content string      `json:"content"` //js语法
	Rules   string      `json:"rules"`
	System  utils.SBool `json:"system"`
}

type resBodyDTO struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
