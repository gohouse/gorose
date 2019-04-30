package gorose

import "reflect"

type MapRow map[string]interface{}
type MapRows []MapRow

type BindType int

const (
	OBJECT_STRUCT       BindType = iota // 结构体 一条数据	(struct)
	OBJECT_STRUCT_SLICE                 // 结构体 多条数据	([]struct)
	OBJECT_MAP                          // map 一条数据		(map[string]interface{})
	OBJECT_MAP_SLICE                    // map 多条数据		([]map[string]interface{})
	OBJECT_STRING                       // 非结构体 表名字符串	("users")
)

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
