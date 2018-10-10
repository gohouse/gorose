package gorose

import (
	"database/sql"
	"errors"
	"github.com/gohouse/gorose/across"
	"github.com/gohouse/gorose/cors"
	"github.com/gohouse/gorose/utils"
	"math/rand"
)

type sqlDb struct {
	SlaveDbs []*sql.DB
	MasterDb *sql.DB
}
type Connection struct {
	DbConfig *DbConfigCluster
	Db       sqlDb
	Logger   cors.LoggerHandler // 持久化sql日志,如果为空则不持久化
}

// Use : cors
func (conn *Connection) Use(options ...func(*Connection)) *Connection {
	for _, option := range options {
		option(conn)
	}

	return conn
}

// Close database
func (conn *Connection) Close() error {
	//conn.SqlLog = []string{}
	return conn.GetExecuteDb().Close()
}

func (c *Connection) Query(arg string, params ...interface{}) ([]map[string]interface{}, error) {
	return c.NewSession().Query(arg, params...)
}

func (c *Connection) Execute(arg string, params ...interface{}) (int64, error) {
	return c.NewSession().Execute(arg, params...)
}

func (c *Connection) NewSession() *Session {
	dba := NewOrm()
	dba.Connection = c
	return dba
}

func (c *Connection) Table(arg interface{}) *Session {
	return c.NewSession().Table(arg)
}

// parseOpenArgs 解析入口参数
func (c *Connection) parseOpenArgs(args ...interface{}) (*DbConfigCluster, error) {
	var dbConf = NewDbConfigCluster()
	var err error
	var fileOrDriverType, dsnOrFile string
	switch len(args) {
	case 1: // 传入的配置struct ( DbConfigCluster )
		switch args[0].(type) {
		case DbConfigCluster, *DbConfigCluster:
			dbConf = args[0].(*DbConfigCluster)
			if dbConf.Master == nil {
				err = errors.New("master配置参数缺失")
				return dbConf,err
			}
			return dbConf,err
		case DbConfigSingle, *DbConfigSingle:
			dbConf.Master = args[0].(*DbConfigSingle)
			return dbConf,err
		default:
			return dbConf, errors.New("参数格式错误: 一个参数时,需要传入DbConfigCluster或DbConfigSingle类型数据")
		}
	case 2: // 第一个是驱动或文件类型, 第二个是dsn或文件路径
		var ok bool
		var ok2 bool
		fileOrDriverType, ok = args[0].(string)
		dsnOrFile, ok2 = args[1].(string)
		if !ok || !ok2 {
			return dbConf, errors.New("参数格式错误: 两个参数时,需要传入string类型数据")
		}
	default:
		return dbConf, errors.New("参数数量有误: 只接收一个或两个参数")
	}

	// 解析配置
	var typeName string
	// 如果是在配置内
	if typeName, err = across.Getter(fileOrDriverType); err == nil {
		switch typeName {
		case "driver":
			dbConf.Master.Driver = fileOrDriverType
			dbConf.Master.Dsn = dsnOrFile

		case "file":
			// 配置文件, 读取配置文件
			if !utils.FileExists(dsnOrFile) {
				return dbConf, errors.New("配置文件不存在")
			}
			// 调用`parser`目录解析器 解析文件
			dbConf, err = NewFileParser(fileOrDriverType, dsnOrFile)
		}
	}
	return dbConf,err
}

// boot sql driver
func (c *Connection) bootDbs(dbConf *DbConfigCluster) (err error) {
	//db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8")
	// 开始驱动, 分别保存不同配置的链接
	c.Db.MasterDb, err = c.bootReal(dbConf.Master)
	//fmt.Println(dbConf.Master)
	//tmp,err33 := json.Marshal(c.Db.MasterDb)
	//fmt.Println(string(tmp), err33)
	// 驱动从数据库
	if len(dbConf.Slave) > 0 {
		var db *sql.DB
		for _, item := range dbConf.Slave {
			db, err = c.bootReal(item)
			if err != nil {
				return
			}
			c.Db.SlaveDbs = append(c.Db.SlaveDbs, db)
		}
	}
	return
}

// boot sql driver
func (c *Connection) bootReal(dbConf *DbConfigSingle) (db *sql.DB, err error) {
	//db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8")
	// 开始驱动
	db, err = sql.Open(dbConf.Driver, dbConf.Dsn)
	if err != nil {
		return
	}
	if dbConf.SetMaxOpenConns > 0 {
		db.SetMaxOpenConns(dbConf.SetMaxOpenConns)
	}
	if dbConf.SetMaxIdleConns > 0 {
		db.SetMaxIdleConns(dbConf.SetMaxIdleConns)
	}

	// 检查是否可以ping通
	err = db.Ping()
	return
}

func (c *Connection) GetQueryDb() (db *sql.DB) {
	lenSlave := len(c.Db.SlaveDbs)

	if lenSlave == 0 {
		db = c.GetExecuteDb()
	} else if lenSlave == 1 {
		db = c.Db.SlaveDbs[0]
	} else {
		db = c.Db.SlaveDbs[rand.Intn(lenSlave-1)]
	}
	return
}

func (c *Connection) GetExecuteDb() (*sql.DB) {
	return c.Db.MasterDb
}
