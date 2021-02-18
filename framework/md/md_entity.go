package md

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"

	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

const (
	md_domain string = "md"
)

type MDEntity struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	Type      string     `gorm:"size:50;not null"` // simple，entity，enum，interface，dto,view
	Domain    string     `gorm:"size:50;name:领域" json:"domain"`
	Code      string     `gorm:"size:100;index:code_idx;not null"`
	Name      string     `gorm:"size:100"`
	TableName string     `gorm:"size:50"`
	Memo      string     `gorm:"size:500"`
	Tags      string     `gorm:"size:500"`
	System    utils.SBool
	Fields    []MDField `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:EntityID"`
	cache     map[string]MDField
}

func (s MDEntity) String() string {
	return fmt.Sprintf("%s-%s-%s", s.Domain, s.Code, s.ID)
}
func (s *MDEntity) MD() *Mder {
	return &Mder{ID: "md.entity", Domain: md_domain, Name: "实体"}
}
func (s *MDEntity) GetField(code string) *MDField {
	if s.cache == nil {
		s.cache = make(map[string]MDField)
	}

	if s.Fields != nil && len(s.Fields) > 0 && len(s.cache) == 0 {
		for i, v := range s.Fields {
			s.cache[strings.ToLower(v.Code)] = s.Fields[i]
			if v.DbName != "" {
				s.cache[strings.ToLower(v.DbName)] = s.Fields[i]
			}
		}
	}
	if v, ok := s.cache[strings.ToLower(code)]; ok {
		return &v
	}
	return nil
}

type MDEntityRelation struct {
	ID        string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt utils.Time `gorm:"name:更新时间" json:"updated_at"`
	Code      string     `gorm:"size:100;index:code_idx;not null"`
	Name      string     `gorm:"size:100"`
	ParentID  string     `gorm:"size:50;name:父实体;not null"`
	ChildID   string     `gorm:"size:50;name:子实体;not null"`
	Kind      string     `gorm:"name:类型;not null"` //inherit，interface，
	ParentKey string     `gorm:"size:36;name:主键" json:"parent_key"`
	ChildKey  string     `gorm:"size:36;name:外键" json:"child_key"`
	Limit     string     `gorm:"size:500;name:限制"`
}

func (s *MDEntityRelation) MD() *Mder {
	return &Mder{ID: "md.entity.relation", Domain: md_domain, Name: "实体关系"}
}

type MDField struct {
	ID             string     `gorm:"primary_key;size:50" json:"id"`
	CreatedAt      utils.Time `gorm:"name:创建时间" json:"created_at"`
	UpdatedAt      utils.Time `gorm:"name:更新时间" json:"updated_at"`
	EntityID       string     `gorm:"size:50;unique_index:uix_code;not null"`
	Entity         *MDEntity  `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Code           string     `gorm:"size:50;unique_index:uix_code;not null"`
	Name           string     `gorm:"size:50"`
	DbName         string     `gorm:"size:50"`
	IsNormal       utils.SBool
	IsPrimaryKey   utils.SBool
	ForeignKey     string    `gorm:"size:50"` //外键
	AssociationKey string    `gorm:"size:50"` //Association
	Kind           string    `gorm:"size:50"`
	TypeID         string    `gorm:"size:50"`
	TypeType       string    `gorm:"size:50"`
	Type           *MDEntity `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Limit          string    `gorm:"size:500;name:限制"`
	Memo           string    `gorm:"size:500"`
	Tags           string    `gorm:"size:500"` // code,name,ent,import
	Sequence       int
	Nullable       utils.SBool
	Length         int
	Precision      int
	DefaultValue   string
	MinValue       string
	MaxValue       string
	SrcID          string `gorm:"size:50" json:"src_id"`
}

func (s MDField) String() string {
	return fmt.Sprintf("%s-%s-%s", s.Code, s.Name, s.TypeType)
}
func (s MDField) CompileValue(value interface{}) interface{} {
	if value == nil || value == "" || s.TypeID == "" {
		return nil
	}
	if s.TypeID == utils.FIELD_TYPE_STRING || s.TypeID == utils.FIELD_TYPE_TEXT || s.TypeID == utils.FIELD_TYPE_XML {
		return value
	}
	if s.TypeID == utils.FIELD_TYPE_INT {
		return utils.ToInt(value)
	}
	if s.TypeID == utils.FIELD_TYPE_BOOL {
		return utils.ToSBool(value)
	}
	if s.TypeID == utils.FIELD_TYPE_JSON {
		return utils.ToSBool(value)
	}
	if s.TypeID == utils.FIELD_TYPE_DATE || s.TypeID == utils.FIELD_TYPE_DATETIME {
		return utils.ToTime(value)
	}
	if s.TypeID == utils.FIELD_TYPE_DECIMAL {
		if v, err := decimal.NewFromString(utils.ToString(value)); err != nil {
			return glog.Error(err)
		} else {
			return v
		}
	}
	return nil
}
func (s *MDField) MD() *Mder {
	return &Mder{ID: "md.field", Domain: md_domain, Name: "属性"}
}
