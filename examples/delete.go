package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/examples/config"
)

func main() {
	db := config.GetConnection()
	// close DB
	defer db.Close()

	//where := map[string]interface{}{
	//	"id": 17,
	//}
	res, err := db.Table("users").
		//Force().
		//Where(where).
		Delete()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

}
