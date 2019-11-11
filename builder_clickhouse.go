package gorose

const (
	// DriverClickhouse ...
	DriverClickhouse = "clickhouse"
)

// BuilderClickhouse ...
type BuilderClickhouse struct {
	//IOrm
	driver string
}

// sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
//		distinct, fields, table, join, where, group, having, order, limit, offset)
// select {distinct} {fields} from {table} {join} {where} {group} {having} {order} {limit} {offset}
// {execute} {table} {data} {where}
func init() {
	var builder = &BuilderClickhouse{}
	NewBuilderDriver().Register(DriverClickhouse, builder)
}

// Clone : a new obj
func (b *BuilderClickhouse) Clone() IBuilder {
	return &BuilderClickhouse{}
}

// BuildQuery : build query sql string
func (b *BuilderClickhouse) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o).SetDriver(DriverClickhouse).BuildQuery()
}

// BuildExecut : build execute sql string
func (b *BuilderClickhouse) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o).SetDriver(DriverClickhouse).BuildExecute(operType)
}
