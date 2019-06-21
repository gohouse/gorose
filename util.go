package gorose

import (
	"fmt"
	"github.com/gohouse/t"
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
	val := reflect.ValueOf(u).Elem()
	v := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		v[i] = valueField.Addr().Interface()
	}
	return v
}

func StructToMap(obj interface{}) map[string]interface{} {
	ty := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < ty.NumField(); i++ {
		data[ty.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// getTagName 获取结构体中Tag的值，如果没有tag则返回字段值
func getTagName(structName interface{}, tagstr string) []string {
	// 获取type
	tag := reflect.TypeOf(structName)
	// 如果是反射Ptr类型, 就获取他的 element type
	if tag.Kind() == reflect.Ptr {
		tag = tag.Elem()
	}

	// 判断是否是struct
	if tag.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	// 获取字段数量
	fieldNum := tag.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		//fieldName := t.Field(i).Name
		// tag 名字
		tagName := tag.Field(i).Tag.Get(tagstr)
		// tag为-时, 不解析
		if tagName == "-" || tagName == "" {
			// 字段名字
			tagName = tag.Field(i).Name
		}
		result = append(result, tagName)
	}
	return result
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

func addQuotes(data interface{}, sep string) string {
	ret := t.New(data).String()
	ret = strings.Replace(ret, `\`, `\\`, -1)
	ret = strings.Replace(ret, `"`, `\"`, -1)
	ret = strings.Replace(ret, `'`, `\'`, -1)
	return fmt.Sprintf("%s%s%s", sep, ret, sep)
}

// InArray :给定元素值 是否在 指定的数组中
func inArray(needle, hystack interface{}) bool {
	nt := t.New(needle)
	for _, item := range t.New(hystack).Slice() {
		if nt.String() == item.String() {
			return true
		}
	}
	return false
}

func withLockContext(fn func()) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	fn()
}
