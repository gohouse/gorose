package gorose

import (
	"github.com/gohouse/t"
	"reflect"
)

// Get : select more rows , relation limit set
func (dba *Orm) Get() error {
	// 构建sql
	sqlStr, args, err := dba.BuildSql()
	if err != nil {
		return err
	}

	// 执行查询
	return dba.ISession.Query(sqlStr, args...)
}
// Get : select more rows , relation limit set
func (dba *Orm) Value(field string) (v t.T, err error) {
	err = dba.Get()
	if err!=nil {
		return
	}
	var binder = dba.ISession.GetBinder()
	switch binder.GetBindType() {
	case OBJECT_MAP,OBJECT_MAP_SLICE,OBJECT_MAP_SLICE_T,OBJECT_MAP_T:
		v = t.New(binder.GetBindResult().MapIndex(reflect.ValueOf(field)).Interface())
	case OBJECT_STRUCT,OBJECT_STRUCT_SLICE:
		v = t.New(binder.GetBindResult().FieldByName(field).Interface())
	}
	return
}
func (dba *Orm) Pluck(field string, fieldKey ...string) (v t.T, err error) {
	err = dba.Get()
	if err!=nil {
		return
	}
	var binder = dba.ISession.GetBinder()
	switch binder.GetBindType() {
	case OBJECT_MAP,OBJECT_MAP_T,OBJECT_STRUCT:	// row
	var key,val t.T
	if len(fieldKey)>0 {
		key,err = dba.Value(fieldKey[0])
		if err!=nil {
			return
		}
		val,err = dba.Value(field)
		if err!=nil {
			return
		}
		v = t.New(map[t.T]t.T{key:val})
	} else {
		v,err = dba.Value(field)
		if err!=nil {
			return
		}
	}
	case OBJECT_MAP_SLICE,OBJECT_MAP_SLICE_T:
		var resMap = make(t.Map)
		var resSlice = make(t.Slice,0)
		for _,item := range t.New(binder.GetBindResultSlice().Interface()).Slice() {
			val := item.MapString()
			if len(fieldKey)>0 {
				resMap[val[fieldKey[0]]] = val[field]
				v = t.New(resMap)
			} else {
				resSlice = append(resSlice, val[field])
				v = t.New(resSlice)
			}
		}
	case OBJECT_STRUCT_SLICE:	// rows
	}
	return
}
// Get : select more rows , relation limit set
func (dba *Orm) Paginate() error {
	// 构建sql
	sqlStr, args, err := dba.BuildSql()
	if err != nil {
		return err
	}

	// 执行查询
	return dba.ISession.Query(sqlStr, args...)
}
