package main

import (
	"./config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	"github.com/gohouse/utils"
)

func main() {
	// connect db
	db := gorose.Open(config.DbConfig, "mysql_dev")
	// close DB
	defer db.Close()

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
