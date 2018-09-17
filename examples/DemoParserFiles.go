package examples

import "github.com/gohouse/gorose/across"

// 临时配置, 方便 test 测试
var DemoParserFiles = map[string]string{
	across.JSON: "/Users/fizz/go/src/github.com/gohouse/gorose/examples/demoParserFiles/mysql_cluster.json",
	//across.JSON: "/Users/fizz/go/src/github.com/gohouse/gorose/examples/demoParserFiles/mysql.json",
	across.TOML: "/Users/fizz/go/src/github.com/gohouse/gorose/examples/demoParserFiles/mysql.toml",
	across.INI:  "/Users/fizz/go/src/github.com/gohouse/gorose/examples/demoParserFiles/mysql.ini",
	//across.JSON: "E:\\go\\src\\github.com\\gohouse\\laboratory\\gorose\\config\\demoParserFiles\\mysql.json",
	//across.TOML: "E:\\go\\src\\github.com\\gohouse\\laboratory\\gorose\\config\\demoParserFiles\\mysql.toml",
}

type TestT struct {
	Name string
}
