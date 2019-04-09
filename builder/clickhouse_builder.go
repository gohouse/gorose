package builder

import (
	"github.com/gohouse/gorose/across"
)

type ClickhouseBuilder struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var builder IBuilder = &ClickhouseBuilder{}

	// 注册驱动
	Register("clickhouse", builder)
}

func (b ClickhouseBuilder) BuildQuery(api across.OrmApi) (sql string, err error) {
	return builders["mysql"].BuildQuery(api)
}

func (b ClickhouseBuilder) BuildExecute(api across.OrmApi, operType string) (string, error) {
	return builders["mysql"].BuildExecute(api, operType)
}
