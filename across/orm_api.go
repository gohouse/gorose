package across

import (
	"database/sql"
	"reflect"
)

var (
	Regex = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "not like", "in", "not in", "between", "not between"}
)

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
	Driver                string
	Prefix                string
	Sforce                bool
	Sfields               []string
	Swhere                [][]interface{} // where
	Sorder                string          // order
	Slimit                int             // limit
	Soffset               int             // offset
	Sjoin                 [][]interface{} // join
	Sdistinct             bool            // distinct
	Sunion                string          // sum/count/avg/max/min
	Sgroup                string          // group
	Shaving               string          // having
	Sdata                 interface{}     // data
	Stx                   *sql.Tx         //Dbstruct Database
	SbeforeParseWhereData [][]interface{}
	Strans                bool
	LastInsertId          int64 // insert last insert id
	SqlLogs               []string
	LastSql               string
}
