package query

import (
	"database/sql"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/ggoop/mdf/context"
	"github.com/ggoop/mdf/glog"
	"github.com/ggoop/mdf/gorm"
	"github.com/ggoop/mdf/md"
	"github.com/ggoop/mdf/repositories"
	"github.com/ggoop/mdf/utils"

	"github.com/shopspring/decimal"
)

type IExector interface {
	PrepareQuery(mysql *repositories.MysqlRepo) (*gorm.DB, error)
	Query(mysql *repositories.MysqlRepo) ([]map[string]interface{}, error)
	Take(mysql *repositories.MysqlRepo) (map[string]interface{}, error)
	Count(mysql *repositories.MysqlRepo) (int, error)
	Select(query string, args ...interface{}) IExector
	Where(query string, args ...interface{}) IQWhere
	OrWhere(query string, args ...interface{}) IQWhere
	Joins(query string, args ...interface{}) IExector
	Order(query string, args ...interface{}) IExector
	Group(query string, args ...interface{}) IExector
	Page(page, pageSize int) IExector
	SetContext(context *context.Context) IExector
	SetFieldDataTypes(dataTypes map[string]string) IExector
	SetFieldDataType(field, dataType string) IExector
	Clone() IExector
	GetMainFrom() IQFrom
}

const REGEXP_FIELD_EXP_STRICT string = `\$\$([A-Za-z._]+[0-9A-Za-z]*)`
const REGEXP_FIELD_EXP string = `([A-Za-z._]+[0-9A-Za-z]*)`
const REGEXP_VAR_EXP string = `{([A-Za-z._]+[0-9A-Za-z]*)}`

type oqlEntity struct {
	Alia     string
	IsMain   bool
	Entity   *md.MDEntity
	Path     string
	Sequence int
}
type oqlField struct {
	Entity *oqlEntity
	Field  *md.MDField
	Path   string
}
type IQFrom interface {
	GetQuery() string
	GetAlia() string
	GetExpr() string
}
type oqlFrom struct {
	Query string
	Alia  string
	Expr  string
}

func (m *oqlFrom) GetQuery() string {
	return m.Query
}
func (m *oqlFrom) GetAlia() string {
	return m.Alia
}
func (m *oqlFrom) GetExpr() string {
	return m.Expr
}

type oqlJoin struct {
	Query string
	Expr  string
	Args  []interface{}
}
type oqlSelect struct {
	Query string
	Expr  string
	Args  []interface{}
}
type oqlOrder struct {
	Query    string
	Sequence int
	Expr     string
	Args     []interface{}
}
type oqlGroup struct {
	Query string
	Expr  string
	Args  []interface{}
}

type exector struct {
	entities  map[string]*oqlEntity
	fields    map[string]*oqlField
	dataTypes map[string]string
	froms     []*oqlFrom
	selects   []*oqlSelect
	joins     []*oqlJoin
	wheres    []*qWhere
	orders    []*oqlOrder
	groups    []*oqlGroup
	page      int
	pageSize  int
	context   *context.Context
}

func NewExector(query string) IExector {
	exec := &exector{
		entities:  make(map[string]*oqlEntity),
		fields:    make(map[string]*oqlField),
		dataTypes: make(map[string]string),
		froms:     make([]*oqlFrom, 0),
		selects:   make([]*oqlSelect, 0),
		joins:     make([]*oqlJoin, 0),
		wheres:    make([]*qWhere, 0),
		orders:    make([]*oqlOrder, 0),
		groups:    make([]*oqlGroup, 0),
	}
	exec.From(query)
	return exec
}
func (m *exector) Clone() IExector {
	n := *m
	return &n
}
func (m *exector) Page(page, pageSize int) IExector {
	m.page = page
	m.pageSize = pageSize
	return m
}

