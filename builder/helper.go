package builder

import (
	"errors"
	"github.com/gohouse/gorose/across"
	"github.com/gohouse/gorose/utils"
	"strings"
)

// parseJoin : parse the join paragraph
func parseJoin(ormApi across.OrmApi) (string, error) {
	var join []interface{}
	var returnJoinArr []string
	joinArr := ormApi.Sjoin

	for _, join = range joinArr {
		var w string
		var ok bool
		var args []interface{}

		if len(join) != 2 {
			return "", errors.New("join conditions are wrong")
		}

		// 获取真正的where条件
		if args, ok = join[1].([]interface{}); !ok {
			return "", errors.New("join conditions are wrong")
		}

		argsLength := len(args)
		switch argsLength {
		case 1:
			w = args[0].(string)
		case 2:
			w = args[0].(string) + " ON " + args[1].(string)
		case 4:
			w = args[0].(string) + " ON " + args[1].(string) + " " + args[2].(string) + " " + args[3].(string)
		default:
			return "", errors.New("join format error")
		}

		returnJoinArr = append(returnJoinArr, " "+join[0].(string)+" JOIN "+w)
	}

	return strings.Join(returnJoinArr, " "), nil
}


// parseWhere : parse where condition
func parseWhere(ormApi across.OrmApi) (string, error) {
	// 取出所有where
	wheres := ormApi.Swhere
	// where解析后存放每一项的容器
	var where []string

	for _, args := range wheres {
		// and或者or条件
		var condition string = args[0].(string)
		// 统计当前数组中有多少个参数
		params := args[1].([]interface{})
		paramsLength := len(params)

		switch paramsLength {
		case 3: // 常规3个参数:  {"id",">",1}
			res, err := parseParams(params)
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)

		case 2: // 常规2个参数:  {"id",1}
			res, err := parseParams(params)
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)
		case 1: // 二维数组或字符串
			switch paramReal := params[0].(type) {
			case string:
				where = append(where, condition+" ("+paramReal+")")
			case map[string]interface{}: // 一维数组
				var whereArr []string
				for key, val := range paramReal {
					whereArr = append(whereArr, key+"="+utils.AddSingleQuotes(val))
				}
				where = append(where, condition+" ("+strings.Join(whereArr, " and ")+")")
			case [][]interface{}: // 二维数组
				var whereMore []string
				for _, arr := range paramReal { // {{"a", 1}, {"id", ">", 1}}
					whereMoreLength := len(arr)
					switch whereMoreLength {
					case 3:
						res, err := parseParams(arr)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					case 2:
						res, err := parseParams(arr)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					default:
						return "", errors.New("where data format is wrong")
					}
				}
				where = append(where, condition+" ("+strings.Join(whereMore, " and ")+")")
			case func():
				// 清空where,给嵌套的where让路,复用这个节点
				ormApi.Swhere = [][]interface{}{}

				// 执行嵌套where放入Database struct
				paramReal()
				// 再解析一遍后来嵌套进去的where
				wherenested, err := parseWhere(ormApi)
				if err != nil {
					return "", err
				}
				// 嵌套的where放入一个括号内
				where = append(where, condition+" ("+wherenested+")")
			default:
				return "", errors.New("where data format is wrong")
			}
		}
	}

	// 还原初始where, 以便后边调用
	ormApi.Swhere = ormApi.SbeforeParseWhereData

	return strings.TrimLeft(
		strings.TrimLeft(strings.TrimLeft(
			strings.Trim(strings.Join(where, " "), " "),
			"and"), "or"),
		" "), nil
}


/**
 * 将where条件中的参数转换为where条件字符串
 * example: {"id",">",1}, {"age", 18}
 */
