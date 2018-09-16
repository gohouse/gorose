package builder

import (
	"github.com/gohouse/gorose/api"
	"github.com/gohouse/gorose/config"
)

type SqliteBuilder struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var builder IBuilder = &SqliteBuilder{}

	// 注册驱动
	Register(config.SQLITE, builder)
}

func (b SqliteBuilder) BuildQuery(api api.OrmApi) (sql string, err error) {
	return "SqliteBuilder BuildQuery", nil
}

func (b SqliteBuilder) BuildExecute(api api.OrmApi, operType string) (string, error) {
	return "SqliteBuilder BuildExecute", nil
}
