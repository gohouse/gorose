# Gorose Orm

### Documentation

- [中文文档](docs/zh-CN/README.md)
- [English document](docs/en/README.md)

### What is Gorose?

Gorose, a mini database ORM for golang, which inspired by the famous php framwork laravle's eloquent. It will be friendly for php developer and python or ruby developer.  
Currently provides five major database drivers: mysql,sqlite3,postgres,oracle,mssql.

### quick scan

```go
// select * from users where id=1 limit 1
db.Table("users").Where("id",1).First()
// select id as uid,name,age from users where id=1 order by id desc limit 10
db.Table("users").Where("id",1).Fields("id as uid,name,age").Order("id desc").Limit(10).Get()
// query string
db.Query("select * from user limit 10")
db.Execute("update users set name='fizzday' where id=?", 1)
```

### character

- Chain operation  
- connection pool  

### Requirement

- Golang 1.6+
- [Glide](https://glide.sh) (Optional, Dependencies management for golang)

### Installation

- `$ go get -u github.com/gohouse/gorose`
- for Glide: `$ glide get github.com/gohouse/gorose`

### Base Usage
```go
package main

import (
	"github.com/gohouse/gorose"        //import Gorose
	_ "github.com/go-sql-driver/mysql" //import DB driver
	"fmt"
)

// DB Config.(Recommend to use configuration file to import)
var DbConfig = map[string]interface{}{
	"default":         "mysql_dev", // 默认数据库配置
	"SetMaxOpenConns": 300,         // (连接池)最大打开的连接数，默认值为0表示不限制
	"SetMaxIdleConns": 10,         // (连接池)闲置的连接数, 默认1

	"mysql_dev": map[string]string{// 定义名为 mysql_dev 的数据库配置
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
}

func main() {
	db, err := gorose.Open(DbConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	res := db.Table("users").First()
	fmt.Println(res)
}

```
For more usage, please read the Documentation.

### License

MIT

### Exchange and Discussion

- QQ group number: 470809220

- [点击加入qq群: 470809220](https://jq.qq.com/?_wv=1027&k=5JJOG9E)  

### Contribution

- [Issues](https://github.com/gohouse/gorose/issues)
- [Pull requests](https://github.com/gohouse/gorose/pulls)
