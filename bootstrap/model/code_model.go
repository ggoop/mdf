package model

import "github.com/ggoop/mdf/framework/md"

type CodeRule struct {
	md.Model
	Name       string `gorm:"size:50"`
	Tag        string `gorm:"size:50;unique_index:udx_tag"`
	Memo       string `gorm:"size:50"`
	Prefix     string `gorm:"size:50"`
	Suffix     string `gorm:"size:50"`
	TimeFormat string `gorm:"size:50"` //yyyy,yy,mm,dd
	SeqLength  int    `gorm:"default:4;name:长度"`
	SeqStep    int    `gorm:"default:1;name:步长"`
}

func (t CodeRule) TableName() string {
	return "sys_code_rules"
}
func (s *CodeRule) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".code.rule", Domain: MD_DOMAIN, Name: "编码规则"}
}

type CodeValue struct {
	md.Model
	RuleID    string `gorm:"size:50;index:idx_tag"`
	EntID     string `gorm:"size:50;index:idx_tag"`
	TimeValue string `gorm:"size:50;index:idx_tag"`
	Code      string `gorm:"size:50"`
	SeqValue  int    `gorm:"default:0;name:序号"`
}

func (t CodeValue) TableName() string {
	return "sys_code_values"
}
func (s *CodeValue) MD() *md.Mder {
	return &md.Mder{ID: MD_DOMAIN + ".code.value", Domain: MD_DOMAIN, Name: "编码数据"}
}
