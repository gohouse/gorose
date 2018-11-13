package gorose

const GOROSE_IMG = `
                                                                            
  ,ad8888ba,                88888888ba                                      
 d8"'    '"8b               88      "8b                                     
d8'                         88      ,8P                                     
88              ,adPPYba,   88aaaaaa8P'  ,adPPYba,   ,adPPYba,   ,adPPYba,  
88      88888  a8"     "8a  88""""88'   a8"     "8a  I8[    ""  a8P_____88  
Y8,        88  8b       d8  88    '8b   8b       d8   '"Y8ba,   8PP"""""""  
 Y8a.    .a88  "8a,   ,a8"  88     '8b  "8a,   ,a8"  aa    ]8I  "8b,   ,aa  
  '"Y88888P"    '"YbbdP"'   88      '8b  '"YbbdP"'   '"YbbdP"'   '"Ybbd8"'  
                                                                             
`

const (
	VERSION_TEXT = "\ngolang orm of gorose's version : "
	VERSION_NO   = "1.0.4"
	VERSION      = VERSION_TEXT + VERSION_NO + GOROSE_IMG
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
