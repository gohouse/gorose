package api

type DbConfig struct {
	Driver          string // 驱动: mysql/sqlite/oracle/mssql/postgres
	EnableQueryLog  bool   // 是否开启sql日志
	SetMaxOpenConns int    // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns int    // (连接池)闲置的连接数, 默认-1
	Dsn             string // 数据库链接
}

// 数据库集群配置
type DbConfigCluster struct {
	Slave  []DbConfig // 多台读服务器
	Master DbConfig   // 一台主服务器负责写数据
}
