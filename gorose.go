package gorose

import (
	"database/sql"
	"errors"
	"github.com/gohouse/gorose/drivers"
)

var (
	DB *sql.DB // origin DB
	Tx *sql.Tx // transaction DB
	//Stmt *sql.Stmt
	Connect Connection
	//this.SetMaxOpenConns int = 0
	//this.SetMaxIdleConns int = -1
)

func init() {
	Connect.SetMaxOpenConns = 0
	Connect.SetMaxIdleConns = -1
}

type Connection struct {
	// all config sets
	DbConfig map[string]interface{}
	// default database
	Default string
	// current config on use
	CurrentConfig map[string]string
	// all sql logs
	SqlLog []string
	// if in transaction, the code auto change
	Trans bool
	// max open connections
	SetMaxOpenConns int
	// max freedom connections leave
	SetMaxIdleConns int
}
// open instance of sql.DB.Oper
func Open(args ...interface{}) (Connection, error) {
	if len(args) == 1 {
		// continue
	} else if len(args) == 2 {
		if confReal, ok := args[1].(string); ok {
			Connect.Default = confReal
		} else {
			// 指定默认数据库只能为字符串!
			return Connect, errors.New("only str allowed of default database name")
		}
	} else {
		// Open方法只接收1个或2个参数!
		return Connect, errors.New("1 or 2 params need in Open() method")
	}
	// 解析config
	err := Connect.parseConfig(args[0])
	if err != nil {
		return Connect, err
	}

	// 驱动数据库
	errs := Connect.boot()

	return Connect, errs
}
// parse input config
func (this *Connection) parseConfig(args interface{}) error {
	if confReal, ok := args.(map[string]string); ok {
		Connect.CurrentConfig = confReal
	} else if confReal, ok := args.(map[string]interface{}); ok {
		Connect.DbConfig = confReal
		if defaultDb, ok := confReal["default"]; ok {
			if Connect.Default == "" {
				Connect.Default = defaultDb.(string)
			}
		}
		if Connect.Default == "" {
			// 配置文件默认数据库链接未设置
			return errors.New("the default database is missing in config!")
		}
		// 获取指定的默认数据库链接信息
		if defaultDbConnection, ok := confReal[Connect.Default]; ok {
			if configs, ok := defaultDbConnection.(map[string]string); ok {
				Connect.CurrentConfig = configs
			} else {
				// 数据库配置格式有误!
				return errors.New("format error in database config!")
			}
		} else {
			// 指定的数据库链接不存在!
			return errors.New("the database for using is missing!")
		}
		// 设置连接池信息
		if mo, ok := confReal["SetMaxOpenConns"]; ok {
			if moInt, ok := mo.(int); ok {
				this.SetMaxOpenConns = moInt
			} else {
				// 连接池信息配置的值只能是数字
				return errors.New("the value of connection pool config need int")
			}
		}
		if mi, ok := confReal["SetMaxIdleConns"]; ok {
			if miInt, ok := mi.(int); ok {
				this.SetMaxIdleConns = miInt
			} else {
				return errors.New("the value of connection pool config need int")
			}
		}
	} else {
		return errors.New("format error in database config!")
	}
	return nil
}
// boot sql driver
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
	DB.SetMaxOpenConns(this.SetMaxOpenConns)
	DB.SetMaxIdleConns(this.SetMaxIdleConns)

	if err != nil {
		return err
	}

	// 检查是否可以ping通
	err2 := DB.Ping()

	return err2
}
// close database
func (this *Connection) Close() error {
	Connect.SqlLog = []string{}
	return DB.Close()
}
// ping db
func (this *Connection) Ping() error {
	return DB.Ping()
}
// set table from database
func (this *Connection) Table(table string) *Database {
	//this.table = table
	var database Database
	return database.Table(table)
}
// transaction begin
func (this *Connection) Begin() {
	Tx, _ = DB.Begin()
	Connect.Trans = true
}
// transaction commit
func (this *Connection) Commit() {
	Tx.Commit()
	Connect.Trans = false
}
// transaction rollback
func (this *Connection) Rollback() {
	Tx.Rollback()
	Connect.Trans = false
}
// simple transaction
func (this *Connection) Transaction(closure func() error) bool {
	//defer func() {
	//	if err := recover(); err != nil {
	//		this.Rollback()
	//		panic(err)
	//	}
	//}()

	this.Begin()
	err := closure()
	if err != nil {
		this.Rollback()
		return false
	}
	this.Commit()

	return true
}
// query str
func (this *Connection) Query(args ...interface{}) ([]map[string]interface{}, error) {
	var database Database
	return database.Query(args...)
}
// execute str
func (this *Connection) Execute(args ...interface{}) (int64, error) {
	var database Database
	return database.Execute(args...)
}
// get last query sql
func (this *Connection) LastSql() string {
	if len(Connect.SqlLog) > 0 {
		return Connect.SqlLog[len(Connect.SqlLog)-1:][0]
	}
	return ""
}
// all sql query logs in this request
func (this *Connection) SqlLogs() []string {
	return Connect.SqlLog
}
// get origin *sql.DB
func GetDB() *sql.DB {
	return DB
}
