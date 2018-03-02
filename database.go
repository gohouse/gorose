package gorose

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gohouse/gorose/utils"
	"strconv"
	"strings"
)

var (
	regex = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "not like", "in", "not in", "between", "not between"}
	//Dbstruct Database
)

//type MapData map[string]interface{}
//type MultiData []map[string]interface{}

//var instance *Database
//var once sync.Once
//func GetInstance() *Database {
//	once.Do(func() {
//		instance = &Database{}
//	})
//	return instance
//}

// Database is data mapper struct
type Database struct {
	table    string          // table
	fields   string          // fields
	where    [][]interface{} // where
	order    string          // order
	limit    int             // limit
	offset   int             // offset
	join     [][]interface{} // join
	distinct bool            // distinct
	count    string          // count
	sum      string          // sum
	avg      string          // avg
	max      string          // max
	min      string          // min
	group    string          // group
	having   string          // having
	data     interface{}     // data
	RowsAffected int		// insert affected rows
	LastInsertId int		// insert last insert id
	//trans    bool
	//sqlLogs  []string
}

// Fields : select fields
func (dba *Database) Fields(fields string) *Database {
	dba.fields = fields
	return dba
}

// Table : select table
func (dba *Database) Table(table string) *Database {
	dba.table = table
	return dba
}

// Data : insert or update data
func (dba *Database) Data(data interface{}) *Database {
	dba.data = data
	return dba
}

// Group : select group by
func (dba *Database) Group(group string) *Database {
	dba.group = group
	return dba
}

// Having : select having
func (dba *Database) Having(having string) *Database {
	dba.having = having
	return dba
}

// Order : select order by
func (dba *Database) Order(order string) *Database {
	dba.order = order
	return dba
}

// Limit : select limit
func (dba *Database) Limit(limit int) *Database {
	dba.limit = limit
	return dba
}

// Offset : select offset
func (dba *Database) Offset(offset int) *Database {
	dba.offset = offset
	return dba
}

// Page : select page
func (dba *Database) Page(page int) *Database {
	dba.offset = (page - 1) * dba.limit
	return dba
}

// Where : query or execute where condition, the relation is and
func (dba *Database) Where(args ...interface{}) *Database {
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组

	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	w := []interface{}{"and", args}

	dba.where = append(dba.where, w)

	return dba
}

// OrWhere : like where , but the relation is or,
func (dba *Database) OrWhere(args ...interface{}) *Database {
	w := []interface{}{"or", args}
	dba.where = append(dba.where, w)

	return dba
}

// Join : select join query
func (dba *Database) Join(args ...interface{}) *Database {
	//dba.parseJoin(args, "INNER")
	dba.join = append(dba.join, []interface{}{"INNER", args})

	return dba
}

// LeftJoin : like join , the relation is left
func (dba *Database) LeftJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "LEFT")
	dba.join = append(dba.join, []interface{}{"LEFT", args})

	return dba
}

// RightJoin : like join , the relation is right
func (dba *Database) RightJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "RIGHT")
	dba.join = append(dba.join, []interface{}{"RIGHT", args})

	return dba
}

// Distinct : select distinct
func (dba *Database) Distinct() *Database {
	dba.distinct = true

	return dba
}

