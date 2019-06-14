package gorose

//import "sync"
//
//type OrmApi struct {
//	table    string
//	fields   []string
//	where    [][]interface{}
//	order    string
//	limit    int
//	offset   int
//	join     [][]interface{}
//	distinct bool
//	union    string
//	group    string
//	having   string
//	data     interface{}
//}


func (o *Orm) GetTable() string {
	return o.table
}

func (o *Orm) GetFields() []string {
	return o.fields
}

func (o *Orm) GetWhere() [][]interface{} {
	return o.where
}

func (o *Orm) GetOrder() string {
	return o.order
}

func (o *Orm) GetLimit() int {
	return o.limit
}

func (o *Orm) GetOffset() int {
	return o.offset
}

func (o *Orm) GetJoin() [][]interface{} {
	return o.join
}

func (o *Orm) GetDistinct() bool {
	return o.distinct
}

func (o *Orm) GetUnion() string {
	return o.union
}

func (o *Orm) GetGroup() string {
	return o.group
}

func (o *Orm) GetHaving() string {
	return o.having
}

func (o *Orm) GetData() interface{} {
	return o.data
}
