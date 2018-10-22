package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	"github.com/gohouse/gorose/examples/config"
)

func main() {
	connection := config.GetConnection()
	connection.Use(gorose.NewLogger())

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

	res, err := db.Table("users").Data(data).Where(where).Update()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)

	connection.Table("users").First()

}
