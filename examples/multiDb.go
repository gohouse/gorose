package main

import (
	"github.com/gohouse/gorose"
	"github.com/gohouse/gorose/examples/config"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

func main() {
	db,err := gorose.Open(config.DbConfig)
	if err!=nil{
		fmt.Println(err)
		return
	}
	//fmt.Println(db.Table("users"))
	//var db2 gorose.Database
	//
	res1 := db.Table("users").Where("id","<",3)
	res2 := db.Table("area")
	//
	r1,_ := res1.Where("id",1).Count()
	r3,_ := res2.Limit(3).Get()
	//r2,_ := res2.Count()
	fmt.Println(r1)
	fmt.Println(r3)
}
