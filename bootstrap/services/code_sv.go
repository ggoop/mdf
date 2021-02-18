package services

import (
	"fmt"
	"github.com/ggoop/mdf/db"
	"math"
	"strings"

	"github.com/ggoop/mdf/bootstrap/model"
	"github.com/ggoop/mdf/utils"
)

type ICodeSv interface {
}

func CodeSv() ICodeSv {
	return _cacheSvInstance
}

type codeSvImpl struct {
}

var _cacheSvInstance ICodeSv = newCodeSvImpl()

func newCodeSvImpl() *codeSvImpl {
	return &codeSvImpl{}
}
func (s *codeSvImpl) ApplyCode(mdID, entID string) string {
	rule := model.CodeRule{}
	db.Default().First(&rule, "tag=?", mdID)
	if rule.ID == "" {
		//如果没有配置规则，又调用此编码服务，则自动生成编码规则
		rule = model.CodeRule{Tag: mdID, Name: "自动生成-" + mdID, Memo: "系统自动生成", TimeFormat: "yyyyMMdd", SeqLength: 4, SeqStep: 1}
		rule.ID = utils.GUID()
		db.Default().Create(&rule)
	}
	timeValue := ""
	if rule.TimeFormat != "" {
		//yyyy,yy,mm,dd=>2006,06,01,02
		rule.TimeFormat = strings.ToLower(rule.TimeFormat)
		rule.TimeFormat = strings.ReplaceAll(rule.TimeFormat, "yyyy", "2006")
		rule.TimeFormat = strings.ReplaceAll(rule.TimeFormat, "yy", "06")
		rule.TimeFormat = strings.ReplaceAll(rule.TimeFormat, "mm", "01")
		rule.TimeFormat = strings.ReplaceAll(rule.TimeFormat, "dd", "02")
		timeValue = utils.TimeNow().Format(rule.TimeFormat)
	}
	lastCode := model.CodeValue{}
	newCode := model.CodeValue{EntID: entID, RuleID: rule.ID, TimeValue: timeValue}
	db.Default().Last(&lastCode, "rule_id=? and ent_id=? and time_value=?", mdID, entID, timeValue)
	if lastCode.ID == "" {
		if rule.SeqLength > 0 {
			newCode.SeqValue = utils.ToInt(math.Pow10(rule.SeqLength - 1))
		}
	} else {
		newCode.SeqValue = lastCode.SeqValue + rule.SeqStep
	}
	newCode.Code = fmt.Sprintf("%s%s%d%s", rule.Prefix, newCode.TimeValue, newCode.SeqValue, rule.Suffix)
	db.Default().Create(&newCode)
	return newCode.Code
}
