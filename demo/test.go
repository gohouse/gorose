package main

import (
	"github.com/gohouse/gorose"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"./config"
)

func main() {
	DB := gorose.Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db gorose.Database

	user := db.Table("sql1").Get()

	for _,item := range user.([]map[string]interface{}) {
		res := fmt.Sprintf("update fd_organization_info set code='%s' where jg_uuid='%s'", item["code"], item["uuid"])
		fmt.Println(res)
	}

	// return json
	//res2 := user.Limit(2).Get()
	//fmt.Println(db.LastSql())
	//fmt.Println(user)

}

