package parser

type IniConfigParser struct {
}

func init()  {
	// 检查解析器是否实现了接口
	var parserTmp IParser = &IniConfigParser{}

	// 注册驱动
	Register("ini", parserTmp)
}

func (c *IniConfigParser) Parse(file string, dbConfCluster interface{}) (err error) {
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
