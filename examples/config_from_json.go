package main

import (
	"fmt"
	"encoding/json"
	"github.com/gohouse/gorose"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// 定义config的struct
	type Configer struct {
		Default         string
		SetMaxOpenConns int
		SetMaxIdleConns int
		Connections     map[string]map[string]string
	}

	// 数据库配置
	var DbConfig = map[string]interface{}{
		"Default":         "mysql_dev", // 默认数据库配置
		"SetMaxOpenConns": 300,         // (连接池)最大打开的连接数，默认值为0表示不限制
		"SetMaxIdleConns": 10,          // (连接池)闲置的连接数, 默认-1
		"Connections": map[string]map[string]string{
			"mysql_dev": {// 定义名为 mysql_dev 的数据库配置
				"host": "192.168.200.248", // 数据库地址
				"username": "gcore",       // 数据库用户名
				"password": "gcore",       // 数据库密码
				"port": "3306",            // 端口
				"database": "test",        // 链接的数据库名字
				"charset": "utf8",         // 字符集
				"protocol": "tcp",         // 链接协议
				"prefix": "",              // 表前缀
				"driver": "mysql",         // 数据库驱动(mysql,sqlite,postgres,oracle,mssql)
			},
		},
	}

	// json 格式配置
	jsons,_ := json.Marshal(DbConfig)
	// json 的转换结果
	// jsons := `{"Default":"mysql_dev","SetMaxOpenConns":300,"SetMaxIdleConns":10,"Connections":{"mysql_dev":{"charset":"utf8","database":"test","driver":"mysql","host":"192.168.200.248","password":"gcore","port":"3306","prefix":"","protocol":"tcp","username":"gcore"}}}`

	var conf Configer

	json.Unmarshal([]byte(jsons), &conf)

	var confReal = map[string]interface{}{
		"Default":         conf.Default,
		"SetMaxOpenConns": conf.SetMaxOpenConns,
		"SetMaxIdleConns": conf.SetMaxIdleConns,
		"Connections":     conf.Connections,
	}

	//	fmt.Println(confReal)
	//return
	db, err := gorose.Open(confReal, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	res, err := db.Table("users").Where("id>2").First()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res["id"])
}
