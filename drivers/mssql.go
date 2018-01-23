package gorose

import (
	"database/sql"
	"fmt"
	//_ "github.com/denisenkom/go-mssqldb"
	"github.com/gohouse/utils"
)

func (this *Connection) MsSQL() {
	dbObj := Connect.CurrentConfig
	var err error

	dsn := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s",
		dbObj["host"], dbObj["port"], dbObj["database"], dbObj["username"], dbObj["password"])
	DB, err = sql.Open("mssql", dsn)

	utils.CheckErr(err)
}
