package gorose

import (
	"database/sql"
	"errors"
	"github.com/gohouse/gorose/across"
	"github.com/gohouse/gorose/helper"
	"math/rand"
)

type sqlDb struct {
	SlaveDbs []*sql.DB
	MasterDb *sql.DB
}
type Connection struct {
	DbConfig *across.DbConfigCluster
	Db       sqlDb
}

func (c *Connection) Query(arg string, params ...interface{}) (result []map[string]interface{}, errs error) {
	return c.NewSession().Query(arg, params...)
}

func (c *Connection) Execute(arg string, params ...interface{}) int64 {
	return 0
}

func (c *Connection) NewSession() *Database {
	dba := NewDatabase()
	dba.Connection = c
	return dba
}

func (c *Connection) Table(arg interface{}) *Database {
	return c.NewSession().Table(arg)
}

// parseOpenArgs 解析入口参数
func (c *Connection) parseOpenArgs(args ...interface{}) (dbConf *across.DbConfigCluster, err error) {
	var fileOrDriverType, dsnOrFile string
	switch len(args) {
	case 1: // 传入的配置struct ( across.DbConfigCluster )
		switch args[0].(type) {
		case across.DbConfigCluster, *across.DbConfigCluster:
			dbConf = args[0].(*across.DbConfigCluster)
			if dbConf.Master == nil {
				err = errors.New("master配置参数缺失")
				return
			}
			return
		case across.DbConfigSingle, *across.DbConfigSingle:
			dbConf.Master = args[0].(*across.DbConfigSingle)
			return
		default:
			return dbConf, errors.New("参数格式错误: 一个参数时,需要传入across.DbConfigCluster或across.DbConfigSingle类型数据")
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
			if !helper.FileExists(dsnOrFile) {
				return dbConf, errors.New("配置文件不存在")
			}
			// 调用`parser`目录解析器 解析文件
			dbConf, err = NewFileParser(fileOrDriverType, dsnOrFile)
		}
	}
	return
}

// boot sql driver
func (c *Connection) bootDbs(dbConf *across.DbConfigCluster) (err error) {
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
func (c *Connection) bootReal(dbConf *across.DbConfigSingle) (db *sql.DB, err error) {
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
	} else {
		db = c.Db.SlaveDbs[rand.Intn(lenSlave-1)]
	}
	return
}

func (c *Connection) GetExecuteDb() (*sql.DB) {
	return c.Db.MasterDb
}