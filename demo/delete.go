package main

import (
	"./config"
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := gorose.Open(config.DbConfig, "mysql_dev")
	// close DB
	defer db.Close()

	where := map[string]interface{}{
		"id":17,
	}
	res := db.Table("users").Where(where).Delete()
	fmt.Println(db.LastSql())
	fmt.Println(res)

}
