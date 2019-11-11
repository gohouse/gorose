package gorose

const (
	// DriverMysql ...
	DriverMysql = "mysql"
)

// BuilderMysql ...
type BuilderMysql struct {
	//IOrm
	driver string
}

var _ IBuilder = (*BuilderMysql)(nil)

// sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
//		distinct, fields, table, join, where, group, having, order, limit, offset)
// select {distinct} {fields} from {table} {join} {where} {group} {having} {order} {limit} {offset}
// {execute} {table} {data} {where}
func init() {
	NewBuilderDriver().Register(DriverMysql, NewBuilderMysql())
}

// NewBuilderMysql ...
func NewBuilderMysql() *BuilderMysql {
	return new(BuilderMysql)
}

// Clone : a new obj
func (b *BuilderMysql) Clone() IBuilder {
	return NewBuilderMysql()
}

// BuildQuery : build query sql string
func (b *BuilderMysql) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	//fmt.Println(o.GetTable(),o.GetWhere())
	sqlStr, args, err = NewBuilderDefault(o).SetDriver(DriverMysql).BuildQuery()
	return
}

// BuildExecut : build execute sql string
func (b *BuilderMysql) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o).SetDriver(DriverMysql).BuildExecute(operType)
}
