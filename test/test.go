package main

import (
	"github.com/gohouse/gorose"
	"github.com/gohouse/gorose/config"
	_ "github.com/gohouse/gorose/driver/mysql"
	"fmt"
)

type users struct {
	Name string `orm:"name"`
	Age int `orm:"age"`
	Job string `orm:"job"`
}
func main() {
	connection,_ := gorose.Open("json", config.DemoParserFiles["json"])

	db := connection.NewDB()

	var u users
	var u2 []users

	db.Table(&u).Select()
	err := db.Table(&u2).Select()

	fmt.Println(u, u2, err)
	//fmt.Println(u.Age)
	db2 := connection.NewDB()
	res,err := db2.Table("users").First()
	fmt.Println(res, err)
}
