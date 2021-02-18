package md

type OQLStatement struct {
	Query    string
	Args     []interface{}
	Affected int64
}
type oqlEntity struct {
	Path     string
	Entity   *MDEntity
	Sequence int
	IsMain   bool
	Alias    string
}
type oqlField struct {
	Entity *oqlEntity
	Field  *MDField
	Path   string
}

type OQLFrom interface {
	Query() string
	Alias() string
	Args() []interface{}
	Expr() string
	setExpr(expr string)
}

type oqlFrom struct {
	query string
	alias string
	args  []interface{}
	expr  string
}

func (m *oqlFrom) Query() string {
	return m.query
}
func (m *oqlFrom) Alias() string {
	return m.alias
}
func (m *oqlFrom) Args() []interface{} {
	return m.args
}
func (m *oqlFrom) Expr() string {
	return m.expr
}

func (m oqlFrom) setExpr(expr string) {
	m.expr = expr
}

type oqlJoin struct {
	Type      OQLJoinType
	Query     string
	Alias     string
	Condition string
	Args      []interface{}
	expr      string
}
type oqlSelect struct {
	Query string
	Alias string
	Args  []interface{}
	expr  string
}
type oqlGroup struct {
	Query string
	Args  []interface{}
	expr  string
}
type oqlOrder struct {
	Query string
	Order OQLOrderType
	Args  []interface{}
	expr  string
}
