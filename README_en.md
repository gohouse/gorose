<base target="main">

# GoRose ORM
[![GoDoc](https://godoc.org/github.com/gohouse/gorose/v2?status.svg)](https://godoc.org/github.com/gohouse/gorose/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/gohouse/gorose/v2)](https://goreportcard.com/report/github.com/gohouse/gorose/v2)
[![GitHub release](https://img.shields.io/github/release/gohouse/gorose.svg)](https://github.com/gohouse/gorose/v2/releases/latest)
[![Gitter](https://badges.gitter.im/gohouse/gorose.svg)](https://gitter.im/gorose/wechat)
![GitHub](https://img.shields.io/github/license/gohouse/gorose?color=blue)
![GitHub All Releases](https://img.shields.io/github/downloads/gohouse/gorose/total?color=blue)
<a target="_blank" href="https://jq.qq.com/?_wv=1027&k=5JJOG9E">
<img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="gorose-orm" title="gorose-orm"></a>

```
  _______   ______   .______        ______        _______. _______ 
 /  _____| /  __  \  |   _  \      /  __  \      /       ||   ____|
|  |  __  |  |  |  | |  |_)  |    |  |  |  |    |   (----`|  |__   
|  | |_ | |  |  |  | |      /     |  |  |  |     \   \    |   __|  
|  |__| | |  `--'  | |  |\  \----.|  `--'  | .----)   |   |  |____ 
 \______|  \______/  | _| `._____| \______/  |_______/    |_______|
```

## translations  
[English readme](README_en.md) |
[中文 readme](README.md) 

## document
[2.x doc](https://www.kancloud.cn/fizz/gorose-2/1135835) | 
[1.x doc](https://www.kancloud.cn/fizz/gorose/769179) | 
[0.x doc](https://gohouse.github.io/gorose/dist/en/index.html)

## introduction
gorose is a golang orm framework, which is Inspired by laravel's eloquent.  
Gorose 2.0 adopts modular architecture, communicates through the API of interface, and strictly relies on the lower layer. Each module can be disassembled, and even can be customized to its preferred appearance.  
The module diagram is as follows:   
![gorose.2.0.jpg](https://i.loli.net/2019/08/06/7R2GlbwUiFKOrNP.jpg)

## installation
- go.mod
```bash
require github.com/gohouse/gorose/v2 v2.1.7
```
> you should use it like `import "github.com/gohouse/gorose/v2"`  
    don't forget the `v2` in the end

- docker
```bash
docker run -it --rm ababy/gorose sh -c "go run main.go"
```
> docker image: [ababy/gorose](https://cloud.docker.com/u/ababy/repository/docker/ababy/gorose), The docker image contains the packages and runtime environment necessary for gorose, [view `Dockerfile`](https://github.com/docker-box/gorose/blob/master/golang-alpine/Dockerfile)   

- go get  
```bash
go get -u github.com/gohouse/gorose
```

## supported drivers
- mysql : https://github.com/go-sql-driver/mysql  
- sqlite3 : https://github.com/mattn/go-sqlite3  
- postgres : https://github.com/lib/pq  
- oracle : https://github.com/mattn/go-oci8  
- mssql : https://github.com/denisenkom/go-mssqldb  
- clickhouse : https://github.com/kshvakov/clickhouse

## api preview
```go
db.Table().Fields().Distinct().Where().GroupBy().Having().OrderBy().Limit().Offset().Select()
db.Table().Data().Insert()
db.Table().Data().Where().Update()
db.Table().Where().Delete()
```

## simple usage example
```go
package main
import (
	"fmt"
	"github.com/gohouse/gorose/v2"
	_ "github.com/mattn/go-sqlite3"
)
var err error
var engin *gorose.Engin
func init() {
    // Global initialization and reuse of databases
    // The engin here needs to be saved globally, using either global variables or singletons
    // Configuration & gorose. Config {} is a single database configuration
    // If you configure a read-write separation cluster, use & gorose. ConfigCluster {}
	engin, err = gorose.Open(&gorose.Config{Driver: "sqlite3", Dsn: "./db.sqlite"})
    // mysql demo, remeber import mysql driver of github.com/go-sql-driver/mysql
	// engin, err = gorose.Open(&gorose.Config{Driver: "mysql", Dsn: "root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true"})
}
func DB() gorose.IOrm {
	return engin.NewOrm()
}
func main() {
    // Native SQL, return results directly 
    res,err := DB().Query("select * from users where uid>? limit 2", 1)
    fmt.Println(res)
    affected_rows,err := DB().Execute("delete from users where uid=?", 1)
    fmt.Println(affected_rows, err)

    // orm chan operation, fetch one row
    res, err := DB().Table("users").First()
    // res's type is map[string]interface{}
    fmt.Println(res)
    
    // rm chan operation, fetch more rows
    res2, _ := DB().Table("users").Get()
    // res2's type is []map[string]interface{}
    fmt.Println(res2)
}
```

## usage advise
Gorose provides data object binding (map, struct), while supporting string table names and map data return. It provides great flexibility.

It is suggested that data binding should be used as a priority to complete query operation, so that the type of data source can be controlled.
Gorose provides default `gorose. Map'and `gorose. Data' types to facilitate initialization of bindings and data

## Configuration and link initialization
Simple configuration
```go
var configSimple = &gorose.Config{
	Driver: "sqlite3", 
	Dsn: "./db.sqlite",
}
```
More configurations, you can configure the cluster, or even configure different databases in a cluster at the same time. The database will randomly select the cluster database to complete the corresponding reading and writing operations, in which master is the writing database, slave is the reading database, you need to do master-slave replication, here only responsible for reading and writing.
```go
var config = &gorose.ConfigCluster{
    Master:       []&gorose.Config{}{configSimple}
    Slave:        []&gorose.Config{}{configSimple}
    Prefix:       "pre_",
    Driver:       "sqlite3",
}
```
Initial usage
```go
var engin *gorose.Engin
engin, err := Open(config)

if err != nil {
    panic(err.Error())
}
```

## Native SQL operation (add, delete, check), session usage
Create user tables of `users`
```sql
DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
	 "uid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "name" TEXT NOT NULL,
	 "age" integer NOT NULL
);

INSERT INTO "users" VALUES (1, 'gorose', 18);
INSERT INTO "users" VALUES (2, 'goroom', 18);
INSERT INTO "users" VALUES (3, 'fizzday', 18);
```
define table struct
```go
type Users struct {
	Uid  int    `gorose:"uid"`
	Name string `gorose:"name"`
	Age  int    `gorose:"age"`
}
// Set the table name. If not, use struct's name by default
func (u *Users) TableName() string {
	return "users"
}
```
Native query operation
```go
// Here is the structure object to be bound
// If you don't define a structure, you can use map, map example directly
// var u = gorose.Data{}
// var u = gorose.Map{}  Both are possible.
var u Users
session := engin.NewSession()
// Here Bind () is used to store results. If you use NewOrm () initialization, you can use NewOrm (). Table (). Query () directly.
_,err := session.Bind(&u).Query("select * from users where uid=? limit 2", 1)
fmt.Println(err)
fmt.Println(u)
fmt.Println(session.LastSql())
```
Native inesrt delete update
```go
session.Execute("insert into users(name,age) values(?,?)(?,?)", "gorose",18,"fizzday",19)
session.Execute("update users set name=? where uid=?","gorose",1)
session.Execute("delete from users where uid=?", 1)
```
## Object Relational Mapping, the Use of ORM  

- 1. Basic Chain Usage  

```go
var u Users
db := engin.NewOrm()
err := db.Table(&u).Fields("name").AddFields("uid","age").Distinct().Where("uid",">",0).OrWhere("age",18).
	Group("age").Having("age>1").OrderBy("uid desc").Limit(10).Offset(1).Select()
```

- 2. If you don't want to define struct and want to bind map results of a specified type, you can define map types, such as
```go
type user gorose.Map
// Or the following type definitions can be parsed properly
type user2 map[string]interface{}
type users3 []user
type users4 []map[string]string
type users5 []gorose.Map
type users6 []gorose.Data
```
Start using map binding
```go
db.Table(&user).Select()
db.Table(&users4).Limit(5).Select()
```
> Note: If the slice data structure is not used, only one piece of data can be obtained.  

---
The gorose. Data used here is actually the `map [string] interface {}'type.

And `gorose. Map'is actually a `t. MapString' type. Here comes a `t'package, a golang basic data type conversion package. See http://github.com/gohouse/t for more details.  


- 3. laravel's `First()`,`Get()`, Used to return the result set   
That is to say, you can even pass in the table name directly without passing in various bound structs and maps, and return two parameters, one is the `[] gorose. Map `result set, and the second is `error', which is considered simple and rude.

Usage is to replace the `Select ()'method above with Get, First, but `Select ()' returns only one parameter.


- 4. orm Select Update Insert Delete  
```go
db.Table(&user2).Limit(10.Select()
db.Table(&user2).Where("uid", 1).Data(gorose.Data{"name","gorose"}).Update()
db.Table(&user2).Data(gorose.Data{"name","gorose33"}).Insert()
db.Table(&user2).Data([]gorose.Data{{"name","gorose33"},"name","gorose44"}).Insert()
db.Table(&user2).Where("uid", 1).Delete()
```

## Final SQL constructor, builder constructs SQL of different databases
Currently supports mysql, sqlite3, postgres, oracle, mssql, Clickhouse and other database drivers that conform to `database/sql` interface support  
In this part, users are basically insensitive, sorted out, mainly for developers can freely add and modify related drivers to achieve personalized needs.  

## binder, Data Binding Objects  
This part is also user-insensitive, mainly for incoming binding object parsing and data binding, and also for personalized customization by developers.  

## Modularization
Gorose 2.0 is fully modular, each module encapsulates the interface api, calling between modules, through the interface, the upper layer depends on the lower layer

- Main module  
    - engin  
    gorose Initialize the configuration module, which can be saved and reused globally  
    - session  
    Really operate the database underlying module, all operations, will eventually come here to obtain or modify data.   
    - orm  
    Object relational mapping module, all ORM operations, are done here    
    - builder  
    Building the ultimate execution SQL module, you can build any database sql, but to comply with the `database / SQL ` package interface  
- sub module  
    - driver  
    The database driver module, which is dependent on engin and builder, does things according to the driver    
    - binder  
    Result Set Binding Module, where all returned result sets are located   

The above main modules are relatively independent and can be customized and replaced individually, as long as the interface of the corresponding modules is realized.    

## Best Practices
sql
```sql
DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
	 "uid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "name" TEXT NOT NULL,
	 "age" integer NOT NULL
);

INSERT INTO "users" VALUES (1, 'gorose', 18);
INSERT INTO "users" VALUES (2, 'goroom', 18);
INSERT INTO "users" VALUES (3, 'fizzday', 18);
```
Actual Code
```go
package main

import (
	"fmt"
	"github.com/gohouse/gorose/v2"
	_ "github.com/mattn/go-sqlite3"
)

type Users struct {
    Uid int64 `gorose:"uid"`
    Name string `gorose:"name"`
    Age int64 `gorose:"age"`
    Xxx interface{} `gorose:"-"` // This field is ignored in ORM
}

func (u *Users) TableName() string {
	return "users"
}

var err error
var engin *gorose.Engin

func init() {
    // Global initialization and reuse of databases
    // The engin here needs to be saved globally, using either global variables or singletons
    // Configuration & gorose. Config {} is a single database configuration
    // If you configure a read-write separation cluster, use & gorose. ConfigCluster {}
	engin, err = gorose.Open(&gorose.Config{Driver: "sqlite3", Dsn: "./db.sqlite"})
}
func DB() gorose.IOrm {
	return engin.NewOrm()
}
func main() {
	// A variable DB is defined here to reuse the DB object, and you can use db. LastSql () to get the SQL that was executed last.
	// If you don't reuse db, but use DB () directly, you create a new ORM object, which is brand new every time.
	// So reusing DB must be within the current session cycle
	db := DB()
	
	// fetch a row
	var u Users
	// bind result to user{}
	err = db.Table(&u).Fields("uid,name,age").Where("age",">",0).OrderBy("uid desc").Select()
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(u, u.Name)
	fmt.Println(db.LastSql())
	
	// fetch multi rows
	// bind result to []Users, db and context condition parameters are reused here
	// If you don't want to reuse, you can use DB() to open a new session, or db.Reset()
	// db.Reset() only removes contextual parameter interference, does not change links, DB() will change links.
	var u2 []Users
	err = db.Table(&u2).Limit(10).Offset(1).Select()
	fmt.Println(u2)
	
	// count
	var count int64
	// Here reset clears the parameter interference of the upper query and can count all the data. If it is not clear, the condition is the condition of the upper query.
	// At the same time, DB () can be called new, without interference.
	count,err = db.Reset().Count()
	// or
	count, err = DB().Table(&u).Count()
	fmt.Println(count, err)
}
```

## Advanced Usage

- Chunk Data Fragmentation, Mass Data Batch Processing (Cumulative Processing)  

   ` When a large amount of data needs to be manipulated, the chunk method can be used if it is unreasonable to take it out at one time and then operate it again.  
        The first parameter of chunk is the amount of data specified for a single operation. According to the volume of business, 100 or 1000 can be selected.  
        The second parameter of chunk is a callback method for writing normal data processing logic  
        The goal is to process large amounts of data senselessly  
        The principle of implementation is that each operation automatically records the current operation position, and the next time the data is retrieved again, the data is retrieved from the current position.
        `
	```go
	User := db.Table("users")
	User.Fields("id, name").Where("id",">",2).Chunk(2, func(data []gorose.Data) error {
	    // for _,item := range data {
	    // 	   fmt.Println(item)
	    // }
	    fmt.Println(data)
        
        // don't forget return error or nil
        return nil
	})

	// print result:  
	// map[id:3 name:gorose]
	// map[id:4 name:fizzday]
	// map[id:5 name:fizz3]
	// map[id:6 name:gohouse]
	[map[id:3 name:gorose] map[name:fizzday id:4]]
	[map[id:5 name:fizz3] map[id:6 name:gohouse]]
	```
    
- Loop Data fragmentation, mass data batch processing (from scratch)   

	` Similar to chunk method, the implementation principle is that every operation is to fetch data from the beginning.
	Reason: When we change data, the result of the change may affect the result of our data taking as where condition, so we can use Loop.`
    ```go
	User := db.Table("users")
	User.Fields("id, name").Where("id",">",2).Loop(2, func(data []gorose.Data) error {
	    // for _,item := range data {
	    // 	   fmt.Println(item)
	    // }
	    // here run update / delete  
        
        // don't forget return error or nil
        return nil
	})
	```
    
- where nested  

	```go
	// SELECT  * FROM users  
	//     WHERE  id > 1 
	//         and ( name = 'fizz' 
	//             or ( name = 'fizz2' 
	//                 and ( name = 'fizz3' or website like 'fizzday%')
	//                 )
	//             ) 
	//     and job = 'it' LIMIT 1
	User := db.Table("users")
	User.Where("id", ">", 1).Where(func() {
	        User.Where("name", "fizz").OrWhere(func() {
	            User.Where("name", "fizz2").Where(func() {
	                User.Where("name", "fizz3").OrWhere("website", "like", "fizzday%")
	            })
	        })
	    }).Where("job", "it").First()
	```

## realease log
- v2.1.x: 
    * update join with auto table prefix  
    * add query return with []map[string]interface{}  
- v2.0.0: new version, new structure  

## Upgrade Guide
### from 2.0.x to 2.1.x  
- change `xxx.Join("pre_tablename")` into `xxx.Join("tablename")`,the join table name auto prefix  
- change `err:=DB().Bind().Query()` into `res,err:=DB().Query()` with multi return,leave the `Bind()` method as well  
### from 1.x to 2.x
install it for new  


## Jetbrains non-commercial sponsorship  
[![](https://www.jetbrains.com/shop/static/images/jetbrains-logo-inv.svg)](https://www.jetbrains.com/?from=gorose)
-----
## pay me a coffee
wechat|alipay|[paypal: click](https://www.paypal.me/fizzday)
---|---|---
<img src="imgs/wechat.png" width="300">|<img src="imgs/alipay.png" width="300"> | <a href="https://www.paypal.me/fizzday"><img src="imgs/paypal.png" width="300"></a> 

- pay list  

total | avator 
---|---
￥100 | [![](https://avatars1.githubusercontent.com/u/53846155?s=96&v=4)](https://github.com/sanjinhub)  

