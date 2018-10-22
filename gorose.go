package gorose

const (
	VERSION_TEXT = "golang orm of gorose's version : "
	VERSION_NO   = "1.0.4"
	VERSION      = VERSION_TEXT + VERSION_NO
)

// Open 链接数据库入口, 传入配置
// args 接收一个或2个参数, 一个参数时:struct配置文件(across.DbConfigCluster{})
//		两个参数时: 第一个是驱动或文件类型, 第二个是dsn或文件路径
func Open(args ...interface{}) (*Connection, error) {
	var c = NewConnection()
	var err error

	// 解析配置获取参数并保存
	c.DbConfig, err = c.parseOpenArgs(args...)

	if err != nil {
		return c, err
	}

	// 驱动数据库获取链接并保存
	err = c.bootDbs(c.DbConfig)
	return c, err
}

func NewConnection() *Connection {
	return &Connection{}
}

func NewOrm() *Session {
	return new(Session)
}
