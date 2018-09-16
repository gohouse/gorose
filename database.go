package gorose

import (
	"database/sql"
	"errors"
	"github.com/gohouse/gorose/helper"
	"fmt"
	"reflect"
	"strings"
)

type ITable interface {
	TableName() string
}

type table struct {
	table       interface{}
	tableName   string
	tableStruct reflect.Value
	tableSlice  reflect.Value
	tableType   TableType
}

type Database struct {
	connection *Connection
	table      table
	fields     []string
	first      interface{}
	limit      int
}

func (dba *Database) Table(arg interface{}) *Database {
	dba.table.table = arg
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
	dba.limit = 1
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

// BuildSql : build sql string , but not execute sql really
// operType : select/insert/update/delete
func (dba *Database) BuildSql(operType ...string) (string, error) {
	switch len(operType) {
	case 0:
		return dba.BuildQuery()
	case 1:
		if operType[0] == "select" {
			return dba.BuildQuery()
		} else {
			return dba.BuildExecut(operType[0])
		}
	default:
		return "", errors.New("参数有误")
	}
}

func (dba *Database) BuildQuery() (sql string, err error) {
	var fields, table, limit, offset string
	// table
	if table, err = dba.parseTable(); err != nil {
		return
	}
	// fields
	fields = strings.Join(dba.fields, ", ")
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

	stmt, err := dba.connection.db.Prepare(arg)
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
	switch dba.table.tableType {
	case TABLE_STRUCT_SLICE:
		err = dba.ScanAll(rows, dba.table.table)
	case TABLE_STRUCT:
		err = dba.ScanRow(rows, dba.table.table)
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
		err := rows.Scan(StrutForScan(dba.table.tableStruct.Interface())...)
		if err != nil {
			return err
		}
		// add to the result slice
		dba.table.tableSlice.Set(reflect.Append(dba.table.tableSlice, dba.table.tableStruct.Elem()))
	}

	return rows.Err()
}
func checkTableStruct() string {
	return "slice"
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

func (dba *Database) parseTable() (string, error) {
	var tableName string
	switch dba.table.table.(type) {
	case string: // 直接传入的是表名
		dba.table.tableType = TABLE_STRING
		return dba.table.table.(string), nil

	default: // 传入的是struct
		// make sure dst is an appropriate type
		dstVal := reflect.ValueOf(dba.table.table)
		if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
			return tableName, fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.table)
		}
		sliceVal := reflect.Indirect(dstVal)
		switch sliceVal.Kind() {
		case reflect.Struct: // struct
			dba.table.tableType = TABLE_STRUCT
			tableName = sliceVal.Type().Name()
			dba.table.tableStruct = sliceVal
		case reflect.Slice: // []struct
			eltType := sliceVal.Type().Elem()
			if eltType.Kind() != reflect.Struct {
				return tableName, fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.table)
			}
			dba.table.tableType = TABLE_STRUCT_SLICE
			tableName = eltType.Name()
			dba.table.tableStruct = reflect.New(eltType)
			dba.table.tableSlice = sliceVal
		default:
			return tableName, fmt.Errorf("table只接收字符串表名和struct, 但是传入的是: %T", dba.table)
		}
		// 是否设置了表名
		if i, ok := dba.table.table.(ITable); ok {
			tableName = i.TableName()
		}

		if len(dba.fields) == 0 {
			dba.fields = helper.GetTagName(dba.table.tableStruct.Interface())
		}
	}
	return tableName, nil
}
