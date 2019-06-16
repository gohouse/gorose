package gorose

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gohouse/t"
	"reflect"
	"strings"
)

type Session struct {
	IEngin
	//IOrm
	//IBuilder
	IBinder
	master       dbObject
	slave        dbObject
	lastInsertId int64
	sqlLogs      []string
	lastSql      string
	union interface{}
}

var _ ISession = &Session{}

// NewSession : 初始化 Session
func NewSession(e IEngin) ISession {
	var s = new(Session)
	s.IEngin = e

	s.master = s.IEngin.GetExecuteDB()
	s.slave = s.IEngin.GetQueryDB()

	s.IBinder = NewBinder()

	return s
}

// Close : 关闭 Session
func (s *Session) Close() {
	s = nil
}

//// GetBinder 获取绑定对象
//func (s *Session) GetBinder() *Binder {
//	return s.GetBinder()
//}

// GetMasterDriver 获取驱动
func (s *Session) GetMasterDriver() string {
	return s.master.driver
}

// GetSlaveDriver 获取驱动
func (s *Session) GetSlaveDriver() string {
	return s.slave.driver
}

// Bind : 传入绑定结果的对象, 参数一为对象, 可以是 struct, gorose.MapRow 或对应的切片
//		如果是做非query操作,第一个参数也可以仅仅指定为字符串表名
func (s *Session) Bind(tab interface{}) ISession {
	s.IBinder.SetBindOrigin(tab)
	return s
}

// GetTableName 获取解析后的名字, 提供给orm使用
func (s *Session) GetTableName() (string, error) {
	err := s.IBinder.BindParse(s.IEngin.GetPrefix())
	return s.IBinder.GetBindName(), err
}

func (s *Session) Begin() (err error) {
	s.master.tx, err = s.master.db.Begin()
	return
}

func (s *Session) Rollback() (err error) {
	err = s.master.tx.Rollback()
	s.master.tx = nil
	return
}

func (s *Session) Commit() (err error) {
	err = s.master.tx.Commit()
	s.master.tx = nil
	return
}

func (s *Session) Transaction(closers ...func(ses ISession) error) (err error) {
	err = s.Begin()
	if err != nil {
		return err
	}

	for _, closer := range closers {
		err = closer(s)
		if err != nil {
			_ = s.Rollback()
			return
		}
	}
	return s.Commit()
}

func (s *Session) Query(sqlstring string, args ...interface{}) error {
	err := s.IBinder.BindParse(s.IEngin.GetPrefix())
	if err != nil {
		return err
	}
	// 记录sqlLog
	s.lastSql = fmt.Sprint(sqlstring, ", ", args)
	if s.IfEnableSqlLog() {
		s.sqlLogs = append(s.sqlLogs, s.lastSql)
	}

	stmt, err := s.slave.db.Prepare(sqlstring)
	if err != nil {
		return err
	}

	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return err
	}

	// make sure we always close rows
	defer rows.Close()

	return s.scan(rows)
}

func (s *Session) Execute(sqlstring string, args ...interface{}) (rowsAffected int64, err error) {
	err = s.IBinder.BindParse(s.IEngin.GetPrefix())
	if err != nil {
		return
	}
	s.lastSql = fmt.Sprintf(sqlstring, args...)
	// 记录sqlLog
	if s.IfEnableSqlLog() {
		s.sqlLogs = append(s.sqlLogs, s.lastSql)
	}

	var operType = strings.ToLower(sqlstring[0:6])
	if operType == "select" {
		return 0, errors.New("Execute does not allow select operations, please use Query")
	}

	var stmt *sql.Stmt
	if s.master.tx == nil {
		stmt, err = s.master.db.Prepare(sqlstring)
	} else {
		stmt, err = s.master.tx.Prepare(sqlstring)
	}

	if err != nil {
		return 0, err
	}
	//return dba.parseExecute(stmt, operType, vals)

	//var err error
	defer stmt.Close()
	result, errs := stmt.Exec(args...)
	if errs != nil {
		return 0, errs
	}

	if operType == "insert" {
		// get last insert id
		lastInsertId, err := result.LastInsertId()
		if err == nil {
			s.lastInsertId = lastInsertId
		}
	}
	// get rows affected
	rowsAffected, err = result.RowsAffected()

	//// 如果是事务, 则重置所有参数
	//if dba.Strans == true {
	//	dba.Reset("transaction")
	//}

	//// 持久化日志
	//if dba.Connection.Logger != nil {
	//	dba.Connection.Logger.Write(dba.lastSql, time.Since(t_start).String(), time.Now().Format("2006-01-02 15:04:05"))
	//}

	return rowsAffected, err
}
func (s *Session) LastInsertId() int64 {
	return s.lastInsertId
}
func (s *Session) LastSql() string {
	return s.lastSql
}

