package gorose

type TableType int
const (
	TABLE_STRING       TableType = iota // 非结构体 表名字符串	("users")
	TABLE_STRUCT                        // 结构体 一条数据		(struct)
	TABLE_STRUCT_SLICE                  // 结构体 多条数据		([]struct)
)
