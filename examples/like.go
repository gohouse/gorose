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

	db := connection.NewSession()
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
