# Gorose ORM

[![GoDoc](https://godoc.org/github.com/gohouse/gorose?status.svg)](https://godoc.org/github.com/gohouse/gorose)
[![Go Report Card](https://goreportcard.com/badge/github.com/gohouse/gorose)](https://goreportcard.com/report/github.com/gohouse/gorose)
[![Gitter](https://badges.gitter.im/gohouse/gorose.svg)](https://gitter.im/gorose/wechat)
<a target="_blank" href="https://jq.qq.com/?_wv=1027&k=5JJOG9E">
<img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="gorose-orm" title="gorose-orm"></a>

### Official Website | (官方网站)

[https://gohouse.github.io/gorose](https://gohouse.github.io/gorose)

### Documentation

- [中文文档](https://gohouse.github.io/gorose/dist/zh-cn)
- [English document](https://gohouse.github.io/gorose/dist/en)

### What is Gorose?

Gorose, a mini database ORM for golang, which inspired by the famous php framwork laravle's eloquent. It will be friendly for php developers and python or ruby developers.  
Currently provides five major database drivers:   
- **mysql** : <https://github.com/go-sql-driver/mysql>  
- **sqlite3** : <https://github.com/mattn/go-sqlite3>  
- **postgres** : <https://github.com/lib/pq>  
- **oracle** : <https://github.com/mattn/go-oci8>  
- **mssql** : <https://github.com/denisenkom/go-mssqldb>  

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

- standard:  
```go
go get -u github.com/gohouse/gorose
```
- or for [Glide](https://glide.sh):  
```go
glide get github.com/gohouse/gorose
```

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
	"Default": "mysql_dev",
	// (Connection pool) Max open connections, default value 0 means unlimit.
	"SetMaxOpenConns": 300,
	// (Connection pool) Max idle connections, default value is 1.
	"SetMaxIdleConns": 10,

	// Define the database configuration character "mysql_dev".
	"Connections":map[string]map[string]string{
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

	res,err := db.Table("users").First()
	if err != nil {
    		fmt.Println(err)
    		return
    }
	fmt.Println(res)
}

```
For more usage, please read the Documentation.

### License

MIT

### Contribution

- [Issues](https://github.com/gohouse/gorose/issues)
- [Pull requests](https://github.com/gohouse/gorose/pulls)

### release notes

> 0.8.2  

- improve config format, new config format support file config like json/toml etc.

> 0.8.1

- improve multi connection and nulti transation

> 0.8.0  

- add connection pool  
- adjust dir for open source standard  
- add glide version control  
- translate for english and chinese docment  
