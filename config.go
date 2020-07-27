package gorose

// Config ...
type Config struct {
	Driver string `json:"driver"` // 驱动: mysql/sqlite3/oracle/mssql/postgres/clickhouse, 如果集群配置了驱动, 这里可以省略
	// mysql 示例:
	// root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true
	Dsn             string `json:"dsn"`             // 数据库链接
	SetMaxOpenConns int    `json:"setMaxOpenConns"` // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns int    `json:"setMaxIdleConns"` // (连接池)闲置的连接数, 默认0
	Prefix          string `json:"prefix"`          // 表前缀, 如果集群配置了前缀, 这里可以省略
}

// ConfigCluster ...
type ConfigCluster struct {
	Master []Config // 主
	Slave  []Config // 从
	Driver string   // 驱动
	Prefix string   // 前缀
}
