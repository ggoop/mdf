package md

type IQWhere interface {
	Where(query string, args ...interface{}) IQWhere
	OrWhere(query string, args ...interface{}) IQWhere
	And() IQWhere
	Or() IQWhere
	GetArgs() []interface{}
	GetQuery() string
	GetDataType() string
	SetDataType(dataType string) IQWhere
	GetLogical() string
	String() string
}
type qWhere struct {
	//字段与操作号之间需要有空格
	//示例1: Org =? ; Org in (?) ;$$Org =?  and ($$Period = ?  or $$Period = ? )
	//示例2：abs($$Qty)>$$TempQty + ?
	Query    string
	DataType string
	Logical  string //and or
	Sequence int
	Children []*qWhere
	Expr     string
	Args     []interface{}
}

func NewQWhere(logical, query string, args ...interface{}) IQWhere {
	return &qWhere{Query: query, Args: args, Logical: logical}
}
func (m qWhere) String() string {
	return m.Query
}
func (m qWhere) GetDataType() string {
	return m.DataType
}
func (m *qWhere) SetDataType(dataType string) IQWhere {
	m.DataType = dataType
	return m
}
func (m qWhere) GetArgs() []interface{} {
	return m.Args
}
func (m qWhere) GetQuery() string {
	return m.Query
}
func (m qWhere) GetLogical() string {
	return m.Logical
}
func (m *qWhere) Where(query string, args ...interface{}) IQWhere {
	if m.Children == nil {
		m.Children = make([]*qWhere, 0)
	}
	item := &qWhere{Query: query, Args: args, Logical: "and"}
	m.Children = append(m.Children, item)
	return m
}
func (m *qWhere) OrWhere(query string, args ...interface{}) IQWhere {
	if m.Children == nil {
		m.Children = make([]*qWhere, 0)
	}
	item := &qWhere{Query: query, Args: args, Logical: "or"}
	m.Children = append(m.Children, item)
	return m
}
func (m *qWhere) And() IQWhere {
	if m.Children == nil {
		m.Children = make([]*qWhere, 0)
	}
	item := &qWhere{Logical: "and"}
	m.Children = append(m.Children, item)
	return item
}
func (m *qWhere) Or() IQWhere {
	if m.Children == nil {
		m.Children = make([]*qWhere, 0)
	}
	item := &qWhere{Logical: "or"}
	m.Children = append(m.Children, item)
	return item
}
