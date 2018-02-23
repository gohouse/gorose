package utils

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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
	case int:
		return strconv.Itoa(data.(int))
	case int64:
		return strconv.FormatInt(data.(int64), 10)
	case float32:
		return strconv.FormatFloat(float64(data.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(data.(float64), 'f', -1, 64)
	case string:
		return data.(string)
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
	return "'" + strings.Replace(ParseStr(data), "'", `\'`, -1) + "'"
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

// SuccessReturn : 接口成功返回
// args: 传入的值,可接收1~3个值,第一个值是返回的数据,第二个值是状态码,第三个值是附加额外数据.
// 		这里第二个值默认缺省为200(成功), 第三个值默认缺省为空
// example: SuccessReturn([]map[string]interface{{"id":1,"name":"fizz"},{"id":2,"name":"fizz2"}}, 200, map[string]int{"page":1,"total":93,"limit":10})
// return: {"data":[{"id":1,"name":"fizz"},{"id":2,"name":"fizz2"}], "status":200, "ext":{"page":1,"total":93,"limit":10}}
func SuccessReturn(args ...interface{}) interface{} {
	argsLength := len(args)

	//var w http.ResponseWriter
	var data = make(map[string]interface{})

	data["msg"] = "success"

	if argsLength > 0 {
		data["data"] = args[0]
	} else {
		data["data"] = ""
	}

	switch argsLength {
	case 0:
		data["status"] = http.StatusOK
	case 1:
		//w.WriteHeader(http.StatusOK)
		// 正确的返回数据
		data["status"] = http.StatusOK
	case 2:
		switch args[1].(type) {
		case int:
			//w.WriteHeader(args[1].(int))
			data["status"] = args[1].(int)
		case string:
			//w.WriteHeader(http.StatusOK)
			data["status"] = http.StatusOK
		default:
			panic("调用返回的状态值应该为int类型")
		}
	case 3:
		switch args[1].(type) {
		case int:
			//w.WriteHeader(args[1].(int))
			data["status"] = args[1].(int)
		case string:
			//w.WriteHeader(http.StatusOK)
			data["status"] = http.StatusOK
		default:
			panic("调用返回的状态值应该为int类型")
		}
		data["ext"] = args[2]
	default:
		panic("调用返回的参数有wu")
	}

	return data
}

// FailReturn : 接口失败返回
func FailReturn(args ...interface{}) interface{} {
	var data []interface{}
	argsLength := len(args)
	if argsLength == 0 {
		data = append(data, "fail")
	} else {
		data = append(data, args[0])
	}
	switch argsLength {
	case 0:
		data = append(data, http.StatusNoContent)
	case 1:
		data = append(data, http.StatusNoContent)
	case 2:
		switch args[1].(type) {
		case int:
			data = append(data, args[1].(int))
		case string:
			data = append(data, http.StatusNoContent)
		default:
			panic("调用返回的状态值应该为int类型")
		}
	case 3:
		switch args[1].(type) {
		case int:
			data = append(data, args[1].(int))
		case string:
			data = append(data, http.StatusNoContent)
		default:
			panic("调用返回的状态值应该为int类型")
		}
	default:
		panic("调用返回的参数有误")
	}

	return SuccessReturn(data...)
}
