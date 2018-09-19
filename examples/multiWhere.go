package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
)

func main() {
	connection, err := gorose.Open(config.DbConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer connection.Close()

	db := connection.NewSession()
	fmt.Println(db)
	res, err := db.Table("users").Where([][]interface{}{{"id", ">", 2}}).First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(res)
}
