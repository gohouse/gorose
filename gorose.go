package gorose

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"utils"
)

var DB *sql.DB
var Tx *sql.Tx
var dbDriver string = "mysql" // sqlite, postgre...
var dbDefault string = "mysql"
var dbConfig = map[string]map[string]interface{}{
	"mysql": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "test",
		"charset":  "utf8",
		"protocol": "tcp",
	},
	"mysql_dev": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "gorose",
		"charset":  "utf8",
		"protocol": "tcp",
	},
}

var regex = []string{"=", ">", "<", "!=", ">=", "<=", "in", "not in", "between", "not between"}

func init() {
	Connect(dbDefault)
}

//var instance *Database
//var once sync.Once
//func GetInstance() *Database {
//	once.Do(func() {
//		instance = &Database{}
//	})
//	return instance
//}
func Connect(arg interface{}) *sql.DB {
	var err error
	var dbObj map[string]interface{}
	if utils.GetType(arg) == "string" {
		dbObj = dbConfig[arg.(string)]
	} else {
		dbObj = arg.(map[string]interface{})
	}

	conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s",
		dbObj["username"], dbObj["password"], dbObj["protocol"], dbObj["host"], dbObj["port"], dbObj["database"], dbObj["charset"])
	//DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	DB, err = sql.Open(dbDriver, conn)
	if err != nil {
		//log.Fatal(err.Error())
		panic(err.Error())
	}
	err = DB.Ping()
	if err != nil {
		//log.Fatal(err.Error())
		panic(err.Error())
	}

	return DB
}

type Database struct {
	Table_    string
	Fields_   string
	Where_    [][]interface{}
	OrWhere_  [][]interface{}
	Order_    string
	Limit_    int
	Offset_   int
	Page_     int
	Join_     []string
	Distinct_ bool
	Count_    string
	Sum_      string
	Avg_      string
	Max_      string
	Min_      string
	Group_    string
	Trans_    bool
	Data_     interface{}
}

func (this *Database) Connect(arg interface{}) *Database {
	Connect(arg)
	return this
}
func (this *Database) Fields(Fields_ string) *Database {
	this.Fields_ = Fields_
	return this
}
func (this *Database) Table(Table_ string) *Database {
	this.Table_ = Table_
	return this
}
func (this *Database) Order(Order_ string) *Database {
	this.Order_ = Order_
	return this
}
func (this *Database) Limit(Limit_ int) *Database {
	this.Limit_ = Limit_
	return this
}
func (this *Database) Offset(Offset_ int) *Database {
	this.Offset_ = Offset_
	return this
}
func (this *Database) Page(Page_ int) *Database {
	this.Page_ = Page_
	return this
}
func (this *Database) First() map[string]interface{} {
	this.Limit_ = 1
	// 构建sql
	sqls := this.buildSql()
	fmt.Println(sqls)
	// 执行查询
	result := this.Query(sqls)

	if len(result) == 0 {
		return nil
	}

	return result[0]
}
func (this *Database) Get() []map[string]interface{} {
	// 构建sql
	sqls := this.buildSql()
	fmt.Println(sqls)
	// 执行查询
	result := this.Query(sqls)

	if len(result) == 0 {
		return nil
	}

	return result
}
func (this *Database) Where(args ...interface{}) *Database {
	argsLen := len(args)

	argsType := "string"

	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组
	if argsLen == 1 {
		argsType = utils.GetType(args[0])
	}

	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为参数类型(三种类型), 第三项为具体传入的参数
	w := []interface{}{"and", argsType, args}

	this.Where_ = append(this.Where_, w)

	return this
}
func (this *Database) OrWhere(args ...interface{}) *Database {
	argsLen := len(args)

	argsType := "string"

	if argsLen == 1 {
		argsType = utils.GetType(args[0])
	}

	w := []interface{}{"or", argsType, args}
	this.Where_ = append(this.Where_, w)
	return this
}
func (this *Database) Join(args ...interface{}) *Database {
	this.parseJoin(args, "INNER")

	return this
}
func (this *Database) LeftJoin(args ...interface{}) *Database {
	this.parseJoin(args, "LEFT")

	return this
}
func (this *Database) RightJoin(args ...interface{}) *Database {
	this.parseJoin(args, "RIGHT")

	return this
}

func (this *Database) Distinct() *Database {
	this.Distinct_ = true

	return this
}
func (this *Database) Count(count string) *Database {
	this.Count_ = "count(" + count + ") as count"

	return this
}
func (this *Database) Sum(sum string) *Database {
	this.Sum_ = "sum(" + sum + ") as sum"

	return this
}
func (this *Database) Avg(avg string) *Database {
	this.Avg_ = "avg(" + avg + ") as avg"

	return this
}
func (this *Database) Max(max string) *Database {
	this.Max_ = "max(" + max + ") as max"

	return this
}
func (this *Database) Min(min string) *Database {
	this.Min_ = "min(" + min + ") as min"

	return this
}

