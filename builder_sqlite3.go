package gorose

import (
	"errors"
	"fmt"
	"github.com/gohouse/t"
	"strconv"
	"strings"
)

type BuilderSqlite3 struct {
	*Builder
	IOrm
}

// sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
//		distinct, fields, table, join, where, group, having, order, limit, offset)
// select {distinct} {fields} from {table} {join} {where} {group} {having} {order} {limit} {offset}
func init() {
	NewDriver().Register("sqlite3", &BuilderSqlite3{})
}

func (b *BuilderSqlite3) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	b.IOrm = o
	join, err := b.BuildJoin()
	if err != nil {
		return
	}
	where, err := b.BuildWhere()
	if err != nil {
		return
	}
	sqlStr = fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
		b.BuildDistinct(), b.BuildFields(), b.BuildTable(), join, where,
		b.BuildGroup(), b.BuildHaving(), b.BuildOrder(), b.BuildLimit(), b.BuildOffset())

	//args = b.bindParams
	args = b.IOrm.GetBindValues()
	return
}

// BuildExecut : build execute query string
func (b *BuilderSqlite3) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	// insert : {"name":"fizz, "website":"fizzday.net"} or {{"name":"fizz2", "website":"www.fizzday.net"}, {"name":"fizz", "website":"fizzday.net"}}}
	// update : {"name":"fizz", "website":"fizzday.net"}
	// delete : ...
	b.IOrm = o
	var update, insertkey, insertval string
	if operType != "delete" {
		if b.IOrm.GetData() == nil {
			err = errors.New("insert,update请传入数据操作")
			return
		}
		update, insertkey, insertval = b.BuildData(operType)
	}

	where, err := b.BuildWhere()
	if err != nil {
		return
	}

	switch operType {
	case "insert":
		sqlStr = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", b.BuildTable(), insertkey, insertval)
	case "update":
		if where == "" && o.GetForce() == false {
			err = errors.New("出于安全考虑, update时where条件不能为空, 如果真的不需要where条件, 请使用force(如: db.xxx.Force().Update())")
			return
		}
		sqlStr = fmt.Sprintf("UPDATE %s SET %s%s", b.BuildTable(), update, where)
	case "delete":
		if where == "" && o.GetForce() == false {
			err = errors.New("出于安全考虑, delete时where条件不能为空, 如果真的不需要where条件, 请使用force(如: db.xxx.Force().Delete())")
			return
		}
		sqlStr = fmt.Sprintf("DELETE FROM %s%s", b.BuildTable(), where)
	}

	args = b.IOrm.GetBindValues()
	return
}

// buildData : build inert or update data
func (b *BuilderSqlite3) BuildData(operType string) (string, string, string) {
	// insert
	var dataFields []string
	var dataValues []string
	// update or delete
	var dataObj []string

	data := b.IOrm.GetData()

	switch data.(type) {
	case string:
		dataObj = append(dataObj, data.(string))
	case []map[string]interface{}, []Data: // insert multi datas ([]map[string]interface{})
		sliceData := t.New(data).Slice()
		for key, _ := range sliceData[0].MapString() {
			if inArray(key, dataFields) == false {
				dataFields = append(dataFields, key)
			}
		}
		for _, itemT := range sliceData {
			item := itemT.MapString()
			var dataValuesSub []string
			for _, key := range dataFields {
				if item[key] == nil {
					dataValuesSub = append(dataValuesSub, "null")
				} else {
					dataValuesSub = append(dataValuesSub, "?")
					b.IOrm.SetBindValues(item[key])
				}
			}
			dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
		}
	case map[string]interface{}, Data: // insert multi datas ([]map[string]interface{})
		var dataValuesSub []string
		for key, val := range t.New(data).MapString() {
			if operType == "insert" {
				// insert
				dataFields = append(dataFields, key)
				if val.Interface() == nil {
					dataValuesSub = append(dataValuesSub, "null")
				} else {
					dataValuesSub = append(dataValuesSub, "?")
					b.IOrm.SetBindValues(val.Interface())
				}
			} else if operType == "update" {
				// update
				if val.Interface() == nil {
					dataObj = append(dataObj, key+"=null")
				} else {
					dataObj = append(dataObj, key+"=?")
					b.IOrm.SetBindValues(val.Interface())
				}
			}
		}
		if operType == "insert" {
			// insert
			dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
		}
	default: // update or insert
		return "", "", ""
	}

	return strings.Join(dataObj, ","), strings.Join(dataFields, ","), strings.Join(dataValues, ",")
}

func (b *BuilderSqlite3) BuildJoin() (string, error) {
	// 用户传入的join参数+join类型
	var join []interface{}
	var returnJoinArr []string
	joinArr := b.GetJoin()

	for _, join = range joinArr {
		var w string
		var ok bool
		// 用户传入 join 的where值, 即第二个参数
		var args []interface{}

		if len(join) != 2 {
			return "", errors.New("join conditions are wrong")
		}

		// 获取真正的用户传入的join参数
		if args, ok = join[1].([]interface{}); !ok {
			return "", errors.New("join conditions are wrong")
		}

		argsLength := len(args)
		switch argsLength {
		case 1: // join字符串 raw
			w = args[0].(string)
		case 2: // join表 + 字符串
			w = args[0].(string) + " ON " + args[1].(string)
		case 4: // join表 + (a字段+关系+a字段)
			w = args[0].(string) + " ON " + args[1].(string) + " " + args[2].(string) + " " + args[3].(string)
		default:
			return "", errors.New("join format error")
		}

		returnJoinArr = append(returnJoinArr, " "+join[0].(string)+" JOIN "+w)
	}

	return strings.Join(returnJoinArr, " "), nil
}

