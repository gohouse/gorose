package gorose

import (
	"database/sql"
	"fmt"
	"github.com/gohouse/utils"
)

var (
	DB *sql.DB	// origin DB
	Tx *sql.Tx	// transaction
	Config map[string]map[string]string	// config
	SqlLogs []string	// all sql logs
)

func Open(arg ...interface{}) *sql.DB{
	if len(arg) == 1 {
		Connect(arg[0])
	} else {
		Config = arg[0].(map[string]map[string]string)
		Connect(arg[1])
	}

	return DB
}

func Connect(arg interface{}) *sql.DB {
	var err error
	var dbObj map[string]string

	if utils.GetType(arg) == "string" {
		dbObj = Config[arg.(string)]
	} else {
		dbObj = arg.(map[string]string)
	}

	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s", dbObj["username"], dbObj["password"], dbObj["protocol"], dbObj["host"], dbObj["port"], dbObj["database"], dbObj["charset"])
	//DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	DB, err = sql.Open("mysql", dsn)
	checkErr(err)

	err = DB.Ping()
	checkErr(err)

	return DB
}

func GetDB() *sql.DB {
	return DB
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
