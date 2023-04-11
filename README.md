<h1 align="center">Go Cinch Auth</h1>

<div align="center">
Cinch Auth是一套权限验证服务.
<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/go-cinch/auth" alt="Go version"/>
<img src="https://img.shields.io/badge/Kratos-v2.5.3-brightgreen" alt="Kratos version"/>
<img src="https://img.shields.io/badge/MySQL-8.0-brightgreen" alt="MySQL version"/>
<img src="https://img.shields.io/badge/go--redis-v8.11.5-brightgreen" alt="Go redis version"/>
<img src="https://img.shields.io/badge/Gorm-1.24.0-brightgreen" alt="Gorm version"/>
<img src="https://img.shields.io/badge/Wire-0.5.0-brightgreen" alt="Wire version"/>
<img src="https://img.shields.io/github/license/go-cinch/auth" alt="License"/>
</p>
</div>

# 起源

你的单体服务架构是否遇到一些问题, 不能满足业务需求? 那么微服务会是好的解决方案.

Cinch是一套轻量级微服务脚手架, 基于[Kratos], 节省基础服务搭建时间, 快速投入业务开发.

我们参考了Go的许多微服务架构, 结合实际需求, 最终选择简洁的[Kratos]作为基石(B站架构), 从架构的设计思路以及代码的书写格式和我们非常匹配.

> cinch意为简单的事, 小菜. 希望把复杂的事变得简单, 提升开发效率.

若你想深入学习微服务每个组件, 建议直接看Kratos官方文档. 本项目整合一些业务常用组件, 开箱即用, 并不会把每个组件都介绍那么详细.


# 特性


- `Proto` - proto协议同时开启gRPC & HTTP支持, 只需开发一次接口, 不用写两套
- `Jwt` - 认证, 用户登入登出一键搞定
- `Action` - 权限, 基于行为的权限校验
- `Redis` - 缓存, 内置防缓存穿透/缓存击穿/缓存雪崩示例
- `Gorm` - 数据库ORM管理框架, 可自行扩展多种数据库类型, 目前使用MySQL, 其他自行扩展
- `SqlMigrate` - 数据库迁移工具, 每次更新平滑迁移
- `Asynq` - 分布式定时任务(异步任务)
- `Log` - 日志, 在Kratos基础上增加一层包装, 无需每个方法传入
- `Embed` - go 1.16文件嵌入属性, 轻松将静态文件打包到编译后的二进制应用中
- `Opentelemetry` - 链路追踪, 跨服务调用可快速追踪完整链路
- `Idempotent` - 接口幂等性(解决重复点击或提交)
- `Pprof` - 内置性能分析开关, 对并发/性能测试友好
- `Wire` - 依赖注入, 编译时完成依赖注入
- `Swagger` - Api文档一键生成, 无需在代码里写注解
- `I18n` - 国际化支持, 简单切换多语言

# Auth服务


