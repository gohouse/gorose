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
	//defer db.Close()
	//var Users = db.Table("users")
	user, err := connection.Query("select * from users")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user)

	db := connection.GetInstance()

	lid,err := db.Execute("insert into users(name,age) values('fizz3', 19)")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(lid)
	fmt.Println(db.LastSql)
	fmt.Println(db.LastInsertId)
}
