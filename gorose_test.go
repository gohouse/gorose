package gorose

import (
	_ "github.com/mattn/go-sqlite3"
)

func initDB() IEngin {
	e, err := Open(&Config{Driver: "sqlite3", Dsn: "./db.sqlite"})

	if err != nil {
		panic(err.Error())
	}

	return e
}
