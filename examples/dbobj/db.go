package dbobj

import (
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	GetSqliteEngin()
}

func GetMysqlEngin() *gorose.Engin {
	var err error
	var engin *gorose.Engin
	engin, err = gorose.Open(&gorose.Config{
		Driver: "mysql",
		Dsn:    "root:123456@tcp(localhost:3306)/test?charset=utf8mb4",
		Prefix: "nv_",
	})
	if err != nil {
		panic(err.Error())
	}
	return engin
}

var engin *gorose.Engin

func GetSqliteEngin() *gorose.Engin {
	var err error
	engin, err = gorose.Open(&gorose.Config{
		Driver: "sqlite3",
		Dsn:    "./db.sqlite",
		Prefix: "",
	})
	if err != nil {
		panic(err.Error())
	}
	return engin
}

func Getdb() gorose.IOrm {
	return engin.NewOrm()
}
