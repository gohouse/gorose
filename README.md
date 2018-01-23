## brief introduction
(gorose, 最风骚的go orm, 开箱即用, 一分钟上手, 链式操作, 让golang操作数据库成为一种享受, 妈妈再也看不到我处理数据的痛苦了)  
gorose(go orm), a mini database orm for golang , which inspired by the famous php framwork laravle's eloquent. it will be friendly for php developer and python or ruby developer  
目前提供5大数据库驱动, mysql,sqlite,postgres,oracle,mssql, 同时可以自由更换驱动

## document

- [English Document](docs/en/README.md)
- [中文文档](docs/en/README.md)

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
