package gorose

import (
	"database/sql"
	"fmt"
	"github.com/gohouse/gorose/api"
	"github.com/gohouse/gorose/builder"
	"github.com/gohouse/gorose/helper"
	"reflect"
	"strings"
)

type ITable interface {
	TableName() string
}

//type table struct {
//	table       interface{}
//	tableName   string
//	tableStruct reflect.Value
//	tableSlice  reflect.Value
//	tableType   api.TableType
//}

type Database struct {
	api.OrmApi
	Connection *Connection
	//table      table
	//fields     []string
	//limit      int
}

func NewDatabase() *Database {
	return &Database{}
}

func (dba *Database) Table(arg interface{}) *Database {
	dba.STable = arg
	//fmt.Println(dba.STable)
	//os.Exit(1)
	return dba
}

func (dba *Database) Get() (result []map[string]interface{}, err error) {
	var sqlStr string
	sqlStr, err = dba.BuildSql()
	fmt.Println(sqlStr)
	if err != nil {
		return
	}
	result, err = dba.Query(sqlStr)

	return
}

func (dba *Database) First() (result map[string]interface{}, err error) {
	dba.SLimit = 1
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

//// BuildSql : build sql string , but not execute sql really
//// operType : select/insert/update/delete
func (dba *Database) BuildSql(operType ...string) (string, error) {
	//dba.Driver = dba.Connection.DbConfig.Driver
	dba.ParseTable()
	dba.Driver = "mysql"
	return builder.BuildSql(dba.OrmApi, operType...)
}

func (dba *Database) BuildQuery() (sql string, err error) {
	var fields, table, limit, offset string
	// table
	if table, err = dba.ParseTable(); err != nil {
		return
	}
	// fields
	fields = strings.Join(dba.SFields, ", ")
	if fields == "" {
		fields = "*"
	}
	// limit
	limit = " limit 3"
	// offset
	offset = " offset 0"

	//sqlstr := "select " + fields + " from " + table + limit + offset
	sqlstr := fmt.Sprintf("SELECT %s FROM %s%s%s", fields, table, limit, offset)

	return sqlstr, nil
}

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

	stmt, err := dba.Connection.Db.Prepare(arg)
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
	case TABLE_STRUCT_SLICE:
		err = dba.ScanAll(rows, dba.STable)
	case TABLE_STRUCT:
		err = dba.ScanRow(rows, dba.STable)
	case TABLE_STRING:
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
	var fields = StrutForScan(dst)

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
	//// make sure dst is an appropriate type
	//dstVal := reflect.ValueOf(dst)
	//if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
	//	return fmt.Errorf("ScanAll called with non-pointer destination: %T", dst)
	//}
	//sliceVal := reflect.Indirect(dstVal)
	//if sliceVal.Kind() != reflect.Slice {
	//	return fmt.Errorf("ScanAll called with pointer to non-slice: %T", dst)
	//}
	//eltType := sliceVal.Type().Elem()
	//if eltType.Kind() != reflect.Struct {
	//	return fmt.Errorf("ScanAll expects element to be pointers to structs, found %T", dst)
	//}
	//
	//var eltVal reflect.Value
	//var elt interface{}
	//// create a new element
	//eltVal = reflect.New(eltType)
	//elt = eltVal.Interface()
	// gather the results
	for rows.Next() {
		// scan it
		err := rows.Scan(StrutForScan(dba.TableStruct.Interface())...)
		if err != nil {
			return err
		}
		// add to the result slice
		dba.TableSlice.Set(reflect.Append(dba.TableSlice, dba.TableStruct.Elem()))
	}

	return rows.Err()
}

func StrutForScan(u interface{}) []interface{} {
	val := reflect.ValueOf(u).Elem()
	v := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		v[i] = valueField.Addr().Interface()
	}
	return v
}

func (dba *Database) BuildExecut(operType string) (string, error) {
	return "", nil
}

func (dba *Database) ParseTable() (string, error) {
	var tableName string
	switch dba.STable.(type) {
	case string: // 直接传入的是表名
		dba.TableType = TABLE_STRING
		dba.TableName = dba.STable.(string)
		return dba.TableName, nil

	default: // 传入的是struct
		// make sure dst is an appropriate type
		dstVal := reflect.ValueOf(dba.STable)
		if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
			return tableName, fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.STable)
		}
		sliceVal := reflect.Indirect(dstVal)
		switch sliceVal.Kind() {
		case reflect.Struct: // struct
			dba.TableType = TABLE_STRUCT
			tableName = sliceVal.Type().Name()
			dba.TableStruct = sliceVal
		case reflect.Slice: // []struct
			eltType := sliceVal.Type().Elem()
			if eltType.Kind() != reflect.Struct {
				return tableName, fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.STable)
			}
			dba.TableType = TABLE_STRUCT_SLICE
			tableName = eltType.Name()
			dba.TableStruct = reflect.New(eltType)
			dba.TableSlice = sliceVal
		default:
			return tableName, fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.STable)
		}
		// 是否设置了表名
		if i, ok := dba.STable.(ITable); ok {
			tableName = i.TableName()
		}

		if len(dba.SFields) == 0 {
			dba.SFields = helper.GetTagName(dba.TableStruct.Interface())
		}
	}
	return tableName, nil
}
