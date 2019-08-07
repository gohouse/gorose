package gorose

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func initDB() *Engin {
	e, err := Open(&Config{Driver: "sqlite3", Dsn: "./db.sqlite"})

	if err != nil {
		panic(err.Error())
	}

	e.TagName("orm")
	e.IgnoreName("ignore")

	initTable(e)

	return e
}

func initTable(e *Engin) {
	var sql = `CREATE TABLE IF NOT EXISTS "users" (
	 "uid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "name" TEXT NOT NULL default "",
	 "age" integer NOT NULL default 0
)`
	var s = e.NewSession()
	var err error
	var aff int64

	aff, err = s.Execute(sql)
	if err != nil {
		return
	}
	if aff == 0 {
		return
	}

	aff, err = s.Execute("insert into users(name,age) VALUES(?,?),(?,?),(?,?)",
		"fizz", 18, "gorose", 19, "fizzday", 20)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("初始化数据和表成功:", aff)
}
