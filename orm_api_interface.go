package gorose

type IOrmApi interface {
	//SetTable(arg string)
	GetTable() string
	//SetFields(arg []string)
	GetFields() []string
	//SetWhere(arg [][]interface{})
	GetWhere() [][]interface{}
	//SetOrder(arg string)
	GetOrder() string
	//SetLimit(arg int)
	GetLimit() int
	//SetOffset(arg int)
	GetOffset() int
	//SetJoin(arg [][]interface{})
	GetJoin() [][]interface{}
	//SetDistinct(arg bool)
	GetDistinct() bool
	//SetUnion(arg string)
	GetUnion() string
	//SetGroup(arg string)
	GetGroup() string
	//SetHaving(arg string)
	GetHaving() string
	//SetData(arg interface{})
	GetData() interface{}
}
