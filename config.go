package gorose

type Config struct {
	Driver          string // 驱动: mysql/sqlite3/oracle/mssql/postgres/clickhouse
	Dsn             string // 数据库链接
	SetMaxOpenConns int    // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns int    // (连接池)闲置的连接数, 默认0
}

type ConfigCluster struct {
	Master         []Config
	Slave          []Config
	Prefix         string
	EnableSqlLog bool
}