[Auth](https://github.com/go-cinch/auth)是基于[layout](https://github.com/go-cinch/layout)生成的一个通用权限验证微服务, 节省鉴权服务搭建时间, 快速投入业务开发.


# 行为管理


## 为什么不用casbin


Golang主流的权限管理基本上用的是[Casbin](https://casbin.org), 曾经的项目中也用到过, 如[gin-web](https://github.com/piupuer/gin-web), 设计模式很好,
但实际使用过程中遇到一些问题


以RBAC为例以及结合常见的web管理系统, 你大概需要这么设计

- api表 - 接口表, 存储所有的api, 如`GET /user`, `POST /login`(method+path)
- casbin表 - casbin关系表, 存储对象和api之间的关系, 如用户admin有`GET /user`的访问权限, 数据记录可能是`v0=admin, v1=GET, v2=/user`
- menu表 - 菜单表(可能包含关联表, 这里不细说), 存储所有的菜单
    - 和api关联, 如`用户管理`菜单下有`GET /user`, `POST /user`, `PATCH /user`, `DELETE /user`
    - 和role关联, 如角色是超级管理员可以展示所有页面, 访客只能查看首页
- 若还需要增加单个按钮, 那需要btn表, 和menu类似
- 若还需要单个用户权限不一样, 那user表与menu/btn要建立关联
- 若还需要用户组权限, 那用户组与上面的表要建立关联

那么问题来了
1. 新增或删除接口, 你怎么维护上述各个表之间的关联关系?
2. app没有菜单, 怎么来设计

> 或许这并不是casbin的问题, 这是业务层的问题. 你有更好的方案欢迎讨论


## 表结构


### Action


- `name` - 名称, 可以重复, 定义行为名, 描述当前action用途
- `code` - 唯一标识, 主要用于其他表关联
- `word` - 英文关键字, 唯一, 方便前端展示
- `resource` - 资源列表(\n隔开), 可以是http接口地址, 也可以是grpc接口地址
- `menu` - 菜单列表(\n隔开), 和前端路由有关
- `btn` - 按钮列表(\n隔开), 和前端有关

> Tips: 建议resource规则使用kratos自动生成的[Operation](https://github.com/go-cinch/auth/blob/dev/api/auth/auth_http.pb.go#L23) `/auth.v1.Auth/CreateAction`(即`服务名`+`方法名`)  
menu/btn是给前端看的, 前端来决定是否需要显示, 后端不单独存储menu表/btn表


### User/Role/UserGroup


- `action` - 行为code列表(逗号隔开), 包含用户/角色/用户组具有的所有行为


## 优缺点


优点
1. 权限变更高效
   - 程序员只需关心action表, 修改action表对应内容即可
   - 使用者勾选或添加action, action通过name显示, 简单易懂
2. action不区分pc/app
3. 减少关联表

缺点
1. 增加冗余性(这个和减少关联表是相悖的)
2. 菜单等更新由前端管理, 若需更新, 必须重新发布前端代码


## 验证权限


各个微服务通过权限中间件调用`/auth.v1.Auth/Permission`接口进行权限校验, 步骤如下

1. 从header中获取jwt token或grpc接口传入的user code, 最终得到user
2. 校验当前用户是否有参数resource的权限, 有直接返回
3. 校验当前用户所在角色是否有参数resource的权限, 有直接返回
4. 校验当前用户所在用户组是否有参数resource的权限

参见[Permission](https://github.com/go-cinch/auth/blob/dev/internal/service/auth.go#L121)


# 常用接口


- `/auth.v1.Auth/Register` - 用户注册
- `/auth.v1.Auth/Pwd` - 修改密码
- `/auth.v1.Auth/Status` - 获取用户状态, 是否锁定/需要输入验证码等
- `/auth.v1.Auth/Captcha` - 获取验证码base64图片, 默认4位数字
- `/auth.v1.Auth/Login` - 登入, 获取jwt token
- `/auth.v1.Auth/Logout` - 登出
- `/auth.v1.Auth/Refresh` - 刷新jwt token
- `/auth.v1.Auth/Idempotent` - 获取幂等性token(登陆后)
- `/auth.v1.Auth/Info` - 获取用户信息(登陆后)

其他接口是Action/User/Role/UserGroup的CRUD, 完整接口参见[auth.proto](https://github.com/go-cinch/auth-proto/blob/master/auth.proto)


# 调用Auth


## 调用链路


[layout](https://github.com/go-cinch/layout)内置auth中间件, 以game服务为例, 调用链路如下

1. `req` - 接收请求, 前端发起请求到game服务
2. `permissionWhitelist` - 白名单过滤, 判断该服务是否需要鉴权
3. `authClient.Permission` - 调用Auth鉴权, ctx包含用户jwt, 传入resource, Auth服务校验当前用户是否有权限访问resource
4. `handler` - 业务逻辑

```go
func Permission(authClient auth.AuthClient) middleware.Middleware {
	...
	return func(ctx context.Context, req interface{}) (rp interface{}, err error) {
		...
        res, err := authClient.Permission(ctx,
            &auth.PermissionRequest{
                Resource: resource,
            },
            grpc.Header(&reply),
        )
		...
        if !res.Pass {
            err = reason.ErrorForbidden(i18n.FromContext(ctx).T(biz.NoPermission))
            return
        }
        return handler(ctx, req)
	}
}
```

> 完整代码参见[internal/server/middleware/permission](https://github.com/go-cinch/auth/blob/dev/internal/server/middleware/permission.go#L21)


## 权限白名单


暂时关闭鉴权

为了测试方便, 暂时关闭鉴权
```bash
vim auth/internal/server/middleware/whitelist.go
```

增加以下内容
```go
	whitelist[auth.OperationAuthIdempotent] = struct{}{}
```

最终代码

```go
func permissionWhitelist() selector.MatchFunc {
	whitelist := make(map[string]struct{})
	whitelist["/grpc.health.v1.Health/Check"] = struct{}{}
	whitelist["/grpc.health.v1.Health/Watch"] = struct{}{}
	whitelist[auth.OperationAuthIdempotent] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whitelist[operation]; ok {
			return false
		}
		return true
	}
}
```

> 完整代码参见[internal/server/middleware/whitelist](https://github.com/go-cinch/auth/blob/dev/internal/server/middleware/whitelist.go#L9)

> 除了健康检查, 其他接口默认都需要鉴权


# 演示


正在规划腾讯UI组件库[tdesign](https://tdesign.tencent.com)作为演示, Vue3+, 敬请期待~


[Kratos]: (https://go-kratos.dev/docs/)
