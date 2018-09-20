package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/examples/config"
)

func main() {
	connection := config.GetConnection()
	// close DB
	defer connection.Close()

	db := connection.NewSession()
	fmt.Println(db)
	res, err := db.Table("users").Where("id", "<", 1).First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(res))
	fmt.Println(db.LastSql)
	fmt.Println(res)

	var db2 = connection.NewSession()
	res2, err := db2.Table("users").Limit(2).Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res2)
	fmt.Println(db.JsonEncode(res2))

	//============== result ======================

	//SELECT * FROM users WHERE  id > '2' LIMIT 1
	//map[id:3 name:gorose age:18 website:go-rose.com job:go orm]
	//SELECT * FROM users LIMIT 2
	//[map[id:1 name:fizz age:18 website:fizzday.net job:it] map[id:2 name:fizzday age:18 website:fizzday.net job:engineer]]

}
