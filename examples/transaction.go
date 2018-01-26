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

	//var data datas
	data2 := map[string]interface{}{
		"age":  17,
		"job":  "it3",
		"name": "fizz5",
	}
	where := map[string]interface{}{
		"id": 17,
	}

	trans := db.Transaction(func() (error) {
		_,err := db.Table("users").Data(data2).Where(where).Update()
		if err != nil {
			return err
		}

		_,err2 := db.Table("users").Data(data2).Insert()
		if err2 != nil {
			return err2
		}

		return nil
	})

	fmt.Println(trans)
}
