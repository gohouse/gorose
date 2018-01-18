package main

import (
	"github.com/gohouse/gorose"
	"fmt"
	"encoding/json"
	_ "github.com/mattn/go-oci8"
	"./config"
)

func main() {
	//gorose.T
	// open a db connection
	DB := gorose.Connect.Open(config.DbConfig, "mysql_dev")
	// close db
	defer DB.Close()
	// get the db chaining object
	var db gorose.Database

	//res := db.Table("users").First()
	//fmt.Println(res)

	// return json
	res2 := db.Table("users").Limit(2).Get()
	jsons, _ := json.Marshal(res2)
	fmt.Println(string(jsons))

}