func (s *Session) scan(rows *sql.Rows) (err error) {
	// 检查实多维数组还是一维数组
	switch s.GetBindType() {
	case OBJECT_STRUCT:
		err = s.scanRow(rows, s.GetBindOrigin())
	case OBJECT_STRUCT_SLICE:
		err = s.scanAll(rows, s.GetBindResultSlice())
	//case OBJECT_MAP, OBJECT_MAP_T:
	//	err = s.scanMap(rows, s.GetBindResult())
	case OBJECT_MAP, OBJECT_MAP_T, OBJECT_MAP_SLICE, OBJECT_MAP_SLICE_T:
		err = s.scanMapAll(rows, s.GetBindResultSlice())
	default:
		err = errors.New("Bind value error")
	}
	return
}

//func (s *Session) scanMap(rows *sql.Rows, dst interface{}) (err error) {
//	return s.scanMapAll(rows, dst)
//}

func (s *Session) scanMapAll(rows *sql.Rows, dst interface{}) (err error) {
	//var result = make([]map[string]interface{}, 0)
	var columns []string
	if columns, err = rows.Columns(); err != nil {
		return
	}

	count := len(columns)
	for rows.Next() {
		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := 0; i < count; i++ {
			scanArgs[i] = &values[i]
		}
		_ = rows.Scan(scanArgs...)
		//entry := make(map[string]interface{})
		var bindResultTmp = reflect.MakeMap(s.GetBindResult().Type())
		//sliceVal := reflect.Indirect(dstVal)
		for i, col := range columns {
			var v interface{}
			val := values[i]
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}
			// 是否union操作, 不是的话,就绑定数据
			if s.GetUnion()!=nil {
				s.union = v
				return
			}
			switch s.GetBindType() {
			case OBJECT_MAP_T, OBJECT_MAP_SLICE_T:
				s.GetBindResult().SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(t.New(v)))
				if s.GetBindType() == OBJECT_MAP_SLICE || s.GetBindType() == OBJECT_MAP_SLICE_T {
					bindResultTmp.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(t.New(v)))
				}
			default:
				s.GetBindResult().SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(v))
				if s.GetBindType() == OBJECT_MAP_SLICE || s.GetBindType() == OBJECT_MAP_SLICE_T {
					bindResultTmp.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(v))
				}
			}
		}
		//result = append(result, entry)
		if s.GetBindType() == OBJECT_MAP_SLICE || s.GetBindType() == OBJECT_MAP_SLICE_T {
			s.GetBindResultSlice().Set(reflect.Append(s.GetBindResultSlice(), bindResultTmp))
		}
	}
	return
}

func (s *Session) SetUnion(u interface{})  {
	s.union = u
}

func (s *Session) GetUnion() interface{} {
	return s.union
}

// scan a single row of data into a struct.
func (s *Session) scanRow(rows *sql.Rows, dst interface{}) error {
	// check if there is data waiting
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	// get a list of targets
	var fields = strutForScan(dst)

	// perform the scan
	if err := rows.Scan(fields...); err != nil {
		//if err := rows.Scan(strutForScan(s.BindResult.Interface())...); err != nil {
		return err
	}

	return rows.Err()
}

// ScanAll scans all sql result rows into a slice of structs.
// It reads all rows and closes rows when finished.
// dst should be a pointer to a slice of the appropriate type.
// The new results will be appended to any existing data in dst.
func (s *Session) scanAll(rows *sql.Rows, dst interface{}) error {
	for rows.Next() {
		// scan it
		err := rows.Scan(strutForScan(s.GetBindResult().Interface())...)
		if err != nil {
			return err
		}
		// add to the result slice
		s.GetBindResultSlice().Set(reflect.Append(s.GetBindResultSlice(), s.GetBindResult().Elem()))
	}

	return rows.Err()
}

func (s *Session) GetBinder() IBinder {
	return s.IBinder
}
