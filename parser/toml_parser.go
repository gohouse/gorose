package parser

import (
	"github.com/BurntSushi/toml"
	"github.com/gohouse/gorose/across"
	"io/ioutil"
)

type TomlConfigParser struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var parser IParser = &TomlConfigParser{}

	// 注册驱动
	Register(across.TOML, parser)
}

func (c *TomlConfigParser) Parse(file string) (conf *across.DbConfigCluster, err error) {
	var fp []byte
	fp, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = toml.Unmarshal([]byte(fp), &conf)
	return
}
