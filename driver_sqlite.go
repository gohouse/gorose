package gorose

import (
	"database/sql"
	//_ "github.com/mattn/go-sqlite3"
	"github.com/gohouse/utils"
)

func (this *Connection) Sqlite() {
	dbObj := Connect.CurrentConfig
	var err error

	dsn := dbObj["database"]
	DB, err = sql.Open("sqlite3", dsn)

	utils.CheckErr(err)
}