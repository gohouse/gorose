package parser

import (
	"fmt"
	"github.com/gohouse/gorose/across"
	"testing"
)

func TestFileParser_Toml(test *testing.T) {
	var file = "/Users/fizz/go/src/github.com/gohouse/gorose/examples/demoParserFiles/mysql_cluster.toml"
	//var file = across.DemoParserFiles["toml"]

	var confP = &TomlConfigParser{}

	var v across.DbConfigCluster
	err := confP.Parse(file, &v)

	if err != nil {
		test.Error("FAIL: toml parser failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: toml parser %v", v.Master))
}
