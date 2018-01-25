package gorose

import (
	"database/sql"
	"fmt"
	"github.com/gohouse/gorose/utils"
	"strconv"
	"strings"
	"errors"
)

var (
	regex    = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "in", "not in", "between", "not between"}
	Dbstruct Database
)

type MapData map[string]interface{}
type MultiData []map[string]interface{}

//var instance *Database
//var once sync.Once
//func GetInstance() *Database {
//	once.Do(func() {
//		instance = &Database{}
//	})
//	return instance
//}

type Database struct {
	table    string
	fields   string
	where    [][]interface{}
	order    string
	limit    int
	offset   int
	join     [][]interface{}
	distinct bool
	count    string
	sum      string
	avg      string
	max      string
	min      string
	group    string
	trans    bool
	data     interface{}
	sqlLogs  []string
}

func (this *Database) Close() error {
	this.sqlLogs = []string{}
	return DB.Close()
}
func (this *Database) Ping() error {
	return DB.Ping()
}

func (this *Database) Fields(fields string) *Database {
	this.fields = fields
	return this
}
func (this *Database) Table(table string) *Database {
	this.table = table
	return this
}
func (this *Database) Data(data interface{}) *Database {
	this.data = data
	return this
}
func (this *Database) Group(group string) *Database {
	this.group = group
	return this
}
func (this *Database) Order(order string) *Database {
	this.order = order
	return this
}
func (this *Database) Limit(limit int) *Database {
	this.limit = limit
	return this
}
func (this *Database) Offset(offset int) *Database {
	this.offset = offset
	return this
}
func (this *Database) Page(page int) *Database {
	this.offset = (page - 1) * this.limit
	return this
}
func (this *Database) Where(args ...interface{}) *Database {
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组

	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	w := []interface{}{"and", args}

	this.where = append(this.where, w)

	return this
}
func (this *Database) OrWhere(args ...interface{}) *Database {
	w := []interface{}{"or", args}
	this.where = append(this.where, w)

	return this
}
func (this *Database) Join(args ...interface{}) *Database {
	//this.parseJoin(args, "INNER")
	this.join = append(this.join, []interface{}{"INNER",args})

	return this
}
func (this *Database) LeftJoin(args ...interface{}) *Database {
	//this.parseJoin(args, "LEFT")
	this.join = append(this.join, []interface{}{"LEFT",args})

	return this
}
func (this *Database) RightJoin(args ...interface{}) *Database {
	//this.parseJoin(args, "RIGHT")
	this.join = append(this.join, []interface{}{"RIGHT",args})

	return this
}
func (this *Database) Distinct() *Database {
	this.distinct = true

	return this
}

func (this *Database) First(args ...interface{}) (map[string]interface{}, error) {
	//var result map[string]interface{}
	//func (this *Database) First() interface{} {
	this.limit = 1
	// 构建sql
	sqls, err := this.buildQuery()
	if err != nil {
		return nil, err
	}

	// 执行查询
	res, err := this.Query(sqls)
	if err != nil {
		return nil, err
	}

	// 之所以不在 Query 中统一Reset, 是因为chunk会复用到查询相关条件
	this.Reset()

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}
func (this *Database) Get() ([]map[string]interface{}, error) {
	//func (this *Database) Get() interface{} {
	// 构建sql
	sqls, err := this.buildQuery()
	if err != nil {
		return nil, err
	}
fmt.Println(sqls)
	// 执行查询
	result, err := this.Query(sqls)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	//if JsonEncode == true {
	//	jsons, _ := json.Marshal(result)
	//	return string(jsons)
	//}

	this.Reset()

	return result, nil
}
func (this *Database) Count() (int, error) {
	res, err := this.buildUnion("count", "*")
	if err != nil {
		return 0, err
	}
	return int(res.(int64)), nil
}
func (this *Database) Sum(sum string) (interface{}, error) {
	return this.buildUnion("sum", sum)
}
func (this *Database) Avg(avg string) (interface{}, error) {
	return this.buildUnion("avg", avg)
}
func (this *Database) Max(max string) (interface{}, error) {
	return this.buildUnion("max", max)
}
func (this *Database) Min(min string) (interface{}, error) {
	return this.buildUnion("min", min)
}
func (this *Database) Chunk(limit int, callback func([]map[string]interface{})) {
	var step = 0
	for {
		this.limit = limit
		this.offset = step * limit

		// 查询当前区块的数据
		sqls, _ := this.buildQuery()
		data, _ := this.Query(sqls)

		if len(data) == 0 {
			this.Reset()
			break
		}

		callback(data)

		//fmt.Println(res)
		if len(data) < limit {
			this.Reset()
			break
		}
		step++
	}
}

