package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 读写分离集群
	var config1 = gorose.Config{Dsn: "./db.sqlite"}
	var config2 = gorose.Config{Dsn: "./db2.sqlite"}
	var config3 = gorose.Config{Dsn: "./db3.sqlite"}
	var config4 = gorose.Config{Dsn: "./db4.sqlite"}
	var configCluster = &gorose.ConfigCluster{
		Master: []gorose.Config{config3, config4},
		Slave:  []gorose.Config{config1, config2},
		Driver: "sqlite3",
	}
	engin, err := gorose.Open(configCluster)
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < 5; i++ {
		db := engin.NewOrm()
		res, err := db.Table("users").First()
		fmt.Println(err)
		fmt.Println(res)
	}
}
