package gorose

import (
	"fmt"
	"database/sql"
	//_ "github.com/go-sql-driver/mysql"
)

func (this *Connection) MySQL() {
	dbObj := CurrentConfig
	var err error

	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s",
		dbObj["username"], dbObj["password"], dbObj["protocol"], dbObj["host"], dbObj["port"], dbObj["database"], dbObj["charset"])
	DB, err = sql.Open("mysql", dsn)

	CheckErr(err)
}