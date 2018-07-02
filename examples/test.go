package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
)

func main() {
	connection, err := gorose.Open(config.DbConfig2, "mysql_dev2")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer connection.Close()

	db := connection.GetInstance()
	fmt.Println(db)
	res, err := db.Table("users").Where("id", ">", 2).First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(res)

	connection2, err2 := gorose.Open(config.DbConfig2, "mysql_dev3")
	if err2 != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer connection2.Close()

	db2 := connection.GetInstance()
	fmt.Println(db2)
	res2, err2 := db2.Table("fd_logs").Where("id", ">", 2).First()
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println(db2.LastSql)
	fmt.Println(res2)
}
