package parser

import (
	"fmt"
	"github.com/gohouse/gorose/config"
	"testing"
)

func TestFileParser_Ini(test *testing.T) {
	//var file = "/Users/fizz/go/src/github.com/gohouse/dp/config/mysql.ini"
	var file = config.DemoParserFiles["ini"]

	var confP = &IniConfigParser{}

	pr, err := confP.Parse(file)

	if err != nil {
		test.Error("FAIL: ini parser failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: ini parser %v", pr))
}
