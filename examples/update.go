package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
)

func main() {
	connection, err := gorose.Open(config.DbConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	// close DB
	defer connection.Close()

	db := connection.NewSession()

	data := map[string]interface{}{
		"age":  17,
		"job":  "it33",
		"name": "fizz5",
	}
	where := map[string]interface{}{
		"id": 75,
	}

	fmt.Println(db.Table("users").Data(data).Where(where).BuildExecut("update"))
	res, err := db.Table("users").Data(data).Where(where).Update()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)

}
