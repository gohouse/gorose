package gorose

import (
	"fmt"
	"testing"
)
type users struct {
	Name string
}
var db1 *Session
func InitOrm() {
	db1 = NewOrm()
	db1.Connection = &Connection{DbConfig:&DbConfigCluster{Master:&DbConfigSingle{
		"mysql",
		true,
		0,
		0,
		"",
		"",
	}}}
}
func TestSession_QueryApi(test *testing.T) {
	InitOrm()
	var b = users{"fizz"}
	sql,err := db1.Table(&b).
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
func TestSession_Insert(test *testing.T) {
	InitOrm()
	var b = users{"fizz"}
	sql,err := db1.Table(&b).
		Data(map[string]interface{}{"name":"fizz333", "age":19}).
		BuildSql("insert")
	if err != nil {
		test.Error("FAIL: orm failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: orm: %v", sql))
}
func TestSession_Update(test *testing.T) {
	InitOrm()
	//var b = "users"
	var b = users{"fizz"}
	sql,err := db1.Table(&b).
		Data(map[string]interface{}{"name":"fizz333", "age":18}).
		Where("a","1").
		BuildSql("update")
	if err != nil {
		test.Error("FAIL: orm failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: orm: %v", sql))
}
func TestSession_Delete(test *testing.T) {
	InitOrm()
	//var b = "users"
	var b = users{"fizz"}
	sql,err := db1.Table(&b).
		Where("a",1).
		BuildSql("delete")
	if err != nil {
		test.Error("FAIL: orm failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: orm: %v", sql))
}
func TestSession_InsertMulti(test *testing.T) {
	InitOrm()
	//var b = "users"
	var b = users{"fizz"}
	sql,err := db1.Table(&b).
		Data([]map[string]interface{}{{"name":"fizz333","age":10},{"name":"fizz222","age":20}}).
		BuildSql("insert")
	if err != nil {
		test.Error("FAIL: orm failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: orm: %v", sql))
}
