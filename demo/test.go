package main

import (
	"github.com/gohouse/gorose"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"./config"
)

func main() {
	db := gorose.Open(config.DbConfig, "mysql_dev")
	// close DB
	defer db.Close()

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

func TTT()  {
	fmt.Sprintf("TTT")
}
