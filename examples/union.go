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

	res, err := db.Table("users").Count()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

	max, err := db.Table("users").Max("money")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(max)

	min, err := db.Table("users").Min("age")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(min)

	avg, err := db.Table("users").Avg("age")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(avg)

	sum, err := db.Table("users").Sum("age")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sum)

	fmt.Println(db.LastSql())

}
