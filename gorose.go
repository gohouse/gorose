package gorose

import (
	"database/sql"
	"errors"
	"github.com/gohouse/gorose/drivers"
	"kuaixinwen/utils"
)

// Connection is the database pre handle
type Connection struct {
	// DB is origin DB
	DB *sql.DB
	// all config sets
	DbConfig map[string]interface{}
	// default database
	Default string
	// current config on use
	CurrentConfig map[string]string
	//// all sql logs
	//SqlLog []string
	//// if in transaction, the code auto change
	//Trans bool
	// max open connections
	SetMaxOpenConns int
	// max freedom connections leave
	SetMaxIdleConns int
}

// Open instance of sql.DB.Oper
// if args has 1 param , it will be derect connection or with default config set
// if args has 2 params , the second param will be the default dirver key
func Open(args ...interface{}) (Connection, error) {
	var conn = Connection{}
	//fmt.Println(args)
	//return conn, errors.New("dsf")
	if len(args) == 1 {
		// continue
	} else if len(args) == 2 {
		if confReal, ok := args[1].(string); ok {
			conn.Default = confReal
		} else {
			// 指定默认数据库只能为字符串!
			return conn, errors.New("only str allowed of default database name")
		}
	} else {
		// Open方法只接收1个或2个参数!
		return conn, errors.New("1 or 2 params need in Open() method")
	}
	// 解析config
	err := conn.parseConfig(args[0])
	if err != nil {
		return conn, err
	}

	// 驱动数据库
	errs := conn.boot()

	return conn, errs
}

// Parse input config
func (conn *Connection) parseConfig(args interface{}) error {
	if confReal, ok := args.(map[string]string); ok { // direct connection
		conn.CurrentConfig = confReal
	} else if confReal, ok := args.(map[string]interface{}); ok {
		// store the full connection
		conn.DbConfig = confReal
		// if set the Default conf, store it
		if defaultDb, ok := confReal["Default"]; ok {
			// judge if seted
			if conn.Default == "" {
				conn.Default = defaultDb.(string)
			}
		}
		if conn.Default == "" {
			// 配置文件默认数据库链接未设置
			return errors.New("the default database is missing in config!")
		}
		// 获取指定的默认数据库链接信息
		var connections map[string]map[string]string
		if connectionsInterface, ok := confReal["Connections"]; ok {
			switch connectionsInterface.(type) {
			case map[string]map[string]string:
				connections = connectionsInterface.(map[string]map[string]string)
			case map[string]interface{}:
				connectionsTmp := connectionsInterface.(map[string]interface{})
				if connectionsTmpReal, ok := connectionsTmp[conn.Default]; ok {
					switch connectionsTmpReal.(type) {
					case map[string]string:
						connections = map[string]map[string]string{conn.Default: connectionsTmpReal.(map[string]string)}
					default:
						return errors.New("the database connections format error !")
					}
				}
			default:
				return errors.New("the database connections format error !")
			}
		} else {
			return errors.New("the database connections missing !")
		}
		if defaultDbConnection, ok := connections[conn.Default]; ok {
			conn.CurrentConfig = defaultDbConnection
		} else {
			// 指定的数据库链接不存在!
			return errors.New("the database for using is missing!")
		}
		// 设置连接池信息
		if mo, ok := confReal["SetMaxOpenConns"]; ok {
			moInt := utils.ParseInt(mo)
			if moInt>0 {
				conn.SetMaxOpenConns = moInt
			}
		}
		if mi, ok := confReal["SetMaxIdleConns"]; ok {
			miInt := utils.ParseInt(mi)
			if miInt>0 {
				conn.SetMaxIdleConns = miInt
			}
		}
	} else {
		return errors.New("format error in database config!")
	}
	return nil
}

// Boot sql driver
func (conn *Connection) boot() error {
	//dbObj := conn.CurrentConfig
	var driver, dsn string
	var err error

	//DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	driver,dsn = drivers.GetDsnByDriverName(conn.CurrentConfig)

	// 开始驱动
	conn.DB, err = sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	if conn.SetMaxOpenConns>0 {
		conn.DB.SetMaxOpenConns(conn.SetMaxOpenConns)
	}
	if conn.SetMaxIdleConns>0 {
		conn.DB.SetMaxIdleConns(conn.SetMaxIdleConns)
	}

	// 检查是否可以ping通
	err2 := conn.DB.Ping()

	return err2
}

// Close database
func (conn *Connection) Close() error {
	//conn.SqlLog = []string{}
	return conn.DB.Close()
}

// Ping db
func (conn *Connection) Ping() error {
	return conn.DB.Ping()
}

// Table is set table from database
func (conn *Connection) Table(table string) *Database {
	return conn.GetInstance().Table(table)
}

//// Begin transaction begin
//func (conn *Connection) Begin() {
//	Tx, _ = DB.Begin()
//	conn.Trans = true
//}
//
//// Commit is transaction commit
//func (conn *Connection) Commit() {
//	Tx.Commit()
//	conn.Trans = false
//}
//
//// Rollback is transaction rollback
//func (conn *Connection) Rollback() {
//	Tx.Rollback()
//	conn.Trans = false
//}
//
//// Transaction is simple transaction
//func (conn *Connection) Transaction(closure func() error) bool {
//	//defer func() {
//	//	if err := recover(); err != nil {
//	//		conn.Rollback()
//	//		panic(err)
//	//	}
//	//}()
//
//	conn.Begin()
//	err := closure()
//	if err != nil {
//		conn.Rollback()
//		return false
//	}
//	conn.Commit()
//
//	return true
//}

// Query str
func (conn *Connection) Query(args ...interface{}) ([]map[string]interface{}, error) {
	return conn.GetInstance().Query(args...)
}

// Execute str
func (conn *Connection) Execute(args ...interface{}) (int64, error) {
	return conn.GetInstance().Execute(args...)
}

// GetInstance , get the database object
func (conn *Connection) GetInstance() *Database {
	//var database *Database
	//return database
	return &Database{connection:conn}
}

// JsonEncode : parse json
func (conn *Connection) JsonEncode(arg interface{}) string {
	return conn.GetInstance().JsonEncode(arg)
}

//// LastSql is get last query sql
//func (conn *Connection) LastSql() string {
//	if len(conn.SqlLog) > 0 {
//		return conn.SqlLog[len(conn.SqlLog)-1:][0]
//	}
//	return ""
//}
//
//// SqlLogs is all sql query logs in this request
//func (conn *Connection) SqlLogs() []string {
//	return conn.SqlLog
//}

// GetDB is get origin *sql.DB
func (conn *Connection) GetDB() *sql.DB {
	return conn.DB
}
