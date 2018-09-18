package gorose

import (
	"fmt"
	_ "github.com/gohouse/gorose/driver/mysql"
	"github.com/gohouse/gorose/across"
	"testing"
)
var conn *Connection
var err error

func InitConn() (*Connection,error) {
	conn,err = Open("json", across.DemoParserFiles["json"])
	return conn,err
}

func TestGorose_Query(test *testing.T) {
	conn,err = InitConn()
	if err != nil {
		test.Error("FAIL: open failed.", err)
		return
	}
	res, err := conn.Query("select * from users limit 1")
	if err != nil {
		test.Error("FAIL: Query failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: Query: %v", res))
}

func TestGorose_Execute(test *testing.T) {
	conn,err = InitConn()
	if err != nil {
		test.Error("FAIL: open failed.", err)
		return
	}
	res, err := conn.Execute("update users set job='it22' where id=47")
	if err != nil {
		test.Error("FAIL: Execute failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: Execute: %v", res))
}


