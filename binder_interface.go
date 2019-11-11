package gorose

import "reflect"

// IBinder 数据绑定对象接口
type IBinder interface {
	SetBindOrigin(arg interface{})
	GetBindOrigin() interface{}
	SetBindName(arg string)
	GetBindName() string
	SetBindResult(arg interface{})
	GetBindResult() interface{}
	SetBindResultSlice(arg reflect.Value)
	GetBindResultSlice() reflect.Value
	SetBindFields(arg []string)
	GetBindFields() []string
	SetBindType(arg BindType)
	GetBindType() BindType
	//SetBindLimit(arg int)
	//GetBindLimit() int
	BindParse(prefix string) error
	SetBindPrefix(arg string)
	GetBindPrefix() string
	ResetBindResultSlice()
	SetBindAll(arg []Data)
	GetBindAll() []Data
	ResetBinder()
}
