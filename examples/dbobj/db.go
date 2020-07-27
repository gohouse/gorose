package dbobj

import (
	"github.com/gohouse/gorose/v2"

	_ "github.com/mattn/go-sqlite3"
)

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

func GetSqliteEngin() *gorose.Engin {
	var err error
	var engin *gorose.Engin
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
