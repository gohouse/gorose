package main

import (
	"./config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
)

func main() {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	data := map[string]interface{}{
		"age":  17,
		"job":  "it2",
		"name": "fizz4",
	}
	res, err := db.Table("users").Data(data).Insert()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql())
	fmt.Println(res)

}