func (b *BuilderSqlite3) BuildWhere() (where string, err error) {
	var beforeParseWhere = b.IOrm.GetWhere()
	where, err = b.parseWhere(b.IOrm)
	b.IOrm.SetWhere(beforeParseWhere)
	return If(where == "", "", " WHERE "+where).(string), err
}

func (b *BuilderSqlite3) BuildDistinct() (dis string) {
	return If(b.IOrm.GetDistinct(), "DISTINCT ", "").(string)
}

func (b *BuilderSqlite3) BuildFields() string {
	if len(b.IOrm.GetFields()) == 0 {
		return "*"
	}
	return strings.Join(b.IOrm.GetFields(), ",")
}

func (b *BuilderSqlite3) BuildTable() string {
	return b.IOrm.GetTable()
}

func (b *BuilderSqlite3) BuildGroup() string {
	return If(b.IOrm.GetGroup() == "", "", " GROUP BY "+b.IOrm.GetGroup()).(string)
}

func (b *BuilderSqlite3) BuildHaving() string {
	return If(b.IOrm.GetHaving() == "", "", " HAVING "+b.IOrm.GetHaving()).(string)
}

func (b *BuilderSqlite3) BuildOrder() string {
	return If(b.IOrm.GetOrder() == "", "", " ORDER BY "+b.IOrm.GetOrder()).(string)
}

func (b *BuilderSqlite3) BuildLimit() string {
	return If(b.IOrm.GetLimit() == 0, "", " LIMIT "+strconv.Itoa(b.IOrm.GetLimit())).(string)
}

func (b *BuilderSqlite3) BuildOffset() string {
	if b.BuildLimit()==""{
		return ""
	}
	return If(b.IOrm.GetOffset() == 0, "", " OFFSET "+strconv.Itoa(b.IOrm.GetOffset())).(string)
}

// parseWhere : parse where condition
func (b *BuilderSqlite3) parseWhere(ormApi IOrm) (string, error) {
	// 取出所有where
	wheres := ormApi.GetWhere()
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
			res, err := b.parseParams(params, ormApi)
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)

		case 2: // 常规2个参数:  {"id",1}
			res, err := b.parseParams(params, ormApi)
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
					//whereArr = append(whereArr, key+"="+addSingleQuotes(val))
					whereArr = append(whereArr, key+"=?")
					b.IOrm.SetBindValues(val)
				}
				where = append(where, condition+" ("+strings.Join(whereArr, " and ")+")")
			case [][]interface{}: // 二维数组
				var whereMore []string
				for _, arr := range paramReal { // {{"a", 1}, {"id", ">", 1}}
					whereMoreLength := len(arr)
					switch whereMoreLength {
					case 3:
						res, err := b.parseParams(arr, ormApi)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					case 2:
						res, err := b.parseParams(arr, ormApi)
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
				ormApi.SetWhere([][]interface{}{})

				// 执行嵌套where放入Database struct
				paramReal()
				// 再解析一遍后来嵌套进去的where
				wherenested, err := b.parseWhere(ormApi)
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
func (b *BuilderSqlite3) parseParams(args []interface{}, ormApi IOrm) (string, error) {
	paramsLength := len(args)
	argsReal := args

	// 存储当前所有数据的数组
	var paramsToArr []string

	switch paramsLength {
	case 3: // 常规3个参数:  {"id",">",1}
		if !inArray(argsReal[1], b.GetRegex()) {
			return "", errors.New("where parameter is wrong")
		}

		paramsToArr = append(paramsToArr, argsReal[0].(string))
		paramsToArr = append(paramsToArr, argsReal[1].(string))

		switch argsReal[1] {
		case "like", "not like":
			paramsToArr = append(paramsToArr, "?")
			b.IOrm.SetBindValues(argsReal[2])
		case "in", "not in":
			var tmp []string
			var ar2 = t.New(argsReal[2]).MapString()
			for _, item := range ar2 {
				tmp = append(tmp, "?")
				b.IOrm.SetBindValues(t.New(item).Interface())
			}
			paramsToArr = append(paramsToArr, "("+strings.Join(tmp, ",")+")")
		case "between", "not between":
			var ar2 = t.New(argsReal[2]).Slice()
			paramsToArr = append(paramsToArr, "? and ?")
			b.IOrm.SetBindValues(ar2[0].Interface())
			b.IOrm.SetBindValues(ar2[1].Interface())
		default:
			paramsToArr = append(paramsToArr, "?")
			b.IOrm.SetBindValues(argsReal[2])
		}
	case 2:
		paramsToArr = append(paramsToArr, argsReal[0].(string))
		paramsToArr = append(paramsToArr, "=")
		paramsToArr = append(paramsToArr, "?")
		b.IOrm.SetBindValues(argsReal[1])
	}

	return strings.Join(paramsToArr, " "), nil
}
