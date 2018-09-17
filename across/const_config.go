package across

import (
	"errors"
)

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
//	MYSQL:    "driver", // 驱动 ...
//	// 配置文件
//	JSON: "file", // 文件 ...
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
