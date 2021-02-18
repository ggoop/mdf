package md

import (
	"fmt"
	"github.com/ggoop/mdf/db"
	"strings"

	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

const (
	RuleType_Widget string = "widget"
	RuleType_Entity string = "entity"
)

// 注册器
type RuleRegister struct {
	Domain    string //领域：
	Code      string //规则编码：save,delete,query,find
	OwnerID   string //规则拥有者：common,widgetID,entityID
	OwnerType string //规则拥有者类型：widget,entity
}

func (s RuleRegister) GetKey() string {
	return fmt.Sprintf("%s:%s", strings.ToLower(s.OwnerType), strings.ToLower(s.Code))
}

/**
规则通用接口
*/
type IActionRule interface {
	Exec(token *utils.TokenContext, req *ReqContext, res *ResContext)
	Register() RuleRegister
}
type IActionSv interface {
	DoAction(token *utils.TokenContext, req *ReqContext) *ResContext
	RegisterRule(rules ...IActionRule)
	RegisterAction(rules ...IActionRule)
}

func ActionSv() IActionSv {
	return actionSv
}

type actionSvImpl struct {
	rules   map[string]IActionRule
	actions map[string]IActionRule
}

var actionSv IActionSv = newActionSvImpl()

func newActionSvImpl() *actionSvImpl {
	return &actionSvImpl{rules: make(map[string]IActionRule), actions: make(map[string]IActionRule)}
}
func (s actionSvImpl) GetRule(reg RuleRegister) (IActionRule, bool) {
	if r, ok := s.rules[reg.GetKey()]; ok {
		return r, ok
	}
	return nil, false
}

func (s actionSvImpl) GetAction(reg RuleRegister) (IActionRule, bool) {
	if r, ok := s.actions[reg.GetKey()]; ok {
		return r, ok
	}
	return nil, false
}

//执行命令
func (s actionSvImpl) DoAction(token *utils.TokenContext, req *ReqContext) *ResContext {
	res := &ResContext{}

	// 查找动作执行
	var action IActionRule
	if a, ok := s.GetAction(RuleRegister{OwnerType: req.OwnerType, OwnerID: req.OwnerID, Code: req.Action}); ok {
		action = a
	}
	if action != nil {
		if action.Exec(token, req, res); res.Error() != nil {
			return res
		}
	}
	//执行规则集合
	rules := make([]IActionRule, 0)
	ownerIds := []string{"common", req.OwnerID}
	ruleDatas := make([]MDActionRule, 0)
	//查询拥有者规则
	db.Default().Order("sequence,code").
		Where("owner_type=? and owner_id in (?) and action=? and enabled=1", req.OwnerType, ownerIds, req.Action).
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
			if rule, ok := s.GetRule(RuleRegister{Domain: r.Domain, OwnerType: r.OwnerType, OwnerID: r.OwnerID, Code: r.Code}); ok {
				rules = append(rules, rule)
			} else {
				glog.Error("找不到规则", glog.Any("rule", r))
			}
		}
	} else {
		if rule, ok := s.GetRule(RuleRegister{Domain: req.Domain, OwnerType: req.OwnerType, OwnerID: req.OwnerID, Code: req.Action}); ok {
			rules = append(rules, rule)
		}
	}
	if len(rules) == 0 {
		res.SetError(glog.Error("没有找到任何规则可执行！"))
	} else {
		for _, rule := range rules {
			if rule.Exec(token, req, res); res.Error() != nil {
				return res
			}
		}
	}
	return res
}
func (s actionSvImpl) RegisterRule(rules ...IActionRule) {
	if len(rules) > 0 {
		for i, _ := range rules {
			rule := rules[i]
			s.rules[rule.Register().GetKey()] = rule
		}
	}
}
func (s actionSvImpl) RegisterAction(rules ...IActionRule) {
	if len(rules) > 0 {
		for i, _ := range rules {
			rule := rules[i]
			s.actions[rule.Register().GetKey()] = rule
		}
	}
}
