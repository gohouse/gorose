package config

import (
	"github.com/gohouse/gorose"
)

// DbConf
var DbConfigSalve = &gorose.DbConfigSingle{
	Driver:          "mysql",                                                   // 驱动: mysql/sqlite/oracle/mssql/postgres
	EnableQueryLog:  false,                                                     // 是否开启sql日志
	SetMaxOpenConns: 0,                                                         // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns: 0,                                                         // (连接池)闲置的连接数, 默认-1
	Prefix:          "",                                                        // 表前缀
	Dsn:             "gcore:gcore@tcp(192.168.200.248:3306)/test?charset=utf8", // 数据库链接
}
var DbConfigMaster = &gorose.DbConfigSingle{
	Driver:          "mysql",                                                   // 驱动: mysql/sqlite/oracle/mssql/postgres
	EnableQueryLog:  false,                                                     // 是否开启sql日志
	SetMaxOpenConns: 0,                                                         // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns: 0,                                                         // (连接池)闲置的连接数, 默认-1
	Prefix:          "",                                                        // 表前缀
	Dsn:             "gcore:gcore@tcp(192.168.200.248:3306)/wcc_service_fooddrug?charset=utf8", // 数据库链接
}

var DbConfCluster = &gorose.DbConfigCluster{
	Master:DbConfigMaster,
	Slave: []*gorose.DbConfigSingle{DbConfigSalve},
}

func GetConnection() *gorose.Connection {
	connection, err := gorose.Open(DbConfigMaster)
	if err != nil {
		panic(err)
	}
	return connection
}

func GetConnectionCluster() *gorose.Connection {
	connection, err := gorose.Open(DbConfCluster)
	if err != nil {
		panic(err)
	}
	return connection
}