// parseParams : 将where条件中的参数转换为where条件字符串
func parseParams(args []interface{}) (string, error) {
	paramsLength := len(args)
	argsReal := args

	// 存储当前所有数据的数组
	var paramsToArr []string

	switch paramsLength {
	case 3: // 常规3个参数:  {"id",">",1}
		if !utils.InArray(argsReal[1], across.Regex) {
			return "", errors.New("where parameter is wrong")
		}

		paramsToArr = append(paramsToArr, argsReal[0].(string))
		paramsToArr = append(paramsToArr, argsReal[1].(string))

		switch argsReal[1] {
		case "like", "not like":
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(argsReal[2]))
			//case "not like":
			//	paramsToArr = append(paramsToArr, utils.AddSingleQuotes(utils.ParseStr(argsReal[2])))
		case "in", "not in":
			var tmp []string
			switch argsReal[2].(type) {
			case []string:
				for _, item := range argsReal[2].([]string) {
					tmp = append(tmp, utils.AddSingleQuotes(item))
				}
			case []int:
				for _, item := range argsReal[2].([]int) {
					tmp = append(tmp, utils.AddSingleQuotes(item))
				}
			case []interface{}:
				for _, item := range argsReal[2].([]interface{}) {
					tmp = append(tmp, utils.AddSingleQuotes(item))
				}
			}
			//for _, item := range argsReal[2].([]interface{}) {
			//	tmp = append(tmp, utils.AddSingleQuotes(item))
			//}
			paramsToArr = append(paramsToArr, "("+strings.Join(tmp, ",")+")")
			//case "not in":
			//	paramsToArr = append(paramsToArr, "("+utils.Implode(argsReal[2], ",")+")")
		case "between", "not between":
			var tmpB []interface{}
			switch argsReal[2].(type) {
			case []string:
				tmp := argsReal[2].([]string)
				tmpB = append(tmpB, tmp[0])
				tmpB = append(tmpB, tmp[1])
			case []int:
				tmp := argsReal[2].([]int)
				tmpB = append(tmpB, tmp[0])
				tmpB = append(tmpB, tmp[1])
			case []interface{}:
				tmp := argsReal[2].([]interface{})
				tmpB = append(tmpB, tmp[0])
				tmpB = append(tmpB, tmp[1])
			}
			//tmpB := argsReal[2].([]string)
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(tmpB[0])+" and "+utils.AddSingleQuotes(tmpB[1]))
			//case "not between":
			//	tmpB := argsReal[2].([]string)
			//	paramsToArr = append(paramsToArr, utils.AddSingleQuotes(tmpB[0])+" and "+utils.AddSingleQuotes(tmpB[1]))
		default:
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(argsReal[2]))
		}
	case 2:
		//if !utils.TypeCheck(args[0], "string") {
		//	panic("where条件参数有误!")
		//}
		//fmt.Println(argsReal)
		paramsToArr = append(paramsToArr, argsReal[0].(string))
		paramsToArr = append(paramsToArr, "=")
		paramsToArr = append(paramsToArr, utils.AddSingleQuotes(argsReal[1]))
	}
	return strings.Join(paramsToArr, " "), nil
}


// buildData : build inert or update data
func buildData(ormApi across.OrmApi) (string, string, string) {
	// insert
	var dataFields []string
	var dataValues []string
	// update or delete
	var dataObj []string

	data := ormApi.Sdata

	switch data.(type) {
	case string:
		dataObj = append(dataObj, data.(string))
	case []map[string]interface{}: // insert multi datas ([]map[string]interface{})
		datas := data.([]map[string]interface{})
		for key, _ := range datas[0] {
			if utils.InArray(key, dataFields) == false {
				dataFields = append(dataFields, key)
			}
		}
		for _, item := range datas {
			var dataValuesSub []string
			for _, key := range dataFields {
				if item[key] == nil {
					dataValuesSub = append(dataValuesSub, "null")
				} else {
					dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(item[key]))
				}
			}
			dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
		}
	default: // update or insert
		var dataValuesSub []string
		for key, val := range data.(map[string]interface{}) {
			// insert
			dataFields = append(dataFields, key)
			//dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(val))
			if val == nil {
				dataValuesSub = append(dataValuesSub, "null")
			} else {
				dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(val))
			}
			// update
			//dataObj = append(dataObj, key+"="+utils.AddSingleQuotes(val))
			if val == nil {
				dataObj = append(dataObj, key+"=null")
			} else {
				dataObj = append(dataObj, key+"="+utils.AddSingleQuotes(val))
			}
		}
		// insert
		dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
	}

	return strings.Join(dataObj, ","), strings.Join(dataFields, ","), strings.Join(dataValues, ",")
}
