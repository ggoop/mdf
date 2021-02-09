package md

type OQLWhere struct {
	//字段与操作号之间需要有空格
	//示例1: Org =? ; Org in (?) ;$$Org =?  and ($$Period = ?  or $$Period = ? )
	//示例2：abs($$Qty)>$$TempQty + ?
	Query   string
	Logical string //and or
	//参数值数据类型
	DataType string
	Sequence int
	Children []*OQLWhere
	Args     []interface{}
	expr     string
}

func NewOQLWhere(query string, args ...interface{}) *OQLWhere {
	return &OQLWhere{Query: query, Args: args, Logical: "and"}
}
func (m OQLWhere) String() string {
	return m.Query
}
func (m *OQLWhere) Where(query string, args ...interface{}) *OQLWhere {
	if m.Children == nil {
		m.Children = make([]*OQLWhere, 0)
	}
	item := &OQLWhere{Query: query, Args: args, Logical: "and"}
	m.Children = append(m.Children, item)
	return m
}
func (m *OQLWhere) OrWhere(query string, args ...interface{}) *OQLWhere {
	if m.Children == nil {
		m.Children = make([]*OQLWhere, 0)
	}
	item := &OQLWhere{Query: query, Args: args, Logical: "or"}
	m.Children = append(m.Children, item)
	return m
}
func (m *OQLWhere) And() *OQLWhere {
	if m.Children == nil {
		m.Children = make([]*OQLWhere, 0)
	}
	item := &OQLWhere{Logical: "and"}
	m.Children = append(m.Children, item)
	return item
}
func (m *OQLWhere) Or() *OQLWhere {
	if m.Children == nil {
		m.Children = make([]*OQLWhere, 0)
	}
	item := &OQLWhere{Logical: "or"}
	m.Children = append(m.Children, item)
	return item
}
