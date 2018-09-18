package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"fmt"
	"regexp"
	"math/rand"
	"errors"
	"net/url"
	"os/exec"
	"bytes"
)

// GetType : 获取数据类型字符串 (string, int, float64, []int, []string, map[string]int ...)
// GetType : (能不用则不用,由于涉及到使用reflect包,性能堪忧)
func GetType(params interface{}) string {
	//数据初始化
	v := reflect.ValueOf(params)
	//获取传递参数类型
	vT := v.Type()

	//类型名称对比
	return vT.String()
}

// InArray :给定元素值 是否在 指定的数组中
func InArray(needle interface{}, hystack interface{}) bool {
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
		//panic("needle only support string,int,int64 type")
	}
	//for _, item := range hystack {
	//	if GetType(needle) == GetType(item) {
	//		if needle == item {
	//			return true
	//		}
	//	}
	//}

	return false
}

// ParseStr 转换为string
func ParseStr(data interface{}) string {
	switch data.(type) {
	case time.Time:
		return data.(time.Time).Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprint(data)
	}
}
func ParseInt(data interface{}) int{
	dataType := GetType(data)
	var res int
	switch dataType {
	case "string":
		res,_ = strconv.Atoi(data.(string))
	case "int":
		res = data.(int)
	default:
		panic("只能转换字符串类型")
	}

	return res
}

