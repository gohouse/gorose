package main

import (
	"./config"
	"encoding/json"
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := gorose.Open(config.DbConfig, "sqlite_dev")
	// close DB
	defer db.Close()

	//res := db.Table("users").First()
	//fmt.Println(res)

	// return json
	res2 := db.Table("users").Limit(2).Get()
	jsons, _ := json.Marshal(res2)
	fmt.Println(string(jsons))

}
