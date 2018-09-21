package parser

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"reflect"
	"strings"
)

type TomlConfigParser struct {
}

func init() {
	// 检查解析器是否实现了接口
	var parserTmp IParser = &TomlConfigParser{}

	// 注册驱动
	Register("toml", parserTmp)
}

func (c *TomlConfigParser) Parse(file string, dbConfCluster interface{}) (err error) {
	var fp []byte
	fp, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}

	// 是否是主从格式
	strFp := string(fp)
	if strings.Contains(strFp, "Slave") &&
		strings.Contains(strFp, "Master") {
		err = toml.Unmarshal(fp, dbConfCluster)
	} else {
		//err = json.Unmarshal([]byte(fp), &conf.Master)
		err = tomlDecoder(fp, dbConfCluster)
	}

	return err
}

func tomlDecoder(str []byte, dbConfCluster interface{}) (err error) {
	srcElem := reflect.Indirect(reflect.ValueOf(dbConfCluster))
	fieldType := srcElem.FieldByName("Master").Type().Elem()
	fieldPtr := reflect.New(fieldType)
	//tmp := fieldPtr.Interface()
	err = toml.Unmarshal(str, fieldPtr.Interface())
	srcElem.FieldByName("Master").Set(fieldPtr)
	return
}
