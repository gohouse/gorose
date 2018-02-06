## brief introduction
gorose(go orm), a mini database orm for golang , which inspired by the famous php framwork laravle's eloquent. it will be friendly for php developer and python or ruby developer  

## document
- [中文文档](https://www.kancloud.cn/fizz/gorose)

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

## config and sample
- multi config 

```go
package main

import (
	"github.com/gohouse/gorose"
	"fmt"
)

func main() {

    var dbConfig = map[string]interface{} {
        "default":         "mysql_dev",// 默认数据库配置
        "SetMaxOpenConns": 0,          // (连接池)最大打开的连接数，默认值为0表示不限制
        "SetMaxIdleConns": 1,          // (连接池)闲置的连接数, 默认1
    
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
        "sqlite_dev": map[string]string{
            "database": "./foo.db",
            "prefix": "",
            "driver": "sqlite3",
        },
    }

    // 初始化数据库链接, 默认会链接配置中 default 指定的值 
    // 也可以在第二个参数中指定对应的数据库链接, 见下边注释的那一行链接示例
    db,err := gorose.Open(dbConfig)
    // db,err := gorose.Open(dbConfig, "sqlite_dev")
    if err!=nil{
        fmt.Println(err)
        return
    }
    res,err := db.Table("users").First()
    if err!=nil{
    	fmt.Println(err)
    	return
    }
    fmt.Println(res)
}
```

- single config
```go
gorose.Open(map[string]string{
                "database": "./foo.db",
                "prefix": "",
                "driver": "sqlite3",
            })
```

- connection poll  
just set the param in config  
```go
"SetMaxOpenConns": 0,        // max open connection, default 0 with no limit
"SetMaxIdleConns": 1,        // connection with free, default 1
```

## example
### query
#### Native string
```go
db.Query("select * from user where id = 1")
db.Query("select * from user where name = ?", "fizz")
```
#### call chaining
```go
db.Table("user").
    Field("id, name, avg(age) as age_avg").  // field
    Where("id",">",1).  // simple where
    Where("head = 3 or rate is not null").  // where string
    Where(map[string]interface{}{"name":"fizzday", "age":18}).  // where object
    Where([]map[string]interface{}{ {"website", "like", "%fizz%"}, {"job", "it"} }).    // multi where
    Where("head = 3 or rate is not null").  // where string
    OrWhere("cash", "1000000"). // or where ...
    OrWhere("score", "between", []string{50, 80}).  // between
    OrWhere("role", "not in", []string{"admin", "read"}).   // in 
    Group("job").   // group
    Having("age_avg>1").   // having
    Order("age asc").   // order 
    Limit(10).  // limit
    Offset(1).  // offset
    Get()   // fetch multi rows
```
parse sql result: 
```go
select id,name from user 
    where (id>1) 
    and (head =3 or rate is not null)
    and (name='fizzday' and age='18') 
    and ((website like '%fizz%') and (job='it'))
    and (head =3 or rate is not null)
    or (cash = '100000') 
    or (score between '50' and '100') 
    or (role not in ('admin', 'read'))
    group by job 
    having age_avg>1
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
```
parse sql result: `select * from user limit 1`  

- fetch a value of one field  
```go
name := User.Value("name")
```

- count
```go
User.Count()
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
db.Table("user a")
    .LeftJoin("card b","a.id","=","b.user_id")
    .First()
```
parse sql result: 
```go
select * from user a left join card b on a.id=b.user_id limit 1
```
> RightJoin : right join

#### where nested (嵌套where)
```go
User.Where("id", ">", 1).Where(func() {
		User.Where("name", "fizz").OrWhere(func() {
			User.Where("name", "fizz2").Where(func() {
				User.Where("name", "fizz3").OrWhere("website", "like", "fizzday%")
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
                and ( name = 'fizz3' or website like 'fizzday%')
                )
            ) 
    and job = 'it' LIMIT 1
```  

#### chunk data block
- Chunk  
```go
User.Fields("id, name").Where("id",">",2).Chunk(2, func(data []map[string]interface{}) {
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
db.Execute("update user set job='it2' where id=?", 3)
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
User.Data([]map[string]interface{}{ {"age":17, "job":"it3"},{"age":17, "job":"it4"} }).Insert()
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
trans := db.Transaction(func() (error) {
	
    res1,err := db.Execute("update area set job='sadf' where id=14")
    if err!=nil {
    	return err
    }
    if res1==0 {
    	return errors.New("update failed")
    }
    
    res2,err := db.Table("area").Data(map[string]interface{}{"names": "fizz3", "age": 3}).Insert()
    if err!=nil {
        return err
    }
    if res2==0 {
    	return errors.New("Insert failed")
    }
        
    res3,err := db.Table("area").Data(map[string]interface{}{"names": "fizz3", "age": 3}).Where("id",10).Update()
    if err!=nil {
        return err
    }
    if res3==0 {
    	return errors.New("Update failed")
    }
    
    return nil
})
```

## get origin connection DB  
```go
gorose.GetDB()
```

## get sql logs or last sql
```go
sqllogs := db.SqlLogs()
lastsql := db.LastSql()
```

## json return
```go
result := db.Table("user").Get()
jsonStr := db.JsonEncode(result)
```

------------
#### [ click for getting the news ](https://github.com/gohouse/gorose)
