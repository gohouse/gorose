package gorose

// IBuilder ...
type IBuilder interface {
	BuildQuery(orm IOrm) (sqlStr string, args []interface{}, err error)
	BuildExecute(orm IOrm, operType string) (sqlStr string, args []interface{}, err error)
	Clone() IBuilder
	//GetIOrm() IOrm
}
