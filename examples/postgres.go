package main

import (
	"fmt"
	//_ "github.com/lib/pq"
)

func main() {
	// 测试时, 不要忘记打开import内的pg包注释
	pgtest()
}
func main2() {
	dsn := "user=postgres dbname=postgres password=123456 sslmode=disable"
	engin, err := gorose.Open(&gorose.Config{Driver: "postgres", Dsn: dsn})
	if err != nil {
		panic(err.Error())
	}
	var orm = engin.NewOrm()
	res, err := orm.Query("select * from users where uid>$1", 1)
	fmt.Println(res, err)

	fmt.Println(engin.NewOrm().Table("users").
		Data(map[string]interface{}{"uname": "fizz22"}).
		//Where("uid",4).BuildSql("insert"))
		Where("uid", 4).BuildSql("update"))
}

func pgtest() {
	dsn := "user=postgres dbname=postgres password=123456 sslmode=disable"
	engin, err := gorose.Open(&gorose.Config{Driver: "postgres", Dsn: dsn})
	if err != nil {
		panic(err.Error())
	}
	var orm = engin.NewOrm()
	res, p, err := orm.Table("users").Where("a", 1).BuildSql()
	fmt.Println(res, p, err)
}
