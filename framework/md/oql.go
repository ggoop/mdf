package md

import "github.com/ggoop/mdf/utils"

const (
	regexp_OQL_FROM    = "([\\S]+)(?i:(?:as|[\\s])+)([\\S]+)|([\\S]+)"
	regexp_OQL_SELECT  = "([\\S]+.*\\S)(?i:\\s+as+\\s)([\\S]+)|([\\S]+.*[\\S]+)"
	regexp_OQL_ORDER   = "(?i)([\\S]+.*\\S)(?:\\s)(desc|asc)|([\\S]+.*[\\S]+)"
	regexp_OQL_VAR_EXP = `{([A-Za-z._]+[0-9A-Za-z]*)}`
)

type OQLJoinType int32
type OQLOrderType int32

const (
	OQL_LEFT_JOIN  OQLJoinType = 0
	OQL_RIGHT_JOIN OQLJoinType = 1
	OQL_FULL_JOIN  OQLJoinType = 2
	OQL_UNION_JOIN OQLJoinType = 3

	OQL_ORDER_DESC OQLOrderType = -1
	OQL_ORDER_ASC  OQLOrderType = 1

	OQL_WHERE_LOGICAL_OR  = "or"
	OQL_WHERE_LOGICAL_AND = "and"
)

type OQLOption struct {
}
type OQL interface {
	Error() error
	AddErr(err error) OQL

	Clone() OQL

	SetContext(context *utils.TokenContext) OQL
	SetActuator(actuator OQLActuator) OQL
	From(query interface{}, args ...interface{}) OQL
	Join(joinType OQLJoinType, query string, condition string, args ...interface{}) OQL
	Select(query interface{}, args ...interface{}) OQL
	Order(query interface{}, args ...interface{}) OQL
	Group(query interface{}, args ...interface{}) OQL
	Or() OQLWhere
	Where(query interface{}, args ...interface{}) OQLWhere
	Having(query interface{}, args ...interface{}) OQLWhere

	Count(value interface{}) OQL
	Pluck(column string, value interface{}) OQL
	Take(out interface{}) OQL
	Find(out interface{}) OQL
	Paginate(value interface{}, page int64, pageSize int64) OQL
	Create(data interface{}) OQL
	Update(data interface{}) OQL
	Delete() OQL
}

func NewOQL(names ...OQLOption) OQL {
	oql := &oqlImpl{}
	oql.errors = make([]error, 0)
	oql.entities = make(map[string]*oqlEntity)
	oql.fields = make(map[string]*oqlField)
	oql.froms = make([]OQLFrom, 0)
	oql.joins = make([]*oqlJoin, 0)
	oql.selects = make([]*oqlSelect, 0)
	oql.orders = make([]*oqlOrder, 0)
	oql.wheres = make([]OQLWhere, 0)
	oql.groups = make([]*oqlGroup, 0)
	oql.having = make([]OQLWhere, 0)
	return oql
}

//公共查询
type oqlImpl struct {
	error    error
	errors   []error
	entities map[string]*oqlEntity
	fields   map[string]*oqlField
	froms    []OQLFrom
	joins    []*oqlJoin
	selects  []*oqlSelect
	orders   []*oqlOrder
	wheres   []OQLWhere
	groups   []*oqlGroup
	having   []OQLWhere
	offset   int64
	limit    int64
	context  *utils.TokenContext
	actuator OQLActuator
}
