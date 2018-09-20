package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	//_ "github.com/mattn/go-oci8"
)

func main() {
	connection := config.GetConnection()
	// close DB
	defer connection.Close()

	db := connection.NewSession()

	//res := db.Table("users").First()
	//fmt.Println(res)

	// return json
	res2, err := db.Table("users").Limit(2).Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(db.JsonEncode(res2))

}
