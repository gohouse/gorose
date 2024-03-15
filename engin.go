package gorose

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"reflect"
)

type SqlItem struct {
	Sql      string
	Bindings []any
	Err      error
}
type Engin struct {
	*GoRose
	tx            *sql.Tx
	autoSavePoint uint8
	lastSql       SqlItem
}

func NewEngin(g *GoRose) *Engin {
	return &Engin{GoRose: g}
}

func (s *Engin) LastSql() SqlItem {
	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		return SqlItem{Err: errors.New("only record when slog level in debug mod")}
	}
	return s.lastSql
}

func (s *Engin) Log(sqls string, bindings ...any) {
	slog.With("bindings", bindings).Debug(sqls)
	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		s.lastSql = SqlItem{Sql: sqls, Bindings: bindings}
	}
}

func (s *Engin) execute(query string, args ...any) (int64, error) {
	exec, err := s.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}
func (s *Engin) Exec(query string, args ...any) (sql.Result, error) {
	slog.With("bindings", args).Debug(query)
	if s.tx != nil {
		return s.tx.Exec(query, args...)
	}
	return s.MasterDB().Exec(query, args...)
}
func (s *Engin) Begin() (err error) {
	if s.tx != nil {
		s.autoSavePoint += 1
		return s.SavePoint(s.autoSavePoint)
	}
	s.tx, err = s.MasterDB().Begin()
	return
}
func (s *Engin) SavePoint(name any) (err error) {
	_, err = s.tx.Exec("SAVEPOINT ?", name)
	return
}
func (s *Engin) RollbackTo(name any) (err error) {
	_, err = s.tx.Exec("ROLLBACK TO ?", name)
	return
}
func (s *Engin) Rollback() (err error) {
	if s.autoSavePoint > 0 {
		// decrease in advance whether rollbackTo fail
		currentPoint := s.autoSavePoint
		s.autoSavePoint -= 1
		return s.RollbackTo(currentPoint)
	}
	err = s.tx.Rollback()
	if err != nil {
		return
	}
	s.tx = nil
	return
}
func (s *Engin) Commit() (err error) {
	if s.autoSavePoint > 0 {
		s.autoSavePoint -= 1
		return
	}
	err = s.tx.Commit()
	if err != nil {
		return
	}
	s.tx = nil
	return
}
func (s *Engin) Transaction(closure ...func(*Engin) error) (err error) {
	if err = s.Begin(); err != nil {
		return
	}
	for _, v := range closure {
		err = v(s)
		if err != nil {
			return s.Rollback()
		}
	}
	return s.Commit()
}

func (s *Engin) Query(query string, args ...any) (rows *sql.Rows, err error) {
	slog.With("bindings", args).Debug(query)
	if s.tx != nil {
		return s.tx.Query(query, args...)
	} else {
		return s.SlaveDB().Query(query, args...)
	}
}

func (s *Engin) QueryRow(query string, args ...any) *sql.Row {
	slog.With("bindings", args).Debug(query)
	if s.tx != nil {
		return s.tx.QueryRow(query, args...)
	} else {
		return s.SlaveDB().QueryRow(query, args...)
	}
}
func (s *Engin) QueryTo(bind any, query string, args ...any) (err error) {
	var rows *sql.Rows
	if rows, err = s.Query(query, args...); err != nil {
		return
	}
	return s.rowsToBind(rows, bind)
}
func (s *Engin) rowsToBind(rows *sql.Rows, bind any) (err error) {
	rfv := reflect.Indirect(reflect.ValueOf(bind))
	switch rfv.Kind() {
	case reflect.Slice:
		switch rfv.Type().Elem().Kind() {
		case reflect.Map:
			return s.rowsToMap(rows, rfv)
		case reflect.Struct:
			return s.rowsToStruct(rows, rfv)
		default:
			return errors.New("only struct(slice) or map(slice) supported")
		}
	case reflect.Map:
		return s.rowsToMap(rows, rfv)
	case reflect.Struct:
		return s.rowsToStruct(rows, rfv)
	default:
		return errors.New("only struct(slice) or map(slice) supported")
	}
}

func (s *Engin) rowsToStruct(rows *sql.Rows, rfv reflect.Value) error {
	//FieldTag, FieldStruct, _ := structsParse(rfv)
	FieldTag, FieldStruct, _ := structsTypeParse(rfv.Type())

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// 列的个数
	count := len(columns)

	for rows.Next() {
		// 要先扫描到map, 再做字段比对, 因为这里不确定具体字段数量
		// 主要针对 select * 或者直接sql语句
		//todo 如果是由struct转换而来, 可以新开一个方法, 不需要做转换比对过程
		//entry, err := s.rowsToMapSingle(rows, columns, count)
		//if err != nil {
		//	return err
		//}

		// 一条数据的各列的值（需要指定长度为列的个数，以便获取地址）
		values := make([]any, count)
		// 一条数据的各列的值的地址
		valPointers := make([]any, count)
		// 获取各列的值的地址
		for i := 0; i < count; i++ {
			valPointers[i] = &values[i]
		}
		// 获取各列的值，放到对应的地址中
		err = rows.Scan(valPointers...)
		if err != nil {
			return err
		}

		if rfv.Kind() == reflect.Slice {
			rfvItem := reflect.Indirect(reflect.New(rfv.Type().Elem()))
			for i, _ := range FieldTag {
				b, ok := values[i].([]byte)
				if ok {
					// 字符切片转为字符串
					rfvItem.FieldByName(FieldStruct[i]).Set(reflect.ValueOf(string(b)))
				} else {
					rfvItem.FieldByName(FieldStruct[i]).Set(reflect.ValueOf(values[i]))
				}
			}
			rfv.Set(reflect.Append(rfv, rfvItem))
		} else {
			for i := range FieldTag {
				b, ok := values[i].([]byte)
				if ok {
					// 字符切片转为字符串
					rfv.FieldByName(FieldStruct[i]).Set(reflect.ValueOf(string(b)))
				} else {
					rfv.FieldByName(FieldStruct[i]).Set(reflect.ValueOf(values[i]))
				}
			}
			//rfv.Set(entry)
		}
	}

	return nil
}

func (s *Engin) rowsToMapSingle(rows *sql.Rows, columns []string, count int) (entry map[string]any, err error) {
	// 一条数据的各列的值（需要指定长度为列的个数，以便获取地址）
	values := make([]any, count)
	// 一条数据的各列的值的地址
	valPointers := make([]any, count)
	// 获取各列的值的地址
	for i := 0; i < count; i++ {
		valPointers[i] = &values[i]
	}
	// 获取各列的值，放到对应的地址中
	err = rows.Scan(valPointers...)
	if err != nil {
		return
	}
	// 一条数据的Map (列名和值的键值对)
	entry = make(map[string]any)

	// Map 赋值
	for i, col := range columns {
		var v any
		// 值复制给val(所以Scan时指定的地址可重复使用)
		val := values[i]
		b, ok := val.([]byte)
		if ok {
			// 字符切片转为字符串
			v = string(b)
		} else {
			v = val
		}
		entry[col] = v
	}
	return
}

func (s *Engin) rowsToMap(rows *sql.Rows, rfv reflect.Value) error {
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// 列的个数
	count := len(columns)

	for rows.Next() {
		entry, err := s.rowsToMapSingle(rows, columns, count)
		if err != nil {
			return err
		}
		if rfv.Kind() == reflect.Slice {
			rfv.Set(reflect.Append(rfv, reflect.ValueOf(entry)))
		} else {
			rfv.Set(reflect.ValueOf(entry))
		}
	}
	return nil
}
