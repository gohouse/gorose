package builder

import "fizzday.com/gohouse/gorose/config"

type MysqlBuilder struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var builder IBuilder = &MysqlBuilder{}

	// 注册驱动
	Register(config.MYSQL, builder)
}

func (sql MysqlBuilder) BuildQuery() (string, error) {
	return "MysqlBuilder BuildQuery", nil
}

func (sql MysqlBuilder) BuildExecute() (string, error) {
	return "MysqlBuilder BuildExecute", nil
}
