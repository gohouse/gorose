package main

import (
	"./config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	//"github.com/devfeel/dotweb"
	"reflect"
)

func main() {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	//user := db.Table("users").Where("id", ">", 1).Where(func() {
	//	db.Where("name", "fizz").OrWhere(func() {
	//		db.Where("name", "fizz2").Where(func() {
	//			db.Where("name", "fizz3").OrWhere("website", "fizzday")
	//		})
	//	})
	//}).Where("job", "it").First()
	//user := GetNewsList(db)
	//	fmt.Println(db.SqlLogs())
	//	fmt.Println(user)
	//
	//	fmt.Println(db)
	//	db.Reset()
	//	fmt.Println(db)
	//var res gorose.MapData
	type Result struct {
		Id int
		Name string
	}
	var res Result
	var res2 gorose.MapData

	news, _ := db.Table("news").
		Where("id","<",3).
		First(&res)
	news2, _ := db.Table("news").
		Where("id","<",3).
		First(&res2)
	fmt.Println(db.LastSql())
	fmt.Println(news)
	fmt.Println(news2)
	//res = news
	fmt.Println(reflect.TypeOf(news))
}

