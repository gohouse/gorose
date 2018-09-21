package parser

import (
	"errors"
	"github.com/gohouse/gorose/across"
)

// 注册解析器
var fileParsers = map[string]IParser{}

// NewFileParser 对外提供配置文件解析器接口
// fileType 文件类型
// file 文件路径
func NewFileParser(fileType, file string, dbConfCluster interface{}) error {
	var ip IParser
	var err error
	if ip, err = Getter(fileType); err != nil {
		return errors.New("不支持的配置类型")
	}
	return ip.Parse(file, dbConfCluster)
}

//func NewFileParser(fileType, file string) (*across.DbConfigCluster, error) {
//	var ip IParser
//	var err error
//	if ip, err = Getter(fileType); err!=nil {
//		return &across.DbConfigCluster{}, errors.New("不支持的配置类型")
//	}
//	return ip.Parse(file)
//}

// Getter 获取解析器
func Getter(p string) (IParser, error) {
	if pr, ok := fileParsers[p]; ok {
		return pr, nil
	}
	return nil, errors.New("解析器不存在")
}

// Register 注册解析器
func Register(p string, ip IParser) {
	if _, ok := fileParsers[p]; ok {
		panic("解析器已存在,请检查是否重复注册了该key值: " + p)
	}
	fileParsers[p] = ip
	// 注册类型,方便Open()解析时区分
	across.Register(p, "file")
}
