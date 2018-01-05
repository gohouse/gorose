package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"strings"
	"kuaixinwen/utils"
	"strconv"
)

var DB *sql.DB

var config = map[string]map[string]interface{} {
	"mysql":{
		"host":"localhost",
		"username":"root",
		"password":"",
		"port":"3306",
		"database": "test",
		"charset": "utf8",
	},
}

func init() {
	var err error
	DB, err = sql.Open("mysql", "gcore:gcore@tcp(192.168.200.248:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}

type Database struct {
	// Where_ := {{"and",{"a", 1}}, {"or",{"a", "=", 2}}}
	Table_  string
	Fields_ string
	Where_ [][]interface{}
	OrWhere_ [][]interface{}
	Order_     string
	Limit_     int
	Offset_    int
	Page_      int
	bindParams []interface{}
}

var regex = []string{"=", ">", "<", "!=", ">=", "<=", "in", "not in", "between", "not between"}

func (this *Database) Table(Table_ string) *Database {
	this.Table_ = Table_
	return this
}
func (this *Database) Fields(Fields_ string) *Database {
	this.Fields_ = Fields_
	return this
}
func (this *Database) GetFields() string {
	return this.Fields_
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

	if (len(result) == 0) {
		return nil
	}

	return result[0]
}
func (this *Database) Get() []map[string]interface{} {
	this.Limit_ = 1
	// 构建sql
	sqls := this.buildSql()
fmt.Println(sqls)
	// 执行查询
	result := this.Query(sqls)

	if (len(result) == 0) {
		return nil
	}

	return result
}
func (this *Database) Where(args ...interface{}) *Database {
	w := []interface{}{"and", args}
	this.Where_ = append(this.Where_, w)
	return this
}
func (this *Database) OrWhere(args ...interface{}) *Database {
	w := []interface{}{"or", args}
	this.Where_ = append(this.Where_, w)
	return this
}
func (this *Database) parseWhere() string {
	// example1.0 := string "a=3 and b>4"
	// example1.1 := map[string]interface{} {"a":1, "b":"bye"}
	// example1.2 := [][]interface{} {{"a",1},{"b",">",2},{"c","=",4}}
	// example2.0 := []interface{} {"a", 3}
	// example3.0 := []interface{} {"a", ">", 4}
	var sqlstr []string

	// where
	wheres := this.Where_

	// 查看args的长度
	//var whereFileds []interface{}

	for _, args := range wheres {
		var sqlstrItem string = " "+args[0].(string)+" "
		// 统计当前数组中有多少个参数
		item := args[1].([]interface{})
		argLen := len(item)
		switch argLen {
		case 3:
			if (!utils.TypeCheck(item[0], "string")) {
				panic("where条件参数有误!")
			}
			if (!utils.TypeCheck(item[1], "string")) {
				panic("where条件参数有误!")
			}
			if (!utils.InArray(item[1], utils.Astoi(regex))) {
				panic("where运算条件参数有误!!")
			}

			sqlstrItem += item[0].(string) +" "+ item[1].(string) + " "

			switch item[1] {
			case "in":
				sqlstrItem += "(" + strings.Join(item[2].([]string), ",") + ")"
			case "not in":
				sqlstrItem += "(" + strings.Join(item[2].([]string), ",") + ")"
			case "between":
				tmpB := item[2].([]string)
				sqlstrItem += tmpB[0] + " and " + tmpB[1]
			case "not between":
				tmpB := item[2].([]string)
				sqlstrItem += tmpB[0] + " and " + tmpB[1]
			default:
				sqlstrItem += utils.ParseStr(item[2])
			}

		case 2:
			if (!utils.TypeCheck(item[0], "string")) {
				panic("where条件参数有误!")
			}

			sqlstrItem += item[0].(string) + "=" + utils.ParseStr(item[1])

		case 1: // 二维数组或字符串
			dataType := utils.GetType(item)

			if dataType == "string" { // sql 语句字符串
				sqlstrItem += item[0].(string)
			} else if dataType == "map[string]interface {}" { // 一维数组
				for key, val := range item[0].(map[string]interface{}) {
					sqlstrItem += key+"="+utils.ParseStr(val)
				}
			} else if dataType == "[]map[string]interface {}" { // 二维数组
				var sqlstrItemOne []string
				for _, arr := range item[0].([][]interface{}) {	// {{"a", 1}}
					arrLen := len(arr)
					switch arrLen {
					case 2:
						if (!utils.TypeCheck(arr[0], "string")) {
							panic("where条件参数有误!")
						}
						sqlstrItemOne = append(sqlstrItemOne, arr[0].(string)+"="+arr[1].(string))
					case 3:
						if (!utils.TypeCheck(arr[0], "string")) {
							panic("where条件参数有误!")
						}
						if (!utils.InArray(arr[1], utils.Astoi(regex))) {
							panic("where运算条件参数有误!!")
						}
						var sqlstrItemOneSub string = arr[0].(string)+" "+arr[1].(string)+" "
						switch arr[1] {
						case "in":
							sqlstrItemOneSub += "(" + strings.Join(arr[2].([]string), ",") + ")"
							sqlstrItemOne = append(sqlstrItemOne,  sqlstrItemOneSub)
						case "not in":
							sqlstrItemOneSub += "(" + strings.Join(arr[2].([]string), ",") + ")"
							sqlstrItemOne = append(sqlstrItemOne,  sqlstrItemOneSub)
						case "between":
							tmpB := arr[2].([]string)
							sqlstrItemOneSub += tmpB[0] + " and " + tmpB[1]
						case "not between":
							tmpB := arr[2].([]string)
							sqlstrItemOneSub += tmpB[0] + " and " + tmpB[1]
						default:
							sqlstrItemOne = append(sqlstrItemOne, utils.ParseStr(arr[2]))
						}
					default:
						panic("where数据格式有误")
					}
				}
				sqlstrItem += strings.Join(sqlstrItemOne, " and ")
			} else { // 不符合的类型
				panic("where条件格式错误")
			}
		}

		sqlstr = append(sqlstr, sqlstrItem)
	}

	where3 := strings.Join(sqlstr, " ")
	where2 := strings.Trim(where3, " ")
	where := strings.TrimLeft(where2, "and")

	return where
}
func (this *Database) buildSql() (string) {
	// fields
	fields := utils.If(this.Fields_ == "", "*", this.Fields_).(string)
	// where
	where := utils.If(this.parseWhere()=="", "", "where "+this.parseWhere()).(string)

	order := utils.If(this.Order_=="", "", "order by "+this.Order_).(string)
	// limit
	limit := utils.If(this.Limit_==0, "limit 100", "limit "+strconv.Itoa(this.Limit_)).(string)
	// offset
	offset := utils.If(this.Offset_==0, "", "offset "+strconv.Itoa(this.Offset_)).(string)
	// count

	sqlstr := "select "+fields+" from "+this.Table_+" "+where+" "+order+" "+limit+" "+offset

	return sqlstr
}

func (this *Database) Query(sqlstring string) ([]map[string]interface{}) {
	defer DB.Close()
	stmt, err := DB.Prepare(sqlstring)
	if err != nil {
		fmt.Println("Query Error", err)
		panic(err.Error())
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("Query Error", err)
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
		//jsonstring += "{"
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

func main() {
	fmt.Println("start")

	var db Database

	query := db.Table("userinfo").Fields("id, lvs").
		Where("id", "<", 100).
		Where("id", ">", 1).Get()
	fmt.Println(query)
}
