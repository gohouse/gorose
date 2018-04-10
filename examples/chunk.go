package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	"github.com/gohouse/gorose/examples/config"
)

func main() {
	connection, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer connection.Close()

	db := connection.GetInstance()

	db.Table("users").Fields("id, name").Where("id", ">", 2).Chunk(2, func(data []map[string]interface{}) {
		fmt.Println(data)
		for _, item := range data {
			fmt.Println(item["id"], item["name"])
		}
	})
}
