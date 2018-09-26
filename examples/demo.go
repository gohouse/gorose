package main

import (
	"github.com/gohouse/gorose"
	_ "github.com/gohouse/gorose/driver/mysql"
	"fmt"
)

// DB Config.(Recommend to use configuration file to import)
var dbConfig1 = &gorose.DbConfigSingle{
	Driver:          "mysql",                                                   // 驱动: mysql/sqlite/oracle/mssql/postgres
	EnableQueryLog:  false,                                                     // 是否开启sql日志
	SetMaxOpenConns: 0,                                                         // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns: 0,                                                         // (连接池)闲置的连接数
	Prefix:          "",                                                        // 表前缀
	Dsn:             "gcore:gcore@tcp(192.168.200.248:3306)/test?charset=utf8", // 数据库链接
}

func main() {
	connection, err := gorose.Open(dbConfig1)
	if err != nil {
		fmt.Println(err)
		return
	}

	db := connection.NewSession()

	res, err := db.Table("users").First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(res)
}
