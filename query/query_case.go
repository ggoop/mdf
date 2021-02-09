package query

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ggoop/mdf/context"
	"github.com/ggoop/mdf/di"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/md"
	"github.com/ggoop/mdf/repositories"
	"github.com/ggoop/mdf/utils"
)

type QueryCase struct {
	md.Model
	EntID     string          `gorm:"size:50" json:"ent_id"`
	UserID    string          `gorm:"size:50" json:"user_id"`
	QueryID   string          `gorm:"size:50;name:查询" json:"query_id"`
	Query     *Query          `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;name:查询" json:"query"`
	Name      string          `gorm:"name:名称" json:"name"`
	ScopeID   string          `gorm:"size:50;name:范围ID" json:"scope_id"`
	ScopeType string          `gorm:"size:50;name:范围类型" json:"scope_type"`
	Memo      string          `gorm:"name:备注" json:"memo"`
	Page      int             `gorm:"name:页码" json:"page"`
	PageSize  int             `gorm:"name:每页显示记录数" json:"page_size"`
	Columns   []QueryColumn   `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:栏目集合" json:"columns"`
	Orders    []QueryOrder    `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:排序集合" json:"orders"`
	Wheres    []QueryWhere    `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:OwnerID;name:条件集合" json:"wheres"`
	Filters   []QueryFilter   `gorm:"association_autoupdate:false;association_autocreate:false;association_save_reference:false;foreignkey:QueryID;association_foreignkey:QueryID;name:过滤集合" json:"filters"`
	Context   context.Context `gorm:"type:text;name:上下文" json:"context"` //上下文参数
	Export    utils.SBool     `gorm:"name:是否导出" json:"export"`
	IsDefault utils.SBool     `gorm:"name:默认" json:"is_default"`
	Condition string          `gorm:"size:200;name:条件" json:"condition"` //条件
	Q         string          `gorm:"size:200;name:关键字" json:"q"`        //关键字
}

