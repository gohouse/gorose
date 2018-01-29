# Gorose ORM
↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓
> 升级提示: 由于0.8.0做出了大的调整, 已经下载过老版本的, 在升级最新版本时, 需要彻底删除老版本的 gohouse/gorose 目录, 方可升级到新版本  
update notes: if you want updating to the new version, you have to delete the old version directory of 'gohouse/gorose'

↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑
### Documentation

- [中文文档](docs/zh-CN/README.md)
- [English document](docs/en/README.md)

### What is Gorose?

Gorose, a mini database ORM for golang, which inspired by the famous php framwork laravle's eloquent. It will be friendly for php developer and python or ruby developer.  
Currently provides five major database drivers: mysql,sqlite3,postgres,oracle,mssql.

### Quick Preview

```go
// select * from users where id=1 limit 1
db.Table("users").Where("id",1).First()
// select id as uid,name,age from users where id=1 order by id desc limit 10
db.Table("users").Where("id",1).Fields("id as uid,name,age").Order("id desc").Limit(10).Get()

// query string
db.Query("select * from user limit 10")
db.Execute("update users set name='fizzday' where id=?", 1)
```

### Features

- Chain Operation
- Connection Pool

### Requirement

- Golang 1.6+
- [Glide](https://glide.sh) (optional, dependencies management for golang)

### Installation

- for [Glide](https://glide.sh): `$ glide get github.com/gohouse/gorose`
- `$ go get -u github.com/gohouse/gorose`

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
	// Default database configuration
	"default": "mysql_dev",
	// (Connection pool) Max open connections, default value 0 means unlimit.
	"SetMaxOpenConns": 300,
	// (Connection pool) Max idle connections, default value is 1.
	"SetMaxIdleConns": 10,

	// Define the database configuration character "mysql_dev".
	"mysql_dev": map[string]string{
		"host":     "192.168.200.248",
		"username": "gcore",
		"password": "gcore",
		"port":     "3306",
		"database": "test",
		"charset":  "utf8",
		"protocol": "tcp",
		"prefix":   "",      // Table prefix
		"driver":   "mysql", // Database driver(mysql,sqlite,postgres,oracle,mssql)
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

- [Join QQ group number: 470809220](https://jq.qq.com/?_wv=1027&k=5JJOG9E)

### Contribution

- [Issues](https://github.com/gohouse/gorose/issues)
- [Pull requests](https://github.com/gohouse/gorose/pulls)
