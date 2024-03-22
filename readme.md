# GoRose ORM V3
PHP Laravel ORM 的 go 实现, 与 laravel 官方文档保持一致 https://laravel.com/docs/10.x/queries .  
分为 go 风格 (struct 结构绑定用法) 和 php 风格 (map 结构用法).  
php 风格用法, 完全可以使用 laravel query builder 的文档做参考, 尽量做到 1:1 还原.  

## 安装
目前还处于beta阶段, 请谨慎使用.
```shell
# go.mod

require github.com/gohouse/gorose/v3 master
```

## 概览
go风格用法
```go
package main

import (
    gorose "github.com/gohouse/gorose/v3"
    // 引入驱动,这里可以将驱动独立出去,避免多驱动的不必要引入,只要实现了 gorose.IDriver 接口即可,理论上可以支持任意数据库
    _ "github.com/gohouse/gorose/v3/drivers/mysql"
)

type User struct {
	Id    int64  `db:"id,pk"`   // 这里的 pk 是指主键
	Name  string `db:"name"`
	Email string `db:"email"`
    
    // 定义表名字,等同于 func (User) TableName() string {return "users"}, 二选一即可
	// TableName string `db:"users" json:"-"` 
}
func (User) TableName() string {
    return "users"
}

var gr = gorose.Open("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true")

func db() *gorose.Database {
	return gr.NewDatabase()
}
func main() {
    // select id,name,email from users limit 1;
    var user User
    db().To(&user)

    // insert into users (name,email) values ("test","test@test.com");
    var user = User{Name: "test", Email: "test@test.com"}}
    db().Insert(&user)

    // update users set name="test2" where id=1;
    var user = User{Id: 1, Name: "test2"}
    db().Update(&user)

    // delete from users where id=1;
    var user = User{Id: 1}
    db().Delete(&user) 
}
```
php风格用法
```go
// select id,name,email from users where id=1 or name="test" group by id having id>1 order by id desc limit 2 offset 2
db().Table("users").
    Select("id","name","email").
    Where("id", "=", 1).OrWhere("name", "test").
    GroupBy("id").Having("id", ">", 1).
    Limit(2).Offset(2).OrderBy("id", "desc").
    Get()
// 等同于
var users []User
db().Where("id", "=", 1).OrWhere("name", "test").
    GroupBy("id").Having("id", ">", 1).
    Limit(2).Offset(2).OrderBy("id", "desc").
    To(&users)
```
由此可以看出, 除了对 表 模型的绑定区别, 其他方法通用

## 配置
单数据库连接, 可以直接同官方接口一样用法
```go
var gr = gorose.Open("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true")
```
也可以用
```go
var conf1 = gorose.Config{
    Driver:          "mysql",
    DSN:             "root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true",
    Prefix:          "tb_",
    Weight:          0,
    MaxIdleConns:    0,
    MaxOpenConns:    0,
    ConnMaxLifetime: 0,
    ConnMaxIdleTime: 0,
}
var gr = gorose.Open(conf)
```
或者使用读写分离集群,事务内,自动强制从写库读数据
```go
var gr = gorose.Open(
    gorose.ConfigCluster{
        WriteConf: []gorose.Config{
            conf1,
            conf2,
        },
        ReadConf: []gorose.Config{
            conf3,
            conf4,
        }
    },
)
```

## 事务
```go
// 全自动事务, 有错误会自动回滚, 无错误会自动提交
// Transaction 方法可以接收多个操作, 同用一个 tx, 方便在不同方法内处理同一个事务
db().Transaction(func(tx gorose.TxHandler) error {
    tx().Insert(&user)
    tx().Update(&user)
    tx().To(&user)
}

// 手动事务
tx = db().Begin()
tx().Insert(&user)
tx().Update(&user)
tx().To(&user)
tx().Rollback()
tx().Commit()

// 全自动嵌套事务
db().Transaction(func(tx gorose.TxHandler) error {
    tx().Insert(&user)
    ...
    // 自动子事务
    tx().Transaction(func(tx2 gorose.TxHandler) error {
        tx2().Update(&user)
        ...
    }
}

// 手动嵌套事务
var tx = db().Begin()
// 自动子事务
tx().Begin() // 自动 savepoint 子事务
...
tx().Rollback()   // 自动回滚到上一个 savepoint
...
// 手动子事务
tx().SavePoint("savepoint1")    // 手动 savepoint 到 savepoint1(自定义名字)
...
tx().RollbackTo("savepoint1") // 手动回滚到自定义的 savepoint
tx().Commit()
```

## update
在插入和更新数据时,如果使用 struct 模型作为数据对象的时候, 默认忽略类型零值,如果想强制写入,则可以从第二个参数开始传入需要强制写入的字段即可,如:
```go
var user = User{Id: 1, Name: "test"}
// 这里不会对 sex 做任何操作, 
//update user set name="test" where id=1
db().Update(&user)
// 这里会强制将sex更改为0
//update user set name="test", sex=0 where id=1
db().Update(&user, "sex")
// 等同于
db().Table(&user).Where("id", 1).Update(map[string]any{"name": "test", "sex": 0}))
```
如果没有where条件,则会自动添加tag中指定了pk的字段作为条件,如: `db:"id,pk"`, 因为指定了 pk,如果 id 的值不为0值, 则 id 会作为主键条件更新

## insert
参考 update

## delete
```go
var user = User{Id: 1}
db().Delete(&user)
// 等同于
db().Table(&user).Where("id", 1).Delete()
// 要加上字段0值条件,只需要传入第二个字段,如:
// delete from users where id=1 and sex=0 and name=""
db().Delete(&user, "sex", "name")
```

## table
- 参数  
  table 参数, 可以是字符串, 也可以是 User 结构体
