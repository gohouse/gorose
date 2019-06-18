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
	IBinder
	master       dbObject
	slave        dbObject
	lastInsertId int64
	sqlLogs      []string
	lastSql      string
	union        interface{}
	transaction  bool
	err          error
}

var _ ISession = &Session{}

// NewSession : 初始化 Session
func NewSession(e IEngin) ISession {
	var s = new(Session)
	s.IEngin = e
	//s.IBinder = b

	s.master = s.IEngin.GetExecuteDB()
	s.slave = s.IEngin.GetQueryDB()

	return s
}

// Close : 关闭 Session
func (s *Session) Close() {
	s = nil
}

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
	s.IBinder = NewBinder(tab)
	s.err = s.IBinder.BindParse(s.IEngin.GetPrefix())
	return s
}

// GetBinder 获取绑定对象
func (s *Session) GetErr() error {
	return s.err
}

// GetBinder 获取绑定对象
func (s *Session) GetBinder() IBinder {
	return s.IBinder
}

//// GetBinder 获取绑定对象
//func (s *Session) SetBinder() ISession {
//	s.Bind(s.IBinder.GetBindOrigin())
//	return s
//}

//// GetBinder 获取绑定对象
//func (s *Session) ResetBinder() {
//	var origin = s.IBinder.GetBindOrigin()
//	s.IBinder = NewBinder(origin)
//}

// GetBinder 获取绑定对象
func (s *Session) ResetBinderResult() {
	_ = s.IBinder.BindParse(s.IEngin.GetPrefix())
}

// GetTableName 获取解析后的名字, 提供给orm使用
// 为什么要在这里重复添加该方法, 而不是直接继承 IBinder 的方法呢?
// 是因为, 这里涉及到表前缀的问题, 只能通过session来传递, 所以IOrm就可以选择直接继承
func (s *Session) GetTableName() (string, error) {
	//err := s.IBinder.BindParse(s.IEngin.GetPrefix())
	return s.IBinder.GetBindName(), s.err
}

func (s *Session) Begin() (err error) {
	s.master.tx, err = s.master.db.Begin()
	s.SetTransaction(true)
	return
}

func (s *Session) Rollback() (err error) {
	err = s.master.tx.Rollback()
	s.master.tx = nil
	s.SetTransaction(false)
	return
}

func (s *Session) Commit() (err error) {
	err = s.master.tx.Commit()
	s.master.tx = nil
	s.SetTransaction(false)
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
	//err := s.IBinder.BindParse(s.IEngin.GetPrefix())
	if s.err != nil {
		return s.err
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
	s.lastSql = fmt.Sprint(sqlstring, ", ", args)
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
	case OBJECT_STRING:
		err = s.scanAll(rows)
	case OBJECT_STRUCT, OBJECT_STRUCT_SLICE:
		err = s.scanStructAll(rows, s.GetBindResultSlice())
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
	var columns []string
	// 获取查询的所有字段
	if columns, err = rows.Columns(); err != nil {
		return
	}
	count := len(columns)

	for rows.Next() {
		// 定义要绑定的结果集
		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := 0; i < count; i++ {
			scanArgs[i] = &values[i]
		}
		// 获取结果
		_ = rows.Scan(scanArgs...)

		// 定义预设的绑定对象
		var bindResultTmp = reflect.MakeMap(s.GetBindResult().Type())
		//// 定义union操作的map返回
		//var unionTmp = map[string]interface{}{}
		for i, col := range columns {
			var v interface{}
			val := values[i]
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}
			// 如果是union操作就不需要绑定数据直接返回, 否则就绑定数据
			//TODO 这里可能有点问题, 比如在group时, 返回的结果不止一条, 这里直接返回的就是第一条
			// 默认其实只是取了第一条, 满足常规的 union 操作(count,sum,max,min,avg)而已
			// 后边需要再行完善, 以便group时使用
			// 具体完善方法: 就是这里断点去掉, 不直接绑定union, 新增一个map,将结果放在map中,在方法最后统一返回
			if s.GetUnion() != nil {
				s.union = v
				return
				// 以下上通用解决方法
				//unionTmp[col] = v
				//s.union = unionTmp
			} else {
				switch s.GetBindType() {
				case OBJECT_MAP_T, OBJECT_MAP_SLICE_T: // t.T类型
					// 绑定到一条数据结果对象上,方便其他地方的调用,永远存储最新一条
					s.GetBindResult().SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(t.New(v)))
					// 跟上一行干的事是一样的, 只不过防止上一行的数据被后续的数据改变, 而无法提供给下边多条数据报错的需要
					if s.GetBindType() == OBJECT_MAP_SLICE || s.GetBindType() == OBJECT_MAP_SLICE_T {
						bindResultTmp.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(t.New(v)))
					}
				default: // 普通类型map[string]interface{}, 具体代码注释参照 上一个 case
					s.GetBindResult().SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(v))
					if s.GetBindType() == OBJECT_MAP_SLICE || s.GetBindType() == OBJECT_MAP_SLICE_T {
						bindResultTmp.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(v))
					}
				}
			}
		}
		// 如果是union操作就不需要绑定数据直接返回, 否则就绑定数据
		if s.GetUnion() == nil {
			// 如果是多条数据集, 就插入到对应的结果集slice上
			if s.GetBindType() == OBJECT_MAP_SLICE || s.GetBindType() == OBJECT_MAP_SLICE_T {
				s.GetBindResultSlice().Set(reflect.Append(s.GetBindResultSlice(), bindResultTmp))
			}
		}
	}
	return
}

