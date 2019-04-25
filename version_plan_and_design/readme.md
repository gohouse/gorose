# 2.0版本规划
发表建议地址: https://github.com/gohouse/gorose/issues/73

## 开发计划
按照 laravel 的数据库操作 query builder 的标准开发, 同时可酌情添加 eloquent 的部分设计思想和 api

## 产品规划
采用模块分离的方式, 每个模块以 接口 的方式相互调用, 做到模块相互独立, 后期可以自由方便的拓展.  
比如驱动, 不同驱动的sql构建

## 设计目标
每个模块要做到可以自由横向扩展, 以实现个性化需求.  
比如返回结果: map, struct, 自定义数据类型等.  

## 模块大致结构图
![gorose_2.0](https://github.com/gohouse/gorose/blob/2.0-dev/version_plan_and_design/gorose_2.0_modules.png?raw=true)
