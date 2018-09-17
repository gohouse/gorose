package gorose

import (
	"database/sql"
	"errors"
	"github.com/gohouse/gorose/utils"
	"fmt"
	"github.com/gohouse/gorose/across"
	"github.com/gohouse/gorose/helper"
	"reflect"
	"strings"
)

type Database struct {
	across.OrmApi
	Connection *Connection
}

func (dba *Database) Table(arg interface{}) *Database {
	dba.STable = arg
	return dba
}

// Fields : select fields
func (dba *Database) Fields(fields ...string) *Database {
	dba.Sfields = fields
	return dba
}

// AddFields : If you already have a query builder instance and you wish to add a column to its existing select clause, you may use the AddFields method:
func (dba *Database) AddFields(fields ...string) *Database {
	dba.Sfields = append(dba.Sfields, fields...)
	return dba
}

// Data : insert or update data
func (dba *Database) Data(data interface{}) *Database {
	dba.Sdata = data
	return dba
}

// Group : select group by
func (dba *Database) Group(group string) *Database {
	dba.Sgroup = group
	return dba
}

// GroupBy : equals Group()
func (dba *Database) GroupBy(group string) *Database {
	return dba.Group(group)
}

// Having : select having
func (dba *Database) Having(having string) *Database {
	dba.Shaving = having
	return dba
}

// Order : select order by
func (dba *Database) Order(order string) *Database {
	dba.Sorder = order
	return dba
}

// OrderBy : equal order
func (dba *Database) OrderBy(order string) *Database {
	return dba.Order(order)
}

// Limit : select limit
func (dba *Database) Limit(limit int) *Database {
	dba.Slimit = limit
	return dba
}

// Offset : select offset
func (dba *Database) Offset(offset int) *Database {
	dba.Soffset = offset
	return dba
}

// Take : select limit
func (dba *Database) Take(limit int) *Database {
	return dba.Limit(limit)
}

// Skip : select offset
func (dba *Database) Skip(offset int) *Database {
	return dba.Offset(offset)
}

// Page : select page
func (dba *Database) Page(page int) *Database {
	dba.Soffset = (page - 1) * dba.Slimit
	return dba
}

// Where : query or execute where condition, the relation is and
func (dba *Database) Where(args ...interface{}) *Database {
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组

	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	w := []interface{}{"and", args}

	dba.Swhere = append(dba.Swhere, w)

	return dba
}

// OrWhere : like where , but the relation is or,
func (dba *Database) OrWhere(args ...interface{}) *Database {
	w := []interface{}{"or", args}
	dba.Swhere = append(dba.Swhere, w)

	return dba
}

// WhereNull : like where , where filed is null,
func (dba *Database) WhereNull(arg string) *Database {
	return dba.Where("arg", "is", "null")
}

// WhereNotNull : like where , where filed is not null,
func (dba *Database) WhereNotNull(arg string) *Database {
	return dba.Where("arg", "is not", "null")
}

// OrWhereNull : like WhereNull , the relation is or,
func (dba *Database) OrWhereNull(arg string) *Database {
	return dba.OrWhere("arg", "is", "null")
}

// OrWhereNotNull : like WhereNotNull , the relation is or,
func (dba *Database) OrWhereNotNull(arg string) *Database {
	return dba.OrWhere("arg", "is not", "null")
}

// WhereIn : a given column's value is contained within the given array
func (dba *Database) WhereIn(field string, arr []interface{}) *Database {
	return dba.Where(field, "in", arr)
}

// WhereNotIn : the given column's value is not contained in the given array
func (dba *Database) WhereNotIn(field string, arr []interface{}) *Database {
	return dba.Where(field, "not in", arr)
}

// OrWhereIn : as WhereIn, the relation is or
func (dba *Database) OrWhereIn(field string, arr []interface{}) *Database {
	return dba.OrWhere(field, "in", arr)
}

// OrWhereNotIn : as WhereNotIn, the relation is or
func (dba *Database) OrWhereNotIn(field string, arr []interface{}) *Database {
	return dba.OrWhere(field, "not in", arr)
}

