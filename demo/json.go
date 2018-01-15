package main

import (
	"github.com/gohouse/gorose"
	_ "github.com/go-sql-driver/mysql"
	"./config"
	"github.com/gohouse/utils"
	"fmt"
)

func main() {
	// connect db
	DB := gorose.Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db gorose.Database

	// return json
	res2 := db.Table("users").Limit(2).Get()
	fmt.Println(db.LastSql())
	fmt.Println(utils.JsonEncode(res2))
	// or
	fmt.Println(db.JsonEncode(res2))

	//============== result ======================

	//SELECT * FROM users WHERE  id > '2' LIMIT 1
	//{"age":18,"id":3,"job":"go orm","name":"gorose","website":"go-rose.com"}
	//SELECT * FROM users LIMIT 2
	//[{"age":18,"id":1,"job":"it","name":"fizz","website":"fizzday.net"},{"age":18,"id":2,"job":"engineer","name":"fizzday","website":"fizzday.net"}]

}

