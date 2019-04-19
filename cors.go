package gorose

import (
	"github.com/gohouse/gorose/cors"
)

func BootLogger() func(*Connection) {
	return func(conn *Connection) {
		conn.Logger = cors.NewDefaultLogger()
	}
}

func NewLogger(filePath ...string) func(*Connection) {
	return func(conn *Connection) {
		conn.Logger = cors.NewDefaultLogger(filePath...)
	}
}

//func NewTableToStruct(c *Connection) *converter.Table2Struct {
//	return converter.NewTable2Struct().DB(c.GetExecuteDb())
//}
