package helper

import (
	"errors"
	"log"
	"os"
	"reflect"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// GetStructName 获取结构体的名字
func GetStructName(structName interface{}) (string,error) {
	// 获取type
	t := reflect.TypeOf(structName)
	// 如果是反射Ptr类型, 就获取他的 element type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 判断是否是struct
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return "",errors.New("非struct")
	}
	return t.Name(),nil
}

// GetTagName 获取结构体中Tag的值，如果没有tag则返回字段值
func GetTagName(structName interface{}) []string {
	// 获取type
	t := reflect.TypeOf(structName)
	// 如果是反射Ptr类型, 就获取他的 element type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 判断是否是struct
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	// 获取字段数量
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		//fieldName := t.Field(i).Name
		// tag 名字
		tagName := t.Field(i).Tag.Get("orm")
		// tag为-时, 不解析
		if tagName=="-" || tagName=="" {
			// 字段名字
			tagName = t.Field(i).Name
		}
		result = append(result, tagName)
	}
	return result
}

func StrutForScan(u interface{}) []interface{} {
	val := reflect.ValueOf(u).Elem()
	v := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		v[i] = valueField.Addr().Interface()
	}
	return v
}
