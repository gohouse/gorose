package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/examples/config"
)

func main() {
	connection := config.GetConnection()
	// close DB
	defer connection.Close()

	db := connection.NewSession()

	db.Table("users").Fields("id, name").Where("id", ">", 2).Chunk(2, func(data []map[string]interface{}) {
		fmt.Println(data)
		for _, item := range data {
			fmt.Println(item["id"], item["name"])
		}
	})
}
