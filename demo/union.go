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

	res := db.Table("users").Count()
	fmt.Println(res)

	max := db.Table("users").Max("money")
	fmt.Println(max)

	min := db.Table("users").Min("age")
	fmt.Println(min)

	avg := db.Table("users").Avg("age")
	fmt.Println(avg)

	sum := db.Table("users").Sum("age")
	fmt.Println(sum)

	fmt.Println(db.LastSql())

}

