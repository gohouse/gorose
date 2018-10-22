# GoRose ORM

[![GoDoc](https://godoc.org/github.com/gohouse/gorose?status.svg)](https://godoc.org/github.com/gohouse/gorose)
[![Go Report Card](https://goreportcard.com/badge/github.com/gohouse/gorose)](https://goreportcard.com/report/github.com/gohouse/gorose)
[![Gitter](https://badges.gitter.im/gohouse/gorose.svg)](https://gitter.im/gorose/wechat)
<a target="_blank" href="https://jq.qq.com/?_wv=1027&k=5JJOG9E">
<img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="gorose-orm" title="gorose-orm"></a>

### [中文-readme](https://github.com/gohouse/gorose/blob/master/README_zh-cn.md) | [english-readme](https://github.com/gohouse/gorose/blob/master/README.md)

### What is GoRose?

GoRose, a mini database ORM for golang, which inspired by the famous php framwork laravel's eloquent. It will be friendly for php developers and python or ruby developers.  
Currently provides five major database drivers:   
- **mysql** : <https://github.com/go-sql-driver/mysql>  
- **sqlite3** : <https://github.com/mattn/go-sqlite3>  
- **postgres** : <https://github.com/lib/pq>  
- **oracle** : <https://github.com/mattn/go-oci8>  
- **mssql** : <https://github.com/denisenkom/go-mssqldb>  

### 1.0.0 update notes
- struct support  
- seperation of write & read cluster  
- New architecture  


### Documentation

[latest document](https://www.kancloud.cn/fizz/gorose) | [最新中文文档](https://www.kancloud.cn/fizz/gorose)  
[0.x version english document](https://gohouse.github.io/gorose/dist/en/index.html) | [0.x版本中文文档](https://gohouse.github.io/gorose/dist/zh-cn/index.html)  
[github](https://github.com/gohouse/gorose)  

### Quick Preview

```go
type users struct {
	Name string
	Age int `orm:"age"`
}

// select * from users where id=1 limit 1
var user users      // a row data
var users []users   // several rows
// use struct
db.Table(&user).Select()
db.Table(&users).Where("id",1).Limit(10).Select()
// use string instead struct
db.Table("users").Where("id",1).First()

// select id as uid,name,age from users where id=1 order by id desc limit 10
db.Table(&user).Where("id",1).Fields("id as uid,name,age").Order("id desc").Limit(10).Get()

// query string
db.Query("select * from user limit 10")
db.Execute("update users set name='fizzday' where id=?", 1)
```

### Features

- Chain Operation
- Connection Pool
- struct/string compatible
- read/write separation cluster
- process a lot of data into slices  
- transaction easily  
- friendly for extended (extend more builders or config parsers)  

### Installation

- standard:  
```go
go get -u github.com/gohouse/gorose
```

### Base Usage
```go
package main

import (
	"github.com/gohouse/gorose"
	_ "github.com/gohouse/gorose/driver/mysql"
	"fmt"
)

type Users struct {
	Name string
	Age  int `orm:"age"`
}

// DB Config.(Recommend to use configuration file to import)
var dbConfig = &gorose.DbConfigSingle{
    Driver:          "mysql", // driver: mysql/sqlite/oracle/mssql/postgres
    EnableQueryLog:  true,    // if enable sql logs
    SetMaxOpenConns: 0,       // connection pool of max Open connections, default zero
    SetMaxIdleConns: 0,       // connection pool of max sleep connections
    Prefix:          "",      // prefix of table
    Dsn:             "root:root@tcp(localhost:3306)/test?charset=utf8", // db dsn
}

func main() {
	connection, err := gorose.Open(dbConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	// start a new session
	session := connection.NewSession()
	// get a row of data
	var user Users
	err2 := session.Table(&user).Select()
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println(session.LastSql)
	fmt.Println(user)
	
	// get several rows of data
	var users []Users
	// use connection derictly instead of NewSession()
	err3 := connection.Table(&users).Limit(3).Select()
	if err3 != nil {
		fmt.Println(err3)
		return
	}
	fmt.Println(users)
}
```

For more usage, please read the Documentation.

### Contribution

- [Issues](https://github.com/gohouse/gorose/issues)
- [Pull requests](https://github.com/gohouse/gorose/pulls)

### Contributors

- `fizzday` : Initiator  
- `wuyumin` : pursuing the open source standard  
- `holdno`  : official website builder  
- `LazyNeo` : bug fix and improve source code  
- `dmhome`  : improve source code 
 
### release notes

> v1.0.4

- add middleware support, add logger cors

> v1.0.3

- add version get by const: gorose.VERSION

> v1.0.2

- improve go mod's bug

> 1.0.0

- New architecture, struct support, seperation of write & read cluster  

> 0.9.2  

- new connection pack for supporting multi connection

> 0.9.1  

- replace the insert result lastInsertId with rowsAffected as default

> 0.9.0  

- new seperate db instance

> 0.8.2  

- improve config format, new config format support file config like json/toml etc.

> 0.8.1

- improve multi connection and nulti transation

> 0.8.0  

- add connection pool  
- adjust dir for open source standard  
- add glide version control  
- translate for english and chinese docment  

### License

MIT
