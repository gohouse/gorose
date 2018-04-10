package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	"github.com/gohouse/gorose/examples/config"
)

func main() {
	db, err := gorose.Open(config.DbConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(db.Table("users"))
	//var db2 gorose.Database
	//
	res1 := db.Table("users").Where("id", "<", 3)
	res2 := db.Table("area")
	//
	r1, _ := res1.Count()
	//r3,_ := res2.Limit(3).Get()
	r2, _ := res2.Get()
	fmt.Println(r1)
	fmt.Println(r2)
	res3, _ := db.Query("select id from users where id=?", 1)

	fmt.Println(res3)
}