func (this *Database) parseJoin(args []interface{}, joinType string) bool {
	var w string
	argsLength := len(args)
	switch argsLength {
	case 1:
		w = args[0].(string)
	case 4:
		w = utils.ParseStr(args[0]) + " ON " + utils.ParseStr(args[1]) + " " + utils.ParseStr(args[2]) + " " + utils.ParseStr(args[3])
	default:
		panic("join格式错误")
	}

	this.Join_ = append(this.Join_, joinType+" JOIN "+w)

	return true
}

/**
 * where解析器
 */
func (this *Database) parseWhere() string {
	wheres := this.Where_

	// where解析后存放每一项的容器
	var where []string

	for _, args := range wheres {
		// and或者or条件
		var condition string = args[0].(string)
		// 数据类型
		var dataType string = args[1].(string)
		// 统计当前数组中有多少个参数
		params := args[2].([]interface{})
		paramsLength := len(params)

		switch paramsLength {
		case 3: // 常规3个参数:  {"id",">",1}
			where = append(where, condition+" ("+this.parseParams(params)+")")
		case 2: // 常规2个参数:  {"id",1}
			where = append(where, condition+" ("+this.parseParams(params)+")")
		case 1: // 二维数组或字符串
			if dataType == "string" { // sql 语句字符串
				where = append(where, condition+" ("+params[0].(string)+")")
			} else if dataType == "map[string]interface {}" { // 一维数组
				var whereArr []string
				for key, val := range params[0].(map[string]interface{}) {
					whereArr = append(whereArr, "("+key+"="+utils.AddSingleQuotes(val)+")")
				}
				where = append(where, condition+" ("+strings.Join(whereArr, " and ")+")")
			} else if dataType == "[][]interface {}" { // 二维数组
				var whereMore []string
				for _, arr := range params[0].([][]interface{}) { // {{"a", 1}, {"id", ">", 1}}
					whereMoreLength := len(arr)
					switch whereMoreLength {
					case 2:
						whereMore = append(whereMore, "("+this.parseParams(arr)+")")
					case 3:
						whereMore = append(whereMore, "("+this.parseParams(arr)+")")
					default:
						panic("where数据格式有误")
					}
				}
				where = append(where, condition+" ("+strings.Join(whereMore, " and ")+")")
			} else { // 不符合的类型
				panic("where条件格式错误")
			}
		}
	}

	return strings.TrimLeft(strings.Trim(strings.Join(where, " "), " "), "and")
}

/**
 * 将where条件中的参数转换为where条件字符串
 * example: {"id",">",1}, {"age", 18}
 */
func (this *Database) parseParams(args []interface{}) string {

	paramsLength := len(args)

	// 存储当前所有数据的数组
	var paramsToArr []string

	switch paramsLength {
	case 3: // 常规3个参数:  {"id",">",1}
		if !utils.TypeCheck(args[0], "string") {
			panic("where条件参数有误!")
		}
		if !utils.TypeCheck(args[1], "string") {
			panic("where条件参数有误!")
		}
		if !utils.InArray(args[1], utils.Astoi(regex)) {
			panic("where运算条件参数有误!!")
		}

		paramsToArr = append(paramsToArr, args[0].(string))
		paramsToArr = append(paramsToArr, args[1].(string))

		switch args[1] {
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
		if !utils.TypeCheck(args[0], "string") {
			panic("where条件参数有误!")
		}
		paramsToArr = append(paramsToArr, args[0].(string))
		paramsToArr = append(paramsToArr, "=")
		paramsToArr = append(paramsToArr, utils.AddSingleQuotes(args[1]))
	}

	return strings.Join(paramsToArr, " ")
}
func (this *Database) buildSql() string {
	// 聚合
	unionArr := []string{
		this.Count_,
		this.Sum_,
		this.Avg_,
		this.Max_,
		this.Min_,
	}
	var union string
	for _, item := range unionArr {
		if item != "" {
			union = item
			break
		}
	}
	// distinct
	distinct := utils.If(this.Distinct_, "distinct", "")
	// fields
	fields := utils.If(this.Fields_ == "", "*", this.Fields_).(string)
	// table
	table := this.Table_
	// join
	join := utils.If(strings.Join(this.Join_, "") == "", "", strings.Join(this.Join_, " "))
	// where
	where := utils.If(this.parseWhere() == "", "", "WHERE "+this.parseWhere()).(string)
	// group
	group := utils.If(this.Group_ == "", "", "GROUP BY "+this.Group_).(string)
	// order
	order := utils.If(this.Order_ == "", "", "ORDER BY "+this.Order_).(string)
	// limit
	limit := utils.If(this.Limit_ == 0, "LIMIT 1000", "LIMIT "+strconv.Itoa(this.Limit_)).(string)
	// offset
	offset := utils.If(this.Offset_ == 0, "", "OFFSET "+strconv.Itoa(this.Offset_)).(string)

	//sqlstr := "select " + fields + " from " + table + " " + where + " " + order + " " + limit + " " + offset
	sqlstr := fmt.Sprintf("SELECT %s %s FROM %s %s %s %s %s %s %s", distinct, utils.If(union != "", union, fields), table, join, where, group, order, limit, offset)

	return sqlstr
}

func (this *Database) Query(sqlstring string) []map[string]interface{} {
	stmt, err := DB.Prepare(sqlstring)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// 结果
	var result []map[string]interface{}
	var result_map = make(map[string]interface{})
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			result_map[columns[i]] = value
		}
		result = append(result, result_map)
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return result
}

