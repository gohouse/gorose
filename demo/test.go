package main

import (
	"github.com/gohouse/gorose"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"./config"
)

func main() {
	DB := gorose.Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db gorose.Database

	user := db.Table("users").Where("id",">",2)
	fmt.Println(user)
	res := db.First()
	fmt.Println(db.LastSql())
	fmt.Println(res)

	// return json
	//res2 := user.Limit(2).Get()
	//fmt.Println(db.LastSql())
	fmt.Println(user)

}

