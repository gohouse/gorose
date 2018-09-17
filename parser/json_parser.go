package parser

import (
	"encoding/json"
	"github.com/gohouse/gorose/across"
	"io/ioutil"
	"strings"
)

type JsonConfigParser struct {
}

func init() {
	// 检查解析器是否实现了接口
	var parser IParser = &JsonConfigParser{}

	// 注册驱动
	Register(across.JSON, parser)
}

func (c *JsonConfigParser) Parse(file string) (*across.DbConfigCluster, error) {
	var conf = &across.DbConfigCluster{}
	var err error
	var fp []byte
	fp, err = ioutil.ReadFile(file)
	if err != nil {
		return conf,err
	}
	// 是否是主从格式
	strFp := string(fp)
	if strings.Contains(strFp, "Slave") &&
		strings.Contains(strFp, "Master") {
		err = json.Unmarshal([]byte(fp), &conf)
	} else {
		err = json.Unmarshal([]byte(fp), &conf.Master)
	}
	//fmt.Println(conf.Master)

	return conf,err
}
