## brief introduction
(gorose, 最风骚的go orm, 开箱即用, 一分钟上手, 链式操作, 让golang操作数据库成为一种享受, 妈妈再也看不到我处理数据的痛苦了)  
gorose(go orm), a mini database orm for golang , which inspired by the famous php framwork laravle's eloquent. it will be friendly for php developer and python or ruby developer  
目前提供5大数据库驱动, mysql,sqlite,postgres,oracle,mssql, 同时可以自由更换驱动

## document
- [中文文档](https://www.kancloud.cn/fizz/gorose)

## 随时在线交流心得
[点击加入qq群: 470809220](https://jq.qq.com/?_wv=1027&k=5JJOG9E)  

## install
```go
go get github.com/gohouse/gorose
```

## quick scan
```go
db.Table("tablename").First()
db.Table("tablename").Distinct().Where("id", ">", 5).Get()
db.Table("tablename").Fields("id, name, age, job").Group("job").Limit(10).Offset(20).Order("id desc").Get()
// query str
db.Query("select * from user")
db.Query("select * from user where id=?", 1)
db.Execute("update users set name=? where id=?", "fizz", 1)
```

## install
- install gorose
```go
go get github.com/gohouse/gorose
```
## config and sample
- multi config 
```go
import (
	"github.com/gohouse/gorose"
	"fmt"
)

var dbConfig = map[string]map[string]string {
	"mysql_master": {
        "host":     "localhost",
        "username": "root",
        "password": "root",
        "port":     "3306",
        "database": "test",
        "charset":  "utf8",
        "protocol": "tcp",
        "driver":   "mysql", // 数据库驱动(mysql,sqlite,postgres,oracle,mssql)
	},
	"mysql_dev": {
        "host":     "localhost",
        "username": "root",
        "password": "root",
        "port":     "3306",
        "database": "test",
        "charset":  "utf8",
        "protocol": "tcp",
        "driver":   "mysql",
	},
}

db := gorose.Open(config.DbConfig, "mysql_dev")
// close DB
defer db.Close()

func main() {
    res := db.Table("users").First()
    
    fmt.Println(res)
}
```
- single config
```go
gorose.Open(map[string]string {
                "host":     "localhost",
                "username": "root",
                "password": "",
                "port":     "3306",
                "database": "test",
                "charset":  "utf8",
                "protocol": "tcp",
                "driver":   "mysql",
            })
```
## example
### query
#### Native string
```go
db.Query("select * from user where id = 1")
```
#### call chaining
```go
db.Table("user").
    Field("id, name").  // field
    Where("id",">",1).  // simple where
    Where(map[string]interface{}{"name":"fizzday", "age":18}).  // where object
    Where([]map[string]interface{}{{"website", "like", "fizz"}, {"job", "it"}}).    // multi where
    Where("head = 3 or rate is not null").  // where string
    OrWhere("cash", "1000000"). // or where ...
    OrWhere("score", "between", []string{50, 80}).  // between
    OrWhere("role", "not in", []string{"admin", "read"}).   // in 
    Group("job").   // group
    Order("age asc").   // order 
    Limit(10).  // limit
    Offset(1).  // offset
    Get()   // fetch multi rows
```
parse sql result: 
```go
select id,name from user 
    where (id>1) 
    and (name='fizzday' and age='18') 
    and ((website like '%fizz%') and (job='it'))
    and (head =3 or rate is not null)
    or (cash = '100000') 
    or (score between '50' and '100') 
    or (role not in ('admin', 'read'))
    group by job 
    order by age asc 
    limit 10 offset 1
```  
#### more query usage
- get user obj
```go
User := db.Table("user")
```
- fetch one row
```go
User.First()
// or
db.Fisrt()
```
parse sql result: `select * from user limit 1`  

- count
```go
User.Count("*")
// or 
db.Count("*")
```
parse sql result: `select count(*) as count from user`  

- max
```go
User.Max("age")
```
parse sql result: `select max(age) as max from user`  

- min
```go
User.Min("age")
```
parse sql result: `select min(age) as min from user`  

- avg
```go
User.Avg("age")
```
parse sql result: `select avg(age) as avg from user`  

- distinct
```go
User.Fields("id, name").Distinct()
```
parse sql result: `select distinct id,name from user`  

#### join
```go
db.Table("user")
    .Join("card","user.id","=","card.user_id")
    .Limit(10)
    .Get()
```
parse sql result: 
```go
select * from user inner join card on user.id=card.user_id limit 10
```
```go
db.Table("user")
    .LeftJoin("card","user.id","=","card.user_id")
    .First()
```
parse sql result: 
```go
select * from user left join card on user.id=card.user_id limit 1
```
> RightJoin : right join

#### where nested (嵌套where)
```go
db.Table("user").Where("id", ">", 1).Where(func() {
		db.Where("name", "fizz").OrWhere(func() {
			db.Where("name", "fizz2").Where(func() {
				db.Where("name", "fizz3").OrWhere("website", "fizzday")
			})
		})
	}).Where("job", "it").First()
```
parse sql result: 
```go
SELECT  * FROM user  
    WHERE  id > '1' 
        and ( name = 'fizz' 
            or ( name = 'fizz2' 
                and ( name = 'fizz3' or website like '%fizzday%')
                )
            ) 
    and job = 'it' LIMIT 1
```  

#### chunk data block
```go
db.Table("users").Fields("id, name").Where("id",">",2).Chunk(2, func(data []map[string]interface{}) {
    // for _,item := range data {
    // 	   fmt.Println(item)
    // }
    fmt.Println(data)
})
```
result:  
```go
// map[id:3 name:gorose]
// map[id:4 name:fizzday]
// map[id:5 name:fizz3]
// map[id:6 name:gohouse]
[map[id:3 name:gorose] map[name:fizzday id:4]]
[map[id:5 name:fizz3] map[id:6 name:gohouse]]
```

### execute
#### Native string
```go
db.Execute("update user set job='it2' where id=3")
```
#### Call chaining
```go
db.Table("user").
	Data(map[string]interface{}{"age":17, "job":"it3"}).
    Where("id", 1).
    OrWhere("age",">",30).
    Update()
```
parse sql result: `update user set age=17, job='ite3' where (id=1) or (age>30)`  

#### more execute usage
- insert  
```go
User.Data(map[string]interface{}{"age":17, "job":"it3"}).Insert()
User.Data([]map[string]interface{}{{"age":17, "job":"it3"},{"age":17, "job":"it4"}).Insert()
```
parse sql result: 
```go
insert into user (age, job) values (17, 'it3')
insert into user (age, job) values (17, 'it3') (17, 'it4')
```

- delete  
```go
User.Where("id", 5).Delete()
```
parse sql result: `delete from user where id=5`

## transaction
- standard using
```go
db.Begin()

res := db.Table("user").Where("id", 1).Data(map[string]interface{}{"age":18}).Update()
if (res == 0) {
	db.Rollback()
}

res2 := db.Table("user").Data(map[string]interface{}{"age":18}).Insert()
if (res2 == 0) {
	db.Rollback()
}

db.Commit()
```
- simple using
```go
db.Transaction(func() {
    db.Execute("update area set job='sadf' where id=14")
    db.Table("area").Data(map[string]interface{}{"names": "fizz3", "age": 3}).Insert()
    db.Table("area").Data(map[string]interface{}{"names": "fizz3", "age": 3}).Where("id",10).Update()
})
```

## Temporary connection
```go
db.Connect("mysql_dev").Table().First()
// or
db.Connect(map[string]string {
                "host":     "localhost",
                "username": "root",
                "password": "",
                "port":     "3306",
                "database": "test",
                "charset":  "utf8",
                "protocol": "tcp",
            }).Table().First()
```  

## get origin connection DB  
```go
gorose.GetDB()
```

## get sql logs or last sql
```go
db.SqlLogs()
db.LastSql()
```

## json return
```go
result := db.Table("user").Get()
jsonStr := db.JsonEncode(result)
```

## TODO (finish)

[] connection pool


------------
#### [ click for getting the news ](https://github.com/gohouse/gorose)
