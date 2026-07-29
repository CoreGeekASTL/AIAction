# 用户绑定管理

> 功能域：用户绑定管理　接口数：3（仅内部）　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

用户绑定关系查询/过期/更新，仅 innerServer（HTTP）暴露，LoginController 承载（routers/beego_router.go:34，与登录接口同 controller）。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| GetUserBind | 查询用户绑定 | controllers/login_controller.go | GET /user-bind/v1/:sessionID |
| ExpiredUserBind | 标记绑定过期 | controllers/login_controller.go | PUT /user-bind/v1/:sessionId |
| UpdateUserBind | 更新绑定 | controllers/login_controller.go | POST /user-bind/v1/update |

## 3. 数据结构说明

- **GetUserBind**
  - 请求：path 参数 `:sessionID`（string）
  - 响应 `db.UserBind`（models/db/user.go，t_user_bind 表）：BrowserInstance；MediaEndpoint；ControlEndpoint；MediaTlsEndpoint；ControlTlsEndpoint；InnerMediaEndpoint；InnerBrowserEndpoint；Token；Heartbeats。无记录返回 404
- **ExpiredUserBind**
  - 请求：path 参数 `:sessionId`（string）
  - 响应：retcode 标准结构 BaseResponse{code,msg}；无记录返回 404
- **UpdateUserBind**
  - 请求 `req.UpdateUserBindRequest`（models/req/request_entity.go）：SessionID（必填，非空校验）；BrowserInstance；MediaEndpoint；ControlEndpoint；MediaTlsEndpoint；ControlTlsEndpoint；InnerMediaEndpoint；InnerBrowserEndpoint
  - 响应：retcode 标准结构

## 4. 风险与注意点

- **路由参数大小写不一致**：controllers/login_controller.go:33-34（`:sessionID` 与 `:sessionId` 同 controller 内大小写不一，Beego 路由参数匹配可能受影响）
- **UserBind 响应含 Token 明文**：models/db/user.go:44（GetUserBind 直接把含 Token 的 db 实体序列化返回，内部 HTTP 无鉴权，需注意内网暴露面）
