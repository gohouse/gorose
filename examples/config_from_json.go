package main

import (
	"fmt"
	"encoding/json"
	"github.com/gohouse/gorose"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// 定义config的struct
	type Configer struct {
		Default         string
		SetMaxOpenConns int
		SetMaxIdleConns int
		Connections     map[string]map[string]string
	}

	jsons := `{"Default":"mysql_dev","SetMaxOpenConns":300,"SetMaxIdleConns":10,"Connections":{"mysql_dev":{"charset":"utf8","database":"test","driver":"mysql","host":"192.168.200.248","password":"gcore","port":"3306","prefix":"","protocol":"tcp","username":"gcore"}}}`

	var conf Configer

	json.Unmarshal([]byte(jsons), &conf)

	var confReal = map[string]interface{}{
		"Default":         conf.Default,
		"SetMaxOpenConns": conf.SetMaxOpenConns,
		"SetMaxIdleConns": conf.SetMaxIdleConns,
		"Connections":     conf.Connections,
	}

	//	fmt.Println(confReal)
	//return
	db, err := gorose.Open(confReal, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	res, err := db.Table("users").Where("id>2").First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql())
	fmt.Println(res["id"])
}
