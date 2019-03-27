package main

import (
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	connection, err := gorose.Open([]interface{}{"sqlite3","./db.sqlite"}...)
	if err != nil {
		panic(err)
	}

	res,err := connection.Table("users").First()
	fmt.Println(res,err)
}
