package md

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/ggoop/mdf/framework/glog"
	"github.com/ggoop/mdf/utils"
)

func (s *oqlImpl) parse() error {
	for _, v := range s.froms {
		s.parseFromField(v)
	}
	for _, v := range s.joins {
		s.parseJoinField(v)
	}
	for _, v := range s.selects {
		s.parseSelectField(v)
	}
	for _, v := range s.wheres {
		s.parseWhereField(v)
	}
	for _, v := range s.having {
		s.parseWhereField(v)
	}
	for _, v := range s.orders {
		s.parseOrderField(v)
	}
	for _, v := range s.groups {
		s.parseGroupField(v)
	}
	return s.error
}

//=============== build
func (s *oqlImpl) buildFroms() *OQLStatement {
	statement := &OQLStatement{}
	queries := make([]string, 0)
	args := make([]interface{}, 0)
	for _, v := range s.froms {
		queries = append(queries, v.Expr())
		args = append(args, v.Args()...)
		statement.Affected++
	}
	statement.Query = strings.Join(queries, ",")
	statement.Args = args
	return statement
}

func (s *oqlImpl) buildJoins() *OQLStatement {
	statement := &OQLStatement{}
	queries := make([]string, 0)
	args := make([]interface{}, 0)
	tables := make([]*oqlEntity, 0)
	for _, v := range s.entities {
		tables = append(tables, v)
	}
	sort.Slice(tables, func(i, j int) bool {
		return tables[i].Sequence < tables[j].Sequence
	})
	for _, t := range tables {
		if t.Path == "" || t.IsMain {
			continue
		}
		relationship, _ := s.parseEntityField(t.Path)
		if relationship == nil {
			glog.Errorf("找不到关联字段")
			continue
		}
		if relationship.Field.Kind == KIND_TYPE_BELONGS_TO || relationship.Field.Kind == KIND_TYPE_HAS_ONE {
			fkey := relationship.Entity.Entity.GetField(relationship.Field.ForeignKey)
			lkey := t.Entity.GetField(relationship.Field.AssociationKey)
			condition := ""
			tag := false
			if relationship.Field.TypeType == utils.TYPE_ENUM {
				if relationship.Field.Limit != "" {
					condition = fmt.Sprintf(" and %v.entity_id=?", t.Alias)
					queries = append(queries, fmt.Sprintf("left join %v  %v on %v.%v=%v.%v%v", t.Entity.TableName, t.Alias, t.Alias, "id", relationship.Entity.Alias, fkey.DbName, condition))
					args = append(args, relationship.Field.Limit)
					statement.Affected++
					tag = true
				}
			}
			if !tag {
				queries = append(queries, fmt.Sprintf("left join %v  %v on %v.%v=%v.%v%v", t.Entity.TableName, t.Alias, t.Alias, lkey.DbName, relationship.Entity.Alias, fkey.DbName, condition))
				statement.Affected++
			}
		} else if relationship.Field.Kind == "has_many" {
			fkey := t.Entity.GetField(relationship.Field.ForeignKey)
			lkey := relationship.Entity.Entity.GetField(relationship.Field.AssociationKey)
			if fkey != nil && lkey != nil {
				queries = append(queries, fmt.Sprintf("left join %v  %v on %v.%v=%v.%v", t.Entity.TableName, t.Alias, t.Alias, fkey.DbName, relationship.Entity.Alias, lkey.DbName))
				statement.Affected++
			} else {
				glog.Error("构建join 联系出错，", glog.String("ForeignKey", relationship.Field.ForeignKey), glog.String("AssociationKey", relationship.Field.AssociationKey))
			}
		}
	}
	return statement
}
func (s *oqlImpl) buildSelects() *OQLStatement {
	statement := &OQLStatement{}
	queries := make([]string, 0)
	args := make([]interface{}, 0)
	for _, v := range s.selects {
		if v.expr != "" {
			queries = append(queries, v.expr)
			args = append(args, v.Args...)
			statement.Affected++
		}
	}
	statement.Query = strings.Join(queries, ",")
	statement.Args = args
	return statement
}
func (s *oqlImpl) buildWheres() *OQLStatement {
	return s.buildWheresItem(s.wheres)
}

