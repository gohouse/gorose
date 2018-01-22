package main

import (
	"github.com/gohouse/gorose"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"./config"
	//"github.com/devfeel/dotweb"
)

func main() {
	db := gorose.Open(config.DbConfig, "mysql_dev")
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
	news := db.Table("news").
		Where("status", 1).
		Order("id desc").
		Limit(10).
		Page(1).
		First()

	fmt.Println(db.LastSql())
		fmt.Println(news)
}


// 获取列表
func GetNewsList(db gorose.Database) interface{} {
	return db.Table("news").
		Where("status", 1).
		Order("id desc").
		Limit(10).
		Page(1).
		Get()
}

func TTT()  {
	fmt.Sprintf("TTT")
}
