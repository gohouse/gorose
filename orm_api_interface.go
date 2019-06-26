package gorose

type IOrmApi interface {
	GetTable() string
	GetFields() []string
	SetWhere(arg [][]interface{})
	GetWhere() [][]interface{}
	GetOrder() string
	GetLimit() int
	GetOffset() int
	GetJoin() [][]interface{}
	GetDistinct() bool
	GetGroup() string
	GetHaving() string
	GetData() interface{}
	ExtraExecCols(args ...string) IOrm
	ResetExtraExecCols() IOrm
	GetExtraExecCols() []string
}