func (s *oqlImpl) buildHaving() *OQLStatement {
	return s.buildWheresItem(s.having)
}
func (s *oqlImpl) buildWheresItem(wheres []OQLWhere) *OQLStatement {
	statement := &OQLStatement{}
	queries := make([]string, 0)
	args := make([]interface{}, 0)
	for _, v := range wheres {
		sub := s.buildWheresItem(v.Children())
		statement.Affected += sub.Affected
		//当前节点加上子节点条件
		if v.Expr() != "" {
			if len(queries) > 0 {
				queries = append(queries, " ", v.Logical(), " ")
			}
			//如果有子条件，则需要把子条件也加入到条件集合中
			if sub.Query != "" {
				// (a=b and (a=1 or a=2))
				queries = append(queries, "((", v.Expr(), ") ", v.Logical(), " ", sub.Query, ")")
				args = append(args, s.getWhereArgs(v)...)
				args = append(args, sub.Args...)
			} else {
				//没有子条件时，就只加当前条件
				queries = append(queries, "(", v.Expr(), ")")
				args = append(args, s.getWhereArgs(v)...)
			}
			statement.Affected++
		} else if sub.Query != "" { //仅仅有子节点
			if len(queries) > 0 {
				queries = append(queries, " ", v.Logical(), " ")
			}
			//如果子条件多于一个，则需要用括号包裹起来
			if sub.Affected > 1 {
				queries = append(queries, "(", sub.Query, ")")
			} else {
				queries = append(queries, sub.Query)
			}
			args = append(args, sub.Args...)
		}
	}
	statement.Args = args
	statement.Query = strings.Join(queries, "")
	return statement
}
func (s *oqlImpl) getWhereArgs(where OQLWhere) []interface{} {
	if len(where.Args()) <= 0 {
		return nil
	}
	args := where.Args()
	for i, item := range args {
		if where.DataType() == utils.FIELD_TYPE_ENTITY || where.DataType() == utils.FIELD_TYPE_ENUM || where.DataType() == "" {
			if v, ok := item.(map[string]interface{}); ok {
				if v["_isRefObject"] != nil && v["id"] != nil {
					args[i] = v["id"]
				} else if v["_isEnumObject"] != nil && v["id"] != nil {
					args[i] = v["id"]
				} else if vv, ok := v["id"]; ok {
					args[i] = vv
				}
			} else if v, ok := item.(utils.SJson); ok {
				args[i] = v.GetValue()
			}
		} else if where.DataType() == utils.FIELD_TYPE_DATE {
			args[i] = utils.ToTime(item).Format(utils.Layout_YYYYMMDD)
		} else if where.DataType() == utils.FIELD_TYPE_DATETIME {
			args[i] = utils.ToTime(item).Format(utils.Layout_YYYYMMDDHHIISS)
		} else if where.DataType() == utils.FIELD_TYPE_BOOL {
			args[i] = utils.ToSBool(item)
		}
	}
	return args
}

func (s *oqlImpl) buildGroups() *OQLStatement {
	statement := &OQLStatement{}
	queries := make([]string, 0)
	args := make([]interface{}, 0)
	for _, v := range s.groups {
		if v.expr != "" {
			queries = append(queries, v.expr)
			args = append(args, v.Args...)
			statement.Affected++
		}
	}
	statement.Query = strings.Join(queries, ",")
	statement.Args = args
	return statement
}
func (s *oqlImpl) buildOrders() *OQLStatement {
	statement := &OQLStatement{}
	queries := make([]string, 0)
	args := make([]interface{}, 0)
	for _, v := range s.orders {
		if v.expr != "" {
			queries = append(queries, v.expr)
			args = append(args, v.Args...)
			statement.Affected++
		}
	}
	statement.Query = strings.Join(queries, ",")
	statement.Args = args
	return statement
}

//=============== parse
func (s *oqlImpl) parseFromField(value OQLFrom) {
	//主表，使用别名作路径
	form := s.parseEntity(value.Query(), value.Alias())
	parts := make([]string, 0)
	if form == nil {
		form = &oqlEntity{}

		path := strings.ToLower(value.Alias())
		form.Entity = &MDEntity{ID: value.Query(), TableName: value.Query()}
		form.Sequence = len(s.entities) + 1
		form.Alias = fmt.Sprintf("a%v", form.Sequence)
		form.Path = path
		s.entities[path] = form
	}
	form.IsMain = true
	parts = append(parts, form.Entity.TableName)
	if form.Alias != "" {
		parts = append(parts, form.Alias)
	}
	value.setExpr(strings.Join(parts, " "))
}
func (s *oqlImpl) parseJoinField(value *oqlJoin) error {
	joins := make([]string, 0)
	switch value.Type {
	case OQL_LEFT_JOIN:
		joins = append(joins, "left join")
	case OQL_RIGHT_JOIN:
		joins = append(joins, "left join")
	case OQL_FULL_JOIN:
		joins = append(joins, "join")
	}
	//主表，使用别名作路径
	form := s.parseEntity(value.Query, value.Alias)
	if form == nil {
		form = &oqlEntity{}
		path := strings.ToLower(value.Alias)
		form.Entity = &MDEntity{ID: value.Query, TableName: value.Query}
		form.Sequence = len(s.entities) + 1
		form.Alias = fmt.Sprintf("a%v", form.Sequence)
		form.Path = path
		s.entities[path] = form
	}
	joins = append(joins, form.Entity.TableName)
	if form.Alias != "" {
		joins = append(joins, form.Alias)
	}
	if value.Condition != "" {
		condition := s.parseFieldExpr(value.Condition)
		if condition != "" {
			joins = append(joins, "on ", condition)
		} else {
			joins = append(joins, "on ", value.Condition)
		}
	}
	value.expr = strings.Join(joins, " ")
	return nil
}
func (s *oqlImpl) parseWhereField(value OQLWhere) {
	value.setExpr(s.parseFieldExpr(value.Query()))
	if len(value.Children()) > 0 {
		for _, v := range value.Children() {
			s.parseWhereField(v)
		}
	}
}
func (s *oqlImpl) parseSelectField(value *oqlSelect) {
	expr := s.parseFieldExpr(value.Query)
	if value.Alias != "" {
		value.expr = fmt.Sprintf("%v as %v", expr, value.Alias)
	} else {
		value.expr = expr
	}
}
func (s *oqlImpl) parseGroupField(value *oqlGroup) {
	value.expr = s.parseFieldExpr(value.Query)
}
func (s *oqlImpl) parseOrderField(value *oqlOrder) {
	expr := s.parseFieldExpr(value.Query)
	if value.Order == OQL_ORDER_DESC {
		value.expr = fmt.Sprintf("%v desc", expr)
	} else {
		value.expr = expr
	}
}

