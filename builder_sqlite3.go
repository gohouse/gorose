package gorose

type BuilderSqlite3 struct {
	//IOrm
	driver string
}

// sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
//		distinct, fields, table, join, where, group, having, order, limit, offset)
// select {distinct} {fields} from {table} {join} {where} {group} {having} {order} {limit} {offset}
// {execute} {table} {data} {where}
func init() {
	var driver = "sqlite3"
	var builder = &BuilderSqlite3{driver: driver}
	NewBuilderDriver().Register(driver, builder)
}

// BuildQuery : build query sql string
func (b *BuilderSqlite3) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o).SetDriver(b.driver).BuildQuery()
}

// BuildExecut : build execute sql string
func (b *BuilderSqlite3) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o).SetDriver(b.driver).BuildExecute(operType)
}
