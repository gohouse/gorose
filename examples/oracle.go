package main

import (
	"./config"
	"encoding/json"
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/mattn/go-oci8"
)

func main() {
	db,err := gorose.Open(config.DbConfig, "oracle_dev")
	if err != nil{
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	//res := db.Table("users").First()
	//fmt.Println(res)

	// return json
	res2,err := db.Table("users").Limit(2).Get()
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(db.JsonEncode(res2))

}
