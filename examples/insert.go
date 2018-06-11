package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
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

	data := map[string]interface{}{
		"age":  17,
		"job":  "it2",
		"name": "fizz4",
	}
	res, err := db.Table("users").Data(data).Insert()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("RowsAffected: %d \n", res)
	fmt.Printf("LastInsertId: %d", db.LastInsertId)
}
