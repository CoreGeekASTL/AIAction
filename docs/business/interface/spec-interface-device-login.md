# 终端登录鉴权与用户绑定

> 功能域：device-login　接口数：6　所属 server：外部(HTTPS/HTTP) + 内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

终端（App）登录鉴权、浏览器实例路由分配与预开浏览器；以及实例侧上报/查询会话绑定的用户-浏览器实例关系。3 条登录路径在 externalServer 与 innerServer 双暴露（两套 controller 实现），user-bind 仅内部暴露。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| GridLoginAuth | 网格登录鉴权，返回实例连接信息 | src/controllers/exlogin_controller.go；src/controllers/login_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth |
| GridLoginAuthOpenBrowser | 登录鉴权并预开浏览器 | src/controllers/exlogin_controller.go；src/controllers/login_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser |
| DeviceLoginAuth | 设备登录鉴权（TikTok 应用额外走 Muen 二次登录取 token） | src/controllers/exlogin_controller.go；src/controllers/login_controller.go | POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| GetUserBind | 按 sessionID 查询用户绑定信息 | src/controllers/login_controller.go | GET /user-bind/v1/:sessionID |
| ExpiredUserBind | 按 sessionID 使绑定关系过期 | src/controllers/login_controller.go | PUT /user-bind/v1/:sessionId |
| UpdateUserBind | 更新用户绑定（实例端点信息） | src/controllers/login_controller.go | POST /user-bind/v1/update |

## 3. 数据结构说明

- **GridLoginAuth / GridLoginAuthOpenBrowser / DeviceLoginAuth**
  - 请求 `req.LoginAuthRequest`（src/models/req/request_entity.go）：IMEI/IMSI（内嵌 UserIdentity，业务要求 15 位纯数字）；Manufacturer、Model、AppType（等于 `TikTokAppType` 时 DeviceLoginAuth 走 Muen 二次登录，src/controllers/login_controller.go:129）、ExtendModel、TotalKb、FreeKb 等，均为 string；Validate 为空实现
  - 响应 `resp.DeviceLoginAuthResponse`（src/models/resp/response_entity.go）：`BaseResponse{code,msg}` + `Data LoginInfo`；LoginInfo 内嵌 AuthInfo（Token、ExpiresTime、TimeAxis）与 AssignInfo（TcpAddr、TlsTcpAddr、NodeGateWayURL、HttpsNodeGateWayUrl、ShortAddr、NodeIntranetWayURL 等）；Grid 两个接口会清空 TcpAddr/TlsTcpAddr/ShortAddr 等字段后返回；路由分配失败时仍返回成功、Data 为空 LoginInfo（src/controllers/login_controller.go:161-164）
- **GetUserBind**
  - 请求：路径参数 `sessionID`
  - 响应 `db.UserBind`（src/models/db/user.go）：BrowserInstance、MediaEndpoint、ControlEndpoint、MediaTlsEndpoint、ControlTlsEndpoint、InnerMediaEndpoint、InnerBrowserEndpoint、Token、Heartbeats；不存在返回 HTTP 404
- **ExpiredUserBind**
  - 请求：路径参数 `sessionId`；响应 `resp.BaseResponse`；不存在返回 HTTP 404
- **UpdateUserBind**
  - 请求 `req.UpdateUserBindRequest`（src/models/req/request_entity.go）：SessionID（必填，Validate 校验非空）、BrowserInstance、MediaEndpoint、ControlEndpoint、MediaTlsEndpoint、ControlTlsEndpoint、InnerMediaEndpoint、InnerBrowserEndpoint
  - 响应 `resp.BaseResponse`

