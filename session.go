package gorose

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gohouse/gorose/across"
	"github.com/gohouse/gorose/utils"
	"reflect"
	"strings"
	"time"
)

type Session struct {
	across.OrmApi
	Connection *Connection
}

func (dba *Session) Table(arg interface{}) *Session {
	dba.STable = arg
	return dba
}

// Fields : select fields
func (dba *Session) Fields(fields ...string) *Session {
	dba.Sfields = fields
	return dba
}

// AddFields : If you already have a query builder instance and you wish to add a column to its existing select clause, you may use the AddFields method:
func (dba *Session) AddFields(fields ...string) *Session {
	dba.Sfields = append(dba.Sfields, fields...)
	return dba
}

// Data : insert or update data
func (dba *Session) Data(data interface{}) *Session {
	dba.Sdata = data
	return dba
}

// Group : select group by
func (dba *Session) Group(group string) *Session {
	dba.Sgroup = group
	return dba
}

// GroupBy : equals Group()
func (dba *Session) GroupBy(group string) *Session {
	return dba.Group(group)
}

// Having : select having
func (dba *Session) Having(having string) *Session {
	dba.Shaving = having
	return dba
}

// Order : select order by
func (dba *Session) Order(order string) *Session {
	dba.Sorder = order
	return dba
}

// OrderBy : equal order
func (dba *Session) OrderBy(order string) *Session {
	return dba.Order(order)
}

// Limit : select limit
func (dba *Session) Limit(limit int) *Session {
	dba.Slimit = limit
	return dba
}

// Offset : select offset
func (dba *Session) Offset(offset int) *Session {
	dba.Soffset = offset
	return dba
}

// Take : select limit
func (dba *Session) Take(limit int) *Session {
	return dba.Limit(limit)
}

// Skip : select offset
func (dba *Session) Skip(offset int) *Session {
	return dba.Offset(offset)
}

// Page : select page
func (dba *Session) Page(page int) *Session {
	dba.Soffset = (page - 1) * dba.Slimit
	return dba
}

// Where : query or execute where condition, the relation is and
func (dba *Session) Where(args ...interface{}) *Session {
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组

	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	w := []interface{}{"and", args}

	dba.Swhere = append(dba.Swhere, w)

	return dba
}

// OrWhere : like where , but the relation is or,
func (dba *Session) OrWhere(args ...interface{}) *Session {
	w := []interface{}{"or", args}
	dba.Swhere = append(dba.Swhere, w)

	return dba
}

// WhereNull : like where , where filed is null,
func (dba *Session) WhereNull(arg string) *Session {
	return dba.Where("arg", "is", "null")
}

// WhereNotNull : like where , where filed is not null,
func (dba *Session) WhereNotNull(arg string) *Session {
	return dba.Where("arg", "is not", "null")
}

// OrWhereNull : like WhereNull , the relation is or,
func (dba *Session) OrWhereNull(arg string) *Session {
	return dba.OrWhere("arg", "is", "null")
}

// OrWhereNotNull : like WhereNotNull , the relation is or,
func (dba *Session) OrWhereNotNull(arg string) *Session {
	return dba.OrWhere("arg", "is not", "null")
}

// WhereIn : a given column's value is contained within the given array
func (dba *Session) WhereIn(field string, arr interface{}) *Session {
	return dba.Where(field, "in", arr)
}

// WhereNotIn : the given column's value is not contained in the given array
func (dba *Session) WhereNotIn(field string, arr interface{}) *Session {
	return dba.Where(field, "not in", arr)
}

// OrWhereIn : as WhereIn, the relation is or
func (dba *Session) OrWhereIn(field string, arr interface{}) *Session {
	return dba.OrWhere(field, "in", arr)
}

// OrWhereNotIn : as WhereNotIn, the relation is or
func (dba *Session) OrWhereNotIn(field string, arr interface{}) *Session {
	return dba.OrWhere(field, "not in", arr)
}

// WhereBetween : a column's value is between two values:
func (dba *Session) WhereBetween(field string, arr interface{}) *Session {
	return dba.Where(field, "between", arr)
}

// WhereNotBetween : a column's value lies outside of two values:
func (dba *Session) WhereNotBetween(field string, arr interface{}) *Session {
	return dba.Where(field, "not between", arr)
}

// OrWhereBetween : like WhereNull , the relation is or,
func (dba *Session) OrWhereBetween(field string, arr interface{}) *Session {
	return dba.OrWhere(field, "not between", arr)
}

// OrWhereNotBetween : like WhereNotNull , the relation is or,
func (dba *Session) OrWhereNotBetween(field string, arr interface{}) *Session {
	return dba.OrWhere(field, "not in", arr)
}

// Join : select join query
func (dba *Session) Join(args ...interface{}) *Session {
	//dba.parseJoin(args, "INNER")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"INNER", args})

	return dba
}