//// scan a single row of data into a struct.
//func (s *Session) scanRow(rows *sql.Rows, dst interface{}) error {
//	// check if there is data waiting
//	if !rows.Next() {
//		if err := rows.Err(); err != nil {
//			return err
//		}
//		return sql.ErrNoRows
//	}
//
//	// get a list of targets
//	var fields = strutForScan(dst)
//
//	// perform the scan
//	if err := rows.Scan(fields...); err != nil {
//		//if err := rows.Scan(strutForScan(s.BindResult.Interface())...); err != nil {
//		return err
//	}
//
//	return rows.Err()
//}

// ScanAll scans all sql result rows into a slice of structs.
// It reads all rows and closes rows when finished.
// dst should be a pointer to a slice of the appropriate type.
// The new results will be appended to any existing data in dst.
func (s *Session) scanStructAll(rows *sql.Rows, dst interface{}) error {
	// check if there is data waiting
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}
	for rows.Next() {
		// scan it
		err := rows.Scan(strutForScan(s.GetBindResult().Interface())...)
		if err != nil {
			return err
		}

		// 如果是union操作就不需要绑定数据直接返回, 否则就绑定数据
		if s.GetUnion() == nil {
			// 如果是多条数据集, 就插入到对应的结果集slice上
			if s.GetBindType() == OBJECT_STRUCT_SLICE {
				// add to the result slice
				s.GetBindResultSlice().Set(reflect.Append(s.GetBindResultSlice(), s.GetBindResult().Elem()))
			}
		}
	}

	return rows.Err()
}

func (s *Session) scanAll(rows *sql.Rows) (err error) {
	var columns []string
	// 获取查询的所有字段
	if columns, err = rows.Columns(); err != nil {
		return
	}
	count := len(columns)

	var result = []Map{}
	for rows.Next() {
		// 定义要绑定的结果集
		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := 0; i < count; i++ {
			scanArgs[i] = &values[i]
		}
		// 获取结果
		_ = rows.Scan(scanArgs...)

		// 定义预设的绑定对象
		var resultTmp = Map{}
		//// 定义union操作的map返回
		//var unionTmp = map[string]interface{}{}
		for i, col := range columns {
			var v interface{}
			val := values[i]
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}
			resultTmp[col] = t.New(v)
		}
		result = append(result, resultTmp)
	}
	s.IBinder.SetBindAll(result)
	return
}

func (s *Session) SetUnion(u interface{}) {
	s.union = u
}

func (s *Session) GetUnion() interface{} {
	return s.union
}

func (s *Session) SetTransaction(b bool) {
	s.transaction = b
}

// GetTransaction 提供给 orm 使用的, 方便reset操作
func (s *Session) GetTransaction() bool {
	return s.transaction
}
