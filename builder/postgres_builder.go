package builder

import "github.com/gohouse/gorose/across"

type PostgresBuilder struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var builder IBuilder = &PostgresBuilder{}

	// 注册驱动
	Register("postgres", builder)
}

func (b PostgresBuilder) BuildQuery(api across.OrmApi) (sql string, err error) {
	return builders["mysql"].BuildQuery(api)
}

func (b PostgresBuilder) BuildExecute(api across.OrmApi, operType string) (string, error) {
	return builders["mysql"].BuildExecute(api, operType)
}
