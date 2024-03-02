package gorose

import (
	"database/sql"
	"time"
)

type Config struct {
	Driver          string
	DSN             string
	Prefix          string
	Weight          int8 // 权重越高,选中的概率越大
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type ConfigCluster struct {
	WriteConf []Config
	ReadConf  []Config
}

//type ConfigMulti map[string]ConfigCluster

func (c ConfigCluster) init() (master []*sql.DB, slave []*sql.DB) {
	if len(c.WriteConf) > 0 {
		for _, v := range c.WriteConf {
			master = append(master, c.initDB(&v))
		}
	}
	if len(c.ReadConf) > 0 {
		for _, v := range c.ReadConf {
			slave = append(master, c.initDB(&v))
		}
	}
	return
}
func (c ConfigCluster) initDB(v *Config) *sql.DB {
	db, err := sql.Open(v.Driver, v.DSN)
	if err != nil {
		panic(err.Error())
	}

	if v.MaxIdleConns > 0 {
		db.SetMaxIdleConns(v.MaxIdleConns)
	}
	if v.MaxOpenConns > 0 {
		db.SetMaxOpenConns(v.MaxOpenConns)
	}
	if v.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(v.ConnMaxLifetime)
	}
	if v.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(v.ConnMaxIdleTime)
	}

	return db
}
