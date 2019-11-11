package gorose

// IOrmApi ...
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
	ExtraCols(args ...string) IOrm
	ResetExtraCols() IOrm
	GetExtraCols() []string
	GetPessimisticLock() string
}