// Force : delete or update without where condition
func (dba *Session) Force(arg ...bool) *Session {
	if len(arg) > 0 {
		dba.Sforce = arg[0]
	} else {
		dba.Sforce = true
	}
	return dba
}

// InnerJoin : equals join
func (dba *Session) InnerJoin(args ...interface{}) *Session {
	//dba.parseJoin(args, "INNER")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"INNER", args})

	return dba
}

// LeftJoin : like join , the relation is left
func (dba *Session) LeftJoin(args ...interface{}) *Session {
	//dba.parseJoin(args, "LEFT")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"LEFT", args})

	return dba
}

// RightJoin : like join , the relation is right
func (dba *Session) RightJoin(args ...interface{}) *Session {
	//dba.parseJoin(args, "RIGHT")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"RIGHT", args})

	return dba
}

// CrossJoin : like join , the relation is cross
func (dba *Session) CrossJoin(args ...interface{}) *Session {
	//dba.parseJoin(args, "RIGHT")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"CROSS", args})

	return dba
}

// UnionJoin : like join , the relation is union
func (dba *Session) UnionJoin(args ...interface{}) *Session {
	//dba.parseJoin(args, "RIGHT")
	dba.Sjoin = append(dba.Sjoin, []interface{}{"UNION", args})

	return dba
}

// Distinct : select distinct
func (dba *Session) Distinct() *Session {
	dba.Sdistinct = true

	return dba
}

