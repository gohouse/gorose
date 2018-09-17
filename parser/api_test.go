package parser

import (
	"fmt"
	"github.com/gohouse/gorose/examples"
	"testing"
)

func TestFileParser_New(test *testing.T) {
	pr, err := NewFileParser("json", examples.DemoParserFiles["json"])

	if err != nil {
		test.Error("FAIL: read file failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: json %v", pr))
}
