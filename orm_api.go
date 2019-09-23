package gorose

type OrmApi struct {
	table           string
	fields          []string
	where           [][]interface{}
	order           string
	limit           int
	offset          int
	join            [][]interface{}
	distinct        bool
	union           string
	group           string
	having          string
	data            interface{}
	force           bool
	extraCols       []string
	// 悲观锁
	pessimisticLock string
}

func (o *Orm) GetTable() string {
	return o.table
}

func (o *Orm) GetFields() []string {
	return o.fields
}

func (o *Orm) SetWhere(arg [][]interface{}) {
	o.where = arg
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

func (o *Orm) GetGroup() string {
	return o.group
}

func (o *Orm) GetHaving() string {
	return o.having
}

func (o *Orm) GetData() interface{} {
	return o.data
}

func (dba *Orm) GetForce() bool {
	return dba.force
}

func (dba *Orm) GetExtraCols() []string {
	return dba.extraCols
}

func (dba *Orm) GetPessimisticLock() string {
	return dba.pessimisticLock
}
