package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
)

func main() {
	connection, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer connection.Close()

	db := connection.GetInstance()

	User := db.Table("users")

	User.Where("name", ">", 1)
	User.Where(func() {
		User.OrWhere("name", "fizz").OrWhere(func() {
			User.Where("name", "fizz2").Where(func() {
				User.Where("name", "fizz53").OrWhere("website", "like", "fizzday%")
			})
		})
	})

	res,_ := User.First()
	fmt.Println(db.LastSql)
	fmt.Println(res)

	res2,_ := User.Where("job", "it").First()
	fmt.Println(db.LastSql)
	fmt.Println(res2)
}
