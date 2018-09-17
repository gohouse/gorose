package parser

import (
	"github.com/gohouse/gorose/across"
)

type IniConfigParser struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var parser IParser = &IniConfigParser{}

	// 注册驱动
	Register(across.TOML, parser)
}

func (c *IniConfigParser) Parse(file string) (conf *across.DbConfigCluster, err error) {
	//conf = &across.DbConfigCluster{}
	//var iniConfig *ini.Config
	//iniConfig, err = ini.ReadDefault(file)
	//if err != nil {
	//	return
	//}
	//conf.Driver, err = iniConfig.String("MYSQL", "Driver")
	//if err != nil {
	//	return
	//}
	//conf.EnableQueryLog, err = iniConfig.Bool("MYSQL", "EnableQueryLog")
	//if err != nil {
	//	return
	//}
	//conf.SetMaxOpenConns, err = iniConfig.Int("MYSQL", "SetMaxOpenConns")
	//if err != nil {
	//	return
	//}
	//conf.SetMaxIdleConns, err = iniConfig.Int("MYSQL", "SetMaxIdleConns")
	//if err != nil {
	//	return
	//}
	//conf.Dsn, err = iniConfig.String("MYSQL", "Dsn")
	//if err != nil {
	//	return
	//}
	return
}
