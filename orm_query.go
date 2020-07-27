package gorose

import (
	"github.com/gohouse/t"
	"math"
	"reflect"
	"strings"
)

// Select : select one or more rows , relation limit set
func (dba *Orm) Select() error {
	switch dba.GetIBinder().GetBindType() {
	case OBJECT_STRUCT, OBJECT_MAP, OBJECT_MAP_T:
		dba.Limit(1)
	}
	// 构建sql
	sqlStr, args, err := dba.BuildSql()
	if err != nil {
		return err
	}

	// 执行查询
	_, err = dba.GetISession().Query(sqlStr, args...)
	return err
}

// First : select one row , relation limit set
func (dba *Orm) First() (result Data, err error) {
	dba.GetIBinder().SetBindType(OBJECT_STRING)
	err = dba.Limit(1).Select()
	if err != nil {
		return
	}
	res := dba.GetISession().GetBindAll()
	if len(res) > 0 {
		result = res[0]
	}
	return
}

// Get : select more rows , relation limit set
func (dba *Orm) Get() (result []Data, err error) {
	dba.GetIBinder().SetBindType(OBJECT_STRING)
	tabname := dba.GetISession().GetIBinder().GetBindName()
	prefix := dba.GetISession().GetIBinder().GetBindPrefix()
	tabname2 := strings.TrimPrefix(tabname, prefix)
	dba.ResetTable()
	dba.Table(tabname2)
	err = dba.Select()
	result = dba.GetISession().GetBindAll()
	return
}

// Count : select count rows
func (dba *Orm) Count(args ...string) (int64, error) {
	fields := "*"
	if len(args) > 0 {
		fields = args[0]
	}
	count, err := dba._unionBuild("count", fields)
	if count == nil {
		return 0, err
	}
	return t.New(count).Int64(), err
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
	fields := union + "(" + field + ") as " + union
	dba.fields = []string{fields}

	res, err := dba.First()
	if r, ok := res[union]; ok {
		return r, err
	}
	return 0, err
}

//func (dba *Orm) _unionBuild_bak(union, field string) (interface{}, error) {
//	var tmp interface{}
//
//	dba.union = union + "(" + field + ") as " + union
//	// 缓存fields字段,暂时由union占用
//	fieldsTmp := dba.fields
//	dba.fields = []string{dba.union}
//	dba.GetISession().SetUnion(true)
//
//	// 构建sql
//	sqls, args, err := dba.BuildSql()
//	if err != nil {
//		return tmp, err
//	}
//
//	// 执行查询
//	_, err = dba.GetISession().Query(sqls, args...)
//	if err != nil {
//		return tmp, err
//	}
//
//	// 重置union, 防止复用的时候感染
//	dba.union = ""
//	// 返还fields
//	dba.fields = fieldsTmp
//
//	// 语法糖获取union值
//	if dba.GetISession().GetUnion() != nil {
//		tmp = dba.GetISession().GetUnion()
//		// 获取之后, 释放掉
//		dba.GetISession().SetUnion(nil)
//	}
//
//	return tmp, nil
//}

// Pluck 获取一列数据, 第二个字段可以指定另一个字段的值作为这一列数据的key
func (dba *Orm) Pluck(field string, fieldKey ...string) (v interface{}, err error) {
	var resMap = make(map[interface{}]interface{}, 0)
	var resSlice = make([]interface{}, 0)

	res, err := dba.Get()

	if err != nil {
		return
	}

	if len(res) > 0 {
		for _, val := range res {
			if len(fieldKey) > 0 {
				resMap[val[fieldKey[0]]] = val[field]
			} else {
				resSlice = append(resSlice, val[field])
			}
		}
	}
	if len(fieldKey) > 0 {
		v = resMap
	} else {
		v = resSlice
	}
	return
}

