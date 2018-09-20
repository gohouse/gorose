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

	res2, err2 := db.Table("test").
		Fields("id", "test_json->'$.name' as name").
		Where("test_json->'$.age'", 18).First()
	fmt.Println(res2,err2)


	db2 := connection.NewSession()
	res3, err3 := db2.Table("test").
		Where("id", 1).
		Data(map[string]interface{}{
			"test_json->'$.age'":20,
		}).
		Update()
	fmt.Println(db2.LastSql)
	fmt.Println(res3,err3)
}
