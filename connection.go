package gorose

import (
	"database/sql"
	"github.com/gohouse/utils"
)

var (
	DB *sql.DB	// origin DB
	Tx *sql.Tx	// transaction
	Conf map[string]map[string]string	// config
	SqlLogs []string	// all sql logs
	CurrentConfig map[string]string
	//conn Connection
	Connect Connection
	JsonEncode bool
)

type Connection struct {

}

func (this *Connection) Open(arg ...interface{}) *sql.DB{
	if len(arg) == 1 {
		this.Connect(arg[0])
	} else {
		Conf = arg[0].(map[string]map[string]string)
		this.Connect(arg[1])
	}

	return DB
}

func (this *Connection) Connect(arg interface{}) *sql.DB {
	if utils.GetType(arg) == "string" {
		CurrentConfig = Conf[arg.(string)]
	} else {
		CurrentConfig = arg.(map[string]string)
	}

	// get driver
	this.getDriver()

	var err error = DB.Ping()
	CheckErr(err)

	return DB
}

func (this *Connection) getDriver() {
	var err error
	dbObj := CurrentConfig

	//DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	switch dbObj["driver"] {
	case "mysql":
		this.MySQL()
	case "sqlite":
		this.Sqlite()
	case "postgre":
		this.Postgres()
	case "oracle":
		this.Oracle()
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
