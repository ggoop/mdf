package md

import (
	"github.com/ggoop/mdf/utils"
)

/**
枚举类型
*/
type MDEnumType struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	Name      string     `gorm:"size:100"`
	Domain    string     `gorm:"size:50" json:"domain"`
	Enums     []MDEnum   `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:EntityID"`
}

/**
枚举值
*/
type MDEnum struct {
	EntityID  string     `gorm:"size:50;primary_key:uix;morph:limit" json:"entity_id"`
	ID        string     `gorm:"size:50;primary_key:uix" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	Name      string     `gorm:"size:50" json:"name"`
	Sequence  int        `json:"sequence"`
	SrcID     string     `gorm:"size:50" json:"src_id"`
}

func (t MDEnum) TableName() string {
	return "md_enums"
}
func (s *MDEnum) MD() *Mder {
	return &Mder{ID: "md.enum", Domain: md_domain, Name: "枚举", Type: utils.TYPE_ENUM}
}
