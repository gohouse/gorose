package builder

import "github.com/gohouse/gorose/across"

type MssqlBuilder struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var builder IBuilder = &MssqlBuilder{}

	// 注册驱动
	Register("mssql", builder)
}

func (b MssqlBuilder) BuildQuery(api across.OrmApi) (sql string, err error) {
	return builders["mysql"].BuildQuery(api)
}

func (b MssqlBuilder) BuildExecute(api across.OrmApi, operType string) (string, error) {
	return builders["mysql"].BuildExecute(api, operType)
}
