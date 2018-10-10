package builder

import "github.com/gohouse/gorose/across"

type OracleBuilder struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var builder IBuilder = &OracleBuilder{}

	// 注册驱动
	Register("oracle", builder)
}

func (b OracleBuilder) BuildQuery(api across.OrmApi) (sql string, err error) {
	return builders["mysql"].BuildQuery(api)
}

func (b OracleBuilder) BuildExecute(api across.OrmApi, operType string) (string, error) {
	return builders["mysql"].BuildExecute(api, operType)
}
