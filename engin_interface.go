package gorose

type IEngin interface {
	GetExecuteDB() dbObject
	GetQueryDB() dbObject
	EnableSqlLog(e ...bool)
	IfEnableSqlLog() (e bool)
	Prefix(pre string)
	GetPrefix() (pre string)
	//NewSession() ISession
}
