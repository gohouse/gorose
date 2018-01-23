package gorose

import (
	"fmt"
	"database/sql"
	//_ "github.com/lib/pq"
	"github.com/gohouse/utils"
)

func (this *Connection) Postgres() {
	dbObj := Connect.CurrentConfig
	var err error

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbObj["host"], dbObj["port"], dbObj["username"], dbObj["password"], dbObj["database"])

	DB, err = sql.Open("postgres", dsn)

	utils.CheckErr(err)
}
