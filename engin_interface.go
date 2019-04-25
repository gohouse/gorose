package gorose

import "database/sql"

type IEngin interface {
	GetExecuteDB() *sql.DB
	GetQueryDB() *sql.DB
	EnableQueryLog(e ...bool)
	IfEnableQueryLog() (e bool)
	Prefix(pre string)
	GetPrefix() (pre string)
	NewSession() ISession
}
