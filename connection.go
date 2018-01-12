package gorose

import (
	"database/sql"
	"github.com/gohouse/utils"
)

var (
	DB *sql.DB	// origin DB
	Tx *sql.Tx	// transaction
	Config map[string]map[string]string	// config
	SqlLogs []string	// all sql logs
	CurrentConfig map[string]string
)

type Connection struct {

}

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
	if utils.GetType(arg) == "string" {
		CurrentConfig = Config[arg.(string)]
	} else {
		CurrentConfig = arg.(map[string]string)
	}

	// get driver
	getDriver()

	var err error = DB.Ping()
	CheckErr(err)

	return DB
}

func getDriver() {
	var err error
	dbObj := CurrentConfig

	//DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	switch dbObj["driver"] {
	case "mysql":
		MySQL()
	case "sqlite":
		Sqlite()
	case "postgre":
		Postgres()
	case "oracle":
		Oracle()
	}

	CheckErr(err)
}

func GetDB() *sql.DB {
	return DB
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
