package main

import (
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/gohouse/gorose/driver/mysql"
)

type users3 struct {
	Name string `orm:"name"`
	Age int `orm:"age"`
	Job string `orm:"job"`
}

func (u *users3) TableName() string {
	return "users"
}

// DB Config.(Recommend to use configuration file to import)
var dbConfig3 = &gorose.DbConfigSingle {
	Driver:          "mysql", // 驱动: mysql/sqlite/oracle/mssql/postgres
	EnableQueryLog:  false,   // 是否开启sql日志
	SetMaxOpenConns: 0,    // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns: 0,    // (连接池)闲置的连接数
	Prefix:          "", // 表前缀
	Dsn:             "gcore:gcore@tcp(192.168.200.248:3306)/test?charset=utf8", // 数据库链接
}

func main() {
	connection, err := gorose.Open(dbConfig3)
	if err != nil {
		fmt.Println(err)
		return
	}

	db := connection.NewSession()

	var user []users3

	err2 := db.Table(&user).Limit(3).Select()
	if err2 != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(user)
}
