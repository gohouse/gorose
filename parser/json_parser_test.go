package parser

import (
	"fmt"
	"fizzday.com/gohouse/gorose/config"
	"testing"
)

func TestFileParser_Json(test *testing.T) {
	//var file = "/Users/fizz/go/src/github.com/gohouse/laboratory/dp/config/mysql.json"
	var file = config.DemoParserFiles["json"]

	var confP = &JsonConfigParser{}

	pr, err := confP.Parse(file)

	if err != nil {
		test.Error("FAIL: json parser failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: json parser %v", pr))
}
