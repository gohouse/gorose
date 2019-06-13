package gorose

import (
	"fmt"
	"testing"
)

func TestEngin(t *testing.T) {
	e := initDB()
	e.EnableSqlLog()
	e.Prefix("pre_")

	db := e.GetQueryDB()

	err := db.db.Ping()

	fmt.Println(err, db.driver)
}