func (this *Database) buildQuery() (string, error) {
	// 聚合
	unionArr := []string{
		this.count,
		this.sum,
		this.avg,
		this.max,
		this.min,
	}
	var union string
	for _, item := range unionArr {
		if item != "" {
			union = item
			break
		}
	}
	// distinct
	distinct := utils.If(this.distinct, "distinct ", "")
	// fields
	fields := utils.If(this.fields == "", "*", this.fields).(string)
	// table
	table := this.table
	// join
	parseJoin, err := this.parseJoin()
	if err != nil {
		return "", err
	}
	join := parseJoin
	// where
	parseWhere, err := this.parseWhere()
	if err != nil {
		return "", err
	}
	where := utils.If(parseWhere == "", "", " WHERE "+parseWhere).(string)
	// group
	group := utils.If(this.group == "", "", " GROUP BY "+this.group).(string)
	// order
	order := utils.If(this.order == "", "", " ORDER BY "+this.order).(string)
	// limit
	limit := utils.If(this.limit == 0, "", " LIMIT "+strconv.Itoa(this.limit))
	// offset
	offset := utils.If(this.offset == 0, "", " OFFSET "+strconv.Itoa(this.offset))

	//sqlstr := "select " + fields + " from " + table + " " + where + " " + order + " " + limit + " " + offset
	sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s",
		distinct, utils.If(union != "", union, fields), table, join, where, group, order, limit, offset)

	//fmt.Println(sqlstr)
	// Reset Database struct

	return sqlstr, nil
}
func (this *Database) buildExecut(operType string) (string, error) {
	// insert : {"name":"fizz, "website":"fizzday.net"} or {{"name":"fizz2", "website":"www.fizzday.net"}, {"name":"fizz", "website":"fizzday.net"}}}
	// update : {"name":"fizz", "website":"fizzday.net"}
	// delete : ...
	var update, insertkey, insertval string
	if operType != "delete" {
		update, insertkey, insertval = this.buildData()
	}
	res,err := this.parseWhere()
	if err!=nil{
		return res,err
	}
	where := utils.If( res== "", "", " WHERE "+res).(string)
	var sqlstr string

	switch operType {
	case "insert":
		sqlstr = fmt.Sprintf("insert into %s (%s) values %s", this.table, insertkey, insertval)
	case "update":
		sqlstr = fmt.Sprintf("update %s set %s%s", this.table, update, where)
	case "delete":
		sqlstr = fmt.Sprintf("delete from %s%s", this.table, where)
	}
	//fmt.Println(sqlstr)
	this.Reset()

	return sqlstr,nil
}
func (this *Database) buildData() (string, string, string) {
	// insert
	var dataFields []string
	var dataValues []string
	// update or delete
	var dataObj []string

	data := this.data

	switch data.(type) {
	case MultiData: // insert multi datas ([]map[string]interface{})
		datas := data.(MultiData)
		for _, item := range datas {
			var dataValuesSub []string
			for key, val := range item {
				if utils.InArray(key, dataFields) == false {
					dataFields = append(dataFields, key)
				}
				dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(val))
			}
			dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
		}
		//case "map[string]interface {}":
	default: // update or insert
		datas := make(map[string]string)
		switch data.(type) {
		case map[string]interface{}:
			for key, val := range data.(map[string]interface{}) {
				datas[key] = utils.ParseStr(val)
			}
		case map[string]int:
			for key, val := range data.(map[string]int) {
				datas[key] = utils.ParseStr(val)
			}
		case map[string]string:
			for key, val := range data.(map[string]string) {
				datas[key] = val
			}
		}

		//datas := data.(map[string]interface{})
		var dataValuesSub []string
		for key, val := range datas {
			// insert
			dataFields = append(dataFields, key)
			dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(val))
			// update
			dataObj = append(dataObj, key+"="+utils.AddSingleQuotes(val))
		}
		// insert
		dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
	}

	return strings.Join(dataObj, ","), strings.Join(dataFields, ","), strings.Join(dataValues, "")
}
func (this *Database) buildUnion(union, field string) (interface{}, error) {
	unionStr := union + "(" + field + ") as " + union
	switch union {
	case "count":
		this.count = unionStr
	case "sum":
		this.sum = unionStr
	case "avg":
		this.avg = unionStr
	case "max":
		this.max = unionStr
	case "min":
		this.min = unionStr
	}

	// 构建sql
	sqls, err := this.buildQuery()
	if err != nil {
		return nil, err
	}

	// 执行查询
	result, err := this.Query(sqls)
	if err != nil {
		return nil, err
	}

	this.Reset()

	//fmt.Println(union, reflect.TypeOf(union), " ", result[0][union])
	return result[0][union], nil
}


