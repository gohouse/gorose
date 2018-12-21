package parser

import (
	"fmt"
	"github.com/gohouse/gorose/across"
	"testing"
)

func TestFileParser_New(test *testing.T) {
	var v interface{}
	err := NewFileParser("json", across.DemoParserFiles["json"], &v)

	if err != nil {
		test.Error("FAIL: read file failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: json %v", v))
}