// First : select one row
func (dba *Database) First(args ...interface{}) (map[string]interface{}, error) {
	//var result map[string]interface{}
	//func (dba *Database) First() interface{} {
	dba.limit = 1
	// 构建sql
	sqls, err := dba.buildQuery()
	if err != nil {
		return nil, err
	}

	// 执行查询
	res, err := dba.Query(sqls)
	if err != nil {
		return nil, err
	}

	// 之所以不在 Query 中统一Reset, 是因为chunk会复用到查询相关条件
	//dba.Reset()

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

// Get : select more rows , relation limit set
func (dba *Database) Get() ([]map[string]interface{}, error) {
	//func (dba *Database) Get() interface{} {
	// 构建sql
	sqls, err := dba.buildQuery()
	if err != nil {
		return nil, err
	}

	// 执行查询
	result, err := dba.Query(sqls)
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

	//dba.Reset()

	return result, nil
}

// Value : select one field's value
func (dba *Database) Value(arg string) (interface{}, error) {
	res, err := dba.First()
	if err != nil {
		return nil, err
	}
	if val, ok := res[arg]; ok {
		return val, nil
	}
	return nil, errors.New("the field is not exists")
}

// Count : select count rows
func (dba *Database) Count(args ...interface{}) (int, error) {
	fields := "*"
	if len(args) >0 {
		fields = utils.ParseStr(args[0])
	}
	res, err := dba.buildUnion("count", fields)
	if err != nil {
		return 0, err
	}
	return int(res.(int64)), nil
}

// Sum : select sum field
func (dba *Database) Sum(sum string) (interface{}, error) {
	return dba.buildUnion("sum", sum)
}

// Avg : select avg field
func (dba *Database) Avg(avg string) (interface{}, error) {
	return dba.buildUnion("avg", avg)
}

// Max : select max field
func (dba *Database) Max(max string) (interface{}, error) {
	return dba.buildUnion("max", max)
}

// Min : select min field
func (dba *Database) Min(min string) (interface{}, error) {
	return dba.buildUnion("min", min)
}

// Chunk : select chunk more data to piceses block
func (dba *Database) Chunk(limit int, callback func([]map[string]interface{})) {
	var step = 0
	var offset = dba.offset
	for {
		dba.limit = limit
		dba.offset = offset + step*limit

		// 查询当前区块的数据
		sqls, _ := dba.buildQuery()
		data, _ := dba.Query(sqls)

		if len(data) == 0 {
			//dba.Reset()
			break
		}

		callback(data)

		//fmt.Println(res)
		if len(data) < limit {
			//dba.Reset()
			break
		}
		step++
	}
}

// buildQuery : build query string
func (dba *Database) buildQuery() (string, error) {
	// 聚合
	unionArr := []string{
		dba.count,
		dba.sum,
		dba.avg,
		dba.max,
		dba.min,
	}
	var union string
	for _, item := range unionArr {
		if item != "" {
			union = item
			break
		}
	}
	// distinct
	distinct := utils.If(dba.distinct, "distinct ", "")
	// fields
	fields := utils.If(dba.fields == "", "*", dba.fields).(string)
	// table
	table := Connect.CurrentConfig["prefix"] + dba.table
	// join
	parseJoin, err := dba.parseJoin()
	if err != nil {
		return "", err
	}
	join := parseJoin
	// where
	parseWhere, err := dba.parseWhere()
	if err != nil {
		return "", err
	}
	where := utils.If(parseWhere == "", "", " WHERE "+parseWhere).(string)
	// group
	group := utils.If(dba.group == "", "", " GROUP BY "+dba.group).(string)
	// having
	having := utils.If(dba.having == "", "", " HAVING "+dba.having).(string)
	// order
	order := utils.If(dba.order == "", "", " ORDER BY "+dba.order).(string)
	// limit
	limit := utils.If(dba.limit == 0, "", " LIMIT "+strconv.Itoa(dba.limit))
	// offset
	offset := utils.If(dba.offset == 0, "", " OFFSET "+strconv.Itoa(dba.offset))

	//sqlstr := "select " + fields + " from " + table + " " + where + " " + order + " " + limit + " " + offset
	sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
		distinct, utils.If(union != "", union, fields), table, join, where, group, having, order, limit, offset)

	//fmt.Println(sqlstr)
	// Reset Database struct

	return sqlstr, nil
}

