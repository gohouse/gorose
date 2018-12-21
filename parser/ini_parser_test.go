package parser

import (
	"fmt"
	"github.com/gohouse/gorose/across"
	"testing"
)

func TestFileParser_Ini(test *testing.T) {
	//var file = "/Users/fizz/go/src/github.com/gohouse/dp/config/mysql.ini"
	var file = across.DemoParserFiles["ini"]

	var confP = &IniConfigParser{}

	var v across.DbConfigCluster
	err := confP.Parse(file, &v)

	if err != nil {
		test.Error("FAIL: ini parser failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: ini parser %v", v))
}
