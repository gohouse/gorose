package gorose

import (
	"database/sql"
	"errors"
	"github.com/gohouse/gorose/drivers"
)

var (
	DB *sql.DB // origin DB
	Tx *sql.Tx // transaction
	//Stmt *sql.Stmt
	Connect Connection
)

type Connection struct {
	DbConfig      map[string]map[string]string
	Default       string
	CurrentConfig map[string]string
}

//type aaa string

func Open(args ...interface{}) (Database, error) {
	if len(args) == 1 {
		if confReal, ok := args[0].(map[string]string); ok {
			Connect.CurrentConfig = confReal
			//Connect.Boot(confReal)
		} else {
			return Database{}, errors.New("配置文件格式有误!")
		}
	} else if len(args) == 2 {
		if confReal, ok := args[0].(map[string]string); ok {
			Connect.CurrentConfig = confReal
			//Connect.Boot(confReal)
		} else if confReal, ok := args[0].(map[string]map[string]string); ok {
			Connect.DbConfig = confReal
			if confReal, ok := args[1].(string); ok {
				Connect.CurrentConfig = Connect.DbConfig[confReal]
				//Connect.Boot(confReal)
			} else {
				return Database{}, errors.New("配置文件格式有误!")
			}
		} else {
			return Database{}, errors.New("配置文件格式有误!")
		}
	}

	err := Connect.boot()

	return Dbstruct, err
}

func (this *Connection) boot() error {
	dbObj := Connect.CurrentConfig
	var driver, dsn string
	var err error

	//DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	switch dbObj["driver"] {
	case "mysql":
		driver, dsn = drivers.MySQL(dbObj)
	case "sqlite3":
		driver, dsn = drivers.Sqlite3(dbObj)
	case "postgres":
		driver, dsn = drivers.Postgres(dbObj)
	case "oracle":
		driver, dsn = drivers.Oracle(dbObj)
	case "mssql":
		driver, dsn = drivers.MsSQL(dbObj)
	}

	// 开始驱动
	DB, err = sql.Open(driver, dsn)

	if err != nil {
		return err
	}

	// 检查是否可以ping通
	err2 := DB.Ping()

	return err2
}

func GetDB() *sql.DB {
	return DB
}
