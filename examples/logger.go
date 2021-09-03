package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() *gorose.Engin {
	e, err := gorose.Open(&gorose.Config{Driver: "sqlite3", Dsn: "./db.sqlite"})

	if err != nil {
		panic(err.Error())
	}

	// 这里可以设置日志相关配置
	//e.Use(func(eg *gorose.Engin) {
	//	eg.SetLogger(gorose.NewLogger(&gorose.LogOption{
	//		FilePath:       "./log",
	//		EnableSqlLog:   true,
	//		EnableSlowLog:  5,
	//		EnableErrorLog: true,
	//	}))
	//})

	return e
}
func main() {
	db := initDB().NewOrm()
	res, err := db.Table("users").First()
	fmt.Println(err)
	fmt.Println(db.LastSql())
	fmt.Println(res)
}
