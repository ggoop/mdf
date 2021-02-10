package md

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

type IQFrom interface {
	GetQuery() string
	GetAlias() string
	GetExpr() string
}

func (m *oqlFrom) GetQuery() string {
	return m.Query
}
func (m *oqlFrom) GetAlias() string {
	return m.Alias
}
func (m *oqlFrom) GetExpr() string {
	return m.expr
}

type oqlFrom struct {
	Query string
	Alias string
	Args  []interface{}
	expr  string
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
