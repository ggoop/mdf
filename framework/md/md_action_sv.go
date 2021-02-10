package md

import (
	"fmt"
	"strings"

	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

/**
请求通用类
*/
type ReqContext struct {
	OwnerID   string                 `json:"owner_id"  form:"owner_id"`
	OwnerType string                 `json:"owner_type"  form:"owner_type"`
	Domain    string                 `json:"domain" form:"domain"`
	Widget    string                 `json:"widget" form:"widget"`  //页面ID
	Params    map[string]interface{} `json:"params"  form:"params"` //一般指 页面 URI 参数
	ID        string                 `json:"id" form:"id"`
	IDS       []string               `json:"ids" form:"ids"`
	UserID    string                 `json:"user_id" form:"user_id"` //用户ID
	EntID     string                 `json:"ent_id" form:"ent_id"`   //企业ID
	OrgID     string                 `json:"org_id" form:"org_id"`   //组织ID
	Page      int                    `json:"page" form:"page"`
	PageSize  int                    `json:"page_size" form:"page_size"`
	Action    string                 `json:"action" form:"action"` // 动作编码
	Rule      string                 `json:"rule" form:"rule"`     //规则编码
	URI       string                 `json:"uri" form:"uri"`
	Method    string                 `json:"method" form:"method"`
	Q         string                 `json:"q" form:"q"`                 //模糊查询条件
	Condition interface{}            `json:"condition" form:"condition"` //附件条件
	Entity    string                 `json:"entity" form:"entity"`
	Data      interface{}            `json:"data" form:"data"` //数据
	Tag       string                 `json:"tag" form:"tag"`
	Orders    []QueryOrder           `json:"orders" form:"orders"`
	Wheres    []QueryWhere           `json:"wheres" form:"wheres"`
	Columns   []QueryColumn          `json:"columns" form:"columns"`
}
type ResContext struct {
	Data  utils.Map
	Error error
}

func (s *ResContext) SetData(name string, value interface{}) *ResContext {
	if s.Data == nil {
		s.Data = utils.Map{}
	}
	s.Data[name] = value
	return s
}

func (s ResContext) New() ResContext {
	return ResContext{}
}

type PageViewDTO struct {
	Data       interface{} `json:"data"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	EntityID   string      `json:"entity_id"`
	PrimaryKey string      `json:"primary_key"`
	Multiple   utils.SBool `json:"multiple"`
	Nullable   utils.SBool `json:"nullable"`
	IsMain     utils.SBool `json:"is_main"`
}

func (s ReqContext) New() ReqContext {
	return ReqContext{
		UserID: s.UserID, EntID: s.EntID, OrgID: s.OrgID,
	}
}
func (s ReqContext) Copy() ReqContext {
	return ReqContext{
		OwnerType: s.OwnerType, OwnerID: s.OwnerID, Domain: s.Domain,
		ID: s.ID, IDS: s.IDS, UserID: s.UserID, EntID: s.EntID, OrgID: s.OrgID,
		Widget: s.Widget,
		Page:   s.Page, PageSize: s.PageSize, Q: s.Q,
		Action: s.Action, Rule: s.Rule,
		URI: s.URI, Method: s.Method,
		Condition: s.Condition, Entity: s.Entity, Data: s.Data,
	}
}

/**
规则通用接口
*/
type IActionRule interface {
	Exec(context *ReqContext, res *ResContext) error
	GetRule() RuleRegister
}

var _action_rule = initActionRule()

//基础规则
type SimpleRule struct {
	Rule   RuleRegister
	Handle func(context *ReqContext, res *ResContext) error
}

func (s *SimpleRule) Exec(context *ReqContext, res *ResContext) error {
	if s.Handle != nil {
		return s.Handle(context, res)
	}
	return glog.Error("没有实现")
}

func (s *SimpleRule) GetRule() RuleRegister {
	return s.Rule
}

func initActionRule() map[string]IActionRule {
	return make(map[string]IActionRule)
}

func GetActionRule(reg RuleRegister) (IActionRule, bool) {
	if r, ok := _action_rule[reg.GetKey()]; ok {
		return r, ok
	}
	return nil, false
}

//执行命令
func DoAction(req ReqContext) ResContext {
	res := ResContext{}
	rules := make([]IActionRule, 0)
	ruleCodes := []string{"common", req.OwnerID}
	ruleDatas := make([]MDActionRule, 0)
	//查询拥有者规则
	repositories.Default().
		Model(MDActionRule{}).
		Order("sequence,code").
		Where("owner_type=? and owner_id in (?) and action=? and enabled=1", req.OwnerType, ruleCodes, req.Action).
		Take(&ruleDatas)

	if len(ruleDatas) > 0 {
		replacedList := make(map[string]MDActionRule)
		for _, r := range ruleDatas {
			if r.Replaced != "" {
				replacedList[strings.ToLower(r.Replaced)] = r
			}
		}
		for _, r := range ruleDatas {
			if replaced, ok := replacedList[strings.ToLower(r.Code)]; ok {
				glog.Error("规则被替换", glog.Any("replaced", replaced.Code))
				continue
			}
			if rule, ok := GetActionRule(RuleRegister{Domain: r.Domain, OwnerType: r.OwnerType, OwnerID: r.OwnerID, Code: r.Code}); ok {
				rules = append(rules, rule)
			} else {
				glog.Error("找不到规则", glog.Any("rule", r))
			}
		}
	}
	if len(rules) == 0 {
		res.Error = glog.Error("没有找到任何规则可执行！")
	}
	for _, rule := range rules {
		if err := rule.Exec(&req, &res); err != nil {
			res.Error = err
			return res
		}
	}
	return res
}
func RegisterActionRule(rules ...IActionRule) {
	if len(rules) > 0 {
		for i, _ := range rules {
			rule := rules[i]
			_action_rule[rule.GetRule().GetKey()] = rule
		}
	}
}

// 注册器
type RuleRegister struct {
	Domain    string //领域：
	Code      string //规则编码：save,delete,query,find
	OwnerID   string //规则拥有者：common,widgetID,entityID
	OwnerType string //规则拥有者类型：widget,entity
}

const (
	RuleType_Widget string = "widget"
	RuleType_Entity string = "entity"
)

func (s RuleRegister) GetKey() string {
	return fmt.Sprintf("%s", strings.ToLower(s.Code))
}
