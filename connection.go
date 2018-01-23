package gorose

import (
	"database/sql"
	"github.com/gohouse/utils"
)

var (
	DB *sql.DB // origin DB
	Tx *sql.Tx // transaction
	//Stmt *sql.Stmt
	Connect Connection
)

type Connection struct {
	DbConfig      map[string]map[string]string
	CurrentConfig map[string]string
}

//type aaa string

func Open(args ...interface{}) Database {
	if len(args) == 1 {
		if confReal, ok := args[0].(map[string]string); ok {
			Connect.Boot(confReal)
		} else {
			panic("配置文件格式有误!")
		}
	} else if len(args) == 2 {
		if confReal, ok := args[0].(map[string]string); ok {
			Connect.Boot(confReal)
		} else if confReal, ok := args[0].(map[string]map[string]string); ok {
			Connect.DbConfig = confReal
			if confReal, ok := args[1].(string); ok {
				Connect.Boot(confReal)
			} else {
				panic("选择默认数据库格式有误!")
			}
		} else {
			panic("配置文件格式有误!")
		}
	}

	return Dbstruct
}

func (this *Connection) Boot(arg interface{}) *sql.DB {
	if argReal, ok := arg.(string); ok {
		Connect.CurrentConfig = Connect.DbConfig[argReal]
	} else if argReal, ok := arg.(map[string]string); ok {
		Connect.CurrentConfig = argReal
	}

	// get driver
	this.getDriver()

	err := DB.Ping()
	utils.CheckErr(err)

	return DB
}

func (this *Connection) getDriver() {
	var err error
	dbObj := Connect.CurrentConfig

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
	case "mssql":
		this.MsSQL()
	}

	utils.CheckErr(err)
}

func GetDB() *sql.DB {
	return DB
}