// WhereBetween : a column's value is between two values:
func (dba *Database) WhereBetween(field string, arr []interface{}) *Database {
	return dba.Where(field, "between", arr)
}

// WhereNotBetween : a column's value lies outside of two values:
func (dba *Database) WhereNotBetween(field string, arr []interface{}) *Database {
	return dba.Where(field, "not between", arr)
}

// OrWhereBetween : like WhereNull , the relation is or,
func (dba *Database) OrWhereBetween(field string, arr []interface{}) *Database {
	return dba.OrWhere(field, "not between", arr)
}

// OrWhereNotBetween : like WhereNotNull , the relation is or,
func (dba *Database) OrWhereNotBetween(field string, arr []interface{}) *Database {
	return dba.OrWhere(field, "not in", arr)
}

// Join : select join query
func (dba *Database) Join(args ...interface{}) *Database {
	//dba.parseJoin(args, "INNER")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"INNER", args})

	return dba
}

// InnerJoin : equals join
func (dba *Database) InnerJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "INNER")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"INNER", args})

	return dba
}

// LeftJoin : like join , the relation is left
func (dba *Database) LeftJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "LEFT")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"LEFT", args})

	return dba
}

// RightJoin : like join , the relation is right
func (dba *Database) RightJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "RIGHT")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"RIGHT", args})

	return dba
}

// CrossJoin : like join , the relation is cross
func (dba *Database) CrossJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "RIGHT")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"CROSS", args})

	return dba
}

// UnionJoin : like join , the relation is union
func (dba *Database) UnionJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "RIGHT")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"UNION", args})

	return dba
}

// Distinct : select distinct
func (dba *Database) Distinct() *Database {
	dba.Sdistinct = true

	return dba
}

// Pluck : Retrieving A List Of Column Values
func (dba *Database) Pluck(args ...string) (interface{}, error) {
	res, err := dba.Get()
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	switch len(args) {
	case 1:
		var pluckTmp []interface{}
		for _, item := range res {
			pluckTmp = append(pluckTmp, item[args[0]])
		}
		return pluckTmp, nil
	case 2:
		var pluckTmp = make(map[interface{}]interface{})
		for _, item := range res {
			pluckTmp[item[args[1]]] = item[args[0]]
		}
		return pluckTmp, nil
	default:
		return nil, errors.New("params error")
	}
}

// Value : If you don't even need an entire row, you may extract a single value from a record using the  value method. This method will return the value of the column directly:
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
func (dba *Database) Count(args ...string) (int64, error) {
	fields := "*"
	if len(args) > 0 {
		fields = utils.ParseStr(args[0])
	}
	return dba.UnionAct("count", fields)
}

// Sum : select sum field
func (dba *Database) Sum(sum string) (int64, error) {
	return dba.UnionAct("sum", sum)
}

// Avg : select avg field
func (dba *Database) Avg(avg string) (int64, error) {
	return dba.UnionAct("avg", avg)
}

// Max : select max field
func (dba *Database) Max(max string) (int64, error) {
	return dba.UnionAct("max", max)
}

// Min : select min field
func (dba *Database) Min(min string) (int64, error) {
	return dba.UnionAct("min", min)
}

