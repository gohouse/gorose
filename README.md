# Gorose

### Documentation
- [English document](docs/en/README.md)
- [中文文档](docs/zh-CN/README.md)

### What is Gorose?

Gorose, a mini database ORM for golang, which inspired by the famous php framwork laravle's eloquent. It will be friendly for php developer and python or ruby developer.  
Currently provides five major database drivers: mysql,sqlite3,postgres,oracle,mssql.

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
var DbConfig = map[string]map[string]string{
	"mysql_dev": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "test",
		"charset":  "utf8",
		"protocol": "tcp",
		"driver":   "mysql", //DB driver
	},
}

func main() {
	db, err := gorose.Open(DbConfig, "mysql_dev")
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

### Contribution

- [Issues](https://github.com/gohouse/gorose/issues)
- [Pull requests](https://github.com/gohouse/gorose/pulls)
