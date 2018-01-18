package main

import (
	"github.com/gohouse/gorose"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"./config"
)

func main() {
	DB := gorose.Connect.Open(config.DbConfig, "mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db gorose.Database

	user := db.Table("users").Where("id", ">", 1).Where(func() {
		db.Where("name", "fizz").OrWhere(func() {
			db.Where("name", "fizz2").Where(func() {
				db.Where("name", "fizz3").OrWhere("website", "fizzday")
			})
		})
	}).Where("job", "it").First()

	fmt.Println(db.LastSql())
	fmt.Println(user)
}