func (m *exector) SetFieldDataTypes(dataTypes map[string]string) IExector {
	if dataTypes == nil {
		dataTypes = make(map[string]string)
	}
	m.dataTypes = dataTypes
	return m
}
func (m *exector) SetFieldDataType(field, dataType string) IExector {
	m.dataTypes[field] = dataType
	return m
}
func (m *exector) formatEntity(entity *md.MDEntity) *oqlEntity {
	e := oqlEntity{Entity: entity}
	return &e
}
func (m *exector) formatField(entity *oqlEntity, field *md.MDField) *oqlField {
	e := oqlField{Entity: entity, Field: field}
	return &e
}
func (m *exector) Count(mysql *repositories.MysqlRepo) (int, error) {
	if q, err := m.PrepareQuery(mysql); err != nil {
		return 0, err
	} else {
		count := 0
		if err := q.Count(&count).Error; err != nil {
			return 0, err
		}
		return count, nil
	}
}
func (m *exector) Take(mysql *repositories.MysqlRepo) (map[string]interface{}, error) {
	m.page = 1
	m.pageSize = 1
	datas, err := m.Query(mysql)
	if err != nil {
		return nil, err
	}
	if len(datas) > 0 {
		return datas[0], nil
	}
	return nil, nil
}
func (m *exector) Query(mysql *repositories.MysqlRepo) ([]map[string]interface{}, error) {
	q, err := m.PrepareQuery(mysql)
	if err != nil {
		return nil, err
	}
	if m.page > 0 && m.pageSize > 0 {
		q = q.Limit(m.pageSize).Offset((m.page - 1) * m.pageSize)
	}
	rows, err := q.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	columnTypes, _ := rows.ColumnTypes()
	columnTypeMap := map[string]string{}
	for _, c := range columnTypes {
		columnTypeMap[c.Name()] = c.DatabaseTypeName()
		if ct, ok := m.dataTypes[c.Name()]; ok && ct != "" {
			columnTypeMap[c.Name()] = ct
		}
	}

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for index, column := range columns {
			dbType := strings.ToUpper(columnTypeMap[column])
			switch dbType {
			case "VARCHAR", "TEXT", "NVARCHAR", "STRING", "VARCHAR2", "CHAR", "CLOB":
				var ignored sql.NullString
				values[index] = &ignored
				break
			case "BOOL", "BOOLEAN":
				var ignored utils.SBool
				values[index] = &ignored
				break
			case "INT", "BIGINT", "TINYINT", "INTEGER":
				var ignored sql.NullInt64
				values[index] = &ignored
				break
			case "DECIMAL", "NUMBER":
				var ignored decimal.Decimal
				values[index] = &ignored
				break
			case "TIMESTAMP", "DATE", "DATETIME":
				var ignored utils.Time
				values[index] = &ignored
				break
			default:
				var ignored interface{}
				values[index] = &ignored
			}
		}
		if err := rows.Scan(values...); err != nil {
			glog.Error(err)
		}
		resultItem := make(map[string]interface{})
		for index, column := range columns {
			if v, ok := values[index].(*sql.NullString); ok {
				resultItem[column] = v.String
			} else if v, ok := values[index].(*utils.SBool); ok {
				resultItem[column] = v
			} else if v, ok := values[index].(*sql.NullInt64); ok {
				resultItem[column] = v.Int64
			} else if v, ok := values[index].(*decimal.Decimal); ok {
				resultItem[column] = *v
			} else if v, ok := values[index].(*utils.Time); ok && !v.IsZero() {
				resultItem[column] = v.Format(utils.Layout_YYYYMMDDHHIISS)
			} else if v == nil {
				resultItem[column] = v
			} else {
				resultItem[column] = *v
			}
		}
		results = append(results, resultItem)
	}
	return results, nil
}
func (m *exector) PrepareQuery(mysql *repositories.MysqlRepo) (*gorm.DB, error) {
	//parse
	for _, v := range m.froms {
		if err := m.parseFromField(v); err != nil {
			return nil, err
		}
	}
	for _, v := range m.selects {
		m.parseSelectField(v)
	}
	for _, v := range m.wheres {
		m.parseWhereField(v)
	}
	for _, v := range m.orders {
		v.Expr = m.replaceFieldString(v.Query)
	}
	for _, v := range m.groups {
		v.Expr = m.replaceFieldString(v.Query)
	}
	//build

	queryDB := m.buildFroms(mysql)
	queryDB = m.buildSelects(queryDB)
	queryDB = m.buildJoins(queryDB)

	if whereExpr, whereArgs, tag := m.buildWheres(m.wheres); tag > 0 {
		// 上下文替换
		if m.context != nil {
			r, _ := regexp.Compile(REGEXP_VAR_EXP)
			matched := r.FindAllStringSubmatch(whereExpr, -1)
			for _, match := range matched {
				v := m.context.GetValue(utils.SnakeString(match[1]))
				whereExpr = strings.ReplaceAll(whereExpr, match[0], "'"+v+"'")
			}
		}
		queryDB = queryDB.Where(whereExpr, whereArgs...)
	}
	queryDB = m.buildGroups(queryDB)
	queryDB = m.buildOrders(queryDB)
	return queryDB, nil
}

