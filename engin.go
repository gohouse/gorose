package gorose

import (
	"database/sql"
)

type dbObject struct {
	driver string
	db     *sql.DB
	tx     *sql.Tx
}
type cluster struct {
	master     []dbObject
	masterSize int
	slave      []dbObject
	slaveSize  int
}
type Engin struct {
	config       *ConfigCluster
	enableSqlLog bool
	prefix       string
	dbs          *cluster
}

var _ IEngin = &Engin{}

// NewEngin : init Engin struct pointer
// NewEngin : 初始化 Engin 结构体对象指针
func NewEngin() *Engin {
	return new(Engin)
}

// EnableSqlLog : wither record sql logs, if no args input, the arg value default true
// EnableSqlLog : 是否启用sql日志记录, 如果不传递参数, 则参数值默认为true
func (c *Engin) EnableSqlLog(arg ...bool) {
	if len(arg) == 0 {
		c.enableSqlLog = true
	} else {
		c.enableSqlLog = arg[0]
	}
}

// IfEnableSqlLog 是否启用sql日志
func (c *Engin) IfEnableSqlLog() bool {
	return c.enableSqlLog
}

// Prefix 设置前缀
func (c *Engin) Prefix(pre string) {
	c.prefix = pre
}

// GetPrefix 获取前缀
func (c *Engin) GetPrefix() string {
	return c.prefix
}

// GetQueryDB : get a slave db for using query operation
// GetQueryDB : 获取一个从库用来做查询操作
func (c *Engin) GetQueryDB() dbObject {
	if c.dbs.slaveSize == 0 {
		return c.GetExecuteDB()
	}
	var randint = getRandomInt(c.dbs.slaveSize)
	return c.dbs.slave[randint]
}

// GetExecuteDB : get a master db for using execute operation
// GetExecuteDB : 获取一个主库用来做查询之外的操作
func (c *Engin) GetExecuteDB() dbObject {
	if c.dbs.masterSize == 0 {
		return dbObject{}
	}
	var randint = getRandomInt(c.dbs.masterSize)
	return c.dbs.master[randint]
}

func (c *Engin) bootSingle(conf *Config) error {
	var cc = new(ConfigCluster)
	cc.Master = append(cc.Master, *conf)
	c.config = cc
	return c.bootCluster()
}

func (c *Engin) bootCluster() error {
	if len(c.config.Slave) > 0 {
		for _, item := range c.config.Slave {
			db, err := c.bootReal(item)
			if err != nil {
				return err
			}
			if c.dbs == nil {
				c.dbs = new(cluster)
			}
			c.dbs.slave = append(c.dbs.slave, dbObject{driver: item.Driver, db: db})
			c.dbs.slaveSize++
		}
	}
	if len(c.config.Master) > 0 {
		for _, item := range c.config.Master {
			db, err := c.bootReal(item)
			if err != nil {
				return err
			}
			if c.dbs == nil {
				c.dbs = new(cluster)
			}
			c.dbs.master = append(c.dbs.master, dbObject{driver: item.Driver, db: db})
			c.dbs.masterSize++
		}
	}

	return nil
}

// boot sql driver
func (c *Engin) bootReal(dbConf Config) (db *sql.DB, err error) {
	//db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8")
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
	db.SetMaxOpenConns(dbConf.SetMaxOpenConns)
	db.SetMaxIdleConns(dbConf.SetMaxIdleConns)

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
