package gorose

import (
	"fmt"
	"testing"
)

func TestDatabase_Session(test *testing.T) {
	var db = NewOrm()
	sql,err := db.Table("users").BuildSql()
	if err != nil {
		test.Error("FAIL: orm failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: orm: %v", sql))
}
