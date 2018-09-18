package main

import (
	"github.com/gohouse/gorose"
	_ "github.com/gohouse/gorose/driver/mysql"
	"fmt"
	"reflect"
)

type users2 struct {
	Name string `orm:"name"`
	Age int `orm:"age"`
	Job string `orm:"job"`
}

func (u *users2) TableName() string {
	return "users"
}

// DB Config.(Recommend to use configuration file to import)
var dbConfig2 = &gorose.DbConfigSingle {
	Driver:          "mysql", // 驱动: mysql/sqlite/oracle/mssql/postgres
	EnableQueryLog:  false,   // 是否开启sql日志
	SetMaxOpenConns: 0,    // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns: 0,    // (连接池)闲置的连接数, 默认-1
	Prefix:          "", // 表前缀
	Dsn:             "gcore:gcore@tcp(192.168.200.248:3306)/test?charset=utf8", // 数据库链接
}

func main() {
	connection, err := gorose.Open(dbConfig2)
	if err != nil {
		fmt.Println(err)
		return
	}

	db := connection.NewSession()

	var user users2
	dstVal := reflect.ValueOf(user)
	if tn := dstVal.MethodByName("TableName"); tn.IsValid() {
		tableName := tn.Call(nil)[0].String()
		fmt.Println(tableName)
	}
	_,err2 := db.Table(&user).First()
	if err2 != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(user)
}
