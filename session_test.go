package gorose

import (
	"fmt"
	"testing"
)
type users struct {
	Name string
}
func TestSession_Api(test *testing.T) {
	var db = &Database{}
	//var b = "users"
	var b = users{"fizz"}
	sql,err := db.Table(&b).
		Distinct().
		//Fields("id as uid","name").
		Where("a",1).
		OrderBy("id desc").
		GroupBy("name").
		Having("count(uid)>0").
		Limit(1).
		Offset(2).
		BuildSql()
	if err != nil {
		test.Error("FAIL: orm failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: orm: %v", sql))
}
