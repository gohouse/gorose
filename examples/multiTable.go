package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
)

func main() {
	connection, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer connection.Close()

	res, err := connection.Table("users").Where("id", "<", 1).First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

	res2, err := connection.Table("users").Limit(2).Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res2)
}
