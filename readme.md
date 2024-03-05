# GoRose ORM V3
PHP Laravel ORM 的 go 实现, 与 laravel 官方文档保持一致 https://laravel.com/docs/10.x/queries .  
分为 go 风格 (struct 结构绑定用法) 和 php 风格 (map 结构用法).  
php 风格用法, 完全可以使用 laravel query builder 的文档做参考, 尽量做到 1:1 还原.  

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
	Id    int64  `db:"id,pk" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
    
	TableName string `db:"users" json:"-"` // 定义表名字,等同于 func (User) TableName() string {return "users"}, 二选一即可
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
db().Table("users").Select("id","name","email").Where("id", "=", 1).OrWhere("name", "test").GroupBy("id").Having("id", ">", 1).Limit(2).Offset(2).OrderBy("id", "desc").Get()
// 等同于
var users []User
db().Where("id", "=", 1).OrWhere("name", "test").GroupBy("id").Having("id", ">", 1).Limit(2).Offset(2).OrderBy("id", "desc").To(&users)
```
由此可以看出, 除了对 表 模型的绑定区别, 其他方法通用, table 参数, 可以是字符串, 也可以是 User 结构体(db().Table(User{}, "u"))

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
db().Transaction(func(tx dbgo.Database) error {
    tx.Insert(&user)
    tx.Update(&user)
    tx.To(&user)
}

// 手动事务
var tx = db()
tx.Begin()
tx.Rollback()
tx.Commit()

// 全自动嵌套事务
var tx = db()
tx.Transaction(func(db1 dbgo.Database) error {
    db1.Insert(&user)
    ...
    // 自动子事务
    tx.Transaction(func(db2 dbgo.Database) error {
        db2.Update(&user)
        ...
    }
}

// 手动嵌套事务
var tx = db()
tx.Begin()
// 自动子事务
tx.Begin() // 自动 savepoint 子事务
tx.Rollback()   // 自动回滚到上一个 savepoint
// 手动子事务
tx.SavePoint("savepoint1")    // 手动 savepoint 到 savepoint1(自定义名字)
tx.RollbackTo("savepoint1") // 手动回滚到自定义的 savepoint
tx.Commit()
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
- [x] Page  
- [x] LastSql  
- [x] To  

## update
在插入和更新数据时,如果使用 struct 模型作为数据对象的时候, 默认忽略类型零值,如果想强制写入,则可以从第二个参数开始传入需要强制写入的字段即可,如:
```go
var user = User{Id: 1, Name: "test"}
// 这里不会对 sex 做任何操作, update user set name="test" where id=1
db().Update(&user)
// 这里会强制将sex更改为0, update user set name="test", sex=0 where id=1
db().Update(&user, "sex")
```

## join
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


