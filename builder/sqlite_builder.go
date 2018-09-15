package builder

import "fizzday.com/gohouse/gorose/config"

type SqliteBuilder struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var builder IBuilder = &SqliteBuilder{}

	// 注册驱动
	Register(config.SQLITE, builder)
}

func (sql SqliteBuilder) BuildQuery() (string, error) {
	return "SqliteBuilder BuildQuery", nil
}

func (sql SqliteBuilder) BuildExecute() (string, error) {
	return "SqliteBuilder BuildExecute", nil
}
