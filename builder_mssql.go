package gorose

const (
	// DriverMsSql ...
	DriverMsSql = "mssql"
)

// BuilderMsSql ...
type BuilderMsSql struct {
	//IOrm
	driver string
}

// sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
//		distinct, fields, table, join, where, group, having, order, limit, offset)
// select {distinct} {fields} from {table} {join} {where} {group} {having} {order} {limit} {offset}
// {execute} {table} {data} {where}
func init() {
	var builder = &BuilderMsSql{driver: DriverMsSql}
	NewBuilderDriver().Register(DriverMsSql, builder)
}

// Clone : a new obj
func (b *BuilderMsSql) Clone() IBuilder {
	return &BuilderMsSql{driver: DriverMsSql}
}

// BuildQuery : build query sql string
func (b *BuilderMsSql) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o).SetDriver(b.driver).BuildQuery()
}

// BuildExecut : build execute sql string
func (b *BuilderMsSql) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o).SetDriver(b.driver).BuildExecute(operType)
}
