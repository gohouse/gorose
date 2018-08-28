package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
)

func main() {
	fmt.Println(config.DbConfig)
	connection, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer connection.Close()

	db := connection.GetInstance()
	fmt.Println(db)
	res, err := db.Table("users").
		Where("name", "like", "fizz%").
		First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(res))
	fmt.Println(db.LastSql)
	fmt.Println(res)
}
