package builder

import (
	"errors"
	"github.com/gohouse/gorose/across"
)

// 注册解析器
var builders = map[string]IBuilder{}
// NewBuilder : build sql
func NewBuilder(api across.OrmApi,operType ...string) (string, error) {
	// 获取驱动类型
	driver := api.Driver
	builder,err := Getter(driver)
	if err!=nil{
		return "",err
	}
	switch len(operType) {
	case 0:
		return builder.BuildQuery(api)
	case 1:
		if operType[0] == "select" {
			return builder.BuildQuery(api)
		} else {
			return builder.BuildExecute(api, operType[0])
		}
	default:
		return "", errors.New("参数有误")
	}
}

// Getter 获取解析器
func Getter(p string) (IBuilder, error) {
	if pr, ok := builders[p]; ok {
		return pr,nil
	}
	return nil, errors.New("解析器不存在")
}

// Register 注册解析器
func Register(p string, ip IBuilder) {
	builders[p] = ip
	// 注册类型,方便Open()解析时区分
	across.Register(p, "driver")
}
