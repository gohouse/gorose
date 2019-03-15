package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/examples/config"
)

func main() {
	connection := config.GetConnection()
	// close DB
	defer connection.Close()

	db := connection.NewSession()
	var db2 = connection.NewSession()
	res2, err := db2.Table("users").Fields("id").Limit(2).Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res2)
	fmt.Println(db.JsonEncode(res2))

	var db3 = connection.NewSession()
	res3,err := db3.Table("users").Fields("id,age,job").
		//Where("id", "in", []interface{}{1,55}).
		Where([][]interface{}{{"id", "in", []interface{}{1, 55}}}).
		Get()
	fmt.Println(db3.LastSql)
	fmt.Println(db.JsonEncode(res3))
}
