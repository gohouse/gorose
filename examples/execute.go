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

	user := db.Execute("update users set name=? where id=?", "fizz8", 4)

	fmt.Println(db.LastSql())
	fmt.Println(user)
}

