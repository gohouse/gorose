package gorose

import (
	"fmt"
	"testing"
)

func TestEngin(t *testing.T) {
	e := initDB()
	e.EnableQueryLog()
	e.Prefix("pre_")

	db := e.GetQueryDB()

	err := db.Ping()

	fmt.Println(err)
}
