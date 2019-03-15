package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/examples/config"
	"github.com/gohouse/gorose/utils"
)

func main() {
	connection := config.GetConnection()
	connection.DbConfig.Master.EnableQueryLog = true
	// close DB
	defer connection.Close()

	db := connection.NewSession()
	res, err := db.Table("gp_user").Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _,item := range res {
		//fmt.Println(item)
		//return
		db2 := connection.NewSession()
		res,err:=db2.Table("gp_user").Where("Id",item["Id"]).Data(map[string]interface{}{
			"Utm": fmt.Sprintf("%s%v",utils.GetRandomAlarm(2),utils.GetRandomNum(4)),
		}).Update()
		//fmt.Println(db2.LastSql)
		fmt.Println(res,err)
		//return
	}
}
