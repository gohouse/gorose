package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
)

var conn = make(map[string]gorose.Connection)

func init()  {
	// 链接第一个
	connection, err := gorose.Open(config.DbConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn["mysql_dev"] = connection

	// 链接第二个
	connection2, err2 := gorose.Open(config.DbConfig, "mysql_dev2")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	conn["mysql_dev2"] = connection2
}

func main() {
	// close DB
	conn1 := conn["mysql_dev"]
	conn2 := conn["mysql_dev2"]
	// 第一个
	defer conn1.Close()
	db := conn1.NewDB()
	res, err := db.Table("users").Fields("id,name").First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(res)

	// 第二个
	defer conn2.Close()
	r,_ := conn2.Table("users").Fields("id,name").First()
	fmt.Println(r)
}
