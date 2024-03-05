package gorose

import (
	"database/sql"
	"log/slog"
	"math"
	"os"
)

type GoRose struct {
	ILogger

	Cluster  *ConfigCluster
	master   []*sql.DB
	slave    []*sql.DB
	driver   string
	prefix   string
	handlers HandlersChain
}

type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc

const abortIndex int8 = math.MaxInt8 >> 1

//type handlers struct {
//	handlers HandlersChain
//	index    int8
//}
//
//func (c *handlers) Next() {
//	c.index++
//	for c.index < int8(len(c.handlers)) {
//		c.handlers[c.index](c)
//		c.index++
//	}
//}

func (g *GoRose) Use(h ...HandlerFunc) *GoRose {
	g.handlers = append(g.handlers, h...)
	return g
}

//func (dg *GoRose) Use(h ...HandlerFunc) *GoRose {
//}

// Open db
// examples
//
//	Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true")
//	Open(&ConfigCluster{...})
func Open(conf ...any) *GoRose {
	var g = GoRose{ILogger: DefaultLogger(slog.LevelDebug, os.Stdout)}
	switch len(conf) {
	case 1:
		if single, ok := conf[0].(*Config); ok {
			g.driver = single.Driver
			g.prefix = single.Prefix
			g.Cluster = &ConfigCluster{WriteConf: []Config{*single}}
			if single == nil { // build sql only
				//g.ConfigCluster = &ConfigCluster{Prefix: "test_"}
				return &g
			}
			g.master, g.slave = g.Cluster.init()
		} else if cluster, ok := conf[0].(*ConfigCluster); ok {
			g.driver = cluster.WriteConf[0].Driver
			g.prefix = cluster.WriteConf[0].Prefix
			g.Cluster = cluster
			if cluster == nil { // build sql only
				//g.ConfigCluster = &ConfigCluster{Prefix: "test_"}
				return &g
			}
			g.master, g.slave = g.Cluster.init()
		} else {
			g.driver = "mysql" // for toSql test
		}
	case 2:
		g.driver = conf[0].(string)
		db, err := sql.Open(g.driver, conf[1].(string))
		if err != nil {
			panic(err.Error())
		}
		g.master = append(g.master, db)
	default:
		panic("config must be *gorose.ConfigCluster or sql.Open() origin params")
	}
	return &g
}

func (g *GoRose) Close() (err error) {
	if len(g.master) > 0 {
		for _, db := range g.master {
			err = db.Close()
		}
	}
	if len(g.slave) > 0 {
		for _, db := range g.slave {
			err = db.Close()
		}
	}
	return
}

func (g *GoRose) MasterDB() *sql.DB {
	if len(g.master) == 0 {
		return nil
	}
	return g.master[GetRandomInt(len(g.master))]
}
func (g *GoRose) SlaveDB() *sql.DB {
	if len(g.slave) == 0 {
		return g.MasterDB()
	}
	return g.slave[GetRandomInt(len(g.slave))]
}

func (g *GoRose) NewDatabase() *Database {
	return NewDatabase(g)
}

func (g *GoRose) NewEngin() *Engin {
	return NewEngin(g)
}
