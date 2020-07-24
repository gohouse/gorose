package gorose

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gohouse/t"
	"reflect"
	"strings"
	"time"
)

// Session ...
type Session struct {
	IEngin
	IBinder
	master       *sql.DB
	tx           *sql.Tx
	slave        *sql.DB
	lastInsertId int64
	sqlLogs      []string
	lastSql      string
	union        interface{}
	transaction  bool
	err          error
}

var _ ISession = (*Session)(nil)

// NewSession : 初始化 Session
func NewSession(e IEngin) *Session {
	var s = new(Session)
	s.IEngin = e
	// 初始化 IBinder
	s.SetIBinder(NewBinder())

	s.master = e.GetExecuteDB()
	s.slave = e.GetQueryDB()

	return s
}

func (s *Session) Close() {
	s.master.Close()
	s.slave.Close()
}

// GetIEngin 获取engin
func (s *Session) GetIEngin() IEngin {
	return s.IEngin
}

// GetDriver 获取驱动
func (s *Session) SetIEngin(ie IEngin) {
	s.IEngin = ie
}

// Bind : 传入绑定结果的对象, 参数一为对象, 可以是 struct, gorose.MapRow 或对应的切片
//		如果是做非query操作,第一个参数也可以仅仅指定为字符串表名
func (s *Session) Bind(tab interface{}) ISession {
	//fmt.Println(tab, NewBinder(tab))
	//s.SetIBinder(NewBinder(tab))
	s.GetIBinder().SetBindOrigin(tab)
	s.err = s.IBinder.BindParse(s.GetIEngin().GetPrefix())
	return s
}

// GetBinder 获取绑定对象
func (s *Session) GetErr() error {
	return s.err
}

// GetBinder 获取绑定对象
func (s *Session) SetIBinder(ib IBinder) {
	s.IBinder = ib
}

// GetBinder 获取绑定对象
func (s *Session) GetIBinder() IBinder {
	return s.IBinder
}

// GetBinder 获取绑定对象
func (s *Session) ResetBinderResult() {
	_ = s.IBinder.BindParse(s.GetIEngin().GetPrefix())
}

// GetTableName 获取解析后的名字, 提供给orm使用
// 为什么要在这里重复添加该方法, 而不是直接继承 IBinder 的方法呢?
// 是因为, 这里涉及到表前缀的问题, 只能通过session来传递, 所以IOrm就可以选择直接继承
func (s *Session) GetTableName() (string, error) {
	//err := s.IBinder.BindParse(s.GetIEngin().GetPrefix())
	//fmt.Println(s.GetIBinder())
	return s.GetIBinder().GetBindName(), s.err
}

// Begin ...
func (s *Session) Begin() (err error) {
	s.tx, err = s.master.Begin()
	s.SetTransaction(true)
	return
}

// Rollback ...
func (s *Session) Rollback() (err error) {
	err = s.tx.Rollback()
	s.tx = nil
	s.SetTransaction(false)
	return
}

// Commit ...
func (s *Session) Commit() (err error) {
	err = s.tx.Commit()
	s.tx = nil
	s.SetTransaction(false)
	return
}

// Transaction ...
func (s *Session) Transaction(closers ...func(ses ISession) error) (err error) {
	err = s.Begin()
	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return err
	}

	for _, closer := range closers {
		err = closer(s)
		if err != nil {
			s.GetIEngin().GetLogger().Error(err.Error())
			_ = s.Rollback()
			return
		}
	}
	return s.Commit()
}

// Query ...
func (s *Session) Query(sqlstring string, args ...interface{}) (result []Data, err error) {
	// 记录开始时间
	start := time.Now()
	//withRunTimeContext(func() {
	if s.err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		err = s.err
	}
	// 记录sqlLog
	s.lastSql = fmt.Sprint(sqlstring, ", ", args)
	//if s.IfEnableSqlLog() {
	//	s.sqlLogs = append(s.sqlLogs, s.lastSql)
	//}

	var stmt *sql.Stmt
	// 如果是事务, 则从主库中读写
	if s.tx == nil {
		stmt, err = s.slave.Prepare(sqlstring)
	} else {
		stmt, err = s.tx.Prepare(sqlstring)
	}

	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}

	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}

	// make sure we always close rows
	defer rows.Close()

	err = s.scan(rows)
	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}
	//}, func(duration time.Duration) {
	//	//if duration.Seconds() > 1 {
	//	//	s.GetIEngin().GetLogger().Slow(s.LastSql(), duration)
	//	//} else {
	//	//	s.GetIEngin().GetLogger().Sql(s.LastSql(), duration)
	//	//}
	//})

	timeduration := time.Since(start)
	//if timeduration.Seconds() > 1 {
	s.GetIEngin().GetLogger().Slow(s.LastSql(), timeduration)
	//} else {
	s.GetIEngin().GetLogger().Sql(s.LastSql(), timeduration)
	//}

	result = s.GetIBinder().GetBindAll()
	return
}

