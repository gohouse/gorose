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

	user := db.Execute("update users set name=? where id=?", "fizz8", 4)

	fmt.Println(db.LastSql())
	fmt.Println(user)
}

