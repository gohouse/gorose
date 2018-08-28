## 简介
gorose(go orm), 一个小巧强悍的go语言数据库操作orm, 灵感来源于laravel的数据库操作orm, 也就是eloquent, 风骚的链式操作, 会让php、python、ruby开发者毫无抵抗力的爱上gorose

## github
- [https://github.com/gohouse/gorose](https://github.com/gohouse/gorose)

## 先睹为快
```go
db.Table("tablename").First()
db.Table("tablename").Distinct().Where("id", ">", 5).Get()
db.Table("tablename").Fields("id, name, age, job").Group("job").Limit(10).Offset(20).Order("id desc").Get()
// query str
db.Query("select * from user")
db.Query("select * from user where id=?", 1)
db.Execute("update users set name=? where id=?", "fizz", 1)
```

## 用法示例

- 多个数据库连接配置和连接使用简单示例  

```go
package main

import (
	"github.com/gohouse/gorose"
	"fmt"
)

func main() {

    var dbConfig = map[string]interface{} {
        "Default":         "mysql_dev",// 默认数据库配置
        "SetMaxOpenConns": 0,          // (连接池)最大打开的连接数，默认值为0表示不限制
        "SetMaxIdleConns": 1,          // (连接池)闲置的连接数, 默认1
    
        "Connections":map[string]map[string]string{
            "mysql_dev": {// 定义名为 mysql_dev 的数据库配置
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
            "sqlite_dev": {
             "database": "./foo.db",
             "prefix": "",
             "driver": "sqlite3",
            },
    	},
    }

    // 初始化数据库链接, 默认会链接配置中 default 指定的值 
    // 也可以在第二个参数中指定对应的数据库链接, 见下边注释的那一行链接示例
	connection, err := gorose.Open(DbConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer connection.Close()
	
	db := connection.NewDB()
	
    res,err := db.Table("users").First()
    if err!=nil{
    	fmt.Println(err)
    	return
    }
    fmt.Println(res)
}
```

- 简单的数据库配置  

```go
gorose.Open(map[string]string{
                "database": "./foo.db",
                "prefix": "",
                "driver": "sqlite3",
            })
```

- 连接池  
直接指定配置文件中对应的数据即可
```go
"SetMaxOpenConns": 0,        // (连接池)最大打开的连接数，默认值为0表示不限制
"SetMaxIdleConns": -1,        // (连接池)闲置的连接数, 默认-1
```

## 用法示例
### 查询
#### 原生sql语句查询
```go
db.Query("select * from user where id = 1")
db.Query("select * from user where name = ?", "fizz")
```
#### 链式调用查询  

```go
db.Table("user").
    Field("id, name, avg(age) as age_avg").  // field
    Where("id",">",1).  // simple where
    Where("head = 3 or rate is not null").  // where string
    Where(map[string]interface{}{"name":"fizzday", "age":18}).  // where object
    Where([][]interface{}{ {"website", "like", "%fizz%"}, {"job", "it"} }).    // multi where
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
得到sql结果: 
```go
select id,name,count(age) from user 
    where (id>1) 
    and (head =3 or rate is not null)
    and (name='fizzday' and age='18') 
    and ((website like '%fizz%') and (job='it'))
    or (cash = '100000') 
    or (score between '50' and '100') 
    or (role not in ('admin', 'read'))
    group by job 
    having age_avg>1
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
res,err := User.First()
```
最终执行的sql为: `select * from user limit 1`  

- 查询某一个字段的值  
```go
name := User.Value("name")
```

- count统计  
```go
res,err := User.Count("*")
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

#### join  

- 普通示例  
```go
db.Table("user")
    .Join("card","user.id","=","card.user_id")
    .Limit(10)
    .Get()
```
最终执行的sql为: 
```go
select * from user inner join card on user.id=card.user_id limit 10  
```

- 左链接  
```go
db.Table("user a")
    .LeftJoin("card b","a.id","=","b.user_id")
    .First()
```
最终执行的sql为: 
```go
select * from user a left join card b on a.id=b.user_id limit 1
```

> 右链接 : right join  

#### 嵌套where的查询 (where nested)
```go
User.Where("id", ">", 1).Where(func() {
		User.Where("name", "fizz").OrWhere(func() {
			User.Where("name", "fizz2").Where(func() {
				User.Where("name", "fizz3").OrWhere("website", "like", "fizzday%")
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
                and ( name = 'fizz3' or website like 'fizzday%')
                )
            ) 
    and job = 'it' LIMIT 1
```  

#### 分块操作所有数据
- Chunk()  
    > 当需要操作大量数据的时候, 一次性取出再操作, 不太合理, 就可以使用chunk方法  
        chunk的第一个参数是指定一次操作的数据量, 根据业务量, 取100条或者1000条都可以  
        chunk的第二个参数是一个回调方法, 用于书写正常的数据处理逻辑  
        目的是做到, 无感知处理大量数据  
        
    ```go
    User.Fields("id, name").Where("id",">",2).Chunk(2, func(data []map[string]interface{}) {
        // for _,item := range data {
        // 	   fmt.Println(item)
        // }
        fmt.Println(data)
    })
    ```
    打印结果:  
    ```go
    // map[id:3 name:gorose]
    // map[id:4 name:fizzday]
    // map[id:5 name:fizz3]
    // map[id:6 name:gohouse]
    [map[id:3 name:gorose] map[name:fizzday id:4]]
    [map[id:5 name:fizz3] map[id:6 name:gohouse]]
    ```

### 增删改操作

#### 原生sql字符串

```go
res,err := db.Execute("update user set job='it2' where id=3")
res,err := db.Execute("update user set job='it2' where id=?", 3)
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
res,err := User.Data(map[string]interface{}{"age":17, "job":"it3"}).Insert()
User.Data([]map[string]interface{}{ {"age":17, "job":"it3"},{"age":17, "job":"it4"} }).Insert()
```

最终执行的sql为:  

```go
insert into user (age, job) values (17, 'it3')
insert into user (age, job) values (17, 'it3') (17, 'it4')
```
> 获取影响行数和插入id  
    - 获取最后插入id: User.LastInsertId  
    - 获取影响行数(默认返回, 也可以通过此方法获取): User.RowsAffected  

- delete   
 
```go
User.Where("id", 5).Delete()
```
最终执行的sql为: `delete from user where id=5`

## 事务

- 标准用法  

```go
db.Begin()

res,err := db.Table("user").Where("id", 1).Data(map[string]interface{}{"age":18}).Update()
if (res == 0 || err!=nil) {
	db.Rollback()
}

res2,err := db.Table("user").Data(map[string]interface{}{"age":18}).Insert()
if (res2 == 0 || err!=nil) {
	db.Rollback()
}

db.Commit()
```

- 简单用法, 用闭包实现, 自动开始,回滚和提交 事务  

```go
trans,err := db.Transaction(func() (error) (bool,error) {
	
    res1,err := db.Execute("update area set job='sadf' where id=14")
    if err!=nil {
    	return false,err
    }
    if res1==0 {
    	return false,errors.New("update failed")
    }
    
    res2,err := db.Table("area").Data(map[string]interface{}{"names": "fizz3", "age": 3}).Insert()
    if err!=nil {
        return false,err
    }
    if res2==0 {
    	return false,errors.New("Insert failed")
    }
        
    res3,err := db.Table("area").Data(map[string]interface{}{"names": "fizz3", "age": 3}).Where("id",10).Update()
    if err!=nil {
        return false,err
    }
    if res3==0 {
    	return false,errors.New("Update failed")
    }
    
    return true,nil
})
```

## 获取原始连接 DB

```go
DB := gorose.GetDB()
```

### 获取所有sql记录, 或者获取最后一条sql语句

```go
sqllogs := db.SqlLogs
lastsql := db.LastSql
```

## json返回
```go
result := db.Table("user").Get()
jsonStr := db.JsonEncode(result)
```

## 目录说明
```sh
/docs/      ---- 文档目录, 这里包含多个语言的不同使用文档
/drivers/   ---- 不同数据库的驱动目录, 可以自由增加任何其他数据库的目录
/examples/  ---- 使用示例目录, 可以在这里找到大部分的用例
/test/      ---- go testing 自动测试, 包括简单的压力测试
/utils/     ---- 工具包, 放置常用工具函数
/vendor/    ---- 采用glide管理的依赖包
database.go ---- 数据库映射操作的核心文件
glide.yaml  ---- 项目依赖管理的配置文件
gorose.go   ---- 数据库链接,数据库驱动加载核心文件
README.md   ---- 文档说明文件
```

------------
#### [点击查看最新更新动态](https://github.com/gohouse/gorose)
