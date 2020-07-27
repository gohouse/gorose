package gorose

import (
	"database/sql"
	"fmt"
)

// TAGNAME ...
var TAGNAME = "gorose"

// IGNORE ...
var IGNORE = "-"

type cluster struct {
	master     []*sql.DB
	masterSize int
	slave      []*sql.DB
	slaveSize  int
}

// Engin ...
type Engin struct {
	config *ConfigCluster
	driver string
	prefix string
	dbs    *cluster
	logger ILogger
}

var _ IEngin = (*Engin)(nil)

// NewEngin : init Engin struct pointer
// NewEngin : 初始化 Engin 结构体对象指针
func NewEngin(conf ...interface{}) (e *Engin, err error) {
	engin := new(Engin)
	if len(conf) == 0 {
		return
	}

	// 使用默认的log, 如果自定义了logger, 则只需要调用 Use() 方法即可覆盖
	engin.Use(DefaultLogger())

	switch conf[0].(type) {
	// 传入的是单个配置
	case *Config:
		err = engin.bootSingle(conf[0].(*Config))
	// 传入的是集群配置
	case *ConfigCluster:
		engin.config = conf[0].(*ConfigCluster)
		err = engin.bootCluster()
	default:
		panic(fmt.Sprint("Open() need *gorose.Config or *gorose.ConfigCluster param, also can empty for build sql string only, but ",
			conf, " given"))
	}

	return engin, err
}

// Use ...
func (c *Engin) Use(closers ...func(e *Engin)) {
	for _, closer := range closers {
		closer(c)
	}
}

// Ping ...
func (c *Engin) Ping() error {
	//for _,item := range c.dbs.master {
	//
	//}
	return c.GetQueryDB().Ping()
}

// TagName 自定义结构体对应的orm字段,默认gorose
func (c *Engin) TagName(arg string) {
	//c.tagName = arg
	TAGNAME = arg
}

// IgnoreName 自定义结构体对应的orm忽略字段名字,默认-
func (c *Engin) IgnoreName(arg string) {
	//c.ignoreName = arg
	IGNORE = arg
}

// SetPrefix 设置表前缀
func (c *Engin) SetPrefix(pre string) {
	c.prefix = pre
}

// GetPrefix 获取前缀
func (c *Engin) GetPrefix() string {
	return c.prefix
}

// GetDriver ...
func (c *Engin) GetDriver() string {
	return c.driver
}

// GetQueryDB : get a slave db for using query operation
// GetQueryDB : 获取一个从库用来做查询操作
func (c *Engin) GetQueryDB() *sql.DB {
	if c.dbs.slaveSize == 0 {
		return c.GetExecuteDB()
	}
	var randint = getRandomInt(c.dbs.slaveSize)
	return c.dbs.slave[randint]
}

// GetExecuteDB : get a master db for using execute operation
// GetExecuteDB : 获取一个主库用来做查询之外的操作
func (c *Engin) GetExecuteDB() *sql.DB {
	if c.dbs.masterSize == 0 {
		return nil
	}
	var randint = getRandomInt(c.dbs.masterSize)
	return c.dbs.master[randint]
}

// GetLogger ...
func (c *Engin) GetLogger() ILogger {
	return c.logger
}

// SetLogger ...
func (c *Engin) SetLogger(lg ILogger) {
	c.logger = lg
}

func (c *Engin) bootSingle(conf *Config) error {
	// 如果传入的是单一配置, 则转换成集群配置, 方便统一管理
	var cc = new(ConfigCluster)
	cc.Master = append(cc.Master, *conf)
	c.config = cc
	return c.bootCluster()
}

func (c *Engin) bootCluster() error {
	//fmt.Println(len(c.config.Slave))
	if len(c.config.Slave) > 0 {
		for _, item := range c.config.Slave {
			if c.config.Driver != "" {
				item.Driver = c.config.Driver
			}
			if c.config.Prefix != "" {
				item.Prefix = c.config.Prefix
			}
			db, err := c.bootReal(item)
			if err != nil {
				return err
			}
			if c.dbs == nil {
				c.dbs = new(cluster)
			}
			c.dbs.slave = append(c.dbs.slave, db)
			c.dbs.slaveSize++
			c.driver = item.Driver
		}
	}
	var pre, dr string
	if len(c.config.Master) > 0 {
		for _, item := range c.config.Master {
			if c.config.Driver != "" {
				item.Driver = c.config.Driver
			}
			if c.config.Prefix != "" {
				item.Prefix = c.config.Prefix
			}
			db, err := c.bootReal(item)

			if err != nil {
				return err
			}
			if c.dbs == nil {
				c.dbs = new(cluster)
			}
			c.dbs.master = append(c.dbs.master, db)
			c.dbs.masterSize = c.dbs.masterSize + 1
			c.driver = item.Driver
			//fmt.Println(c.dbs.masterSize)
			if item.Prefix != "" {
				pre = item.Prefix
			}
			if item.Driver != "" {
				dr = item.Driver
			}
		}
	}
	// 如果config没有设置prefix,且configcluster设置了prefix,则使用cluster的prefix
	if pre != "" && c.prefix == "" {
		c.prefix = pre
	}
	// 如果config没有设置driver,且configcluster设置了driver,则使用cluster的driver
	if dr != "" && c.driver == "" {
		c.driver = dr
	}

	return nil
}

// boot sql driver
func (c *Engin) bootReal(dbConf Config) (db *sql.DB, err error) {
	//db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8mb4")
	// 开始驱动
	db, err = sql.Open(dbConf.Driver, dbConf.Dsn)
	if err != nil {
		return
	}

	// 检查是否可以ping通
	err = db.Ping()
	if err != nil {
		return
	}

	// 连接池设置
	if dbConf.SetMaxOpenConns > 0 {
		db.SetMaxOpenConns(dbConf.SetMaxOpenConns)
	}
	if dbConf.SetMaxIdleConns > 0 {
		db.SetMaxIdleConns(dbConf.SetMaxIdleConns)
	}

	return
}

// NewSession 获取session实例
// 这是一个语法糖, 为了方便使用(engin.NewSession())添加的
// 添加后会让engin和session耦合, 如果不想耦合, 就删掉此方法
// 删掉这个方法后,可以使用 gorose.NewSession(gorose.IEngin)
// 通过 gorose.IEngin 依赖注入的方式, 达到解耦的目的
func (c *Engin) NewSession() ISession {
	return NewSession(c)
}

// NewOrm 获取orm实例
// 这是一个语法糖, 为了方便使用(engin.NewOrm())添加的
// 添加后会让engin和 orm 耦合, 如果不想耦合, 就删掉此方法
// 删掉这个方法后,可以使用 gorose.NewOrm(gorose.NewSession(gorose.IEngin))
// 通过 gorose.ISession 依赖注入的方式, 达到解耦的目的
func (c *Engin) NewOrm() IOrm {
	return NewOrm(c)
}