func (s *QueryCase) MD() *md.Mder {
	return &md.Mder{ID: "query.case", Name: "查询方案"}
}
func (s *QueryCase) prepareCase() {
	if s.Query == nil && s.QueryID != "" {
		if err := di.Global.Invoke(func(db *repositories.MysqlRepo) {
			q := Query{}
			if err := db.Preload("Columns").Preload("Orders").Preload("Wheres").Order("id").Take(&q, "id=? or code=?", s.QueryID, s.QueryID).Error; err != nil {
				glog.Errorf("query error:%s", err)
			}
			if q.ID != "" {
				s.Query = &q
			}
		}); err != nil {
			glog.Errorf("di Provide error:%s", err)
		}
	}
	if s.Query != nil {
		s.Condition = s.Query.Condition
	}
	if len(s.Columns) == 0 && s.Query != nil {
		s.Columns = make([]QueryColumn, 0)
		for _, v := range s.Query.Columns {
			if v.IsDefault.IsTrue() {
				s.Columns = append(s.Columns, v)
			}
		}
	}
	//默认栏目名称
	r, _ := regexp.Compile("^[A-Za-z._]$")
	for _, v := range s.Columns {
		if v.Name == "" && v.Field != "" && r.MatchString(v.Field) {
			v.Name = strings.Replace(strings.ToLower(v.Field), ".", "_", -1)
		}
	}
	if len(s.Orders) == 0 && s.Query != nil {
		s.Orders = make([]QueryOrder, 0)
		for _, v := range s.Query.Orders {
			if v.IsDefault.IsTrue() {
				s.Orders = append(s.Orders, v)
			}
		}
	}
	if s.Page <= 0 {
		s.Page = 1
	}
	if s.PageSize <= 0 {
		s.PageSize = s.Query.PageSize
	}
	if s.PageSize <= 0 {
		s.PageSize = 30
	}
}
func (s *QueryCase) GetExector() IExector {
	s.prepareCase()
	if s.Query == nil {
		return nil
	}
	exector := NewExector(s.Query.Entry)
	if s.Page > 0 && s.PageSize > 0 {
		exector.Page(s.Page, s.PageSize)
	}
	if s.Condition != "" {
		exector.Where(s.Condition)
	}
	if s.Q != "" {
		q := "%" + s.Q + "%"
		qwhere := exector.Where("")
		for _, v := range s.Columns {
			if v.CodeField.IsTrue() || v.NameField.IsTrue() || strings.Contains(v.Tags, "quick_filter") {
				if v.Expr != "" {
					qwhere.OrWhere(v.Expr+" like ?", q)
				} else if v.Field != "" {
					qwhere.OrWhere(v.Field+" like ?", q)
				}
			}
		}
	}
	exector.SetContext(&s.Context)
	for _, v := range s.Columns {
		if v.Expr != "" && v.Name != "" {
			exector.Select(v.Expr + " as \"" + v.Name + "\"")
		}
		if v.Expr != "" {
			exector.Select(v.Expr)
		} else if v.Field != "" && v.Name != "" {
			exector.Select("$$" + v.Field + " as \"" + v.Name + "\"")
		} else if v.Field != "" {
			exector.Select("$$" + v.Field)
		}
		if v.Name != "" && v.DataType != "" {
			exector.SetFieldDataType(v.Name, v.DataType)
		}
	}
	for _, v := range s.Wheres {
		iw := s.queryWhereToIWhere(v)
		if iw != nil {
			if iw.GetLogical() == "or" {
				iw = exector.OrWhere(iw.GetQuery(), iw.GetArgs()...).SetDataType(iw.GetDataType())
			} else {
				iw = exector.Where(iw.GetQuery(), iw.GetArgs()...).SetDataType(iw.GetDataType())
			}
			if v.Children != nil && len(v.Children) > 0 {
				for _, item := range v.Children {
					s.addSubItemToIWhere(iw, item)
				}
			}
		}
	}
	if len(s.Orders) > 0 {
		for _, v := range s.Orders {
			if v.Expr != "" && v.Order != "" {
				exector.Order(fmt.Sprintf("%s %s", v.Expr, v.Order))
			} else if v.Expr != "" {
				exector.Order(v.Expr)
			} else if v.Field != "" && v.Order != "" {
				exector.Order(fmt.Sprintf("$$%s %s", v.Field, v.Order))
			} else if v.Field != "" {
				exector.Order(fmt.Sprintf("$$%s", v.Field))
			}
		}
	} else {
		exector.Order("$$ID")
	}

	return exector
}
func (s *QueryCase) addSubItemToIWhere(iw IQWhere, subValue QueryWhere) {
	newIw := s.queryWhereToIWhere(subValue)
	if newIw != nil && iw.GetLogical() == "or" {
		newIw = iw.OrWhere(newIw.GetQuery(), newIw.GetArgs()).SetDataType(newIw.GetDataType())
	} else if newIw != nil {
		newIw = iw.Where(newIw.GetQuery(), newIw.GetArgs()).SetDataType(newIw.GetDataType())
	}
	if newIw != nil && subValue.Children != nil && len(subValue.Children) > 0 {
		for _, item := range subValue.Children {
			s.addSubItemToIWhere(newIw, item)
		}
	}
}
func (s *QueryCase) queryWhereToIWhere(value QueryWhere) IQWhere {
	if value.Enabled.IsFalse() {
		return nil
	}
	if value.Operator == "" {
		value.Operator = "="
	}
	item := qWhere{Logical: value.Logical, DataType: value.DataType}
	//简单模式
	if value.Expr == "" && value.Field != "" && value.Value.Valid() && value.Operator == "contains" {
		item.Query = fmt.Sprintf("$$%v like ?", value.Field)
		item.Args = []interface{}{"%" + value.Value.GetString() + "%"}
		if value.Value.GetString() == "" {
			return nil
		}
	} else if value.Expr == "" && value.Field != "" && value.Value.Valid() && value.Operator == "like" {
		item.Query = fmt.Sprintf("$$%v like ?", value.Field)
		item.Args = []interface{}{"%" + value.Value.GetString() + "%"}
		if value.Value.GetString() == "" {
			return nil
		}
	} else if value.Expr == "" && value.Field != "" && value.Value.Valid() && value.Operator == "not like" {
		item.Query = fmt.Sprintf("$$%v not like ?", value.Field)
		item.Args = []interface{}{"%" + value.Value.GetString() + "%"}
		if value.Value.GetString() == "" {
			return nil
		}
	} else if value.Expr == "" && value.Field != "" && value.Value.Valid() && (value.Operator == "between") {
		item.Args = value.Value.GetInterfaceSlice()
		item.Query = fmt.Sprintf("$$%v between ? and ?", value.Field)
		if len(item.Args) != 2 {
			return nil
		}
	} else if value.Expr == "" && value.Field != "" && value.Value.Valid() && (value.Operator == "in" || value.Operator == "not in") {
		item.Args = value.Value.GetInterfaceSlice()
		item.Query = fmt.Sprintf("$$%v %s (?)", value.Field, value.Operator)
		if len(item.Args) == 0 {
			return nil
		}

	} else if value.Expr == "" && value.Field != "" && value.Value.Valid() && (value.Operator == "=" || value.Operator == "<>" || value.Operator == ">" || value.Operator == ">=" || value.Operator == "<" || value.Operator == "<=") {
		item.Query = fmt.Sprintf("$$%v %s ?", value.Field, value.Operator)
		item.Args = value.Value.GetInterfaceSlice()
		if len(item.Args) == 0 || value.Value.GetString() == "" {
			return nil
		}
	} else if value.Expr == "" && value.Field != "" && value.Value.Valid() && value.Value.GetString() != "" && s.Context.IsValid() && (value.Operator == "=p" || value.Operator == "<>p" || value.Operator == ">p" || value.Operator == ">=p" || value.Operator == "<p" || value.Operator == "<=p") {
		item.Query = fmt.Sprintf("$$%v %s ?", value.Field, strings.Replace(value.Operator, "p", "", -1))
		item.Args = []interface{}{s.Context.GetValue(strings.Replace(value.Value.GetString(), "@", "", -1))}
	} else if value.Expr != "" {
		//表达式 模式
		item.Query = value.Expr
		item.Args = value.Value.GetInterfaceSlice()
		if strings.Contains(item.Query, "?") && len(item.Args) == 0 {
			return nil
		}
	} else {
		if len(value.Children) == 0 {
			return nil
		}
	}
	return &item
}