///=============== build
func (m *exector) buildFroms(mysql *repositories.MysqlRepo) *gorm.DB {
	parts := make([]string, 0)
	for _, v := range m.froms {
		parts = append(parts, v.Expr)
	}
	return mysql.Table(strings.Join(parts, ","))
}
func (m *exector) buildSelects(queryDB *gorm.DB) *gorm.DB {
	selects := make([]string, 0)
	for _, v := range m.selects {
		if v.Expr != "" {
			selects = append(selects, v.Expr)
		}
	}
	queryDB = queryDB.Select(selects)
	return queryDB
}
func (m *exector) getWhereArgs(where IQWhere) []interface{} {
	if where == nil || len(where.GetArgs()) <= 0 {
		return nil
	}
	args := where.GetArgs()
	dType := where.GetDataType()
	for i, item := range args {
		if dType == WHERE_TYPE_ENUM || dType == WHERE_TYPE_REF || dType == "" {
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
		} else if dType == WHERE_TYPE_DATE {
			args[i] = utils.CreateTime(item).Format(utils.Layout_YYYYMMDD)
		} else if dType == WHERE_TYPE_DATETIME {
			args[i] = utils.CreateTime(item).Format(utils.Layout_YYYYMMDDHHIISS)
		} else if dType == WHERE_TYPE_BOOL {
			args[i] = utils.SBool_Parse(item)
		}
	}
	return args
}
func (m *exector) buildWheres(wheres []*qWhere) (string, []interface{}, int) {
	if len(wheres) == 0 {
		return "", nil, 0
	}
	tag := 0
	exprs := make([]string, 0)
	args := make([]interface{}, 0)

	for _, v := range wheres {
		subExpr, subArgs, subTag := m.buildWheres(v.Children)
		tag += subTag
		//当前节点加上子节点条件
		if v.Expr != "" {
			if len(exprs) > 0 {
				exprs = append(exprs, " ", v.Logical, " ")
			}
			//如果有子条件，则需要把子条件也加入到条件集合中
			if subExpr != "" {
				// (a=b and (a=1 or a=2))
				exprs = append(exprs, "((", v.Expr, ") ", v.Logical, " ", subExpr, ")")
				args = append(args, m.getWhereArgs(v)...)
				args = append(args, subArgs...)
			} else {
				//没有子条件时，就只加当前条件
				exprs = append(exprs, "(", v.Expr, ")")
				args = append(args, m.getWhereArgs(v)...)
			}
			tag += 1
		} else if subExpr != "" { //仅仅有子节点
			if len(exprs) > 0 {
				exprs = append(exprs, " ", v.Logical, " ")
			}
			//如果子条件多于一个，则需要用括号包裹起来
			if subTag > 1 {
				exprs = append(exprs, "(", subExpr, ")")
			} else {
				exprs = append(exprs, subExpr)
			}
			args = append(args, subArgs...)
		}
	}
	return strings.Join(exprs, ""), args, tag
}
func (m *exector) buildGroups(queryDB *gorm.DB) *gorm.DB {
	selects := make([]string, 0)
	for _, v := range m.groups {
		if v.Expr != "" {
			selects = append(selects, v.Expr)
		}
	}
	queryDB = queryDB.Group(strings.Join(selects, ","))
	return queryDB
}
func (m *exector) buildOrders(queryDB *gorm.DB) *gorm.DB {
	selects := make([]string, 0)
	for _, v := range m.orders {
		if v.Expr != "" {
			selects = append(selects, v.Expr)
		}
	}
	queryDB = queryDB.Order(strings.Join(selects, ","))
	return queryDB
}
func (m *exector) buildJoins(queryDB *gorm.DB) *gorm.DB {
	tables := make([]*oqlEntity, 0)
	for _, v := range m.entities {
		tables = append(tables, v)
	}
	sort.Slice(tables, func(i, j int) bool {
		return tables[i].Sequence < tables[j].Sequence
	})
	for _, t := range tables {
		if t.Path == "" || t.IsMain {
			continue
		}
		relationship, _ := m.parseField(t.Path)
		if relationship == nil {
			glog.Errorf("找不到关联字段")
			continue
		}
		if relationship.Field.Kind == "belongs_to" || relationship.Field.Kind == "has_one" {
			fkey := relationship.Entity.Entity.GetField(relationship.Field.ForeignKey)
			lkey := t.Entity.GetField(relationship.Field.AssociationKey)
			args := make([]interface{}, 0)
			condition := ""
			tag := false
			if relationship.Field.TypeType == md.TYPE_ENUM {
				if relationship.Field.Limit != "" {
					condition = fmt.Sprintf(" and %v.entity_id=?", t.Alia)
					args = append(args, relationship.Field.Limit)
					queryDB = queryDB.Joins(fmt.Sprintf("left join %v  %v on %v.%v=%v.%v%v",
						t.Entity.TableName, t.Alia, t.Alia, "id", relationship.Entity.Alia, fkey.DbName, condition),
						args...)
					tag = true
				}

			}
			if !tag {
				queryDB = queryDB.Joins(fmt.Sprintf("left join %v  %v on %v.%v=%v.%v%v",
					t.Entity.TableName, t.Alia, t.Alia, lkey.DbName, relationship.Entity.Alia, fkey.DbName, condition),
					args...)
			}
		} else if relationship.Field.Kind == "has_many" {
			fkey := t.Entity.GetField(relationship.Field.ForeignKey)
			lkey := relationship.Entity.Entity.GetField(relationship.Field.AssociationKey)
			if fkey != nil && lkey != nil {
				queryDB = queryDB.Joins(fmt.Sprintf("left join %v  %v on %v.%v=%v.%v", t.Entity.TableName, t.Alia, t.Alia, fkey.DbName, relationship.Entity.Alia, lkey.DbName))
			} else {
				glog.Error("构建join 联系出错，", glog.String("ForeignKey", relationship.Field.ForeignKey), glog.String("AssociationKey", relationship.Field.AssociationKey))
			}
		}
	}
	return queryDB
}
func (m *exector) replaceFieldString(expr string) string {
	if expr == "" {
		return expr
	}
	tag := false
	//先使用 严格模式，如:$$ID > 0
	r, _ := regexp.Compile(REGEXP_FIELD_EXP_STRICT)
	matched := r.FindAllStringSubmatch(expr, -1)
	for _, match := range matched {
		tag = true
		field, _ := m.parseField(match[1])
		if field != nil {
			expr = strings.ReplaceAll(expr, match[0], fmt.Sprintf("%s.%s", field.Entity.Alia, field.Field.DbName))
		}
	}
	if !tag {
		//使用字段模式，如:ID >0
		r, _ = regexp.Compile(REGEXP_FIELD_EXP)
		matched := r.FindAllStringSubmatch(expr, -1)
		for _, match := range matched {
			tag = true
			field, _ := m.parseField(match[1])
			if field != nil {
				expr = strings.ReplaceAll(expr, match[0], fmt.Sprintf("%s.%s", field.Entity.Alia, field.Field.DbName))
			}
		}
	}
	return expr
}

