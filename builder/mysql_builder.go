package builder

import (
	"fmt"
	"github.com/gohouse/gorose/across"
	"strings"
)

type MysqlBuilder struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var builder IBuilder = &MysqlBuilder{}

	// 注册驱动
	Register(across.MYSQL, builder)
}

func (m MysqlBuilder) BuildQuery(api across.OrmApi) (sql string, err error) {
	var fields, table, limit, offset string
	// table
	//if table, err = api.ParseTable(); err != nil {
	//	return
	//}
	table = api.TableName
	// fields
	fields = strings.Join(api.SFields, ", ")
	if fields == "" {
		fields = "*"
	}
	// limit
	limit = " LIMIT 3"
	// offset
	offset = " OFFSET 0"

	//sqlstr := "select " + fields + " from " + table + limit + offset
	sqlstr := fmt.Sprintf("SELECT %s FROM %s%s%s", fields, table, limit, offset)

	return sqlstr, nil
}

func (m MysqlBuilder) BuildExecute(api across.OrmApi, operType string) (string, error) {
	return "MysqlBuilder BuildExecute", nil
}
