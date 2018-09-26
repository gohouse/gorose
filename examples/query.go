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

	//user, err := db.Query("select * from users where id in (?)", "47,55")
	//tmp := "47,55"
	//tmpArr := strings.Split(tmp, ",")
	//user, err := db.Table("users").WhereIn("id",[]int{47,55}).Get()
	user, err := db.Table("users").WhereBetween("id",[]interface{}{"47",55}).Get()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(db.LastSql)
	fmt.Println(user)

	// return json
	//res2 := user.Limit(2).Get()
	//fmt.Println(db.LastSql())
	//fmt.Println(user)

}
