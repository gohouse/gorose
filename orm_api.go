package gorose

// OrmApi ...
type OrmApi struct {
	table    string
	fields   []string
	where    [][]interface{}
	order    string
	limit    int
	offset   int
	join     [][]interface{}
	distinct bool
	//union     string
	group     string
	having    string
	data      interface{}
	force     bool
	extraCols []string
	// 悲观锁
	pessimisticLock string
}

// GetTable ...
func (o *Orm) GetTable() string {
	return o.table
}

// GetFields ...
func (o *Orm) GetFields() []string {
	return o.fields
}

// SetWhere ...
func (o *Orm) SetWhere(arg [][]interface{}) {
	o.where = arg
}

// GetWhere ...
func (o *Orm) GetWhere() [][]interface{} {
	return o.where
}

// GetOrder ...
func (o *Orm) GetOrder() string {
	return o.order
}

// GetLimit ...
func (o *Orm) GetLimit() int {
	return o.limit
}

// GetOffset ...
func (o *Orm) GetOffset() int {
	return o.offset
}

// GetJoin ...
func (o *Orm) GetJoin() [][]interface{} {
	return o.join
}

// GetDistinct ...
func (o *Orm) GetDistinct() bool {
	return o.distinct
}

// GetGroup ...
func (o *Orm) GetGroup() string {
	return o.group
}

// GetHaving ...
func (o *Orm) GetHaving() string {
	return o.having
}

// GetData ...
func (o *Orm) GetData() interface{} {
	return o.data
}

// GetForce ...
func (o *Orm) GetForce() bool {
	return o.force
}

// GetExtraCols ...
func (o *Orm) GetExtraCols() []string {
	return o.extraCols
}

// GetPessimisticLock ...
func (o *Orm) GetPessimisticLock() string {
	return o.pessimisticLock
}
