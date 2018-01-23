最近迷恋上了go语言, 对go有种特别的好感.   
可是, 由于之前做了很久的php开发, 刚开始用go做web开发有点不太习惯, 也许是之前的 laravel 框架的 eloquent orm过于好用, 所以, 导致了使用go做web开发的各种不适应.  
于是, 想着找一个类似的orm用用, 找遍了go生态圈, 发现了很多知名的 go orm, 诸如: gorm, xorm, sqlx等, 发现没有一个是我的feel, 在体验了百般不爽之于, 痛定思痛, 就自己撸了个 go orm, gorose 就这么诞生了.  
gorose, 是一个mini的 go orm, 也可以说是 golang 版本的 laravel eloquent, 因为喜欢这种feel, 就着手撸了起来, 经过一个礼拜的调教, 初版上了线, 看看效果:  
## gorose链接数据库
```go
// 开启一个链接
db := gorose.Open(config.DbConfig, "mysql_dev")
// 执行完毕后关闭数据库 DB
defer db.Close()
```
## laravel般的简单查询
```go
db.Table("userinfo").First()
```
解析的sql为: `select * from userinfo limit 1`  
是不是很熟悉的感觉, 更熟悉的还在后边

## 多条件链式查询
```go
db.Table("userinfo").Where("id","<",10).Order("id desc").Get()
```
解析的sql为: `select * from userinfo where id<10 order by id desc`   

## 原生查询
```go
db.Query("select * from userinfo")
db.Query("select * from userinfo where id>?", 1)
```   
是不是php orm 的feel又回来了, 没错, 不仅仅如此, eloquent 的大多用法, 都可以在这里直接使用, 更多用法  
- 请看 [github.com/gohouse/gorose](https://github.com/gohouse/gorose)  
- 或者 [点击加入qq群: 470809220](https://jq.qq.com/?_wv=1027&k=5JJOG9E) 慢慢撩~~~

------
powered by [fizzday](http://fizzday.net)(星期八)