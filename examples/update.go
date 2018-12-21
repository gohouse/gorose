package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	connection := config.GetConnection()

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

}
