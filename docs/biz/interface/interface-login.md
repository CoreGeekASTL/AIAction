# 登录鉴权与用户绑定

> 功能域：登录鉴权与用户绑定　接口数：6　所属 server：外部 + 内部
> 子文档 of [README.md](README.md)

## 1. 定位

终端登录鉴权（网格登录/预开浏览器/设备登录）与浏览器实例用户绑定关系管理。3 个登录鉴权接口经 externalServer 与 innerServer 双暴露（`ExLoginController` 注册在外部 server，`LoginController` 注册在内部 server，路径相同）；user-bind 系列仅内部 server。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| GridLoginAuth | 网格登录鉴权 | controllers/exlogin_controller.go、controllers/login_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth |
| GridLoginAuthOpenBrowser | 登录鉴权并预开浏览器 | controllers/exlogin_controller.go、controllers/login_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser |
| DeviceLoginAuth | 设备登录鉴权 | controllers/exlogin_controller.go、controllers/login_controller.go | POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| GetUserBind | 按 sessionID 查询用户绑定 | controllers/login_controller.go | GET /user-bind/v1/:sessionID |
| ExpiredUserBind | 使指定 session 的用户绑定过期 | controllers/login_controller.go | PUT /user-bind/v1/:sessionId |
| UpdateUserBind | 更新用户绑定的浏览器实例与端点 | controllers/login_controller.go | POST /user-bind/v1/update |

## 3. 数据结构说明

- **GridLoginAuth / GridLoginAuthOpenBrowser / DeviceLoginAuth**
  - 请求 `req.LoginAuthRequest`（models/req/request_entity.go）：内嵌 `UserIdentity`（IMSI、IMEI，15 位纯数字）；Manufacturer、Model、AppType（TikTok 时走 Muen 云二次鉴权）、DeviceType 等设备信息字段
  - 响应 `resp.DeviceLoginAuthResponse`（models/resp/response_entity.go）：`BaseResponse`（code/msg）+ `Data LoginInfo`，LoginInfo 组合 `AuthInfo`（Token、ExpiresTime、TimeAxis）与 `AssignInfo`（NodeGateWayURL、HttpsShortAddr、NodeCapacity 等分配地址）；Grid 系接口返回前清空 TcpAddr/TlsTcpAddr/VideoMode/ShortAddr 等字段
  - 双实现说明：同一路径在 externalServer 由 `ExLoginController` 处理、在 innerServer 由 `LoginController` 处理，按 server 监听地址选择
- **GetUserBind**
  - 请求：路径参数 sessionID
  - 响应 `db.UserBind`（models/db/user.go）：BrowserInstance、MediaEndpoint、ControlEndpoint、MediaTlsEndpoint、ControlTlsEndpoint、InnerMediaEndpoint、InnerBrowserEndpoint、Token、Heartbeats；无记录返回 404
- **ExpiredUserBind**
  - 请求：路径参数 sessionId（注意与 GetUserBind 路径参数大小写不同）
  - 响应 `resp.BaseResponse`（code/msg）
- **UpdateUserBind**
  - 请求 `req.UpdateUserBindRequest`（models/req/request_entity.go）：SessionID（必填，空则校验失败）、BrowserInstance、MediaEndpoint、ControlEndpoint、MediaTlsEndpoint、ControlTlsEndpoint、InnerMediaEndpoint、InnerBrowserEndpoint
  - 响应 `resp.BaseResponse`
