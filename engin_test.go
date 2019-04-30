package gorose

import (
	"fmt"
	"testing"
)

func TestEngin(t *testing.T) {
	e := initDB()
	e.EnableQueryLog()
	e.Prefix("pre_")

	db,dr := e.GetQueryDB()

	err := db.Ping()

	fmt.Println(err,dr)
}
