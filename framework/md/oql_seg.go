package md

import (
	"fmt"
	"github.com/ggoop/mdf/utils"
	"regexp"
	"strings"
)

func (m *oqlImpl) Clone() OQL {
	n := *m
	return &n
}
func (s *oqlImpl) Error() error {
	return s.error
}
func (s *oqlImpl) SetContext(context *utils.TokenContext) OQL {
	s.context = context
	return s
}
func (s *oqlImpl) SetActuator(actuator OQLActuator) OQL {
	s.actuator = actuator
	return s
}
func (s *oqlImpl) GetActuator() OQLActuator {
	if s.actuator == nil {
		s.actuator = GetOQLActuator()
	}
	return s.actuator
}

// 设置 主 from ，示例：
//  tableA
//	tableA as a
//	tableA a
func (s *oqlImpl) From(query interface{}, args ...interface{}) OQL {
	if v, ok := query.(string); ok {
		seg := &oqlFrom{query: v, args: args}
		r := regexp.MustCompile(regexp_OQL_FROM)
		matches := r.FindStringSubmatch(v)
		if matches != nil && len(matches) == 4 {
			if matches[2] != "" {
				seg.query = matches[1]
				seg.alias = matches[2]
			} else {
				seg.query = matches[3]
			}
		}
		s.froms = append(s.froms, seg)
	} else if v, ok := query.(OQLFrom); ok {
		s.froms = append(s.froms, v)
	}
	return s
}
func (s *oqlImpl) Join(joinType OQLJoinType, query string, condition string, args ...interface{}) OQL {
	seg := &oqlJoin{Type: joinType, Query: query, Condition: condition, Args: args}
	r := regexp.MustCompile(regexp_OQL_FROM)
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
func (s *oqlImpl) Select(query interface{}, args ...interface{}) OQL {
	if v, ok := query.(string); ok {
		items := make([]string, 0)
		if ok, _ := regexp.MatchString(`\(.*,`, v); ok {
			items = append(items, v)
		} else {
			items = strings.Split(strings.TrimSpace(v), ",")
		}
		for _, item := range items {
			seg := &oqlSelect{Query: item, Args: args}
			r := regexp.MustCompile(regexp_OQL_SELECT)
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
func (s *oqlImpl) Order(query interface{}, args ...interface{}) OQL {
	if v, ok := query.(string); ok {
		items := make([]string, 0)
		if ok, _ := regexp.MatchString(`\(.*,`, v); ok {
			items = append(items, v)
		} else {
			items = strings.Split(strings.TrimSpace(v), ",")
		}
		for _, item := range items {
			seg := &oqlOrder{Query: item, Args: args}
			r := regexp.MustCompile(regexp_OQL_ORDER)
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
func (s *oqlImpl) Group(query interface{}, args ...interface{}) OQL {
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
func (s *oqlImpl) Where(query interface{}, args ...interface{}) OQLWhere {
	var seg OQLWhere
	if v, ok := query.(string); ok {
		seg = NewOQLWhere(v, args...)
	} else if v, ok := query.(OQLWhere); ok {
		seg = v
	} else if v, ok := query.(bool); ok && !v {
		s.wheres = make([]OQLWhere, 0)
	}
	s.wheres = append(s.wheres, seg)
	return seg
}
func (s *oqlImpl) Or() OQLWhere {
	seg := &oqlWhere{logical: OQL_WHERE_LOGICAL_OR}
	s.wheres = append(s.wheres, seg)
	return seg
}
func (s *oqlImpl) Having(query interface{}, args ...interface{}) OQLWhere {
	var seg OQLWhere
	if v, ok := query.(string); ok {
		seg = NewOQLWhere(v, args...)
	} else if v, ok := query.(OQLWhere); ok {
		seg = v
	} else if v, ok := query.(bool); ok && !v {
		s.having = make([]OQLWhere, 0)
	}
	s.having = append(s.having, seg)
	return seg
}

//============= exec
func (s *oqlImpl) Count(value interface{}) OQL {
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
func (s *oqlImpl) Pluck(column string, value interface{}) OQL {
	return s.GetActuator().Pluck(s, column, value)
}
func (s *oqlImpl) Take(out interface{}) OQL {
	return s.GetActuator().Take(s, out)
}
func (s *oqlImpl) Find(out interface{}) OQL {
	return s.GetActuator().Find(s, out)
}
func (s *oqlImpl) Paginate(value interface{}, page int64, pageSize int64) OQL {
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
func (s *oqlImpl) Create(data interface{}) OQL {
	return s.GetActuator().Create(s, data)
}

//update table set aa=bb
//field 从select 取， value 从 data 取
func (s *oqlImpl) Update(data interface{}) OQL {
	return s.GetActuator().Update(s, data)
}
func (s *oqlImpl) Delete() OQL {
	return s.GetActuator().Delete(s)
}
func (s *oqlImpl) AddErr(err error) OQL {
	s.errors = append(s.errors, err)
	s.error = err
	return s
}
