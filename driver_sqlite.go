package gorose

import (
	"database/sql"
	//_ "github.com/mattn/go-sqlite3"
)

func (this *Connection) Sqlite() {
	dbObj := CurrentConfig
	var err error

	dsn := dbObj["database"]
	DB, err = sql.Open("sqlite3", dsn)

	CheckErr(err)
}