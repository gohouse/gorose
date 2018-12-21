package gorose

import (
	"fmt"
	"github.com/gohouse/gorose/across"
	"testing"
)

func TestGorose_Open(test *testing.T) {
	conn, err = InitConn()
	if err != nil {
		test.Error("FAIL: open failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: open: %v", conn.Db.MasterDb))
}

func TestGorose_Open_Struct(test *testing.T) {
	conn,err = Open(&DbConfigSingle{
		Driver:          "mysql",                                                   // 驱动: mysql/sqlite/oracle/mssql/postgres
		EnableQueryLog:  false,                                                     // 是否开启sql日志
		SetMaxOpenConns: 0,                                                         // (连接池)最大打开的连接数，默认值为0表示不限制
		SetMaxIdleConns: 0,                                                         // (连接池)闲置的连接数
		Prefix:          "",                                                        // 表前缀
		Dsn:             "gcore:gcore@tcp(192.168.200.248:3306)/test?charset=utf8", // 数据库链接
	})
	if err != nil {
		test.Error("FAIL: open failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: open: %v", conn.Db.MasterDb))
}

func TestGorose_Open_Json(test *testing.T) {
	conn,err = Open("json", across.DemoParserFiles["json"])
	if err != nil {
		test.Error("FAIL: open failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: open: %v", conn.Db.MasterDb))
}

func TestGorose_Version(test *testing.T) {
	test.Log(fmt.Sprintf("PASS: %v", VERSION))
}
