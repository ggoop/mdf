package md

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ggoop/mdf/framework/context"
)

type OQLStatement struct {
	Query    string
	Args     []interface{}
	Affected int64
}

//公共查询
type OQL struct {
	Error    error
	errors   []error
	entities map[string]*oqlEntity
	fields   map[string]*oqlField
	froms    []*oqlFrom
	joins    []*oqlJoin
	selects  []*oqlSelect
	orders   []*oqlOrder
	wheres   []*OQLWhere
	groups   []*oqlGroup
	having   []*OQLWhere
	offset   int64
	limit    int64
	context  *context.Context
	actuator OQLActuator
}

func (s *OQL) SetContext(context *context.Context) *OQL {
	s.context = context
	return s
}
func (s *OQL) SetActuator(actuator OQLActuator) *OQL {
	s.actuator = actuator
	return s
}
func (s *OQL) GetActuator() OQLActuator {
	if s.actuator == nil {
		s.actuator = GetOQLActuator()
	}
	return s.actuator
}

// 设置 主 from ，示例：
//  tableA
//	tableA as a
//	tableA a
func (s *OQL) From(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		seg := &oqlFrom{Query: v, Args: args}
		r := regexp.MustCompile(REGEXP_OQL_FROM)
		matches := r.FindStringSubmatch(v)
		if matches != nil && len(matches) == 4 {
			if matches[2] != "" {
				seg.Query = matches[1]
				seg.Alias = matches[2]
			} else {
				seg.Query = matches[3]
			}
		}
		s.froms = append(s.froms, seg)
	} else if v, ok := query.(oqlFrom); ok {
		s.froms = append(s.froms, &v)
	}
	return s
}
func (s *OQL) Join(joinType OQLJoinType, query string, condition string, args ...interface{}) *OQL {
	seg := &oqlJoin{Type: joinType, Query: query, Condition: condition, Args: args}
	r := regexp.MustCompile(REGEXP_OQL_FROM)
	matches := r.FindStringSubmatch(query)
	if matches != nil && len(matches) == 4 {
		if matches[2] != "" {
			seg.Query = matches[1]
			seg.Alias = matches[2]
		} else {
			seg.Query = matches[3]
		}
	}
	s.joins = append(s.joins, seg)
	return s
}

// 添加字段，示例：
//	单字段：fieldA ，fieldA as A
//	复合字段：sum(fieldA) AS A，fieldA+fieldB as c
//
func (s *OQL) Select(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		items := make([]string, 0)
		if ok, _ := regexp.MatchString(`\(.*,`, v); ok {
			items = append(items, v)
		} else {
			items = strings.Split(strings.TrimSpace(v), ",")
		}
		for _, item := range items {
			seg := &oqlSelect{Query: item, Args: args}
			r := regexp.MustCompile(REGEXP_OQL_SELECT)
			matches := r.FindStringSubmatch(item)
			if matches != nil && len(matches) == 4 {
				if matches[2] != "" {
					seg.Query = matches[1]
					seg.Alias = matches[2]
				} else {
					seg.Query = matches[3]
				}
			}
			s.selects = append(s.selects, seg)
		}
	} else if v, ok := query.(oqlSelect); ok {
		s.selects = append(s.selects, &v)
	} else if v, ok := query.(bool); ok && !v {
		s.selects = make([]*oqlSelect, 0)
	}
	return s
}

