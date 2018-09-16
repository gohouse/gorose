package gorose

import (
	"database/sql"
	"errors"
	"github.com/gohouse/gorose/config"
	"github.com/gohouse/gorose/helper"
	"github.com/gohouse/gorose/parser"
)

type Connection struct {
	DbConfig *config.DbConfig
	Db       *sql.DB
}

// Open 链接数据库入口, 传入配置
// args 接收一个或2个参数, 一个参数时:struct配置文件(config.DbConfig{})
//		两个参数时: 第一个是驱动或文件类型, 第二个是dsn或文件路径
func Open(args ...interface{}) (*Connection, error) {
	var c = &Connection{}
	var err error

	// 解析配置获取参数并保存
	c.DbConfig, err = c.parseOpenArgs(args...)

	if err != nil {
		return c, err
	}

	// 驱动数据库获取链接并保存
	c.Db, err = c.bootDb(c.DbConfig)
	return c, err
}

func (dba *Connection) Query(arg string, params ...interface{}) (result []map[string]interface{}, errs error) {
	return dba.NewDB().Query(arg, params...)
}

func (c *Connection) Execute(arg string, params ...interface{}) int64 {
	return 0
}

func (c *Connection) NewDB() *Database {
	return &Database{Connection: c}
}

func (c *Connection) Table(arg interface{}) *Database {
	return c.NewDB().Table(arg)
}

// parseOpenArgs 解析入口参数
func (c *Connection) parseOpenArgs(args ...interface{}) (dbConf *config.DbConfig, err error) {
	var fileOrDriverType, dsnOrFile string
	switch len(args) {
	case 1: // 传入的配置struct ( config.DbConfig )
		switch args[0].(type) {
		case config.DbConfig, *config.DbConfig:
			dbConf = args[0].(*config.DbConfig)
			return
		default:
			return dbConf, errors.New("参数格式错误: 一个参数时,需要传入config.DbConfig类型数据")
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
	if typeName, err = config.Getter(fileOrDriverType); err == nil {
		switch typeName {
		case "driver":
			dbConf.Driver = fileOrDriverType
			dbConf.Dsn = dsnOrFile

		case "file":
			// 配置文件, 读取配置文件
			if !helper.FileExists(dsnOrFile) {
				return dbConf, errors.New("配置文件不存在")
			}
			// 调用`parser`目录解析器 解析文件
			dbConf, err = parser.NewFileParser(fileOrDriverType, dsnOrFile)
		}
	}
	return
}

// boot sql driver
func (c *Connection) bootDb(dbConf *config.DbConfig) (db *sql.DB, err error) {
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
