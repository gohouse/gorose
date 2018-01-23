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

	user := db.Query("select * from users where id>? limit ?", 1, 2)

	fmt.Println(db.LastSql())
	fmt.Println(user)

	// return json
	//res2 := user.Limit(2).Get()
	//fmt.Println(db.LastSql())
	//fmt.Println(user)

}

