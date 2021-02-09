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
