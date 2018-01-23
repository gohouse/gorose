gorose, 最风骚的go orm, 拥有链式操作, 开箱即用, 一分钟上手等八大风骚, 让golang操作数据库成为一种享受, 妈妈再也看不到我处理数据的痛苦了, 下面就让我一一讲解gorose的风情

### 风骚一 : 开箱即用, 一分钟上手
```go
gorose.Open("xxxxxx这里是配置文件中的数据库配置")

var db gorose.Database

db.Query("select * from user")  // 原生sql执行, 返回格式化后的结果
```

### 风骚二 : 链式操作, 尽显妩媚之姿
```go
db.Table("user").First()
```
get sql : `select * from user where id=1`  

### 风骚三 : 直接查询想要的字段, 无需预先声明字段类型
```go
db.Table("user").Fields("id as uid,name").Where("id", ">", 1).Get()
```

### 风骚四 : JSON返回自由切换

- 指定json
```go
result := db.Table("user").Get()
jsonStr := db.JsonEncode(result)
```

或者
```go
result := db.Table("user").Get()
jsonStr := utils.JsonEncode(result) // 这里的utils是github.com/gohouse/utils工具包, 可以在任何地方调用
```

### 风骚五 : 一键处理事务, 全自动 启动/回滚/提交 事务, 我只需专注于代码本身
```go
db.Transaction(func(){
	db.Table("user").Data(map[string]interface{}{"name":"fizz"}).Insert()
	db.Table("user").Data(map[string]interface{}{"name":"fizz2"}).Where("id",1).Update()
})
```

### 风骚六 : 一键切换数据库连接, 自由畅想在任何数据库之间
```go
db.Connect("mysql2").Table("goods").First()
```

### 风骚七 : 大量数据自动分块处理, 我还是只需要专注于代码本身  

user表中的所有数据, 我每次取出100条, 然后处理完, 自动取下一个100条, 继续处理, 如此反复, 直到处理完指定条件的数据
```go
db.Table("user").Where("id","<",10000).Chunk(100, func(data []map[string]interface{}){
	for _, item := range data {
		fmt.Println(item["name"])
	}
})
```

### 风骚八 : 无感知读写分离, 依然只专注于代码本身, 只需要在配置中设置读库和写库即可

## 更多风骚之处
- 请看 [https://github.com/gohouse/gorose](https://github.com/gohouse/gorose)  
  
- 或者 [点击加入qq群: 470809220](https://jq.qq.com/?_wv=1027&k=5JJOG9E) 慢慢撩~~~  

---

powered by [fizzday](http://fizzday.net)(星期八)