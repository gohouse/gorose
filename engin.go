package gorose

import (
	"database/sql"
)

type cluster struct {
	master     []map[string]*sql.DB
	masterSize int
	slave      []map[string]*sql.DB
	slaveSize  int
}
type Engin struct {
	config         *ConfigCluster
	enableQueryLog bool
	prefix         string
	dbs            *cluster
}

var _ IEngin = &Engin{}

// NewEngin : init Engin struct pointer
// NewEngin : 初始化 Engin 结构体对象指针
func NewEngin() *Engin {
	return new(Engin)
}

// EnableQueryLog : wither record sql logs, if no args input, the arg value default true
// EnableQueryLog : 是否启用sql日志记录, 如果不传递参数, 则参数值默认为true
func (c *Engin) EnableQueryLog(arg ...bool) {
	if len(arg) == 0 {
		c.enableQueryLog = true
	} else {
		c.enableQueryLog = arg[0]
	}
}

func (c *Engin) IfEnableQueryLog() bool {
	return c.enableQueryLog
}

func (c *Engin) Prefix(pre string) {
	c.prefix = pre
}

func (c *Engin) GetPrefix() string {
	return c.prefix
}

// GetQueryDB : get a slave db for using query operation
// GetQueryDB : 获取一个从库用来做查询操作
func (c *Engin) GetQueryDB() (db *sql.DB, driver string) {
	if c.dbs.slaveSize == 0 {
		return c.GetExecuteDB()
	}
	var randint = getRandomInt(c.dbs.slaveSize)
	for driver, db = range c.dbs.slave[randint] {
		return
	}
	return
}

// GetExecuteDB : get a master db for using execute operation
// GetExecuteDB : 获取一个主库用来做查询之外的操作
func (c *Engin) GetExecuteDB() (db *sql.DB, driver string) {
	if c.dbs.masterSize == 0 {
		return
	}
	var randint = getRandomInt(c.dbs.masterSize)
	for driver, db = range c.dbs.master[randint] {
		return
	}
	return
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
			c.dbs.slave = append(c.dbs.slave, map[string]*sql.DB{item.Driver: db})
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
			c.dbs.master = append(c.dbs.master, map[string]*sql.DB{item.Driver: db})
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

func (c *Engin) NewSession() ISession {
	return NewSession(c)
}