///=============== parse
func (m *exector) parseWhereField(value *qWhere) {
	value.Expr = m.replaceFieldString(value.Query)
	if len(value.Children) > 0 {
		for _, v := range value.Children {
			m.parseWhereField(v)
		}
	}
}
func (m *exector) parseFromField(value *oqlFrom) error {
	items := strings.Split(strings.TrimSpace(value.Query), ",")
	strs := make([]string, 0)
	for _, item := range items {
		parts := strings.Split(strings.TrimSpace(item), " ")
		if len(parts) == 1 {
			parts = append(parts, "")
		}
		form, err := m.parseEntity(parts[0], parts[len(parts)-1])
		if err != nil {
			return err
		}
		form.IsMain = true
		strs = append(strs, fmt.Sprintf("%s  %s", form.Entity.TableName, form.Alia))
	}
	value.Expr = strings.Join(strs, ",")
	return nil
}
func (m *exector) parseSelectField(value *oqlSelect) {
	value.Expr = m.replaceFieldString(value.Query)
}

// 解析实体
func (m *exector) parseEntity(id, path string) (*oqlEntity, error) {
	path = strings.ToLower(strings.TrimSpace(path))
	if v, ok := m.entities[path]; ok {
		return v, nil
	}
	entity := md.GetEntity(id)
	if entity == nil {
		return nil, glog.Errorf("找不到实体 %v", id)
	}
	v := m.formatEntity(entity)
	v.Sequence = len(m.entities) + 1
	v.Alia = fmt.Sprintf("a%v", v.Sequence)
	v.Path = path
	m.entities[path] = v
	return v, nil
}