/**
 *　执行增删改 ｓｑｌ 语句
 */
func (this *Database) Execute(sqlstring string) int64 {
	var operType string = strings.ToLower(sqlstring[0:6])
	if operType == "select" {
		panic("该方法不允许select操作, 请使用Query")
	}
	if this.Trans_ == true {
		stmt, err := Tx.Prepare(sqlstring)
		checkErr(err)
		return this.parseExecute(stmt, operType)
	} else {
		stmt, err := DB.Prepare(sqlstring)
		checkErr(err)
		return this.parseExecute(stmt, operType)
	}
}

func (this *Database) parseExecute(stmt *sql.Stmt, operType string) int64 {
	var res int64
	var err error
	result, err := stmt.Exec()
	checkErr(err)

	switch operType {
	case "insert":
		res, err = result.LastInsertId()
	case "update":
		res, err = result.RowsAffected()
	case "delete":
		res, err = result.RowsAffected()
	}
	checkErr(err)
	return res
}

func (this *Database) buildExecut(operType string) string {
	// insert : {"name":"fizz, "website":"fizzday.net"} or {{"name":"fizz2", "website":"www.fizzday.net"}, {"name":"fizz", "website":"fizzday.net"}}}
	// update : {"name":"fizz", "website":"fizzday.net"}
	// delete : ...
	upOrDel, insertkey, insertval := this.buildData()
	where := utils.If(this.parseWhere() == "", "", "WHERE "+this.parseWhere()).(string)
	var sqlstr string

	switch operType {
	case "insert":
		sqlstr = fmt.Sprintf("insert into %s (%s) values %s", this.Table_, insertkey, insertval)
	case "update":
		sqlstr = fmt.Sprintf("update %s set %s %s", this.Table_, upOrDel, where)
	case "delete":
		sqlstr = fmt.Sprintf("update %s set %s %s", this.Table_, upOrDel, where)
	}
	fmt.Println(sqlstr)
	return sqlstr
}

func (this *Database) buildData() (string, string, string) {
	// insert
	var dataFields []string
	var dataValues []string
	// update or delete
	var dataObj []string

	data := this.Data_
	fmt.Println(data, utils.GetType(data))
	dataType := utils.GetType(data)
	switch dataType {
	case "[]map[string]interface {}":
		datas := data.([]map[string]interface{})
		for _, item := range datas {
			var dataValuesSub []string
			for key, val := range item {
				if utils.InArray(key, utils.Astoi(dataFields)) == false {
					dataFields = append(dataFields, key)
				}
				dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(val))
			}
			dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
		}
	//case "map[string]interface {}":
	default:
		datas := make(map[string]string)
		switch dataType {
		case "map[string]interface {}":
			for key, val := range data.(map[string]interface{}) {
				datas[key] = utils.ParseStr(val)
			}
		case "map[string]int":
			for key, val := range data.(map[string]int) {
				datas[key] = utils.ParseStr(val)
			}
		case "map[string]string":
			for key, val := range data.(map[string]string) {
				datas[key] = val
			}
		}

		//datas := data.(map[string]interface{})
		var dataValuesSub []string
		for key, val := range datas {
			dataFields = append(dataFields, key)
			dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(val))
			// update or delete
			dataObj = append(dataObj, key+"="+utils.AddSingleQuotes(val))
		}
		dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
	}

	return strings.Join(dataObj, ","), strings.Join(dataFields, ","), strings.Join(dataValues, "")
}
func (this *Database) Data(data interface{}) *Database {
	//var tmp []interface{}
	//tmp = append(tmp, utils.GetType(data))
	//tmp = append(tmp, data)
	this.Data_ = data
	return this
}
func (this *Database) Insert() int64 {
	sqlstr := this.buildExecut("insert")
	return this.Execute(sqlstr)
}
func (this *Database) Update() int64 {
	sqlstr := this.buildExecut("insert")
	return this.Execute(sqlstr)
}
func (this *Database) Delete(sqlstring string) int64 {
	sqlstr := this.buildExecut("insert")
	return this.Execute(sqlstr)
}
func (this *Database) Begin() *sql.Tx {
	tx, _ := DB.Begin()
	this.Trans_ = true
	Tx = tx
	return tx
}
func (this *Database) Commit() *Database {
	Tx.Commit()
	this.Trans_ = false
	return this
}
func (this *Database) Rollback() *Database {
	Tx.Rollback()
	this.Trans_ = false
	return this
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("start")

	var db Database
	//init()
	query := db.Table("users").Fields("id, name").
		Where("id", "<", 100).
		Where("id", ">", 1).Get()
	fmt.Println(query)
}
