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

	res, err := db.Table("users").Where("id", ">", 1).First()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(db.LastSql)
	fmt.Println(res)


	db.Table("users").Where("id", 47).Decrement("age",2)

	fmt.Println(db.LastSql)

	res3, _ := db.Table("users").Where("id", ">", 1).First()
	fmt.Println(res3)
}
