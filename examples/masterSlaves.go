package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/examples/config"
)

func main() {
	connection := config.GetConnectionCluster()

	defer connection.Close()

	res,err := connection.Table("users").First()
	fmt.Println(res,err)
	res2,err2 := connection.Table("fd_test").Data(map[string]interface{}{
		"name":"fizz",
	}).Insert()
	fmt.Println(res2,err2)
}
