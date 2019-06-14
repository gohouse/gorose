package gorose

import "reflect"

type IBinder interface {
	SetBindOrigin(arg interface{})
	GetBindOrigin() interface{}
	SetBindName(arg string)
	GetBindName() string
	SetBindResult(arg reflect.Value)
	GetBindResult() reflect.Value
	SetBindResultSlice(arg reflect.Value)
	GetBindResultSlice() reflect.Value
	SetBindFields(arg []string)
	GetBindFields() []string
	SetBindType(arg BindType)
	GetBindType() BindType
	SetBindLimit(arg int)
	GetBindLimit() int
	BindParse(prefix string) error
}