/**
 * 将where条件中的参数转换为where条件字符串
 * example: {"id",">",1}, {"age", 18}
 */
func (this *Database) parseParams(args []interface{}) (string, error) {
	paramsLength := len(args)

	// 存储当前所有数据的数组
	var paramsToArr []string

	switch paramsLength {
	case 3: // 常规3个参数:  {"id",">",1}
		if !utils.InArray(args[1], regex) {
			return "",errors.New("where运算条件参数有误!!")
		}

		paramsToArr = append(paramsToArr, args[0].(string))
		paramsToArr = append(paramsToArr, args[1].(string))

		switch args[1] {
		case "like":
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes("%"+utils.ParseStr(args[2])+"%"))
		case "in":
			paramsToArr = append(paramsToArr, "("+utils.Implode(args[2], ",")+")")
		case "not in":
			paramsToArr = append(paramsToArr, "("+utils.Implode(args[2], ",")+")")
		case "between":
			tmpB := args[2].([]string)
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(tmpB[0])+" and "+utils.AddSingleQuotes(tmpB[1]))
		case "not between":
			tmpB := args[2].([]string)
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(tmpB[0])+" and "+utils.AddSingleQuotes(tmpB[1]))
		default:
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(args[2]))
		}
	case 2:
		//if !utils.TypeCheck(args[0], "string") {
		//	panic("where条件参数有误!")
		//}
		paramsToArr = append(paramsToArr, args[0].(string))
		paramsToArr = append(paramsToArr, "=")
		paramsToArr = append(paramsToArr, utils.AddSingleQuotes(args[1]))
	}

	return strings.Join(paramsToArr, " "), nil
}
func (this *Database) parseJoin() (string, error) {
	var join []interface{}
	var returnJoinArr []string
	joinArr := this.join

	for _, join = range joinArr {
		var w string
		var ok bool
		var args []interface{}

		if len(join) != 2 {
			return "", errors.New("join条件有误")
		}

		// 获取真正的where条件
		if args, ok = join[1].([]interface{}); !ok {
			return "", errors.New("join条件有误")
		}

		argsLength := len(args)
		switch argsLength {
		case 1:
			w = args[0].(string)
		case 4:
			w = args[0].(string) + " ON " + args[1].(string) + " " + args[2].(string) + " " + args[3].(string)
		default:
			return "", errors.New("join格式错误")
		}

		returnJoinArr = append(returnJoinArr, " "+join[0].(string)+ " JOIN "+w)
	}

	return strings.Join(returnJoinArr, " "), nil
}
func (this *Database) parseWhere() (string, error) {
	// 取出所有where
	wheres := this.where
	//// where解析后存放每一项的容器
	var where []string

	for _, args := range wheres {
		// and或者or条件
		var condition string = args[0].(string)
		// 统计当前数组中有多少个参数
		params := args[1].([]interface{})
		paramsLength := len(params)

		switch paramsLength {
		case 3: // 常规3个参数:  {"id",">",1}
			res, err := this.parseParams(params)
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)

		case 2: // 常规2个参数:  {"id",1}
			res, err := this.parseParams(params)
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
						res, err := this.parseParams(params)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					case 2:
						res, err := this.parseParams(params)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					default:
						return "", errors.New("where数据格式有误")
					}
				}
				where = append(where, condition+" ("+strings.Join(whereMore, " and ")+")")
			case func():
				// 清空where,给嵌套的where让路,复用这个节点
				this.where = [][]interface{}{}

				// 执行嵌套where放入Database struct
				paramReal()
				// 再解析一遍后来嵌套进去的where
				wherenested,err := this.parseWhere()
				if err!=nil{
					return "",err
				}
				// 嵌套的where放入一个括号内
				where = append(where, condition+" ("+wherenested+")")
			default:
				return "", errors.New("where数据格式有误")
			}
		}
	}

	return strings.TrimLeft(strings.Trim(strings.Join(where, " "), " "), "and"), nil
}
func (this *Database) parseExecute(stmt *sql.Stmt, operType string, vals []interface{}) (int64, error) {
	var res int64
	var err error
	result, errs := stmt.Exec(vals...)
	if errs != nil{
		return 0,errs
	}

	switch operType {
	case "insert":
		res, err = result.LastInsertId()
	case "update":
		res, err = result.RowsAffected()
	case "delete":
		res, err = result.RowsAffected()
	}

	return res, err
}

