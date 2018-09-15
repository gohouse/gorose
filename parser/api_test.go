package parser

import (
	"fmt"
	"fizzday.com/gohouse/gorose/config"
	"testing"
)

func TestFileParser_New(test *testing.T) {
	pr, err := NewFileParser("json",config.DemoParserFiles["json"])

	if err != nil {
		test.Error("FAIL: read file failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: json %v", pr))
}
