package main

import (
	"github.com/gohouse/gorose"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"./config"
	//"github.com/devfeel/dotweb"
	"time"
)

func main() {
	db := gorose.Open(config.DbConfig, "mysql_dev")
	// close DB
	defer db.Close()
	for i := 1; i < 10; i++ {
		go insertUpdate(db, i)
	}
	time.Sleep(10*time.Second)
}

func insertUpdate(db gorose.Database, i int)  {
	db.Transaction(func() {
		data := map[string]interface{}{
			"age":i,
		}
		//where := map[string]interface{}{
		//	"id":i,
		//}

		//res := db.Table("users").Data(data).Where(where).Update()
		res := db.Table("users").Data(data).Insert()
		fmt.Println(gorose.GetDB())
		fmt.Println(res)


		data2 := map[string]interface{}{
			"age":i,
			"name":"222",
		}
		res2 := db.Connect("mysql_dev2").Table("users_copy").Data(data2).Insert()
		//fmt.Println(db.LastSql())
		fmt.Println(res2)
		fmt.Println(gorose.GetDB())



	})

}
