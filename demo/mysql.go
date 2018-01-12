package main

import (
	"github.com/gohouse/gorose"
	"fmt"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"./config"
)

func main() {
	DB := gorose.Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db gorose.Database

	res := db.Table("users").Distinct().First()
	fmt.Println(db.LastSql())
	fmt.Println(res)

	// return json
	res2 := db.Table("users").Limit(2).Get()
	jsons, _ := json.Marshal(res2)
	fmt.Println(string(jsons))

}

