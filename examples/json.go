package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/utils"
)

func main() {
	connection := config.GetConnection()
	// close DB
	defer connection.Close()

	// return json
	res2, err := connection.Table("users").Limit(2).Get()
	if err != nil {
		fmt.Println(err)
		return
	}

	// or
	fmt.Println(utils.JsonEncode(res2))

	//============== result ======================

	//SELECT * FROM users WHERE  id > '2' LIMIT 1
	//{"age":18,"id":3,"job":"go orm","name":"gorose","website":"go-rose.com"}
	//SELECT * FROM users LIMIT 2
	//[{"age":18,"id":1,"job":"it","name":"fizz","website":"fizzday.net"},{"age":18,"id":2,"job":"engineer","name":"fizzday","website":"fizzday.net"}]

}
