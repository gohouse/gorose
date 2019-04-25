package gorose

type IBuilder interface {
	BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error)
	BuildExecute(o IOrm,operType string) (sql string, args []interface{}, err error)
}