// 解析字段表达式，如
//	a.fieldA+fieldB+sum(b.fieldA)   =>a.fieldA ,fieldB, b.fieldA
//	$$a.fieldA + sum( c.fieldA )	=>$$a.fieldA, c.fieldA
// 函数与左括号之间不能有空格
// 多级字段.号不能有空格
func (s *oqlImpl) parseFieldExpr(expr string) string {
	if expr == "" {
		return expr
	}
	r, _ := regexp.Compile(`([\$]?[A-Za-z._]+[0-9A-Za-z|\(])`)
	matches := r.FindAllStringSubmatch(expr, -1)
	for _, match := range matches {
		str := match[1]
		//带有括号的是函数，不需要解析
		if strings.Index(str, utils.PARENTHESIS_LEFT) < 0 {
			field, _ := s.parseEntityField(str)
			if field != nil {
				expr = strings.ReplaceAll(expr, str, fmt.Sprintf("%s.%s", field.Entity.Alias, field.Field.DbName))
			}
		}
	}
	return expr
}

// 解析实体
func (s *oqlImpl) formatEntity(entity *MDEntity) *oqlEntity {
	e := oqlEntity{Entity: entity}
	return &e
}
func (s *oqlImpl) formatEntityField(entity *oqlEntity, field *MDField) *oqlField {
	e := oqlField{Entity: entity, Field: field}
	return &e
}
func (s *oqlImpl) parseEntity(id, path string) *oqlEntity {
	path = strings.ToLower(strings.TrimSpace(path))
	if v, ok := s.entities[path]; ok {
		return v
	}
	entity := MDSv().GetEntity(id)
	if entity == nil {
		err := glog.Errorf("找不到实体 %v", id)
		s.AddErr(err)
		return nil
	}
	v := s.formatEntity(entity)
	v.Sequence = len(s.entities) + 1
	v.Alias = fmt.Sprintf("a%v", v.Sequence)
	v.Path = path
	s.entities[path] = v
	return v
}

// 解析字段
func (s *oqlImpl) parseEntityField(fieldPath string) (*oqlField, error) {
	fieldPath = strings.ToLower(strings.TrimSpace(fieldPath))
	if v, ok := s.fields[fieldPath]; ok {
		return v, nil
	}
	start := 0
	parts := strings.Split(fieldPath, ".")
	var mainFrom OQLFrom
	if len(parts) > 1 {
		//如果主表有别名，则第一个字段为表
		for i, v := range s.froms {
			if v.Alias() != "" && strings.ToLower(v.Alias()) == parts[0] {
				mainFrom = s.froms[i]
				start = 1
				break
			}
		}
	}
	if mainFrom == nil {
		//如果没有找到主表，则说明字段没有 表作为导引
		for i, v := range s.froms {
			if v.Alias() == "" {
				mainFrom = s.froms[i]
				break
			}
		}
	}
	if mainFrom == nil {
		mainFrom = s.froms[0]
	}
	//主实体
	entity := s.entities[strings.ToLower(mainFrom.Alias())]
	if entity == nil {
		return nil, nil
	}
	path := ""
	for i, part := range parts {
		if i > 0 {
			path += "."
		}
		path += part
		if i < start {
			continue
		}
		mdField := entity.Entity.GetField(part)
		if mdField == nil {
			mdField = &MDField{ID: part, Code: part, Name: part, DbName: part}
			s.AddErr(glog.Error(fmt.Sprintf("找不到字段 %v", path)))
		}
		field := s.formatEntityField(entity, mdField)
		field.Path = path
		s.fields[path] = field
		if i < len(parts)-1 {
			entity = s.parseEntity(mdField.TypeID, path)
			if s.Error != nil || entity == nil {
				return nil, nil
			}
		} else {
			return field, nil
		}
	}
	return nil, nil
}
