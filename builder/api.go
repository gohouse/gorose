package builder

import (
	"errors"
	"fizzday.com/gohouse/gorose/config"
)

// 注册解析器
var builders = map[string]IBuilder{}

func BuildQuery(d string) (string, error) {
	//return builders[d].BuildQuery()
	var ip IBuilder
	var err error
	if ip, err = Getter(d); err!=nil {
		return "", err
	}
	return ip.BuildQuery()
}

func BuildExecute(d string) (string, error) {
	//return builders[d].BuildQuery()
	var ip IBuilder
	var err error
	if ip, err = Getter(d); err!=nil {
		return "", err
	}
	return ip.BuildExecute()
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
	config.Register(p, "driver")
}
