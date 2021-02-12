package services

import (
	"fmt"
	"math"
	"strings"

	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/framework/db/repositories"
	"github.com/ggoop/mdf/framework/di"
	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

type CodeSv struct {
	repo *repositories.MysqlRepo
}

var _defaultCodeSv *CodeSv

func ApplyCode(mdID, entID string) string {
	if _defaultCodeSv == nil {
		if err := di.Global.Invoke(func(s1 *CodeSv) {
			_defaultCodeSv = s1
		}); err != nil {
			glog.Errorf("di Provide error:%s", err)
		}
	}
	if _defaultCodeSv == nil {
		if err := di.Global.Invoke(func(s1 *repositories.MysqlRepo) {
			_defaultCodeSv = NewCodeSv(s1)
		}); err != nil {
			glog.Errorf("di Provide error:%s", err)
		}
	}
	if _defaultCodeSv != nil {
		return _defaultCodeSv.ApplyCode(mdID, entID)
	}
	return ""
}
func NewCodeSv(repo *repositories.MysqlRepo) *CodeSv {
	return &CodeSv{repo: repo}
}
func (s *CodeSv) ApplyCode(mdID, entID string) string {
	rule := model.CodeRule{}
	s.repo.First(&rule, "tag=?", mdID)
	if rule.ID == "" {
		//如果没有配置规则，又调用此编码服务，则自动生成编码规则
		rule = model.CodeRule{Tag: mdID, Name: "自动生成-" + mdID, Memo: "系统自动生成", TimeFormat: "yyyyMMdd", SeqLength: 4, SeqStep: 1}
		rule.ID = utils.GUID()
		s.repo.Create(&rule)
	}
	timeValue := ""
	if rule.TimeFormat != "" {
		//yyyy,yy,mm,dd=>2006,06,01,02
		rule.TimeFormat = strings.ToLower(rule.TimeFormat)
		rule.TimeFormat = strings.ReplaceAll(rule.TimeFormat, "yyyy", "2006")
		rule.TimeFormat = strings.ReplaceAll(rule.TimeFormat, "yy", "06")
		rule.TimeFormat = strings.ReplaceAll(rule.TimeFormat, "mm", "01")
		rule.TimeFormat = strings.ReplaceAll(rule.TimeFormat, "dd", "02")
		timeValue = utils.NewTime().Format(rule.TimeFormat)
	}
	lastCode := model.CodeValue{}
	newCode := model.CodeValue{EntID: entID, RuleID: rule.ID, TimeValue: timeValue}
	s.repo.Last(&lastCode, "rule_id=? and ent_id=? and time_value=?", mdID, entID, timeValue)
	if lastCode.ID == "" {
		if rule.SeqLength > 0 {
			newCode.SeqValue = utils.ToInt(math.Pow10(rule.SeqLength - 1))
		}
	} else {
		newCode.SeqValue = lastCode.SeqValue + rule.SeqStep
	}
	newCode.Code = fmt.Sprintf("%s%s%d%s", rule.Prefix, newCode.TimeValue, newCode.SeqValue, rule.Suffix)
	s.repo.Create(&newCode)
	return newCode.Code
}
