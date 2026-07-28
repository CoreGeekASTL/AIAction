# 用户绑定管理

> 功能域：用户绑定管理　接口数：3　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

用户绑定关系查询/过期/更新，innerServer（HTTP）暴露，LoginController 承载。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| GetUserBind | 查询用户绑定 | controllers/login_controller.go | GET /user-bind/v1/:sessionID |
| ExpiredUserBind | 标记绑定过期 | controllers/login_controller.go | PUT /user-bind/v1/:sessionId |
| UpdateUserBind | 更新绑定 | controllers/login_controller.go | POST /user-bind/v1/update |

## 3. 数据结构说明

- **GetUserBind**
  - 请求：path 参数 `:sessionID`（string）
  - 响应 `resp.LoginInfo` 或绑定体；无记录返回 404
- **ExpiredUserBind / UpdateUserBind**
  - 请求：path 参数 `:sessionId` 或绑定体
  - 响应：retcode 标准结构

## 4. 风险与注意点

- **路由参数大小写不一致**：controllers/login_controller.go:33-34（`:sessionID` 与 `:sessionId` 同 controller 内大小写不一，Beego 路由参数匹配可能受影响）