// Pluck : Retrieving A List Of Column Values
func (dba *Session) Pluck(args ...string) (interface{}, error) {
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
func (dba *Session) Value(arg string) (interface{}, error) {
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
func (dba *Session) Count(args ...string) (int64, error) {
	fields := "*"
	if len(args) > 0 {
		fields = utils.ParseStr(args[0])
	}
	count, err := dba.UnionAct("count", fields)
	return count.(int64), err
}

// Sum : select sum field
func (dba *Session) Sum(sum string) (interface{}, error) {
	return dba.UnionAct("sum", sum)
}

// Avg : select avg field
func (dba *Session) Avg(avg string) (interface{}, error) {
	return dba.UnionAct("avg", avg)
}

// Max : select max field
func (dba *Session) Max(max string) (interface{}, error) {
	return dba.UnionAct("max", max)
}

// Min : select min field
func (dba *Session) Min(min string) (interface{}, error) {
	return dba.UnionAct("min", min)
}

// Chunk : select chunk more data to piceses block
func (dba *Session) Chunk(limit int, callback func([]map[string]interface{})) {
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
func (dba *Session) Loop(limit int, callback func([]map[string]interface{})) {
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
func (dba *Session) Execute(sqlstring string, params ...interface{}) (int64, error) {
	t_start := time.Now()
	//defer DB.Close()
	lenParams := len(params)
	var vals []interface{}

	if lenParams > 0 {
		for _, v := range params {
			vals = append(vals, v)
		}
	}
	// 记录sqlLog
	if dba.Connection.DbConfig.Master.EnableQueryLog {
		dba.LastSql = fmt.Sprintf(sqlstring, vals...)
		dba.SqlLogs = append(dba.SqlLogs, dba.LastSql)
	}

	var operType string = strings.ToLower(sqlstring[0:6])
	if operType == "select" {
		return 0, errors.New("Execute does not allow select operations, use Query")
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

	// 持久化日志
	if dba.Connection.Logger != nil {
		dba.Connection.Logger.Write(dba.LastSql, time.Since(t_start).String(), time.Now().Format("2006-01-02 15:04:05"))
	}

	return rowsAffected, err
}

// Insert : insert data and get affected rows
func (dba *Session) Insert() (int64, error) {
	return dba.ExecuteAct("insert")
}

// insertGetId : insert data and get id
func (dba *Session) InsertGetId() (int64, error) {
	_, err := dba.Insert()
	return dba.LastInsertId, err
}

// Update : update data
func (dba *Session) Update() (int64, error) {
	return dba.ExecuteAct("update")
}

// Delete : delete data
func (dba *Session) Delete() (int64, error) {
	return dba.ExecuteAct("delete")
}

func (dba *Session) ExecuteAct(operType string) (int64, error) {
	sqlstr, err := dba.BuildSql(operType)
	if err != nil {
		return 0, err
	}
	return dba.Execute(sqlstr)
}

// Increment : auto Increment +1 default
// we can define step (such as 2, 3, 6 ...) if give the second params
// we can use this method as decrement with the third param as "-"
func (dba *Session) Increment(args ...interface{}) (int64, error) {
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
func (dba *Session) Decrement(args ...interface{}) (int64, error) {
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

func (dba *Session) Begin() {
	dba.Stx, _ = dba.Connection.GetExecuteDb().Begin()
	dba.Strans = true
}
func (dba *Session) Commit() {
	dba.Stx.Commit()
	dba.Strans = false
}
func (dba *Session) Rollback() {
	dba.Stx.Rollback()
	dba.Strans = false
}

// Reset : reset union select
func (dba *Session) Reset(source string) {
	if source == "transaction" {
		//this = new(Session)
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
func (dba *Session) ResetWhere() {
	dba.Swhere = [][]interface{}{}
}

// JsonEncode : parse json
func (dba *Session) JsonEncode(data interface{}) string {
	res, _ := utils.JsonEncode(data)
	return res
}

// UnionAct : build union select real
func (dba *Session) UnionAct(union, field string) (interface{}, error) {
	var tmp interface{}

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
		tmp = result[0][union]
	}

	return tmp, nil
}

func (dba *Session) Get() (result []map[string]interface{}, err error) {
	var sqlStr string
	sqlStr, err = dba.BuildSql()
	if err != nil {
		return
	}
	result, err = dba.Query(sqlStr)

	return
}

func (dba *Session) First() (result map[string]interface{}, err error) {
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

func (dba *Session) Select() (err error) {
	_, err = dba.Get()
	return
}

// Transaction : is a simple usage of trans
func (dba *Session) Transaction(closure func() (error)) (bool, error) {
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
func (dba *Session) BuildSql(operType ...string) (string, error) {
	// table解析
	err := dba.ParseTable()
	if err != nil {
		return "", err
	}
	// 表前缀, 驱动
	dba.Prefix = dba.Connection.DbConfig.Master.Prefix
	dba.Driver = dba.Connection.DbConfig.Master.Driver
	return NewBuilder(dba.OrmApi, operType...)
}

// Query : query instance of sql.DB.Query
func (dba *Session) Query(sqlstring string, params ...interface{}) (result []map[string]interface{}, errs error) {
	t_start := time.Now()
	lenParams := len(params)
	var vals []interface{}

	if lenParams > 0 {
		for _, v := range params {
			vals = append(vals, v)
		}
	}
	// 记录sqlLog
	if dba.Connection.DbConfig.Master.EnableQueryLog {
		dba.LastSql = fmt.Sprintf(sqlstring, vals...)
		dba.SqlLogs = append(dba.SqlLogs, dba.LastSql)
	}

	stmt, err := dba.Connection.GetQueryDb().Prepare(sqlstring)
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

	res,err2 := dba.Scan(rows)

	// 持久化日志
	if dba.Connection.Logger != nil {
		//fmt.Println(dba.Connection.Logger)
		dba.Connection.Logger.Write(dba.LastSql, time.Since(t_start).String(), time.Now().Format("2006-01-02 15:04:05"))
	}

	return res,err2
}
func (dba *Session) Scan(rows *sql.Rows) (result []map[string]interface{}, err error) {
	// 检查实多维数组还是一维数组
	switch dba.TableType {
	case across.TABLE_STRUCT_SLICE:
		err = dba.ScanAll(rows, dba.STable)
	case across.TABLE_STRUCT:
		err = dba.ScanRow(rows, dba.STable)
	//case across.TABLE_STRING:
	default:
		result, err = dba.ScanMap(rows)
	}
	return
}

func (dba *Session) ScanMap(rows *sql.Rows) (result []map[string]interface{}, err error) {
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
func (dba *Session) ScanRow(rows *sql.Rows, dst interface{}) error {
	// check if there is data waiting
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	// get a list of targets
	var fields = utils.StrutForScan(dst)

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
func (dba *Session) ScanAll(rows *sql.Rows, dst interface{}) error {
	for rows.Next() {
		// scan it
		err := rows.Scan(utils.StrutForScan(dba.TableStruct.Interface())...)
		if err != nil {
			return err
		}
		// add to the result slice
		dba.TableSlice.Set(reflect.Append(dba.TableSlice, dba.TableStruct.Elem()))
	}

	return rows.Err()
}

func (dba *Session) ParseTable() (error) {
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
			// 是否设置了表名
			if tn := dstVal.MethodByName("TableName"); tn.IsValid() {
				tableName = tn.Call(nil)[0].String()
			}
		case reflect.Slice: // []struct
			eltType := sliceVal.Type().Elem()
			if eltType.Kind() != reflect.Struct {
				return fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.STable)
			}
			dba.TableType = across.TABLE_STRUCT_SLICE
			tableName = eltType.Name()
			dba.TableStruct = reflect.New(eltType)
			dba.TableSlice = sliceVal
			// 是否设置了表名
			if tn := dba.TableStruct.MethodByName("TableName"); tn.IsValid() {
				tableName = tn.Call(nil)[0].String()
			}
		default:
			return fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.STable)
		}

		if len(dba.Sfields) == 0 {
			dba.Sfields = utils.GetTagName(dba.TableStruct.Interface())
		}
	}
	//fmt.Println("表名: ", tableName)
	dba.TableName = tableName
	return nil
}
