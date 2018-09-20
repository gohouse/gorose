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

	User := db.Table("users")

	//User.Where("id", ">", 1)
	//User.Where("id", "in", []interface{}{47})
	User.WhereNotIn("id", []interface{}{47})

	res,_ := User.First()
	fmt.Println(db.LastSql)
	fmt.Println(res)

	//res2,_ := User.Where("job", "it").First()
	//fmt.Println(db.LastSql)
	//fmt.Println(res2)
}