// Pluck_bak ...
func (dba *Orm) Pluck_bak(field string, fieldKey ...string) (v interface{}, err error) {
	var binder = dba.GetISession().GetIBinder()
	var resMap = make(map[interface{}]interface{}, 0)
	var resSlice = make([]interface{}, 0)

	err = dba.Select()
	if err != nil {
		return
	}

	switch binder.GetBindType() {
	case OBJECT_MAP, OBJECT_MAP_T, OBJECT_STRUCT: // row
		var key, val interface{}
		if len(fieldKey) > 0 {
			key, err = dba.Value(fieldKey[0])
			if err != nil {
				return
			}
			val, err = dba.Value(field)
			if err != nil {
				return
			}
			resMap[key] = val
		} else {
			v, err = dba.Value(field)
			if err != nil {
				return
			}
		}
	case OBJECT_MAP_SLICE, OBJECT_MAP_SLICE_T:
		for _, item := range t.New(binder.GetBindResultSlice().Interface()).Slice() {
			val := item.MapInterfaceT()
			if len(fieldKey) > 0 {
				resMap[val[fieldKey[0]].Interface()] = val[field].Interface()
			} else {
				resSlice = append(resSlice, val[field].Interface())
			}
		}
	case OBJECT_STRUCT_SLICE: // rows
		var brs = binder.GetBindResultSlice()
		for i := 0; i < brs.Len(); i++ {
			val := reflect.Indirect(brs.Index(i))
			if len(fieldKey) > 0 {
				mapkey := dba._valueFromStruct(val, fieldKey[0])
				mapVal := dba._valueFromStruct(val, field)
				resMap[mapkey] = mapVal
			} else {
				resSlice = append(resSlice, dba._valueFromStruct(val, field))
			}
		}
	case OBJECT_STRING:
		res := dba.GetISession().GetBindAll()
		if len(res) > 0 {
			for _, val := range res {
				if len(fieldKey) > 0 {
					resMap[val[fieldKey[0]]] = val[field]
				} else {
					resSlice = append(resSlice, val[field])
				}
			}
		}
	}
	if len(fieldKey) > 0 {
		v = resMap
	} else {
		v = resSlice
	}
	return
}

// Value is get a row of a field value
func (dba *Orm) Value(field string) (v interface{}, err error) {
	res, err := dba.First()
	if v, ok := res[field]; ok {
		return v, err
	}
	return
}

// Value_bak ...
func (dba *Orm) Value_bak(field string) (v interface{}, err error) {
	dba.Limit(1)
	err = dba.Select()
	if err != nil {
		return
	}
	var binder = dba.GetISession().GetIBinder()
	switch binder.GetBindType() {
	case OBJECT_MAP, OBJECT_MAP_SLICE, OBJECT_MAP_SLICE_T, OBJECT_MAP_T:
		v = reflect.ValueOf(binder.GetBindResult()).MapIndex(reflect.ValueOf(field)).Interface()
	case OBJECT_STRUCT, OBJECT_STRUCT_SLICE:
		bindResult := reflect.Indirect(reflect.ValueOf(binder.GetBindResult()))
		v = dba._valueFromStruct(bindResult, field)
	case OBJECT_STRING:
		res := dba.GetISession().GetBindAll()
		if len(res) > 0 {
			v = res[0][field]
		}
	}
	return
}
func (dba *Orm) _valueFromStruct(bindResult reflect.Value, field string) (v interface{}) {
	ostype := bindResult.Type()
	for i := 0; i < ostype.NumField(); i++ {
		tag := ostype.Field(i).Tag.Get(TAGNAME)
		if tag == field || ostype.Field(i).Name == field {
			v = bindResult.FieldByName(ostype.Field(i).Name).Interface()
		}
	}
	return
}

// Chunk : 分块处理数据,当要处理很多数据的时候, 我不需要知道具体是多少数据, 我只需要每次取limit条数据,
// 然后不断的增加offset去取更多数据, 从而达到分块处理更多数据的目的
//TODO 后续增加 gorotine 支持, 提高批量数据处理效率, 预计需要增加获取更多链接的支持
func (dba *Orm) Chunk(limit int, callback func([]Data) error) (err error) {
	var page = 1
	var tabname = dba.GetISession().GetIBinder().GetBindName()
	prefix := dba.GetISession().GetIBinder().GetBindPrefix()
	tabname2 := strings.TrimPrefix(tabname, prefix)
	// 先执行一条看看是否报错, 同时设置指定的limit, offset
	result, err := dba.Table(tabname2).Limit(limit).Page(page).Get()
	if err != nil {
		return
	}
	for len(result) > 0 {
		if err = callback(result); err != nil {
			break
		}
		page++
		// 清理绑定数据, 进行下一次操作, 因为绑定数据是每一次执行的时候都会解析并保存的
		// 而第二次以后执行的, 都会再次解析并保存, 数据结构是slice, 故会累积起来
		dba.ClearBindValues()
		result, _ = dba.Page(page).Get()
	}
	return
}

