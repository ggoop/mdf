package md

type OQLWhere interface {
	And() OQLWhere
	Or() OQLWhere

	Where(query string, args ...interface{}) OQLWhere
	OrWhere(query string, args ...interface{}) OQLWhere

	Query() string
	Logical() string
	DataType() string
	Children() []OQLWhere
	Args() []interface{}
	Expr() string
	setExpr(expr string)
}
type oqlWhere struct {
	//字段与操作号之间需要有空格
	//示例1: Org =? ; Org in (?) ;$$Org =?  and ($$Period = ?  or $$Period = ? )
	//示例2：abs($$Qty)>$$TempQty + ?
	query   string
	logical string //and or
	//参数值数据类型
	dataType string
	sequence int
	children []OQLWhere
	args     []interface{}
	expr     string
}

func NewOQLWhere(query string, args ...interface{}) OQLWhere {
	return &oqlWhere{query: query, args: args, logical: OQL_WHERE_LOGICAL_AND}
}
func (m oqlWhere) Query() string {
	return m.query
}
func (m oqlWhere) Logical() string {
	return m.logical
}
func (m oqlWhere) DataType() string {
	return m.dataType
}
func (m oqlWhere) Children() []OQLWhere {
	return m.children
}

func (m oqlWhere) Args() []interface{} {
	return m.args
}
func (m oqlWhere) Expr() string {
	return m.expr
}

func (m oqlWhere) setExpr(expr string) {
	m.expr = expr
}
func (m oqlWhere) String() string {
	return m.query
}
func (m *oqlWhere) Where(query string, args ...interface{}) OQLWhere {
	if m.children == nil {
		m.children = make([]OQLWhere, 0)
	}
	item := &oqlWhere{query: query, args: args, logical: OQL_WHERE_LOGICAL_AND}
	m.children = append(m.children, item)
	return m
}
func (m *oqlWhere) OrWhere(query string, args ...interface{}) OQLWhere {
	if m.children == nil {
		m.children = make([]OQLWhere, 0)
	}
	item := &oqlWhere{query: query, args: args, logical: OQL_WHERE_LOGICAL_OR}
	m.children = append(m.children, item)
	return m
}
func (m *oqlWhere) And() OQLWhere {
	if m.children == nil {
		m.children = make([]OQLWhere, 0)
	}
	item := &oqlWhere{logical: OQL_WHERE_LOGICAL_AND}
	m.children = append(m.children, item)
	return item
}
func (m *oqlWhere) Or() OQLWhere {
	if m.children == nil {
		m.children = make([]OQLWhere, 0)
	}
	item := &oqlWhere{logical: OQL_WHERE_LOGICAL_OR}
	m.children = append(m.children, item)
	return item
}
