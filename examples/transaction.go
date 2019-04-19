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

	//var data datas
	data2 := map[string]interface{}{
		"age":  17,
		"job":  "it3",
		"name": "fizz4",
	}
	where := map[string]interface{}{
		"id": 17,
	}

	db := connection.NewSession()

	err := db.Transaction(func() error {

		res2, err2 := db.Table("users").Data(data2).Insert()
		if err2 != nil {
			return err2
		}

		res1, err := db.Table("users").Data(data2).Where(where).Update()
		if err != nil {
			return err
		}

		fmt.Println(res1,res2)

		return nil
	})

	fmt.Println(err)
}
