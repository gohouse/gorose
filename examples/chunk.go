package main

import (
	"./config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
)

func main() {
	db := gorose.Open(config.DbConfig, "mysql_dev")
	// close DB
	defer db.Close()

	db.Table("users").Fields("id, name").Where("id", ">", 2).Chunk(2, func(data []map[string]interface{}) {
		fmt.Println(data)
		for _, item := range data {
			fmt.Println(item["id"], item["name"])
		}
	})
}
