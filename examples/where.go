package main

import (
	"fmt"

	"github.com/leeyisoft/gorose/v2/examples/dbobj"
)

func main() {
	var res interface{}
	var err error
	db := dbobj.Getdb()
	res, err = db.Reset().Table("users").Get()
	if err != nil {
		fmt.Println("1 - err:", err.Error())
	}
	fmt.Println("1 - res:", res)

	res, err = db.Reset().Table("users").OrderBy("uid desc").Limit(1).Paginate(1)
	if err != nil {
		fmt.Println("2 - err:", err.Error())
	}
	fmt.Println("2 - res:", res)

	res, err = db.Reset().Table("users").Where("age", ">", 0).Limit(1).Paginate(3)
	if err != nil {
		fmt.Println("2 - err:", err.Error())
	}
	fmt.Println("2 - res:", res)
}
