## gorose orm的由来
### 优势
- 链式调用
- 简单易上手,1分钟上手不是梦
- 动态语言般的使用畅快感
- 友好的 api
- 不需要定义一大堆的 struct 来做表字段类型声明
- 信手拈来, 开箱即用

### 由来
- 用惯了 `laravel eloquent` 的数据库操作, 用其他的语言做数据库操作, 也总是会怀念那种流畅操作, 一气呵成的快感.  
- 突然想到了大道至简, 即简单容易上手, 又强大好用, 方是 `orm` 的归宿.  
- 当接触到 `go` 的时候, 就被他深深的吸引了, 就想着用他做一做微服务, 但是, 用惯了世界上最好的语言 `php` , 总感觉少了点感觉, 找到了 `gin` 框架, 感觉他的路由挺爽的, 但是, 总感觉少了点什么, 找来找去, 终于发现, 是少了一个称手的`orm`, 目前市面上的 `orm` , 比较大众化的有 `gorm`, `xorm`等等. 这些都是比较强大的`orm`, 但是用上去, 总感觉少了点 `feel`, 于是就想着从 php 那里取点经, 就着手自己早了个轮子 .   
- 于是乎, 就想着自己造个轮子吧, 既可以学习交流, 也可以为开源做点小贡献, `gorose` 这个轮子就这么的诞生了  

## github
[https://github.com/gohouse/gorose](https://github.com/gohouse/gorose)

## 随时在线交流心得
[点击加入qq群: 470809220](https://jq.qq.com/?_wv=1027&k=5JJOG9E)  

## 先睹为快
先做个简单的操作, 过过眼瘾
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