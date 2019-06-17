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

// Count : select count rows
func (dba *Orm) Count(args ...string) (int64, error) {
	fields := "*"
	if len(args) > 0 {
		fields = args[0]
	}
	count, err := dba._unionBuild("count", fields)
	if count == nil {
		count = int64(0)
	}
	return count.(int64), err
}

// Sum : select sum field
func (dba *Orm) Sum(sum string) (interface{}, error) {
	return dba._unionBuild("sum", sum)
}

// Avg : select avg field
func (dba *Orm) Avg(avg string) (interface{}, error) {
	return dba._unionBuild("avg", avg)
}

// Max : select max field
func (dba *Orm) Max(max string) (interface{}, error) {
	return dba._unionBuild("max", max)
}

// Min : select min field
func (dba *Orm) Min(min string) (interface{}, error) {
	return dba._unionBuild("min", min)
}

// _unionBuild : build union select real
func (dba *Orm) _unionBuild(union, field string) (interface{}, error) {
	var tmp interface{}

	dba.union = union + "(" + field + ") as " + union
	// 缓存fields字段,暂时由union占用
	fieldsTmp := dba.fields
	dba.fields = []string{dba.union}
	dba.ISession.SetUnion(true)

	// 构建sql
	sqls, args, err := dba.BuildSql()
	if err != nil {
		return tmp, err
	}

	// 执行查询
	err = dba.ISession.Query(sqls, args...)
	if err != nil {
		return tmp, err
	}

	// 重置union, 防止复用的时候感染
	dba.union = ""
	// 返还fields
	dba.fields = fieldsTmp

	// 语法糖获取union值
	if dba.ISession.GetUnion() != nil {
		tmp = dba.ISession.GetUnion()
		dba.ISession.SetUnion(nil)
	}

	return tmp, nil
}

// Get : select more rows , relation limit set
func (dba *Orm) Value(field string) (v t.T, err error) {
	dba.Limit(1)
	err = dba.Get()
	if err != nil {
		return
	}
	var binder = dba.ISession.GetBinder()
	switch binder.GetBindType() {
	case OBJECT_MAP, OBJECT_MAP_SLICE, OBJECT_MAP_SLICE_T, OBJECT_MAP_T:
		//fmt.Println(binder.GetBindResult(), binder.GetBindResult().Type())
		v = t.New(binder.GetBindResult().MapIndex(reflect.ValueOf(field)).Interface())
	case OBJECT_STRUCT, OBJECT_STRUCT_SLICE:
		//v = t.New(reflect.Indirect(binder.GetBindResult()).FieldByName(field).Interface())
		bindResult := reflect.Indirect(binder.GetBindResult())
		//ostype := os.Type()
		//for i := 0; i < ostype.NumField(); i++ {
		//	tag := ostype.Field(i).Tag.Get(TAGNAME)
		//	if tag == field || ostype.Field(i).Name == field {
		//		v = t.New(os.FieldByName(ostype.Field(i).Name))
		//		return
		//	}
		//}
		v = dba._valueFromStruct(bindResult, field)
	}
	return
}
func (dba *Orm) _valueFromStruct(bindResult reflect.Value, field string) (v t.T) {
	//os := val
	ostype := bindResult.Type()
	for i := 0; i < ostype.NumField(); i++ {
		tag := ostype.Field(i).Tag.Get(TAGNAME)
		if tag == field || ostype.Field(i).Name == field {
			v = t.New(bindResult.FieldByName(ostype.Field(i).Name))
		}
	}
	return
}
func (dba *Orm) Pluck(field string, fieldKey ...string) (v t.T, err error) {
	err = dba.Get()
	if err != nil {
		return
	}
	var binder = dba.ISession.GetBinder()
	var resMap = make(t.MapInterface, 0)
	var resSlice = t.Slice{}
	switch binder.GetBindType() {
	case OBJECT_MAP, OBJECT_MAP_T, OBJECT_STRUCT: // row
		var key, val t.T
		if len(fieldKey) > 0 {
			key, err = dba.Value(fieldKey[0])
			if err != nil {
				return
			}
			val, err = dba.Value(field)
			if err != nil {
				return
			}
			v = t.New(t.Map{key: val})
		} else {
			v, err = dba.Value(field)
			if err != nil {
				return
			}
		}
		return
	case OBJECT_MAP_SLICE, OBJECT_MAP_SLICE_T:
		for _, item := range t.New(binder.GetBindResultSlice().Interface()).Slice() {
			val := item.MapInterface()
			if len(fieldKey) > 0 {
				resMap[val[fieldKey[0]].Interface()] = val[field]
				//v = t.New(resMap)
			} else {
				resSlice = append(resSlice, val[field])
				//v = t.New(resSlice)
			}
		}
	case OBJECT_STRUCT_SLICE: // rows
		var brs = binder.GetBindResultSlice()
		for i := 0; i < brs.Len(); i++ {
			val := reflect.Indirect(brs.Index(i))
			//fmt.Println(val)
			if len(fieldKey) > 0 {
				//var resMap = make(t.Map)
				//fmt.Println(dba._valueFromStruct(val, field))
				mapkey := dba._valueFromStruct(val, fieldKey[0])
				mapVal := dba._valueFromStruct(val, field)
				//fmt.Println(mapkey, mapVal)
				resMap[mapkey.Interface()] = mapVal
				//fmt.Println(resMap)
				//v = t.New(resMap)
			} else {
				//var resSlice = make(t.Slice, 0)
				resSlice = append(resSlice, dba._valueFromStruct(val, field))
				//v = t.New(resSlice)
			}
		}
	}
	//fmt.Println(resMap)
	if len(fieldKey) > 0 {
		v = t.New(t.New(resMap).Map())
	} else {
		v = t.New(resSlice)
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
