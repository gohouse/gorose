package parser

import (
	"fmt"
	"github.com/gohouse/gorose/across"
	"testing"
)

func TestFileParser_Json(test *testing.T) {
	//var file = "/Users/fizz/go/src/github.com/gohouse/gorose/examples/demoParserFiles/mysql.json"
	var file = across.DemoParserFiles["json"]

	var confP = &JsonConfigParser{}

	var v across.DbConfigCluster
	err := confP.Parse(file, &v)

	if err != nil {
		test.Error("FAIL: json parser failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: json parser %v", v.Master))
}
