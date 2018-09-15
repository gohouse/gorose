package parser

import (
	"encoding/json"
	"fizzday.com/gohouse/gorose/config"
	"io/ioutil"
)

type JsonConfigParser struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var parser IParser = &JsonConfigParser{}

	// 注册驱动
	Register(config.JSON, parser)
}

func (c *JsonConfigParser) Parse(file string) (conf *config.DbConfig, err error) {
	var fp []byte
	fp, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(fp, &conf)
	return
}
