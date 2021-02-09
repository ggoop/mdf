package md

const (
	REGEXP_OQL_FROM    = "([\\S]+)(?i:(?:as|[\\s])+)([\\S]+)|([\\S]+)"
	REGEXP_OQL_SELECT  = "([\\S]+.*\\S)(?i:\\s+as+\\s)([\\S]+)|([\\S]+.*[\\S]+)"
	REGEXP_OQL_ORDER   = "(?i)([\\S]+.*\\S)(?:\\s)(desc|asc)|([\\S]+.*[\\S]+)"
	REGEXP_OQL_VAR_EXP = `{([A-Za-z._]+[0-9A-Za-z]*)}`
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
)

type OQLOption struct {
}

func GetOQL(names ...OQLOption) *OQL {
	oql := OQL{}
	oql.errors = make([]error, 0)
	oql.entities = make(map[string]*oqlEntity)
	oql.fields = make(map[string]*oqlField)
	oql.froms = make([]*oqlFrom, 0)
	oql.joins = make([]*oqlJoin, 0)
	oql.selects = make([]*oqlSelect, 0)
	oql.orders = make([]*oqlOrder, 0)
	oql.wheres = make([]*OQLWhere, 0)
	oql.groups = make([]*oqlGroup, 0)
	oql.having = make([]*OQLWhere, 0)
	return &oql
}
