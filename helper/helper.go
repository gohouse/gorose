package helper

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strings"
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


//将结构体里的成员按照json名字来赋值
func SetStructFieldByJsonName(ptr interface{}, fields map[string]interface{}) {
	v := reflect.ValueOf(ptr).Elem() // the struct variable

	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("json")

		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}

		if value, ok := fields[name]; ok {
			//给结构体赋值
			//保证赋值时数据类型一致
			//fmt.Println("类型1：", reflect.ValueOf(value).Type(), "类型2：", v.FieldByName(fieldInfo.Name).Type())
			if reflect.ValueOf(value).Type() == v.FieldByName(fieldInfo.Name).Type() {
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value))
			}

		}
	}

	return
}

//func BuildSql(dba *gorose.Database, operType ...string) (string, error) {
//	return gorose.BuildSql(dba, operType...)
//}
