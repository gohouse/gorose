package gorose

import (
	"github.com/gohouse/t"
	"reflect"
	"sync"
)

type MapRow map[string]t.T
type MapRows []MapRow

type BindType int

const (
	OBJECT_STRUCT       BindType = iota // 结构体 一条数据	(struct)
	OBJECT_STRUCT_SLICE                 // 结构体 多条数据	([]struct)
	OBJECT_MAP                          // map 一条数据		(map[string]interface{})
	OBJECT_MAP_SLICE                    // map 多条数据		([]map[string]interface{})
	OBJECT_STRING                       // 非结构体 表名字符串	("users")
	OBJECT_MAP_T                        // map 一条数据		(map[string]t.T)
	OBJECT_MAP_SLICE_T                  // map 多条数据		([]map[string]t.T)
)

var BindString = map[BindType]string {
	OBJECT_STRUCT:"OBJECT_STRUCT",
	OBJECT_STRUCT_SLICE: "OBJECT_STRUCT_SLICE",
	OBJECT_MAP: "OBJECT_MAP",
	OBJECT_MAP_SLICE: "OBJECT_MAP_SLICE",
	OBJECT_STRING: "OBJECT_STRING",
	OBJECT_MAP_T: "OBJECT_MAP_T",
	OBJECT_MAP_SLICE_T: "OBJECT_MAP_SLICE_T",
}

func (b BindType) String() string {
	return BindString[b]
}

type Binder struct {
	// Bind是指传入的对象 [slice]map,[slice]struct
	// 传入的原始对象
	BindOrigin interface{}
	//BindOriginTableName []string
	// 解析出来的对象名字, 或者指定的method(TableName)获取到的名字
	BindName string
	// 一条结果的反射对象
	BindResult reflect.Value
	// 多条
	BindResultSlice reflect.Value
	// 传入结构体解析出来的字段
	BindFields []string
	// 传入的对象类型判定
	BindType BindType
	// 出入传入得是非slice对象, 则只需要取一条, 取多了也是浪费
	BindLimit int
}

var binderOnce sync.Once
var _ IBinder = &Binder{}
var binder *Binder
func NewBinder(o ...interface{}) IBinder {
	binderOnce.Do(func() {
		binder = &Binder{}
		if len(o)>0 {
			binder.SetBindOrigin(o)
		}
	})
	return binder
}

func (o *Binder) SetBindOrigin(arg interface{}) {
	o.BindOrigin = arg
}

func (o *Binder) GetBindOrigin() interface{} {
	return o.BindOrigin
}

func (o *Binder) SetBindName(arg string) {
	o.BindName = arg
}

func (o *Binder) GetBindName() string {
	return o.BindName
}

func (o *Binder) SetBindResult(arg reflect.Value) {
	o.BindResult = arg
}

func (o *Binder) GetBindResult() reflect.Value {
	return o.BindResult
}

func (o *Binder) SetBindResultSlice(arg reflect.Value) {
	o.BindResultSlice = arg
}

func (o *Binder) GetBindResultSlice() reflect.Value {
	return o.BindResultSlice
}

func (o *Binder) SetBindFields(arg []string) {
	o.BindFields = arg
}

func (o *Binder) GetBindFields() []string {
	return o.BindFields
}

func (o *Binder) SetBindType(arg BindType) {
	o.BindType = arg
}

func (o *Binder) GetBindType() BindType {
	return o.BindType
}

func (o *Binder) SetBindLimit(arg int) {
	o.BindLimit = arg
}

func (o *Binder) GetBindLimit() int {
	return o.BindLimit
}
