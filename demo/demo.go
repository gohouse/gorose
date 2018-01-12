package main

import (
	"github.com/gohouse/gorose"
	"fmt"
)

var dbConfig = map[string]map[string]string {
	"mysql_dev": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "test",
		"charset":  "utf8",
		"protocol": "tcp",
		"driver":	"mysql",
	},
	"mysql_postgres": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "test",
		"charset":  "utf8",
		"protocol": "tcp",
		"driver":	"postgres",
	},
	"sqlite_dev": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "./test.db",
		"charset":  "utf8",
		"protocol": "tcp",
		"driver":	"sqlite",
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

