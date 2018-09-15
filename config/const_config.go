package config

import "errors"

// 数据库驱动登记
const (
	MYSQL    = "mysql"
	SQLITE   = "sqlite"
	MSSQL    = "mssql"
	ORACLE   = "oracle"
	POSTGRES = "postgres"
)

// 配置文件类型登记
const (
	JSON = "json"
	TOML = "toml"
	INI  = "ini"
)

// 类型分类
//	// 数据库驱动
//	MYSQL:    "driver", // 驱动
//	// 配置文件
//	JSON: "file", // 文件
var constsType = map[string]string{}

// Getter 获取类型分类
func Getter(p string) (string, error) {
	if pr, ok := constsType[p]; ok {
		return pr, nil
	}
	return "", errors.New("类型分类不存在")
}

// Register 注册类型分类
func Register(p string, ip string) {
	constsType[p] = ip
}

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

var DemoDbConfig = DbConfig{
	Driver:          "mysql",
	EnableQueryLog:  true,
	SetMaxOpenConns: 0,
	SetMaxIdleConns: 0,
	Dsn:             "username:password@tcp(localhost:3306)/database?charset=utf8",
}
