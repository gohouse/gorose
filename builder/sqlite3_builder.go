package builder

import (
	"github.com/gohouse/gorose/across"
)

type SqliteBuilder struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var builder IBuilder = &SqliteBuilder{}

	// 注册驱动
	Register("sqlite3", builder)
}

func (b SqliteBuilder) BuildQuery(api across.OrmApi) (sql string, err error) {
	return builders["mysql"].BuildQuery(api)
}

func (b SqliteBuilder) BuildExecute(api across.OrmApi, operType string) (string, error) {
	return builders["mysql"].BuildExecute(api, operType)
}
