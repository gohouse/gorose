package main

import (
	"github.com/gohouse/gorose"
	_ "github.com/go-sql-driver/mysql"
	"./config"
	"fmt"
)

func main() {
	DB := gorose.Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db gorose.Database

	// 查出总数
	//total := db.Table("codes").Count()
	limit := 2

	//times := int(math.Ceil(float64(total/limit)))
	times := 10

	for i:=0;i<times;i++{
		// 查询需要查询的数据
		//data := db.Table("codes_copy").Fields("code3").Limit(2).Offset(i*limit).Get()
		data := db.Table("codes_copy").Fields("code3").Limit(limit).Offset(i*limit).Get()

		for _,item := range data.([]map[string]interface{}) {
			code1 := item["code3"]

			//code查询是否在code2中1
			exists := db.Table("codes_copy").Where("code4", code1).Count()
			fmt.Println(exists)
			if exists>0{
				db.Table("codes_copy").Where("code3", code1).Data(map[string]int{"tag1":1}).Update()
				fmt.Println(db.LastSql())
			} else {

				res := db.Table("codes_copy").Where("code3", code1).Data(map[string]int{"tag1":1}).Update()
				fmt.Println(db.LastSql())
				fmt.Println(res)

			}


		}
	}
}

