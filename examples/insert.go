package main

import (
	"github.com/gohouse/gorose/examples/config"
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
	User := db.Table("users")
	res, err := User.Data(data).Insert()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql())
	fmt.Printf("RowsAffected: %d \n", User.RowsAffected)
	fmt.Printf("LastInsertId: %d", User.LastInsertId)
}
