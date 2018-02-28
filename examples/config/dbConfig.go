package config

// DbConfig : 所有数据库配置
// 如果想从json读取, 只需要用这个struct读取出来,
// 然后在赋值给类似这个变量类型就可以了: map[string]interface{}
// type Configer struct {
// 	   Default string
// 	   SetMaxOpenConns int
// 	   SetMaxIdleConns int
// 	   Connections map[string]map[string]string
// }
var DbConfig = map[string]interface{}{
	"Default":         "mysql_dev", // 默认数据库配置
	"SetMaxOpenConns": 300,         // (连接池)最大打开的连接数，默认值为0表示不限制
	"SetMaxIdleConns": 10,          // (连接池)闲置的连接数, 默认-1
	"Connections":map[string]map[string]string{
		"mysql_dev": { // 定义名为 mysql_dev 的数据库配置
			"host":     "192.168.200.248", // 数据库地址
			"username": "gcore",           // 数据库用户名
			"password": "gcore",           // 数据库密码
			"port":     "3306",            // 端口
			"database": "test",            // 链接的数据库名字
			"charset":  "utf8",            // 字符集
			"protocol": "tcp",             // 链接协议
			"prefix":   "",                // 表前缀
			"driver":   "mysql",           // 数据库驱动(mysql,sqlite,postgres,oracle,mssql)
		},
		"mssql_dev": {
			"host":     "192.168.200.248",
			"username": "gcore",
			"password": "gcore",
			"port":     "1433",
			"database": "test",
			"prefix":   "",
			//"charset":  "utf8",
			//"protocol": "tcp",
			"driver": "mssql",
		},
		"postgres_dev": {
			"host":     "localhost",
			"username": "postgres",
			"password": "",
			"port":     "5432",
			"database": "test",
			"prefix":   "",
			//"charset":  "utf8",
			//"protocol": "tcp",
			"driver": "postgres",
		},
		"oracle_dev": {
			//"host":     "localhost",
			"username": "root",
			"password": "",
			//"port":     "1521",
			"database": "test",
			"prefix":   "",
			//"charset":  "utf8",
			//"protocol": "tcp",
			"driver": "oracle",
		},
		"sqlite_dev": {
			//"host":     "localhost",
			//"username": "root",
			//"password": "",
			//"port":     "3306",
			"database": "./foo.db",
			"prefix":   "",
			//"charset":  "utf8",
			//"protocol": "tcp",
			"driver": "sqlite3",
		},
	},
}
