package main

import (
	"./config"
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	DB := gorose.Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db gorose.Database

	data := map[string]interface{}{
		"age":17,
		"job":"it3",
		"name":"fizz4",
	}
	res := db.Table("users").Data(data).Insert()
	fmt.Println(db.LastSql())
	fmt.Println(res)

}
