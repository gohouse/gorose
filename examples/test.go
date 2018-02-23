package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	//"github.com/devfeel/dotweb"
	"github.com/gohouse/gorose/utils"
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
	if err != nil {
		fmt.Println(err)
		return
	}
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
	var data = map[string]string{
		"created_at":"2018-02-01 00:00:00",
	}

	res,err := db.Table("users").Data(data).Insert()

	fmt.Println(res)
	fmt.Println(err)

	res2,_ := db.Table("users").Order("id desc").First()

	fmt.Println(res2)
	fmt.Println(res2["created_at"].(string))

	fmt.Println(utils.JsonEncode(utils.SuccessReturn("success")))
	fmt.Println(utils.JsonEncode(utils.SuccessReturn("", 200)))
}
