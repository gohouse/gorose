# gorose
gorose(go orm), a mini database orm for golang , which inspired by the famous php framwork laravle's eloquent. it will be friendly for php developer and python or ruby developer

- [english document](https://github.com/gohouse/gorose)
- [中文文档](https://github.com/gohouse/gorose/blob/master/README-ZH_CN.md)
## install
- install gorose
```go
go get github.com/gohouse/gorose
```
## config and sample
- multi config 
```go
import "github.com/gohouse/gorose"

var dbConfig = map[string]map[string]string {
	"mysql": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "test",
		"charset":  "utf8",
		"protocol": "tcp",
	},
	"mysql_dev": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "gorose",
		"charset":  "utf8",
		"protocol": "tcp",
	},
}

gorose.Open(dbConfig, "mysql")

var db gorose.Database

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
```sql
SELECT  * FROM user  
    WHERE  id > '1' 
        and ( name = 'fizz' 
            or ( name = 'fizz2' 
                and ( name = 'fizz3' or website like '%fizzday%')
                )
            ) 
    and job = 'it' LIMIT 1
```  

#### transaction
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
```sql
insert into user (age, job) values (17, 'it3')
insert into user (age, job) values (17, 'it3') (17, 'it4')
```

- delete  
```go
User.Where("id", 5).Delete()
```
parse sql result: `delete from user where id=5`

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

## TODO (finish)
- list  
[x] where nested  
[x] transaction union (auto begin, rollback or commit) 
- sample  
```go
db.Where(func(){
	db.Where().OrWhere(func() *db{
		return db.Where().OrWhere()
	})
})
```
- transaction
```go
db.Transaction(func(){
	db.Table("user").Data().Where().Update()
	db.Table("card").Data().Insert()
})
```

## TODO  
[] Separation of reading and writing

------------
#### [点击查看最新更新动态](https://github.com/gohouse/gorose)