// Execute ...
func (s *Session) Execute(sqlstring string, args ...interface{}) (rowsAffected int64, err error) {
	// 记录开始时间
	start := time.Now()
	//withRunTimeContext(func() {
	//	err = s.GetIBinder().BindParse(s.GetIEngin().GetPrefix())
	if s.err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}
	s.lastSql = fmt.Sprint(sqlstring, ", ", args)
	//// 记录sqlLog
	//if s.IfEnableSqlLog() {
	//	s.sqlLogs = append(s.sqlLogs, s.lastSql)
	//}

	var operType = strings.ToLower(sqlstring[0:6])
	if operType == "select" {
		s.GetIEngin().GetLogger().Error(err.Error())
		err = errors.New("Execute does not allow select operations, please use Query")
		return
	}

	var stmt *sql.Stmt
	if s.tx == nil {
		stmt, err = s.master.Prepare(sqlstring)
	} else {
		stmt, err = s.tx.Prepare(sqlstring)
	}

	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}

	//var err error
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}

	if operType == "insert" {
		// get last insert id
		lastInsertId, err := result.LastInsertId()
		if err == nil {
			s.lastInsertId = lastInsertId
		} else {
			s.GetIEngin().GetLogger().Error(err.Error())
		}
	}
	// get rows affected
	rowsAffected, err = result.RowsAffected()
	timeduration := time.Since(start)
	//}, func(duration time.Duration) {
	if timeduration.Seconds() > 1 {
		s.GetIEngin().GetLogger().Slow(s.LastSql(), timeduration)
	} else {
		s.GetIEngin().GetLogger().Sql(s.LastSql(), timeduration)
	}
	//})
	return
}

// LastInsertId ...
func (s *Session) LastInsertId() int64 {
	return s.lastInsertId
}

// LastSql ...
func (s *Session) LastSql() string {
	return s.lastSql
}

func (s *Session) scan(rows *sql.Rows) (err error) {
	// 如果不需要绑定, 则需要初始化一下binder
	if s.GetIBinder() == nil {
		s.SetIBinder(NewBinder())
	}
	// 检查实多维数组还是一维数组
	switch s.GetBindType() {
	case OBJECT_STRING:
		err = s.scanAll(rows)
	case OBJECT_STRUCT, OBJECT_STRUCT_SLICE:
		err = s.scanStructAll(rows)
	//case OBJECT_MAP, OBJECT_MAP_T:
	//	err = s.scanMap(rows, s.GetBindResult())
	case OBJECT_MAP, OBJECT_MAP_T, OBJECT_MAP_SLICE, OBJECT_MAP_SLICE_T:
		err = s.scanMapAll(rows)
	case OBJECT_NIL:
		err = s.scanAll(rows)
	default:
		err = errors.New("Bind value error")
	}
	return
}

//func (s *Session) scanMap(rows *sql.Rows, dst interface{}) (err error) {
//	return s.scanMapAll(rows, dst)
//}

func (s *Session) scanMapAll(rows *sql.Rows) (err error) {
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
		//fmt.Println(reflect.TypeOf(s.GetBindResult()).Kind())
		var bindResultTmp = reflect.MakeMap(reflect.Indirect(reflect.ValueOf(s.GetBindResult())).Type())
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
				br := reflect.Indirect(reflect.ValueOf(s.GetBindResult()))
				switch s.GetBindType() {
				case OBJECT_MAP_T, OBJECT_MAP_SLICE_T: // t.T类型
					// 绑定到一条数据结果对象上,方便其他地方的调用,永远存储最新一条
					br.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(t.New(v)))
					// 跟上一行干的事是一样的, 只不过防止上一行的数据被后续的数据改变, 而无法提供给下边多条数据报错的需要
					if s.GetBindType() == OBJECT_MAP_SLICE || s.GetBindType() == OBJECT_MAP_SLICE_T {
						bindResultTmp.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(t.New(v)))
					}
				default: // 普通类型map[string]interface{}, 具体代码注释参照 上一个 case
					br.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(v))
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

// ScanAll scans all sql result rows into a slice of structs.
// It reads all rows and closes rows when finished.
// dst should be a pointer to a slice of the appropriate type.
// The new results will be appended to any existing data in dst.
func (s *Session) scanStructAll(rows *sql.Rows) error {
	// check if there is data waiting
	//if !rows.Next() {
	//	if err := rows.Err(); err != nil {
	//		s.GetIEngin().GetLogger().Error(err.Error())
	//		return err
	//	}
	//	return sql.ErrNoRows
	//}
	var sfs = structForScan(s.GetBindResult())
	for rows.Next() {
		if s.GetUnion() != nil {
			var union interface{}
			err := rows.Scan(&union)
			if err != nil {
				s.GetIEngin().GetLogger().Error(err.Error())
				return err
			}
			s.union = union
			return err
		}
		// scan it
		//fmt.Printf("%#v \n",structForScan(s.GetBindResult()))
		err := rows.Scan(sfs...)
		if err != nil {
			s.GetIEngin().GetLogger().Error(err.Error())
			return err
		}
		// 如果是union操作就不需要绑定数据直接返回, 否则就绑定数据
		if s.GetUnion() == nil {
			// 如果是多条数据集, 就插入到对应的结果集slice上
			if s.GetBindType() == OBJECT_STRUCT_SLICE {
				// add to the result slice
				s.GetBindResultSlice().Set(reflect.Append(s.GetBindResultSlice(),
					reflect.Indirect(reflect.ValueOf(s.GetBindResult()))))
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

	var result = []Data{}
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
		var resultTmp = Data{}
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
			if s.GetUnion() != nil {
				s.union = v
				return
				// 以下上通用解决方法
				//unionTmp[col] = v
				//s.union = unionTmp
			}
			resultTmp[col] = v
		}
		result = append(result, resultTmp)
	}
	s.IBinder.SetBindAll(result)
	return
}

// SetUnion ...
func (s *Session) SetUnion(u interface{}) {
	s.union = u
}

// GetUnion ...
func (s *Session) GetUnion() interface{} {
	return s.union
}

// SetTransaction ...
func (s *Session) SetTransaction(b bool) {
	s.transaction = b
}

// GetTransaction 提供给 orm 使用的, 方便reset操作
func (s *Session) GetTransaction() bool {
	return s.transaction
}
