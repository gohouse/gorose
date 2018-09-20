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

	db := connection.NewSession()

	res, err := db.Table("users").Count()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

	max, err := db.Table("users").Max("age")
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

}
