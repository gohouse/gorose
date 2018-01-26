package utils

import (
	"reflect"
	"strconv"
	"strings"
	"encoding/json"
	"net/http"
)

func CheckErr(err error) error {
	if err != nil {
		return err
	}

	return nil
}

/**
 * 获取数据类型字符串 (string, int, float64, []int, []string, map[string]int ...)
 * (能不用则不用,由于涉及到使用reflect包,性能堪忧)
 */
func GetType(params interface{}) string {
	//数据初始化
	v := reflect.ValueOf(params)
	//获取传递参数类型
	v_t := v.Type()

	//类型名称对比
	return v_t.String()
}

/**
 * 给定元素值 是否在 指定的数组中
 */
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

func ParseStr(data interface{}) string {
	switch data.(type) {
	case int:
		return strconv.Itoa(data.(int))
	case int64:
		return strconv.FormatInt(data.(int64), 10)
	case string:
		return data.(string)
	default:
		return ""
	}
}

/**
 * 三元运算
 */
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

/**
 * 添加单引号
 */
func AddSingleQuotes(data interface{}) string {
	return "'" + strings.Trim(ParseStr(data), " ") + "'"
}

/**
 * 字符串转数组, 接受混合类型, 最终输出的是字符串类型
 */
func Implode(data interface{}, glue string) string {
	var tmp []string
	for _, item := range data.([]interface{}) {
		tmp = append(tmp, ParseStr(item))
	}

	return strings.Join(tmp, glue)
}

/**
 * json转码
 */
func JsonEncode(data interface{}) (string, error) {
	res, err := json.Marshal(data)

	if err != nil {
		return "",err
	}

	return string(res), err
}

// todo
//func JsonDecode(data string) interface{}{
//	type res struct{
//		arr []
//	}
//	err := json.Unmarshal(data, &res)
//	gorose.CheckErr(err)
//
//	return string(res)
//}
func UcFirst(arg string) string {
	if (len(arg) == 0) {
		return arg
	}
	return strings.ToUpper(arg[0:1])+arg[1:]
}
func Empty2(params interface{}) bool {
	//初始化变量
	var (
		flag          bool = true
		default_value reflect.Value
	)

	r := reflect.ValueOf(params)

	//获取对应类型默认值
	default_value = reflect.Zero(r.Type())
	//由于params 接口类型 所以default_value也要获取对应接口类型的值 如果获取不为接口类型 一直为返回false
	if !reflect.DeepEqual(r.Interface(), default_value.Interface()) {
		flag = false
	}
	return flag
}
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

func FailReturn(args ...interface{}) interface{} {
	var data []interface{}
	argsLength := len(args)
	if argsLength==0{
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
