package gorose

import (
	"errors"
	"fmt"
	"reflect"
)

type TypeToSqlIncDec struct {
	Symbol string
	Data   map[string]any
}
type TypeToSqlExists struct {
	Bindings []any
}
type TypeToSqlAggregate struct {
	CallFunc, Column string
}

type TypeToSqlInsertCase struct {
	IgnoreCase                  string
	OnDuplicateKeys, MustFields []string
}
type TypeToSqlPessimisticLocking struct {
	LockType string
}

func (db *Database) ToSqlSelect() (sql4prepare string, binds []any) {
	return db.Driver.ToSqlSelect(db.Context)
}

func (db *Database) ToSqlTable() (sql4prepare string, values []any, err error) {
	return db.Driver.ToSqlTable(db.Context)
}
func (db *Database) ToSqlJoin() (sql4prepare string, binds []any, err error) {
	return db.Driver.ToSqlJoin(db.Context)
}

func (db *Database) ToSqlWhere() (sql4prepare string, values []any, err error) {
	return db.Driver.ToSqlWhere(db.Context)
}

func (db *Database) ToSqlOrderBy() (sql4prepare string) {
	return db.Driver.ToSqlOrderBy(db.Context)
}

func (db *Database) ToSqlLimitOffset() (sqlSegment string, binds []any) {
	return db.Driver.ToSqlLimitOffset(db.Context)
}

func (db *Database) ToSql() (sql4prepare string, values []any, err error) {
	return db.Driver.ToSql(db.Context)
}
func (db *Database) ToSqlIncDec(symbol string, data map[string]any) (sql4prepare string, values []any, err error) {
	return db.Driver.ToSqlIncDec(db.Context, symbol, data)
}
func (db *Database) ToSqlExists(bind ...any) (sql4prepare string, values []any, err error) {
	if len(bind) > 0 {
		sql4prepare, values, err = db.Table(bind[0]).Driver.ToSql(db.Context)
	} else {
		sql4prepare, values, err = db.Driver.ToSql(db.Context)
	}
	if err != nil {
		return
	}
	sql4prepare = fmt.Sprintf("SELECT EXISTS (%s) AS exist", sql4prepare)
	return
}

func (db *Database) ToSqlAggregate(function, column string) (sql4prepare string, values []any, err error) {
	var ctx = *db.Context
	ctx.SelectClause.Columns = append(ctx.SelectClause.Columns, Column{
		Name:  fmt.Sprintf("%s(%s)", function, column),
		Alias: function,
		IsRaw: true,
		Binds: []any{},
	})
	return db.Driver.ToSql(&ctx)
}

func (db *Database) ToSqlTo(obj any, mustFields ...string) (sql4prepare string, binds []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	columns, fieldStruct, _ := StructsParse(obj)
	switch rfv.Kind() {
	case reflect.Struct:
		var data = make(map[string]any)
		data, err = structDataToMap(rfv, columns, fieldStruct, mustFields...)
		if err != nil {
			return
		}
		sql4prepare, binds, err = db.Table(obj).Select(columns...).Where(data).Limit(1).ToSql()
	case reflect.Slice:
		if rfv.Type().Elem().Kind() == reflect.Struct {
			sql4prepare, binds, err = db.Table(obj).Select(columns...).ToSql()
		}
	default:
		err = errors.New("obj must be struct(slice) or map(slice)")
	}
	return
}

func (db *Database) ToSqlInsert(obj any, ignoreCase string, onDuplicateKeys []string, mustFields ...string) (sqlSegment string, binds []any, err error) {
	return db.Driver.ToSqlInsert(db.Context, obj, ignoreCase, onDuplicateKeys, mustFields...)
}

func (db *Database) ToSqlUpdate(obj any, mustFields ...string) (sqlSegment string, binds []any, err error) {
	return db.Driver.ToSqlUpdate(db.Context, obj, mustFields...)
}

func (db *Database) ToSqlDelete(obj any) (sqlSegment string, binds []any, err error) {
	return db.Driver.ToSqlDelete(db.Context, obj)
}
