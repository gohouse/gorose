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
	SetMaxOpenConns	int = 0
	SetMaxIdleConns	int = 1
)

type Connection struct {
	DbConfig      map[string]interface{}
	Default       string
	CurrentConfig map[string]string
}

func Open(args ...interface{}) (Database, error) {
	if len(args) == 1 {
		// continue
	} else if len(args) == 2 {
		if confReal, ok := args[1].(string); ok {
			Connect.Default = confReal
		} else {
			return Database{}, errors.New("指定默认数据库只能位字符串!")
		}
	} else {
		return Database{},errors.New("Open方法只接收1个或2个参数!")
	}
	// 解析config
	err := Connect.parseConfig(args[0])
	if err!=nil{
		return Database{},err
	}

	// 驱动数据库
	errs := Connect.boot()

	return Dbstruct, errs
}

func (this *Connection) parseConfig(args interface{}) error {
	if confReal, ok := args.(map[string]string); ok {
		Connect.CurrentConfig = confReal
	} else if confReal, ok := args.(map[string]interface{}); ok {
		Connect.DbConfig = confReal
		if defaultDb,ok := confReal["default"]; ok{
			if Connect.Default == "" {
				Connect.Default = defaultDb.(string)
			}
		}
		if Connect.Default == ""{
			return errors.New("配置文件默认数据库链接未设置!")
		}
		// 获取指定的默认数据库链接信息
		if defaultDbConnection,ok := confReal[Connect.Default]; ok{
			if configs,ok := defaultDbConnection.(map[string]string);ok{
				Connect.CurrentConfig = configs
			} else {
				return errors.New("数据库配置格式有误!")
			}
		} else {
			return errors.New("指定的数据库链接不存在!")
		}
		// 设置连接池信息
		if mo,ok := confReal["SetMaxOpenConns"];ok{
			if moInt,ok := mo.(int);ok{
				SetMaxOpenConns	= moInt
			} else {
				return errors.New("连接池信息配置的值只能是数字")
			}
		}
		if mi,ok := confReal["SetMaxIdleConns"];ok{
			if miInt,ok := mi.(int);ok{
				SetMaxIdleConns	= miInt
			} else {
				return errors.New("连接池信息配置的值只能是数字")
			}
		}
	} else {
		return errors.New("配置文件格式有误2!")
	}
	return nil
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
	DB.SetMaxOpenConns(SetMaxOpenConns)
	DB.SetMaxIdleConns(SetMaxIdleConns)

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