func (this *Database) Insert() (int, error) {
	sqlstr,err := this.buildExecut("insert")
	if err!=nil{
		return 0,nil
	}

	res,errs := this.Execute(sqlstr)
	if errs!=nil{
		return 0,nil
	}
	return int(res),nil
}
func (this *Database) Update() (int, error) {
	sqlstr,err := this.buildExecut("update")
	if err!=nil{
		return 0,nil
	}

	res,errs := this.Execute(sqlstr)
	if errs!=nil{
		return 0,nil
	}
	return int(res),nil
}
func (this *Database) Delete() (int, error) {
	sqlstr,err := this.buildExecut("delete")
	if err!=nil{
		return 0,nil
	}

	res,errs := this.Execute(sqlstr)
	if errs!=nil{
		return 0,nil
	}
	return int(res),nil
}
func (this *Database) Begin() {
	Tx, _ = DB.Begin()
	this.trans = true
}
func (this *Database) Commit() {
	Tx.Commit()
	this.trans = false
}
func (this *Database) Rollback() {
	Tx.Rollback()
	this.trans = false
}
func (this *Database) Reset() {
	//this = new(Database)
	this.table = ""
	this.fields = ""
	this.where = [][]interface{}{}
	this.order = ""
	this.limit = 0
	this.offset = 0
	this.join = [][]interface{}{}
	this.distinct = false
	this.count = ""
	this.sum = ""
	this.avg = ""
	this.max = ""
	this.min = ""
	this.group = ""
	this.trans = false

	var tmp interface{}
	this.data = tmp
}
func (this *Database) JsonEncode(data interface{}) string {
	res, _ := utils.JsonEncode(data)
	return res
}

/**
 *　执行查询 sql 语句
 */
func (this *Database) Query(args ...interface{}) ([]map[string]interface{}, error) {
	tableData := make([]map[string]interface{}, 0)

	lenArgs := len(args)
	var sqlstring string
	var vals []interface{}

	sqlstring = args[0].(string)

	if lenArgs > 1 {
		for k, v := range args {
			if k > 0 {
				vals = append(vals, v)
			}
		}
	}
	// 记录sqllog
	this.sqlLogs = append(this.sqlLogs, fmt.Sprintf(sqlstring, vals...))

	stmt, err := DB.Prepare(sqlstring)
	if err != nil {
		return tableData, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(vals...)
	if err != nil {
		return tableData, err
	}

	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return tableData, err
	}

	count := len(columns)

	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			scanArgs[i] = &values[i]
		}
		rows.Scan(scanArgs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			//fmt.Println(val, reflect.TypeOf(val))
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}
/**
 *　执行增删改 ｓｑｌ 语句
 */
func (this *Database) Execute(args ...interface{}) (int64, error) {
	lenArgs := len(args)
	var sqlstring string
	var vals []interface{}

	sqlstring = args[0].(string)

	if lenArgs > 1 {
		for k, v := range args {
			if k > 0 {
				vals = append(vals, v)
			}
		}
	}
	// 记录sqllog
	this.sqlLogs = append(this.sqlLogs, fmt.Sprintf(sqlstring, vals...))

	var operType string = strings.ToLower(sqlstring[0:6])
	if operType == "select" {
		return 0, errors.New("该方法不允许select操作, 请使用Query")
	}

	if this.trans == true {
		stmt, err := Tx.Prepare(sqlstring)
		if err!= nil {
			return 0, err
		}
		return this.parseExecute(stmt, operType, vals)
	} else {
		stmt, err := DB.Prepare(sqlstring)
		if err!= nil {
			return 0, err
		}
		return this.parseExecute(stmt, operType, vals)
	}
}

func (this *Database) LastSql() string {
	if len(this.sqlLogs)>0{
		return this.sqlLogs[len(this.sqlLogs)-1:][0]
	}
	return ""
}
func (this *Database) SqlLogs() []string {
	return this.sqlLogs
}

/**
 * simple transaction
 */
func (this *Database) Transaction(closure func() (error)) bool {
	//defer func() {
	//	if err := recover(); err != nil {
	//		this.Rollback()
	//		panic(err)
	//	}
	//}()

	this.Begin()
	err := closure()
	if err!=nil {
		this.Rollback()
		return false
	}
	this.Commit()

	return true
}