// buildExecut : build execute query string
func (dba *Database) buildExecut(operType string) (string, error) {
	// insert : {"name":"fizz, "website":"fizzday.net"} or {{"name":"fizz2", "website":"www.fizzday.net"}, {"name":"fizz", "website":"fizzday.net"}}}
	// update : {"name":"fizz", "website":"fizzday.net"}
	// delete : ...
	var update, insertkey, insertval, sqlstr string
	if operType != "delete" {
		update, insertkey, insertval = dba.buildData()
	}
	res, err := dba.parseWhere()
	if err != nil {
		return res, err
	}
	where := utils.If(res == "", "", " WHERE "+res).(string)

	tableName := Connect.CurrentConfig["prefix"] + dba.table
	switch operType {
	case "insert":
		sqlstr = fmt.Sprintf("insert into %s (%s) values %s", tableName, insertkey, insertval)
	case "update":
		sqlstr = fmt.Sprintf("update %s set %s%s", tableName, update, where)
	case "delete":
		sqlstr = fmt.Sprintf("delete from %s%s", tableName, where)
	}
	//fmt.Println(sqlstr)
	//dba.Reset()

	return sqlstr, nil
}

// buildData : build inert or update data
func (dba *Database) buildData() (string, string, string) {
	// insert
	var dataFields []string
	var dataValues []string
	// update or delete
	var dataObj []string

	data := dba.data

	switch data.(type) {
	case []map[string]interface{}: // insert multi datas ([]map[string]interface{})
		datas := data.([]map[string]interface{})
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

// buildUnion : build union select
func (dba *Database) buildUnion(union, field string) (interface{}, error) {
	unionStr := union + "(" + field + ") as " + union
	switch union {
	case "count":
		dba.count = unionStr
	case "sum":
		dba.sum = unionStr
	case "avg":
		dba.avg = unionStr
	case "max":
		dba.max = unionStr
	case "min":
		dba.min = unionStr
	}

	// 构建sql
	sqls, err := dba.buildQuery()
	if err != nil {
		return nil, err
	}

	// 执行查询
	result, err := dba.Query(sqls)
	if err != nil {
		return nil, err
	}

	dba.Reset()

	//fmt.Println(union, reflect.TypeOf(union), " ", result[0][union])
	return result[0][union], nil
}

/**
 * 将where条件中的参数转换为where条件字符串
 * example: {"id",">",1}, {"age", 18}
 */
// parseParams : 将where条件中的参数转换为where条件字符串
func (dba *Database) parseParams(args []interface{}) (string, error) {
	paramsLength := len(args)

	// 存储当前所有数据的数组
	var paramsToArr []string

	switch paramsLength {
	case 3: // 常规3个参数:  {"id",">",1}
		if !utils.InArray(args[1], regex) {
			return "", errors.New("where parameter is wrong")
		}

		paramsToArr = append(paramsToArr, args[0].(string))
		paramsToArr = append(paramsToArr, args[1].(string))

		switch args[1] {
		case "like":
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(utils.ParseStr(args[2])))
		case "not like":
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(utils.ParseStr(args[2])))
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

// parseJoin : parse the join paragraph
func (dba *Database) parseJoin() (string, error) {
	var join []interface{}
	var returnJoinArr []string
	joinArr := dba.join

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
func (dba *Database) parseWhere() (string, error) {
	// 取出所有where
	wheres := dba.where
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
			res, err := dba.parseParams(params)
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)

		case 2: // 常规2个参数:  {"id",1}
			res, err := dba.parseParams(params)
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
						res, err := dba.parseParams(params)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					case 2:
						res, err := dba.parseParams(params)
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
				dba.where = [][]interface{}{}

				// 执行嵌套where放入Database struct
				paramReal()
				// 再解析一遍后来嵌套进去的where
				wherenested, err := dba.parseWhere()
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

	return strings.TrimLeft(strings.Trim(strings.Join(where, " "), " "), "and"), nil
}

// parseExecute : parse execute condition
func (dba *Database) parseExecute(stmt *sql.Stmt, operType string, vals []interface{}) (int64, error) {
	var rowsAffected int64
	var err error
	result, errs := stmt.Exec(vals...)
	if errs != nil {
		return 0, errs
	}

	switch operType {
	case "insert":
		// get rows affected
		rowsAffected, err = result.RowsAffected()
		dba.RowsAffected = int(rowsAffected)
		// get last insert id
		rowsAffected, err = result.LastInsertId()
		dba.LastInsertId = int(rowsAffected)
	case "update":
		rowsAffected, err = result.RowsAffected()
	case "delete":
		rowsAffected, err = result.RowsAffected()
	}

	return rowsAffected, err
}

// Insert : insert data
func (dba *Database) Insert() (int, error) {
	sqlstr, err := dba.buildExecut("insert")
	if err != nil {
		return 0, err
	}

	res, err := dba.Execute(sqlstr)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

// Update : update data
func (dba *Database) Update() (int, error) {
	sqlstr, err := dba.buildExecut("update")
	if err != nil {
		return 0, err
	}

	res, errs := dba.Execute(sqlstr)
	if errs != nil {
		return 0, err
	}
	return int(res), nil
}

// Delete : delete data
func (dba *Database) Delete() (int, error) {
	sqlstr, err := dba.buildExecut("delete")
	if err != nil {
		return 0, err
	}

	res, errs := dba.Execute(sqlstr)
	if errs != nil {
		return 0, err
	}
	return int(res), nil
}

//func (dba *Database) Begin() {
//	Tx, _ = DB.Begin()
//	dba.trans = true
//}
//func (dba *Database) Commit() {
//	Tx.Commit()
//	dba.trans = false
//}
//func (dba *Database) Rollback() {
//	Tx.Rollback()
//	dba.trans = false
//}

// Reset : reset union select
func (dba *Database) Reset() {
	//this = new(Database)
	//dba.table = ""
	//dba.fields = ""
	//dba.where = [][]interface{}{}
	//dba.order = ""
	//dba.limit = 0
	//dba.offset = 0
	//dba.join = [][]interface{}{}
	//dba.distinct = false
	dba.count = ""
	dba.sum = ""
	dba.avg = ""
	dba.max = ""
	dba.min = ""
	//dba.group = ""
	//dba.having = ""
	//dba.trans = false
	//
	//var tmp interface{}
	//dba.data = tmp
}

// JsonEncode : parse json
func (dba *Database) JsonEncode(data interface{}) string {
	res, _ := utils.JsonEncode(data)
	return res
}

// Query : query instance of sql.DB.Query
func (dba *Database) Query(args ...interface{}) ([]map[string]interface{}, error) {
	//defer DB.Close()
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
	Connect.SqlLog = append(Connect.SqlLog, fmt.Sprintf(sqlstring, vals...))

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

// Execute : query instance of sql.DB.Execute
func (dba *Database) Execute(args ...interface{}) (int64, error) {
	//defer DB.Close()
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
	Connect.SqlLog = append(Connect.SqlLog, fmt.Sprintf(sqlstring, vals...))

	var operType string = strings.ToLower(sqlstring[0:6])
	if operType == "select" {
		return 0, errors.New("this method does not allow select operations, use Query")
	}

	if Connect.Trans == true {
		stmt, err := Tx.Prepare(sqlstring)
		if err != nil {
			return 0, err
		}
		return dba.parseExecute(stmt, operType, vals)
	}

	stmt, err := DB.Prepare(sqlstring)
	if err != nil {
		return 0, err
	}
	return dba.parseExecute(stmt, operType, vals)
}

//func (dba *Database) LastSql() string {
//	if len(Connect.SqlLog) > 0 {
//		return Connect.SqlLog[len(Connect.SqlLog)-1:][0]
//	}
//	return ""
//}
//func (dba *Database) SqlLogs() []string {
//	return Connect.SqlLog
//}
//
///**
// * simple transaction
// */
//func (dba *Database) Transaction(closure func() (error)) bool {
//	//defer func() {
//	//	if err := recover(); err != nil {
//	//		dba.Rollback()
//	//		panic(err)
//	//	}
//	//}()
//
//	dba.Begin()
//	err := closure()
//	if err != nil {
//		dba.Rollback()
//		return false
//	}
//	dba.Commit()
//
//	return true
//}
