    gorose(go orm), 一个小巧强悍的go语言数据库操作orm, 灵感来源于laravel的数据库操作orm, 也就是eloquent, php、python、ruby开发者, 都会喜欢上这个orm的操作方式, 主要是链式操作比较风骚

- [English Document](https://github.com/gohouse/gorose)
- [中文文档](https://github.com/gohouse/gorose/blob/master/README-ZH_CN.md)

## 安装
- 安装 gorose
```go
go get github.com/gohouse/gorose
```
- 安装gorose中使用到的函数工具包
```go
go get github.com/gohouse/utils
```

## 配置和示例
- 多个数据库连接配置  

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
- 简单的但数据库配置  

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

## 用法示例
### 查询
#### 原生sql语句查询
```go
db.Query("select * from user where id = 1")
```
#### 链式调用查询  

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
得到sql结果: 
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

#### 更多链式查询示例
- 获取user表对象
```go
User := db.Table("user")
```
- 查询一条
```go
User.First()
// 或者
db.Fisrt()
```
parse sql result: `select * from user limit 1`  

- count统计
```go
User.Count("*")
// 或(下同) 
db.Count("*")
```
最终执行的sql为: `select count(*) as count from user`  

- max
```go
User.Max("age")
```
最终执行的sql为: `select max(age) as max from user`  

- min
```go
User.Min("age")
```
最终执行的sql为: `select min(age) as min from user`  

- avg
```go
User.Avg("age")
```
最终执行的sql为: `select avg(age) as avg from user`  

- distinct
```go
User.Fields("id, name").Distinct()
```
最终执行的sql为: `select distinct id,name from user`  

#### 嵌套where的查询 (where nested)
```go
db.Table("user").Where("id", ">", 1).Where(func() {
		db.Where("name", "fizz").OrWhere(func() {
			db.Where("name", "fizz2").Where(func() {
				db.Where("name", "fizz3").OrWhere("website", "fizzday")
			})
		})
	}).Where("job", "it").First()
```
最终执行的sql为: 
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

#### 事务
- 标准用法  

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
- 简单用法, 用闭包实现, 自动开始事务, 回滚或提交事务  

```go
db.Transaction(func() {
    db.Execute("update area set job='sadf' where id=14")
    db.Table("area").Data(map[string]interface{}{"names": "fizz3", "age": 3}).Insert()
    db.Table("area").Data(map[string]interface{}{"names": "fizz3", "age": 3}).Where("id",10).Update()
})
```

### 增删改操作
#### 原生sql字符串
```go
db.Execute("update user set job='it2' where id=3")
```
#### 链式调用  

```go
db.Table("user").
	Data(map[string]interface{}{"age":17, "job":"it3"}).
    Where("id", 1).
    OrWhere("age",">",30).
    Update()
```
最终执行的sql为: `update user set age=17, job='ite3' where (id=1) or (age>30)`  

#### 更多增删改的用法
- insert   

```go
User.Data(map[string]interface{}{"age":17, "job":"it3"}).Insert()
User.Data([]map[string]interface{}{{"age":17, "job":"it3"},{"age":17, "job":"it4"}).Insert()
```
最终执行的sql为:  

```go
insert into user (age, job) values (17, 'it3')
insert into user (age, job) values (17, 'it3') (17, 'it4')
```

- delete   
 
```go
User.Where("id", 5).Delete()
```
最终执行的sql为: `delete from user where id=5`

## 切换数据库连接  

```go
// 连接最开始配置的第二个链接(mysql_dev是key)
db.Connect("mysql_dev").Table().First()
// 或者直接输入连接配置
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

## TODO
[] 读写分离

------------
#### [点击查看最新更新动态](https://github.com/gohouse/gorose)