// ChunkStruct : 同Chunk,只不过不用返回map, 而是绑定数据到传入的对象上
// 这里一定要传入绑定struct
func (dba *Orm) ChunkStruct(limit int, callback func() error) (err error) {
	var page = 0
	//var tableName = dba.GetISession().GetIBinder().GetBindName()
	// 先执行一条看看是否报错, 同时设置指定的limit, offset
	err = dba.Limit(limit).Offset(page * limit).Select()
	if err != nil {
		return
	}
	switch dba.GetIBinder().GetBindType() {
	case OBJECT_STRUCT, OBJECT_MAP, OBJECT_MAP_T:
		var ibinder = dba.GetIBinder()
		var result = ibinder.GetBindResult()
		for result != nil {
			if err = callback(); err != nil {
				break
			}
			page++
			// 清空结果
			//result = nil
			var rfRes = reflect.ValueOf(result)
			rfRes.Set(reflect.Zero(rfRes.Type()))
			// 清理绑定数据, 进行下一次操作, 因为绑定数据是每一次执行的时候都会解析并保存的
			// 而第二次以后执行的, 都会再次解析并保存, 数据结构是slice, 故会累积起来
			dba.ClearBindValues()
			_ = dba.Table(ibinder.GetBindOrigin()).Offset(page * limit).Select()
			result = dba.GetIBinder().GetBindResultSlice()
		}
	case OBJECT_STRUCT_SLICE, OBJECT_MAP_SLICE, OBJECT_MAP_SLICE_T:
		var ibinder = dba.GetIBinder()
		var result = ibinder.GetBindResultSlice()
		for result.Interface() != nil {
			if err = callback(); err != nil {
				break
			}
			page++
			// 清空结果
			result.Set(result.Slice(0, 0))
			// 清理绑定数据, 进行下一次操作, 因为绑定数据是每一次执行的时候都会解析并保存的
			// 而第二次以后执行的, 都会再次解析并保存, 数据结构是slice, 故会累积起来
			dba.ClearBindValues()
			_ = dba.Table(ibinder.GetBindOrigin()).Offset(page * limit).Select()
			result = dba.GetIBinder().GetBindResultSlice()
		}
	}
	return
}

// Loop : 同chunk, 不过, 这个是循环的取前limit条数据, 为什么是循环取这些数据呢
// 因为, 我们考虑到一种情况, 那就是where条件如果刚好是要修改的值,
// 那么最后的修改结果因为offset的原因, 只会修改一半, 比如:
// DB().Where("age", 18) ===> DB().Data(gorose.Data{"age":19}).Where().Update()
func (dba *Orm) Loop(limit int, callback func([]Data) error) (err error) {
	var page = 0
	var tabname = dba.GetISession().GetIBinder().GetBindName()
	prefix := dba.GetISession().GetIBinder().GetBindPrefix()
	tabname2 := strings.TrimPrefix(tabname, prefix)
	// 先执行一条看看是否报错, 同时设置指定的limit
	result, err := dba.Table(tabname2).Limit(limit).Get()
	if err != nil {
		return
	}
	for len(result) > 0 {
		if err = callback(result); err != nil {
			break
		}
		page++
		// 同chunk
		dba.ClearBindValues()
		result, _ = dba.Get()
	}
	return
}

// Paginate 自动分页
// @param limit 每页展示数量
// @param current_page 当前第几页, 从1开始
// 以下是laravel的Paginate返回示例
//{
//	"total": 50,
//	"per_page": 15,
//	"current_page": 1,
//	"lastPage": 4,
//	"first_page_url": "http://laravel.app?page=1",
//	"lastPage_url": "http://laravel.app?page=4",
//	"nextPage_url": "http://laravel.app?page=2",
//	"prevPage_url": null,
//	"path": "http://laravel.app",
//	"from": 1,
//	"to": 15,
//	"data":[
//		{
//		// Result Object
//		},
//		{
//		// Result Object
//		}
//	]
//}
func (dba *Orm) Paginate(page ...int) (res Data, err error) {
	if len(page) > 0 {
		dba.Page(page[0])
	}
	var limit = dba.GetLimit()
	if limit == 0 {
		limit = 15
	}
	var offset = dba.GetOffset()
	var currentPage = int(math.Ceil(float64(offset+1) / float64(limit)))
	//dba.ResetUnion()
	// 获取结果
	resData, err := dba.Get()
	if err != nil {
		return
	}
	// 统计总量
	dba.offset = 0
	count, err := dba.Count()
	var lastPage = int(math.Ceil(float64(count) / float64(limit)))
	var nextPage = currentPage + 1
	var prevPage = currentPage - 1
	res = Data{
		"total":          count,
		"per_page":       limit,
		"current_page":   currentPage,
		"last_page":      lastPage,
		"first_page_url": 1,
		"last_page_url":  lastPage,
		"next_page_url":  If(nextPage > lastPage, nil, nextPage),
		"prev_page_url":  If(prevPage < 1, nil, prevPage),
		//"data":           dba.GetIBinder().GetBindResultSlice().Interface(),
		"data": resData,
	}

	return
}
