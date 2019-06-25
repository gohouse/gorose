package gorose

import "database/sql"

type IEngin interface {
	GetExecuteDB() *sql.DB
	GetQueryDB() *sql.DB
	//EnableSqlLog(e ...bool)
	//IfEnableSqlLog() (e bool)
	//SetPrefix(pre string)
	GetPrefix() (pre string)
	//NewSession() ISession
	//NewOrm() IOrm
	GetLogger() ILogger
	GetDriver() string
}
