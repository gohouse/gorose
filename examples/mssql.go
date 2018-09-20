package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	connection := config.GetConnection()
	// close DB
	defer connection.Close()

	db := connection.NewSession()

	res, err := db.Table("users").Where("id", ">", 2).First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(res)

	// return json
	res2, err := db.Table("users").Limit(2).Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(res2)
	//fmt.Println(db.JsonEncode(res2))

	//============== result ======================

	//SELECT * FROM users WHERE  id > '2' LIMIT 1
	//map[id:3 name:gorose age:18 website:go-rose.com job:go orm]
	//SELECT * FROM users LIMIT 2
	//[map[id:1 name:fizz age:18 website:fizzday.net job:it] map[id:2 name:fizzday age:18 website:fizzday.net job:engineer]]

}
