package main

import (
	"./config"
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	DB := gorose.Connect.Open(config.DbConfig, "mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db gorose.Database

	where := map[string]interface{}{
		"id":17,
	}
	res := db.Table("users").Where(where).Delete()
	fmt.Println(db.LastSql())
	fmt.Println(res)

}