//排序，示例：
// fieldA desc，fieldA + fieldB
func (s *OQL) Order(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		items := make([]string, 0)
		if ok, _ := regexp.MatchString(`\(.*,`, v); ok {
			items = append(items, v)
		} else {
			items = strings.Split(strings.TrimSpace(v), ",")
		}
		for _, item := range items {
			seg := &oqlOrder{Query: item, Args: args}
			r := regexp.MustCompile(REGEXP_OQL_ORDER)
			matches := r.FindStringSubmatch(item)
			if matches != nil && len(matches) == 4 {
				if matches[2] != "" {
					seg.Query = matches[1]
					if strings.ToLower(matches[2]) == "desc" {
						seg.Order = OQL_ORDER_DESC
					} else {
						seg.Order = OQL_ORDER_ASC
					}
				} else {
					seg.Query = matches[3]
				}
			}
			s.orders = append(s.orders, seg)
		}
	} else if v, ok := query.(oqlOrder); ok {
		s.orders = append(s.orders, &v)
	} else if v, ok := query.(bool); ok && !v {
		s.orders = make([]*oqlOrder, 0)
	}
	return s
}
func (s *OQL) Group(query interface{}, args ...interface{}) *OQL {
	if v, ok := query.(string); ok {
		items := make([]string, 0)
		if ok, _ := regexp.MatchString(`\(.*,`, v); ok {
			items = append(items, v)
		} else {
			items = strings.Split(strings.TrimSpace(v), ",")
		}
		for _, item := range items {
			seg := &oqlGroup{Query: item, Args: args}
			s.groups = append(s.groups, seg)
		}
	} else if v, ok := query.(oqlGroup); ok {
		s.groups = append(s.groups, &v)
	} else if v, ok := query.(bool); ok && !v {
		s.groups = make([]*oqlGroup, 0)
	}
	return s
}
func (s *OQL) Where(query interface{}, args ...interface{}) *OQLWhere {
	var seg *OQLWhere
	if v, ok := query.(string); ok {
		seg = NewOQLWhere(v, args...)
	} else if v, ok := query.(OQLWhere); ok {
		seg = &v
	} else if v, ok := query.(bool); ok && !v {
		s.wheres = make([]*OQLWhere, 0)
	}
	s.wheres = append(s.wheres, seg)
	return seg
}

func (s *OQL) Having(query interface{}, args ...interface{}) *OQLWhere {
	var seg *OQLWhere
	if v, ok := query.(string); ok {
		seg = NewOQLWhere(v, args...)
	} else if v, ok := query.(OQLWhere); ok {
		seg = &v
	} else if v, ok := query.(bool); ok && !v {
		s.having = make([]*OQLWhere, 0)
	}
	s.having = append(s.having, seg)
	return seg
}

//============= exec
func (s *OQL) Count(value interface{}) *OQL {
	s.parse()
	queries := make([]string, 0)
	args := make([]interface{}, 0)
	queries = append(queries, "select count(*)")
	if statement := s.buildFroms(); statement.Affected > 0 {
		queries = append(queries, fmt.Sprintf("from %s", statement.Query))
		args = append(args, statement.Args...)
	}
	if statement := s.buildWheres(); statement.Affected > 0 {
		queries = append(queries, fmt.Sprintf("where %s", statement.Query))
		args = append(args, statement.Args...)
	}
	if statement := s.buildGroups(); statement.Affected > 0 {
		queries = append(queries, fmt.Sprintf("group by %s", statement.Query))
		args = append(args, statement.Args...)
	}
	if statement := s.buildHaving(); statement.Affected > 0 {
		queries = append(queries, fmt.Sprintf("having %s", statement.Query))
		args = append(args, statement.Args...)
	}
	return s.GetActuator().Count(s, value)
}
func (s *OQL) Pluck(column string, value interface{}) *OQL {
	return s.GetActuator().Pluck(s, column, value)
}
func (s *OQL) Take(out interface{}) *OQL {
	return s.GetActuator().Take(s, out)
}
func (s *OQL) Find(out interface{}) *OQL {
	return s.GetActuator().Find(s, out)
}
func (s *OQL) Paginate(value interface{}, page int64, pageSize int64) *OQL {
	if pageSize > 0 && page <= 0 {
		page = 1
	} else if pageSize <= 0 {
		pageSize = 0
		page = 0
	}
	s.limit = pageSize
	s.offset = (page - 1) * pageSize

	return s.GetActuator().Find(s, value)
}

//insert into table (aaa,aa,aa) values(aaa,aaa,aaa)
//field 从select 取， value 从 data 取
func (s *OQL) Create(data interface{}) *OQL {
	return s.GetActuator().Create(s, data)
}

//update table set aa=bb
//field 从select 取， value 从 data 取
func (s *OQL) Update(data interface{}) *OQL {
	return s.GetActuator().Update(s, data)
}
func (s *OQL) Delete() *OQL {
	return s.GetActuator().Delete(s)
}
func (s *OQL) AddErr(err error) *OQL {
	s.errors = append(s.errors, err)
	s.Error = err
	return s
}
