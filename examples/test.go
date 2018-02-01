package main

import (
	"./config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	//"github.com/devfeel/dotweb"
)

func main() {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	//defer db.Close()
	var Users = db.Table("users")
	user, err := Users.Where("id", ">", 1).Where(func() {
		Users.Where("name", "fizz").OrWhere(func() {
			Users.Where("name", "fizz2").Where(func() {
				Users.Where("name", "fizz3").OrWhere("website", "fizzday")
			})
		})
	}).Where("job", "it").First()
	//user := GetNewsList(db)
	fmt.Println(db.SqlLogs())
	fmt.Println(user)
	//
	//	fmt.Println(db)
	//	db.Reset()
	//	fmt.Println(db)
	//var res gorose.MapData
	//news, _ := db.Table("users").
	//			Where("id","<",3).
	//			First()
	//
	//fmt.Println(db.LastSql())
	//fmt.Println(news)
}
