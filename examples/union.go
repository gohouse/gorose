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

	res := db.Table("users").Count()
	fmt.Println(res)

	max := db.Table("users").Max("money")
	fmt.Println(max)

	min := db.Table("users").Min("age")
	fmt.Println(min)

	avg := db.Table("users").Avg("age")
	fmt.Println(avg)

	sum := db.Table("users").Sum("age")
	fmt.Println(sum)

	fmt.Println(db.LastSql())

}
