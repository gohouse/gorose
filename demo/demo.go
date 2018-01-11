package main

import (
	"github.com/gohouse/gorose"
	"fmt"
)

var dbConfig = map[string]map[string]string {
	"mysql_master": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "test",
		"charset":  "utf8",
		"protocol": "tcp",
	},
	"mysql_dev": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "test",
		"charset":  "utf8",
		"protocol": "tcp",
	},
}

func main() {
	// open a db connection
	DB := gorose.Open(dbConfig, "mysql_dev")
	// close db
	defer DB.Close()
	// get the db chaining object
	var db gorose.Database

	res := db.Table("users").First()
	res2 := db.Table("users").Limit(2).Get()

	fmt.Println(res)
	fmt.Println(res2)
}

