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

	user, err := db.Execute("update users set name=? where id=?", "fizz8", 4)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(user)
}
