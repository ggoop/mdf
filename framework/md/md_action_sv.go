package md

import (
	"fmt"
	"strings"

	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

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

/**
规则通用接口
*/
type IActionRule interface {
	Exec(token *utils.TokenContext, req *ReqContext, res *ResContext) error
	GetRule() RuleRegister
}

var _action_rule = initActionRule()

//基础规则
type SimpleRule struct {
	Rule   RuleRegister
	Handle func(token *utils.TokenContext, req *ReqContext, res *ResContext) error
}

func (s *SimpleRule) Exec(token *utils.TokenContext, req *ReqContext, res *ResContext) error {
	if s.Handle != nil {
		return s.Handle(token, req, res)
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
func DoAction(token *utils.TokenContext, req *ReqContext) *ResContext {
	res := &ResContext{}
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
	} else {
		if rule, ok := GetActionRule(RuleRegister{Domain: req.Domain, OwnerType: req.OwnerType, OwnerID: req.OwnerID, Code: req.Action}); ok {
			rules = append(rules, rule)
		}
	}
	if len(rules) == 0 {
		res.SetError(glog.Error("没有找到任何规则可执行！"))
	}
	for _, rule := range rules {
		if err := rule.Exec(token, req, res); err != nil {
			res.SetError(err)
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
	return fmt.Sprintf("%s:%s", strings.ToLower(s.OwnerType), strings.ToLower(s.Code))
}
