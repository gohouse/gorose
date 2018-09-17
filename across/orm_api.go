package across

import "reflect"

type TableType int

const (
	TABLE_STRING       TableType = iota // 非结构体 表名字符串	("users")
	TABLE_STRUCT                        // 结构体 一条数据		(struct)
	TABLE_STRUCT_SLICE                  // 结构体 多条数据		([]struct)
)

type table struct {
	STable      interface{}
	TableName   string
	TableStruct reflect.Value
	TableSlice  reflect.Value
	TableType   TableType
}
type OrmApi struct {
	table
	Driver  string
	SFields []string
	SLimit  int
	SOffset int
}
