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

	user, err := db.Table("users a").
		Join("area b", "a.id", "=", "b.uid").
		Where("a.id", ">", 1).
		Get()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(db.LastSql())
	fmt.Println(user)

	// return json
	//res2 := user.Limit(2).Get()
	//fmt.Println(db.LastSql())
	//fmt.Println(user)

}
