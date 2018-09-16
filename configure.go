package gorose

import "github.com/gohouse/gorose/api"

const (
	TABLE_STRING       api.TableType = iota // 非结构体 表名字符串	("users")
	TABLE_STRUCT                        // 结构体 一条数据		(struct)
	TABLE_STRUCT_SLICE                  // 结构体 多条数据		([]struct)
)

type Config struct {
	api.DbConfig
}