```go
db().Table(User{})
db().Table("users")
```
- 取别名
```go
db().Table(User{}, "u")
db().Table("users", "u")
```
- 结果集查询
```go
sub := db().Table("users").Select("id", "name")
db().Table(sub).Where("id", ">", 1).Get()
```

## join
- 简单用法
```go
db().Table("users").Join(UserInfo{}, "user.id", "=", "user_info.user_id").Get()
```

- 取别名
```go
type UserInfo struct {
    UserId      int64   `db:"user_id"`
    TableName   string  `db:"user_info"`
}
// select * from users a inner join user_info b on a.id=b.user_id
db().Table("users", "u").Join(gorose.As(UserInfo{}, "b"), "u.id", "=", "b.user_id").Get()
// 等同于
db().Table(User{}, "u").Join(gorose.As("user_info", "b"), "u.id", "=", "b.user_id").Get()
```
`gorose.As(UserInfo{}, "b")` 中, `user_info` 取别名 `b`

- 复杂用法
```go
db().Table("users").Join(UserInfo{}, func(wh gorose.IJoinOn) {
    wh.On("a.id", "b.user_id").OrOn("a.sex", "b.sex")
}).Get()
```

## where
- sub query
```go
// where id in (select user_id from user_info)
sub := db().Table("user_info").Select("user_id")
xxx.Where("id", "in", sub).Get()
```

- where nested
```go
// where id>1 and (sex=1 or sex=2)
xxx.Where("id",">", 1).Where(func(wh gorose.IWhere) {
    wh.Where("sex", 1).OrWhere("sex", 2)
})
```

## Pluck
返回两列数据到一个map中,第一列为value,第二列为key
```go
// select id,name from users
db().Table("users").Pluck("name", "id")
// 返回 map[<id>]<name>
// 得到 map[int64]string{1: "张三", 2: "李四"}
```

## List
返回一列数据到一个数组中
```go
// select id,name from users
db().Table("users").List("id")
// 返回 []<id>
// 得到 []int64{1,2,3}
```

## To 查询结果绑定到对象
使用结构体字段作为 select 字段  
使用结构体字段值作为 where 条件  
查询结果绑定到结构体,支持一条或多条
```go
// 查询一条数据
var user User
db().To(&user)

// 查询条件,一条数据
// select id,name,email from users where id=1
var user = User{Id: 1}
db().To(&user)

// 查询多条数据
var users []User
db().To(&users)

// 查询条件,多条数据
var users []User
db().Where("id", ">", 1).To(&users)
```

## Bind 查询结果绑定到对象
仅仅用作查询结果的绑定   
结构体字段,不作为查询字段和条件  
常用作join或者手动指定字段查询绑定
```go
type Result struct {
    Id    int64 `db:"id"`
    Aname string `db:"aname"`
    Bname string `db:"bname"`
}
var res Result
// select a.id, a.name aname, b.name bname from a inner join b on a.id=b.id where a.id>1
db().Table("a").Join("b", "a.id","b.aid").Select("a.id", "a.name aname","b.name bname").Where("a.id", ">", 1).Bind(&res)
```
查询字段的显示名字一定要跟 结构体的字段 tag(db) 名字相同, 否则不会被赋值  
字段数量可以不一样

## ListTo,PluckTo,ValueTo
```go
var list []int
db().Table("users").ListTo("age", &list)
var pluck map[int64]string
db().Table("users").PluckTo("name","id", &pluck)
var value int
db().Table("users").ValueTo("age", &value)
```
## SumTo,MaxTo,MinTo
```go
var sum int
db().Table("users").SumTo("age", &sum)
var max int
db().Table("users").MaxTo("age", &max)
var min int
db().Table("users").MinTo("age", &min)
```


## 已经支持的 laravel query builder 方法
- [x] Table  
- [x] Select  
- [x] SelectRaw  
- [x] AddSelect  
- [x] Join  
- [x] GroupBy  
- [x] Having  
- [x] HavingRaw  
- [x] OrHaving  
- [x] OrHavingRaw  
- [x] OrderBy  
- [x] Limit  
- [x] Offset

- [x] Where
- [x] WhereRaw
- [x] OrWhere
- [x] OrWhereRaw
- [x] WhereBetween
- [x] OrWhereBetween
- [x] WhereNotBetween
- [x] OrWhereNotBetween
- [x] WhereIn
- [x] OrWhereIn
- [x] WhereNotIn
- [x] OrWhereNotIn
- [x] WhereNull
- [x] OrWhereNull
- [x] WhereNotNull
- [x] OrWhereNotNull
- [x] WhereLike
- [x] OrWhereLike
- [x] WhereNotLike
- [x] OrWhereNotLike
- [x] WhereExists
- [x] WhereNotExists
- [x] WhereNot

- [x] Get  
- [x] First  
- [x] Find  
- [x] Insert  
- [x] Update  
- [x] Delete  
- [x] Max  
- [x] Min  
- [x] Sum  
- [x] Avg  
- [x] Count  

- [x] InsertGetId  
- [x] Upsert  
- [x] InsertOrIgnore  
- [x] Exists
- [x] DoesntExist

- [x] Pluck  
- [x] List  
- [x] Value  

- [x] Increment  
- [x] Decrement  
- [x] IncrementEach  
- [x] DecrementEach  
- [x] Truncate  
- [x] Union  
- [x] UnionAll  
- [x] SharedLock  
- [x] LockForUpdate  

## 额外增加的 api
- [x] Replace
- [x] Page  
- [x] LastSql  
- [x] To  
- [x] Bind  
- [x] ListTo  
- [x] PluckTo  
- [x] ValueTo  
- [x] SumTo  
- [x] MaxTo  
- [x] MinTo  
