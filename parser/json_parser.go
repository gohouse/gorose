package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

type JsonConfigParser struct {
}

func init() {
	// 检查解析器是否实现了接口
	var parserTmp IParser = &JsonConfigParser{}

	// 注册驱动
	Register("json", parserTmp)
}

func (c *JsonConfigParser) Parse(file string, dbConfCluster interface{}) (err error) {
	var fp []byte
	fp, err = ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	// 是否是主从格式
	strFp := string(fp)
	if strings.Contains(strFp, "Slave") &&
		strings.Contains(strFp, "Master") {
		err = json.Unmarshal(fp, dbConfCluster)
	} else {
		//err = json.Unmarshal([]byte(fp), &conf.Master)
		err = jsonDecoder(fp, dbConfCluster)
	}

	return err
}

func jsonDecoder(str []byte, dbConfCluster interface{}) (err error) {
	srcElem := reflect.Indirect(reflect.ValueOf(dbConfCluster))
	fmt.Println(srcElem)
	fieldType := srcElem.FieldByName("Master").Type().Elem()
	fieldElem := reflect.New(fieldType)
	err = json.Unmarshal(str, fieldElem.Interface())
	srcElem.FieldByName("Master").Set(fieldElem)
	return
}
