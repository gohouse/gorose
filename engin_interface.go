package gorose

import "database/sql"

type IEngin interface {
	GetExecuteDB() (db *sql.DB, driver string)
	GetQueryDB() (db *sql.DB, driver string)
	EnableQueryLog(e ...bool)
	IfEnableQueryLog() (e bool)
	Prefix(pre string)
	GetPrefix() (pre string)
	NewSession() ISession
}