// Chunk : select chunk more data to piceses block
func (dba *Database) Chunk(limit int, callback func([]map[string]interface{})) {
	var step = 0
	var offset = dba.Soffset
	for {
		dba.Slimit = limit
		dba.Soffset = offset + step*limit

		// 查询当前区块的数据
		sqls, _ := dba.BuildSql()
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

// Loop : select more data to piceses block from begin
func (dba *Database) Loop(limit int, callback func([]map[string]interface{})) {
	dba.Slimit = limit
	for {
		// 查询当前区块的数据
		sqls, _ := dba.BuildSql()
		data, _ := dba.Query(sqls)
		if len(data) == 0 {
			break
		}

		callback(data)

		if len(data) < limit {
			break
		}
	}
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
	// 记录sqlLog
	dba.LastSql = fmt.Sprintf(sqlstring, vals...)
	dba.SqlLogs = append(dba.SqlLogs, dba.LastSql)

	var operType string = strings.ToLower(sqlstring[0:6])
	if operType == "select" {
		return 0, errors.New("this method does not allow select operations, use Query")
	}

	var stmt *sql.Stmt
	var err error
	if dba.Strans == true {
		stmt, err = dba.Stx.Prepare(sqlstring)
	} else {
		stmt, err = dba.Connection.GetExecuteDb().Prepare(sqlstring)
	}

	if err != nil {
		return 0, err
	}
	//return dba.parseExecute(stmt, operType, vals)

	var rowsAffected int64
	//var err error
	defer stmt.Close()
	result, errs := stmt.Exec(vals...)
	if errs != nil {
		return 0, errs
	}

	if operType == "insert" {
		// get last insert id
		lastInsertId, err := result.LastInsertId()
		if err == nil {
			dba.LastInsertId = lastInsertId
		}
	}
	// get rows affected
	rowsAffected, err = result.RowsAffected()

	// 如果是事务, 则重置所有参数
	if dba.Strans == true {
		dba.Reset("transaction")
	}

	return rowsAffected, err
}
// Insert : insert data and get affected rows
func (dba *Database) Insert() (int, error) {
	sqlstr, err := dba.BuildSql("insert")
	if err != nil {
		return 0, err
	}

	res, err := dba.Execute(sqlstr)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

// insertGetId : insert data and get id
func (dba *Database) InsertGetId() (int, error) {
	_, err := dba.Insert()
	if err != nil {
		return 0, err
	}
	return int(dba.LastInsertId), nil
}

// Update : update data
func (dba *Database) Update() (int, error) {
	sqlstr, err := dba.BuildSql("update")
	if err != nil {
		return 0, err
	}

	res, errs := dba.Execute(sqlstr)
	if errs != nil {
		return 0, errs
	}
	return int(res), nil
}

// Delete : delete data
func (dba *Database) Delete() (int, error) {
	sqlstr, err := dba.BuildSql("delete")
	if err != nil {
		return 0, err
	}

	res, errs := dba.Execute(sqlstr)
	if errs != nil {
		return 0, errs
	}
	return int(res), nil
}

// Increment : auto Increment +1 default
// we can define step (such as 2, 3, 6 ...) if give the second params
// we can use this method as decrement with the third param as "-"
func (dba *Database) Increment(args ...interface{}) (int, error) {
	argLen := len(args)
	var field string
	var value string = "1"
	var mode string = "+"
	switch argLen {
	case 1:
		field = args[0].(string)
	case 2:
		field = args[0].(string)
		switch args[1].(type) {
		case int:
			value = utils.ParseStr(args[1])
		case int64:
			value = utils.ParseStr(args[1])
		case float32:
			value = utils.ParseStr(args[1])
		case float64:
			value = utils.ParseStr(args[1])
		case string:
			value = args[1].(string)
		default:
			return 0, errors.New("第二个参数类型错误")
		}
	case 3:
		field = args[0].(string)
		switch args[1].(type) {
		case int:
			value = utils.ParseStr(args[1])
		case int64:
			value = utils.ParseStr(args[1])
		case float32:
			value = utils.ParseStr(args[1])
		case float64:
			value = utils.ParseStr(args[1])
		case string:
			value = args[1].(string)
		default:
			return 0, errors.New("第二个参数类型错误")
		}
		mode = args[2].(string)
	default:
		return 0, errors.New("参数数量只允许1个,2个或3个")
	}
	dba.Data(field + "=" + field + mode + value)
	return dba.Update()
}

// Decrement : auto Decrement -1 default
// we can define step (such as 2, 3, 6 ...) if give the second params
func (dba *Database) Decrement(args ...interface{}) (int, error) {
	arglen := len(args)
	switch arglen {
	case 1:
		args = append(args, 1)
		args = append(args, "-")
	case 2:
		args = append(args, "-")
	default:
		return 0, errors.New("Decrement参数个数有误")
	}
	return dba.Increment(args...)
}

func (dba *Database) Begin() {
	dba.Stx, _ = dba.Connection.GetExecuteDb().Begin()
	dba.Strans = true
}
func (dba *Database) Commit() {
	dba.Stx.Commit()
	dba.Strans = false
}
func (dba *Database) Rollback() {
	dba.Stx.Rollback()
	dba.Strans = false
}

// Reset : reset union select
func (dba *Database) Reset(source string) {
	if source == "transaction" {
		//this = new(Database)
		dba.STable = ""
		dba.Sfields = []string{}
		dba.Swhere = [][]interface{}{}
		dba.Sorder = ""
		dba.Slimit = 0
		dba.Soffset = 0
		dba.Sjoin = [][]interface{}{}
		dba.Sdistinct = false
		dba.Sgroup = ""
		dba.Shaving = ""
		var tmp interface{}
		dba.Sdata = tmp
	}
	dba.Sunion = ""
}

// ResetWhere : in transaction, when you need update several tables in difference condition
func (dba *Database) ResetWhere() {
	dba.Swhere = [][]interface{}{}
}

// JsonEncode : parse json
func (dba *Database) JsonEncode(data interface{}) string {
	res, _ := utils.JsonEncode(data)
	return res
}

// buildData : build inert or update data
func (dba *Database) buildData() (string, string, string) {
	// insert
	var dataFields []string
	var dataValues []string
	// update or delete
	var dataObj []string

	data := dba.Sdata

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
		//case "map[string]interface {}":
	default: // update or insert
		datas := make(map[string]string)
		switch data.(type) {
		case map[string]interface{}:
			for key, val := range data.(map[string]interface{}) {
				if val == nil {
					datas[key] = "null"
				} else {
					datas[key] = utils.ParseStr(val)
				}
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
			//dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(val))
			if val == "null" {
				dataValuesSub = append(dataValuesSub, "null")
			} else {
				dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(val))
			}
			// update
			//dataObj = append(dataObj, key+"="+utils.AddSingleQuotes(val))
			if val == "null" {
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

// UnionAct : build union select real
func (dba *Database) UnionAct(union, field string) (int64, error) {
	var tmp int64 = 0

	dba.Sunion = union + "(" + field + ") as " + union
	
	// 构建sql
	sqls, err := dba.BuildSql()
	if err != nil {
		return tmp, err
	}

	// 执行查询
	result, err := dba.Query(sqls)
	if err != nil {
		return tmp, err
	}

	dba.Reset("union")

	//fmt.Println(union, reflect.TypeOf(union), " ", result[0][union])
	if len(result) > 0 {
		tmp = result[0][union].(int64)
	}

	return tmp, nil
}

func (dba *Database) Get() (result []map[string]interface{}, err error) {
	var sqlStr string
	sqlStr, err = dba.BuildSql()
	if err != nil {
		return
	}
	result, err = dba.Query(sqlStr)

	return
}

func (dba *Database) First() (result map[string]interface{}, err error) {
	dba.Slimit = 1
	var resultSlice []map[string]interface{}
	if resultSlice, err = dba.Get(); err != nil {
		return
	}
	if len(resultSlice) > 0 {
		result = resultSlice[0]
	}
	return
}

func (dba *Database) Select() (err error) {
	_, err = dba.Get()
	return
}

// Transaction : is a simple usage of trans
func (dba *Database) Transaction(closure func() (error)) (bool, error) {
	dba.Begin()
	err := closure()
	if err != nil {
		dba.Rollback()
		return false, err
	}
	dba.Commit()

	return true, nil
}

//// BuildSql : build sql string , but not execute sql really
//// operType : select/insert/update/delete
func (dba *Database) BuildSql(operType ...string) (string, error) {
	//dba.Driver = dba.Connection.DbConfig.Driver
	err := dba.ParseTable()
	if err!=nil{
		return "",err
	}
	dba.Driver = "mysql"
	return NewBuilder(dba.OrmApi, operType...)
}

//func (dba *Database) BuildQuery() (sql string, err error) {
//	var fields, table, limit, offset string
//	// table
//	if table, err = dba.ParseTable(); err != nil {
//		return
//	}
//	// fields
//	fields = strings.Join(dba.Sfields, ", ")
//	if fields == "" {
//		fields = "*"
//	}
//	// limit
//	limit = " limit 3"
//	// offset
//	offset = " offset 0"
//
//	//sqlstr := "select " + fields + " from " + table + limit + offset
//	sqlstr := fmt.Sprintf("SELECT %s FROM %s%s%s", fields, table, limit, offset)
//
//	return sqlstr, nil
//}

// Query : query instance of sql.DB.Query
func (dba *Database) Query(arg string, params ...interface{}) (result []map[string]interface{}, errs error) {
	lenParams := len(params)
	var vals []interface{}

	if lenParams > 1 {
		for k, v := range params {
			if k > 0 {
				vals = append(vals, v)
			}
		}
	}

	stmt, err := dba.Connection.GetQueryDb().Prepare(arg)
	if err != nil {
		return result, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(vals...)
	if err != nil {
		return result, err
	}

	// make sure we always close rows
	defer rows.Close()

	return dba.Scan(rows)
}
func (dba *Database) Scan(rows *sql.Rows) (result []map[string]interface{}, err error) {
	// 检查实多维数组还是一维数组
	switch dba.TableType {
	case across.TABLE_STRUCT_SLICE:
		err = dba.ScanAll(rows, dba.STable)
	case across.TABLE_STRUCT:
		err = dba.ScanRow(rows, dba.STable)
	case across.TABLE_STRING:
		result, err = dba.ScanMap(rows)
	}
	return
}

func (dba *Database) ScanMap(rows *sql.Rows) (result []map[string]interface{}, err error) {
	var columns []string
	if columns, err = rows.Columns(); err != nil {
		return
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
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else if val == nil {
				v = "NULL"
			} else {
				v = val
			}
			entry[col] = v
		}
		result = append(result, entry)
	}
	return
}

// scan a single row of data into a struct.
func (dba *Database) ScanRow(rows *sql.Rows, dst interface{}) error {
	// check if there is data waiting
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	// get a list of targets
	var fields = helper.StrutForScan(dst)

	// perform the scan
	if err := rows.Scan(fields...); err != nil {
		return err
	}

	return rows.Err()
}

// ScanAll scans all sql result rows into a slice of structs.
// It reads all rows and closes rows when finished.
// dst should be a pointer to a slice of the appropriate type.
// The new results will be appended to any existing data in dst.
func (dba *Database) ScanAll(rows *sql.Rows, dst interface{}) error {
	for rows.Next() {
		// scan it
		err := rows.Scan(helper.StrutForScan(dba.TableStruct.Interface())...)
		if err != nil {
			return err
		}
		// add to the result slice
		dba.TableSlice.Set(reflect.Append(dba.TableSlice, dba.TableStruct.Elem()))
	}

	return rows.Err()
}

func (dba *Database) ParseTable() (error) {
	var tableName string
	switch dba.STable.(type) {
	case string: // 直接传入的是表名
		dba.TableType = across.TABLE_STRING
		tableName = dba.STable.(string)

	default: // 传入的是struct
		// make sure dst is an appropriate type
		dstVal := reflect.ValueOf(dba.STable)
		if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
			return fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.STable)
		}
		sliceVal := reflect.Indirect(dstVal)
		switch sliceVal.Kind() {
		case reflect.Struct: // struct
			dba.TableType = across.TABLE_STRUCT
			tableName = sliceVal.Type().Name()
			dba.TableStruct = sliceVal
			// 默认只查一条
			dba.Slimit = 1
		case reflect.Slice: // []struct
			eltType := sliceVal.Type().Elem()
			if eltType.Kind() != reflect.Struct {
				return fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.STable)
			}
			dba.TableType = across.TABLE_STRUCT_SLICE
			tableName = eltType.Name()
			dba.TableStruct = reflect.New(eltType)
			dba.TableSlice = sliceVal
		default:
			return fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.STable)
		}
		// 是否设置了表名
		if i, ok := dba.STable.(ITable); ok {
			tableName = i.TableName()
		}

		if len(dba.Sfields) == 0 {
			dba.Sfields = helper.GetTagName(dba.TableStruct.Interface())
		}
	}
	fmt.Println("表名: ", tableName)
	dba.TableName = tableName
	return nil
}
