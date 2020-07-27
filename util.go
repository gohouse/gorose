package gorose

import (
	"fmt"
	"github.com/gohouse/t"
	"log"
	"math/rand"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

func getRandomInt(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}

func structForScan(u interface{}) []interface{} {
	val := reflect.Indirect(reflect.ValueOf(u))
	v := make([]interface{}, 0)
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		if val.Type().Field(i).Tag.Get(TAGNAME) != IGNORE {
			if valueField.CanAddr() {
				v = append(v, valueField.Addr().Interface())
			} else {
				//v[i] = valueField
				v = append(v, valueField)
			}
		}
	}
	return v
}

// StructToMap ...
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
		// tag 名字
		tagName := tag.Field(i).Tag.Get(tagstr)
		if tagName != IGNORE {
			// tag为-时, 不解析
			if tagName == "-" || tagName == "" {
				// 字段名字
				tagName = tag.Field(i).Name
			}
			result = append(result, tagName)
		}
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

func addBackticks(arg string) string {
	reg := regexp.MustCompile(`^\w+$`)
	if reg.MatchString(arg) {
		return fmt.Sprintf("`%s`", arg)
	}
	return arg
}

// InArray :给定元素值 是否在 指定的数组中
func inArray(needle, hystack interface{}) bool {
	nt := t.New(needle)
	for _, item := range t.New(hystack).Slice() {
		if strings.ToLower(nt.String()) == strings.ToLower(item.String()) {
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

func withRunTimeContext(closer func(), callback func(time.Duration)) {
	// 记录开始时间
	start := time.Now()
	closer()
	timeduration := time.Since(start)
	//log.Println("执行完毕,用时:", timeduration.Seconds(),timeduration.Seconds()>1.1)
	callback(timeduration)
}

func readFile(filepath string) *os.File {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path.Dir(filepath), os.ModePerm)
		file, _ = os.Create(filepath)
	}
	return file
}
