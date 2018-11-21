package main

import (
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//connection := config.GetConnection()
	connection,err := gorose.Open(&gorose.DbConfigSingle{
		Driver:          "sqlite3",                    // 驱动: mysql/sqlite/oracle/mssql/postgres
		EnableQueryLog:  true,                       // 是否开启sql日志
		SetMaxOpenConns: 0,                          // (连接池)最大打开的连接数，默认值为0表示不限制
		SetMaxIdleConns: 0,                          // (连接池)闲置的连接数
		Prefix:          "",                         // 表前缀
		Dsn:             "resource/sqlite/t.db", // 数据库链接
	})
	if err!=nil {
		panic("数据库链接失败")
	}
	// close DB
	defer connection.Close()

	sql_table := `
    CREATE TABLE IF NOT EXISTS article(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content text NULL,
        created_at TIMESTAMP default (datetime('now', 'localtime'))  
    );
    `
	connection.Execute(sql_table)

	db := connection.NewSession()

	db.Table("article").Data(map[string]interface{}{
		"content":`{"title":"标题啊","author":"fizz"}`,
	}).Insert()

	// return json
	res2, err := db.Table("article").Limit(2).Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db.LastSql)
	fmt.Println(db.JsonEncode(res2))

}
