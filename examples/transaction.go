package main

import (
	"fmt"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/leeyisoft/gorose/v2"
)

func initEngin() *gorose.Engin {
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
func main() {
	err := trans()
	fmt.Println(err)
}

func trans() error {
	var aff int64
	var err error
	var engin = initEngin()
	db := engin.NewOrm()
	db.Begin()
	aff, err = db.Table("logs").Insert(gorose.Data{"username": "xx"})
	if err != nil || aff == 0 {
		db.Rollback()
		return err
	}
	aff, err = db.Table("logs").Where("id", 2).Update(gorose.Data{"username": "xx1232xx1232xx1232xx1232xx1232xx1232xx1232xx1232"})
	if err != nil || aff == 0 {
		db.Rollback()
		return err
	}
	aff, err = db.Table("logs").Insert(gorose.Data{"username": "xx", "pkid": 1})
	if err != nil || aff == 0 {
		db.Rollback()
		return err
	}
	return db.Commit()
}