// ParseStr 转换为string
func ParseStr_bak(data interface{}) string {
	switch data.(type) {
	case int:
		return strconv.Itoa(data.(int))
	case int64:
		return strconv.FormatInt(data.(int64), 10)
	case int32:
		return strconv.FormatInt(int64(data.(int32)), 10)
	case uint32:
		return strconv.FormatUint(uint64(data.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(data.(uint64), 10)
	case float32:
		return strconv.FormatFloat(float64(data.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(data.(float64), 'f', -1, 64)
	case string:
		return data.(string)
	case time.Time:
		return data.(time.Time).Format("2006-01-02 15:04:05")
	default:
		return ""
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

// AddSingleQuotes : 添加单引号
func AddSingleQuotes(data interface{}) string {
	//return "'" + strings.Trim(ParseStr(data), " ") + "'"
	switch data.(type) {
	case int,int64,int32,uint32,uint64:
		return ParseStr(data)
	default:
		return "'" + strings.Replace(ParseStr(data), "'", `\'`, -1) + "'"
	}
}

// Implode : 字符串转数组, 接受混合类型, 最终输出的是字符串类型
func Implode(data interface{}, glue string) string {
	var tmp []string
	for _, item := range data.([]interface{}) {
		tmp = append(tmp, ParseStr(item))
	}

	return strings.Join(tmp, glue)
}

// JsonEncode : json转码
func JsonEncode(data interface{}) (string, error) {
	res, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	return string(res), err
}

// UcFirst : 字符串第一个字母转成大写
func UcFirst(arg string) string {
	if len(arg) == 0 {
		return arg
	}
	return strings.ToUpper(arg[0:1]) + arg[1:]
}

// Empty : 是否位假
func Empty(arg interface{}) bool {
	switch arg.(type) {
	case int:
		return If(arg.(int) == 0, true, false).(bool)
	case int64:
		return If(arg.(int) == 0, true, false).(bool)
	case string:
		return If(arg.(string) == "", true, false).(bool)
	default:
		return true
	}
}

type ApiReturn struct {
	Data interface{}
	Code int
	Msg  interface{}
	Ext  interface{}
}

// SuccessReturn : 接口成功返回
// args: 传入的值,可接收1~3个值,第一个值是返回的数据,第二个值是状态码(默认200),第三个值是附加额外数据.
// 		这里第二个值默认缺省为200(成功), 第三个值默认缺省为空
// example: SuccessReturn([]map[string]interface{{"id":1,"name":"fizz"},{"id":2,"name":"fizz2"}}, 200, map[string]int{"page":1,"total":93,"limit":10})
// return: {"data":[{"id":1,"name":"fizz"},{"id":2,"name":"fizz2"}], "status":200, "ext":{"page":1,"total":93,"limit":10}}
func SuccessReturn(args ...interface{}) ApiReturn {
	data := ApiReturn{
		Msg: "success",
	}

	argsLength := len(args)
	switch argsLength {
	case 0:
		data.Data = "success"
		data.Code = http.StatusOK
	case 1:
		// 正确的返回数据
		data.Data = args[0]
		data.Code = http.StatusOK
	case 2:
		switch args[1].(type) {
		case int:
			data.Data = args[0]
			data.Code = args[1].(int)
		case string:
			data.Data = args[0]
			code, _ := strconv.Atoi(args[1].(string))
			data.Code = If(code > 0, code, http.StatusOK).(int)
		default:
			//panic("调用返回的状态值应该为int类型")
			return FailReturn("SuccessReturn 调用返回的状态值应该为int类型")
		}
	case 3:
		switch args[1].(type) {
		case int:
			data.Data = args[0]
			data.Code = args[1].(int)
		case string:
			data.Data = args[0]
			code, _ := strconv.Atoi(args[1].(string))
			data.Code = If(code > 0, code, http.StatusOK).(int)
		default:
			//panic("调用返回的状态值应该为int类型")
			return FailReturn("SuccessReturn 调用返回的状态值应该为int类型")
		}
		data.Ext = args[2]
	default:
		//panic("调用返回的参数有误")
		return FailReturn("SuccessReturn 调用返回的参数有误")
	}

	return data
}

// FailReturn : 接口失败返回
// 可接收1~3个值,第一个值是返回的数据,第二个值是状态码(默认204),第三个值是附加额外数据.
func FailReturn(args ...interface{}) ApiReturn {
	data := ApiReturn{
		Msg: "fail",
	}

	argsLength := len(args)
	switch argsLength {
	case 0:
		data.Data = "fail"
		data.Code = http.StatusNoContent
	case 1:
		// 正确的返回数据
		data.Data = args[0]
		data.Msg = args[0]
		data.Code = http.StatusNoContent
	case 2:
		data.Msg = args[0]
		data.Data = args[0].(string)
		switch args[1].(type) {
		case int:
			data.Code = args[1].(int)
		case string:
			code, _ := strconv.Atoi(args[1].(string))
			data.Code = If(code > 0, code, http.StatusNoContent).(int)
		default:
			//panic("调用返回的状态值应该为int类型")
			return FailReturn("FailReturn 调用返回的状态值应该为int类型");
		}
	case 3:
		data.Msg = args[0]
		data.Data = args[0]
		switch args[1].(type) {
		case int:
			data.Code = args[1].(int)
		case string:
			code, _ := strconv.Atoi(args[1].(string))
			data.Code = If(code > 0, code, http.StatusNoContent).(int)
		default:
			//panic("调用返回的状态值应该为int类型")
			return FailReturn("FailReturn 调用返回的状态值应该为int类型");
		}
		data.Ext = args[2]
	default:
		//panic("调用返回的参数有误")
		return FailReturn("FailReturn 调用返回的参数有误");
	}

	return data
}
func ArrayReverse(arr []map[string]interface{}) ([]map[string]interface{}, error) {
	lenArr := len(arr)
	if lenArr == 0 {
		return arr, nil
	}

	var newArr []map[string]interface{}

	for i := lenArr - 1; i >= 0; i-- {
		newArr = append(newArr, arr[i])
	}

	return newArr, nil
}
func Ip2long(ipstr string) (ip uint32) {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return
	}
	ips := reg.FindStringSubmatch(ipstr)
	if ips == nil {
		return
	}

	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return
	}

	ip += uint32(ip1 * 0x1000000)
	ip += uint32(ip2 * 0x10000)
	ip += uint32(ip3 * 0x100)
	ip += uint32(ip4)

	return
}
func Long2ip(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip>>24, ip<<8>>24, ip<<16>>24, ip<<24>>24)
}
func GetIp() string {
	ipArr := [15][2]int{
		{607649792, 608174079},
		{975044608, 977272831},
		{999751680, 999784447},
		{1019346944, 1019478015},
		{1038614528, 1039007743},
		{1783627776, 1784676351},
		{1947009024, 1947074559},
		{1987051520, 1988034559},
		{2035023872, 2035154943},
		{2078801920, 2079064063},
		{-1950089216, -1948778497},
		{-1425539072, -1425014785},
		{-1236271104, -1235419137},
		{-770113536, -768606209},
		{-569376768, -564133889},
	}
	randKey := rand.Intn(14)
	ip := Long2ip(uint32(rand.Intn(ipArr[randKey][1]-ipArr[randKey][0]) + ipArr[randKey][0]));
	return ip
}

func UrlQueryStrToMap(urlstr string) (map[string]interface{}, error) {
	formData := make(map[string]interface{})
	if len(urlstr) < 5 {
		return formData, errors.New("url有误")
	}
	u, err := url.Parse(urlstr)
	if err != nil {
		return formData, err
	}

	// 组装map
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return formData, err
	}

	if len(m) > 0 {
		for k, v := range m {
			if len(v) > 0 {
				formData[k] = v[0]
			} else {
				formData[k] = ""
			}
		}
	}

	return formData, nil
}

func ArrayKeys(arr map[string]interface{}) []string {
	var tmp []string
	if len(arr)>0 {
		for k,_ := range arr {
			tmp = append(tmp,k)
		}
	}
	return tmp
}

func ArrayValues(arr map[string]interface{}) []interface{} {
	var tmp []interface{}
	if len(arr)>0 {
		for _,v := range arr {
			tmp = append(tmp,v)
		}
	}
	return tmp
}

func StartWith(originStr string, sepStr string) bool {
	if originStr!="" && sepStr!="" {
		length := len(sepStr)
		if strings.Trim(originStr, " ")[:length]==sepStr {
			return true
		}
	}

	return false
}

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func ExecShell(s string) (string, error){
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()

	return out.String(), err
}


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
