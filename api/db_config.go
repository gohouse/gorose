package api

// 单一数据库配置
type DbConfig struct {
	Driver          string // 驱动: mysql/sqlite/oracle/mssql/postgres
	EnableQueryLog  bool   // 是否开启sql日志
	SetMaxOpenConns int    // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns int    // (连接池)闲置的连接数, 默认-1
	Prefix          string // 表前缀
	Dsn             string // 数据库链接
}

// 数据库集群配置
// 如果不启用集群, 则直接使用 DbConfig 即可
// 如果仍然使用此配置为非集群, 则可以设置 EnableCluster=fasle, 等同于使用 DbConfig
type DbConfigCluster struct {
	EnableCluster bool       // 是否启用主从配置集群,如果不启用,则读取 Master 的值作为默认配置
	Slave         []DbConfig // 多台读服务器, 如果启用则需要放入对应的多台从服务器配置
	Master        DbConfig   // 一台主服务器负责写数据
}
