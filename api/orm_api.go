package api

import "reflect"

type TableType int

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
