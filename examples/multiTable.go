package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/examples/config"
)

func main() {
	connection := config.GetConnection()
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
