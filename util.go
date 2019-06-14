package gorose

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"
)

func getRandomInt(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}

func strutForScan(u interface{}) []interface{} {
	//fmt.Printf("%#v\n",u)
	val := reflect.ValueOf(u).Elem()
	v := make([]interface{}, val.NumField())
	//fmt.Printf("%#v\n",val.Elem())
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		v[i] = valueField.Addr().Interface()
	}
	//fmt.Printf("%#v",v)
	return v
}

func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// getTagName 获取结构体中Tag的值，如果没有tag则返回字段值
func getTagName(structName interface{}) []string {
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
		if tagName == "-" || tagName == "" {
			// 字段名字
			tagName = t.Field(i).Name
		}
		result = append(result, tagName)
	}
	return result
}

// ParseStr 转换为string
func parseStr(data interface{}) string {
	switch data.(type) {
	case time.Time:
		return data.(time.Time).Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprint(data)
	}
}

// If : ternary operator (三元运算)
// condition:比较运算
// trueVal:运算结果为真时的值
// falseVal:运算结果为假时的值
// return: 由于不知道传入值的类型, 所有, 必须在接收结果时, 指定对应的值类型
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// addSingleQuotes : 添加单引号
func addSingleQuotes(data interface{}) string {
	switch data.(type) {
	case int, int64, int32, uint32, uint64:
		return parseStr(data)
	default:
		ret := parseStr(data)
		ret = strings.Replace(ret, `\`, `\\`, -1)
		ret = strings.Replace(ret, `"`, `\"`, -1)
		ret = strings.Replace(ret, `'`, `\'`, -1)
		return "'" + ret + "'"
	}
}

func addQuotes(data interface{}, sep string) string {
	switch data.(type) {
	case int, int64, int32, uint32, uint64:
		return parseStr(data)
	default:
		ret := parseStr(data)
		ret = strings.Replace(ret, `\`, `\\`, -1)
		ret = strings.Replace(ret, `"`, `\"`, -1)
		ret = strings.Replace(ret, `'`, `\'`, -1)
		return fmt.Sprintf("%s%s%s",sep,ret,sep)
	}
}

// InArray :给定元素值 是否在 指定的数组中
func inArray(needle interface{}, hystack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range hystack.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range hystack.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range hystack.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}

	return false
}

func withLockContext(fn func())  {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	fn()
}
