#### gin的一个demo

> 包管理工具：go mod

> 项目目录介绍

```
.
├── app // 项目
│   ├── http
│   │   ├── controllers  // 控制器
│   │   │   └── api
│   │   │       └── v1
│   │   │           ├── login
│   │   │           └── user
│   │   └── middleware  // 中间件目录
│   │       └── logger
│   ├── models //model目录
│   │   ├── entity
│   │   └── user
│   └── services //业务层目录
│       └── user
├── conf // 配置文件
├── databases
├── pkg // 
│   ├── file
│   ├── logging
│   ├── redis
│   ├── setting
│   └── utils
│       └── result
├── routers // 路由
├── runtime
│   └── logs
└── vendor // vendor目录

```
##### 主要对gorm不定条件查询数据时的一个封装【灵感来于laravel项目中对db的封装[laravel版本Eloquent ORM构建](https://www.jianshu.com/p/9896549b7f3d)】
> 封装方法在 app/models/entity/Gorm.go文件里

#### 条件说明
> ["字段名","操作符","查询值","与前一个条件的关系[默认and]"] \
>
> 1.如果是等于，可以省略"操作符" : \
> []interface{}{"username", "chen"} 或 []interface{}{"username","=" , "chen"} 
> 
> 2.大于:
> []interface{}{"createtime", ">", "2019-1-1"}
>
> 3.如果为or，那就得一写全：
[]interface{}{"username", "=", "chen", "or"}
>
> 4.其它的where兼容gorm的where方法

#### 测试

启动项目
```
go run main.go
```
访问测试地址：
[http://127.0.0.1:8100/api/v1/user/test](http://127.0.0.1:8100/api/v1/user/test)

带分页的地址：[http://127.0.0.1:8100/api/v1/user/list](http://127.0.0.1:8100/api/v1/user/list)

> 1、and条件测试
```
where := []interface{}{
	[]interface{}{"id", "=", 1},
	[]interface{}{"username", "chen"},
}
db, err = entity.BuildWhere(db, where)
db.Find(&users)
// SELECT * FROM `users`  WHERE (id = 1)and(username = 'chen')
```

> 2、结构体条件测试
```
where := user.User{ID: 1, UserName: "chen"}
db, err = entity.BuildWhere(db, where)
db.Find(&users)
// SELECT * FROM `users`  WHERE (id = 1) and (username = 'chen')
```
> 3、in,or条件测试
```
where := []interface{}{
	[]interface{}{"id", "in", []int{1, 2}},
	[]interface{}{"username", "=", "chen", "or"},
}
db, err = entity.BuildWhere(db, where)
db.Find(&users)
// SELECT * FROM `users`  WHERE (id in ('1','2')) OR (username = 'chen')
```

> 3.1、not in,or条件测试
```
where := []interface{}{
	[]interface{}{"id", "not in", []int{1}},
	[]interface{}{"username", "=", "chen", "or"},
}
db, err = entity.BuildWhere(db, where)
db.Find(&users)
// SELECT * FROM `users`  WHERE (id not in ('1')) OR (username = 'chen')
```

> 4、map条件测试
```
where := map[string]interface{}{"id": 1, "username": "chen"}
db, err = entity.BuildWhere(db, where)
db.Find(&users)
// SELECT * FROM `users`  WHERE (`users`.`id` = '1') AND (`users`.`username` = 'chen')
```

> 5、and,or混合条件测试
```
where := []interface{}{
	[]interface{}{"id", "in", []int{1, 2}},
	[]interface{}{"username = ? or nickname = ?", "chen", "yond"},
}
db, err = entity.BuildWhere(db, where)
db.Find(&users)
// SELECTSELECT * FROM `users`  WHERE (id in ('1','2')) AND (username = 'chen' or nickname = 'yond')

//注：不要使用下方方法
/*
where := []interface{}{
	[]interface{}{"id", "in", []int{1, 2}},
	[]interface{}{
		[]interface{}{"username", "=", "chen"},
		[]interface{}{"username", "=", "yond", "or"},
	},
}
// 返回sql: SELECT * FROM `users`  WHERE (id in ('1','2')) AND (username = 'chen') OR (username = 'yond')
// 与设想不一样
// 经过测试，gorm底层暂时不支持db.Where(func(db *gorm.DB) *gorm.DB {})闭包方法
*/

```







> demo框架参考 
[eddycjy/go-gin-example](https://github.com/eddycjy/go-gin-example)