// 解析字段
func (m *exector) parseField(fieldPath string) (*oqlField, error) {
	fieldPath = strings.ToLower(strings.TrimSpace(fieldPath))
	if v, ok := m.fields[fieldPath]; ok {
		return v, nil
	}
	var err error
	start := 0
	parts := strings.Split(fieldPath, ".")
	var mainFrom *oqlFrom
	if len(parts) > 1 {
		for i, v := range m.froms {
			if v.Alia != "" && strings.ToLower(v.Alia) == parts[0] {
				mainFrom = m.froms[i]
				start = 1
				break
			}
		}
	}
	if mainFrom == nil {
		for i, v := range m.froms {
			if v.Alia == "" {
				mainFrom = m.froms[i]
				break
			}
		}
	}
	if mainFrom == nil {
		mainFrom = m.froms[0]
	}
	entity := m.entities[strings.ToLower(mainFrom.Alia)]

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
			return nil, nil
		}
		field := m.formatField(entity, mdField)
		field.Path = path
		m.fields[path] = field
		if i < len(parts)-1 {
			entity, err = m.parseEntity(mdField.TypeID, path)
			if err != nil {
				return nil, err
			}

		} else {
			return field, nil
		}
	}
	return nil, nil
}

// fieldA
// fieldA,fieldB,fieldC
// fieldA as a,fieldB
func (m *exector) Select(query string, args ...interface{}) IExector {
	item := &oqlSelect{Query: query, Args: args}
	m.selects = append(m.selects, item)
	return m
}

// fieldA =?
// fieldB =''
// FieldC is null
// ($$FieldA +$$FieldB) >?
func (m *exector) Where(query string, args ...interface{}) IQWhere {
	item := &qWhere{Query: query, Args: args, Logical: "and"}
	m.wheres = append(m.wheres, item)
	return item
}
func (m *exector) OrWhere(query string, args ...interface{}) IQWhere {
	item := &qWhere{Query: query, Args: args, Logical: "or"}
	m.wheres = append(m.wheres, item)
	return item
}

// left join tableA on a=b
// left join tableA as c on c.a=b.b
// inner join tableA
func (m *exector) Joins(query string, args ...interface{}) IExector {
	item := &oqlJoin{Query: query, Args: args}
	m.joins = append(m.joins, item)
	return m
}

//fieldA
//fieldA,fieldB
func (m *exector) Group(query string, args ...interface{}) IExector {
	item := &oqlGroup{Query: query, Args: args}
	m.groups = append(m.groups, item)
	return m
}

//fieldA
//fieldA desc,fieldB
func (m *exector) Order(query string, args ...interface{}) IExector {
	item := &oqlOrder{Query: query, Args: args}
	m.orders = append(m.orders, item)
	return m
}

// tableA
// tableA a
// tableA as a
// tableA,tableB
func (m *exector) From(query string) IExector {
	items := strings.Split(strings.TrimSpace(query), ",")
	for _, item := range items {
		parts := strings.Split(strings.TrimSpace(item), " ")
		item := &oqlFrom{Query: item}
		if len(parts) > 1 {
			item.Alia = parts[len(parts)-1]
		}
		m.froms = append(m.froms, item)
	}
	return m
}
func (m *exector) SetContext(context *context.Context) IExector {
	m.context = context
	return m
}
func (m *exector) GetMainFrom() IQFrom {
	if m.froms == nil || len(m.froms) == 0 {
		return nil
	}
	return m.froms[0]
}
