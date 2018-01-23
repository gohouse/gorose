package main

import (
	"github.com/gohouse/gorose"
	"fmt"
	"encoding/json"
	_ "github.com/mattn/go-oci8"
	"./config"
)

func main() {
	db := gorose.Open(config.DbConfig, "oracle_dev")
	// close DB
	defer db.Close()

	//res := db.Table("users").First()
	//fmt.Println(res)

	// return json
	res2 := db.Table("users").Limit(2).Get()
	jsons, _ := json.Marshal(res2)
	fmt.Println(string(jsons))

}

