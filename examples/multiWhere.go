package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	connection := config.GetConnection()
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
