# 设备登录鉴权 业务规则

> 生成时间：2026-08-11
> 覆盖入口：POST /app-api/devicetcp/app/login/v1/gridLoginAuth → controllers.LoginController.GridLoginAuth / controllers.ExLoginController.GridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser → controllers.LoginController.GridLoginAuthOpenBrowser / controllers.ExLoginController.GridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth → controllers.LoginController.DeviceLoginAuth / controllers.ExLoginController.DeviceLoginAuth；GET /user-bind/v1/:sessionID → controllers.LoginController.GetUserBind；PUT /user-bind/v1/:sessionId → controllers.LoginController.ExpiredUserBind；POST /user-bind/v1/update → controllers.LoginController.UpdateUserBind

## 概述

设备登录鉴权功能域解决终端（IMEI+IMSI 标识）登录云浏览器时的用户建档、实例分配（UserBind 路由）与 TikTok 场景下 Muen 云端二次登录问题。规则提取自三条登录鉴权入口链路与三条 user-bind 管理入口链路，覆盖条件分支、参数校验、阈值常量与错误码返回四类规则点。

## 规则表

| 规则名 | 条件 | 动作 | 依据 | 来源入口 |
| --- | --- | --- | --- | --- |
| login-inject-auth-reject | loginAuth 注入 AuthService.AuthIMEI 白名单联合鉴权判定不通过 | 返回 retcode.ClientFailed(-2)，Message "auth rejected"，终止登录（鉴权细节见 rules-auth.md） | src/controllers/login_controller.go；src/controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| request-body-unmarshal-to | 登录请求体 JSON 反序列化失败 | 返回 retcode.ClientFailed(-2)，拒绝登录 | src/controllers/login_controller.go；src/controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| create-or-update-user | 用户记录不存在（Key = `IMEI_IMSI`，orm.ErrNoRows） | 插入新用户（CreatedAt/UpdatedAt 置当前时间）；已存在则仅刷新 UpdatedAt | src/service/user_service.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| create-or-update-user-failed | CreateOrUpdateUser 返回错误 | 返回 retcode.InternalFailed(-1)，Message 固定 "Login failed"，登录失败 | src/controllers/login_controller.go；src/controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| route-to-instance-fallback | RouteToInstance 分配实例失败 | 鉴权仍通过，返回空 LoginInfo（降级，不阻断登录） | src/controllers/login_controller.go；src/controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| validate-user-bind | UserBind 不存在、绑定实例在 CSE 中查不到、实例 IsHealthy=false 或心跳过期 | 触发 reRouteToInstance 重新分配实例并生成新 Token（uuid.New） | src/service/browser_service.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| default-heartbeats | UserBind.Heartbeats 距今超过 defaultHeartbeats（3 分钟，defaultHeartbeatsMultiplier=3） | 判定 UserBind 过期，重新分配实例 | src/service/browser_service.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| get-all-ready-service-instances | 实例 PluginStatus != db.Complete 或 Cap <= 0 或 IsHealthy=false | 实例不参与分配 | src/service/browser_service.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| assign-instance | 无就绪实例，或排序后最空闲实例 Used >= Cap | 返回 "no idle instances available" 错误，上层降级为空 LoginInfo | src/service/browser_service.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| grid-login-auth-mask | gridLoginAuth / gridLoginAuthOpenBrowser 响应 | 清空 TcpAddr/TlsTcpAddr/VideoMode/ShortAddr/HttpsShortAddr/NodeIntranetWayURL 字段后返回 | src/controllers/login_controller.go；src/controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser |
| pre-open-browser | 入口为 gridLoginAuthOpenBrowser（preOpenBrowser=true） | 向所有就绪实例异步发起 /browsergw/browser/preOpen 预开浏览器 | src/controllers/login_controller.go；src/controllers/exlogin_controller.go；src/service/browser_service.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser |
| tik-tok-app-type | 请求 AppType == constants.TikTokAppType（"2"） | 走 MuenDeviceLogin 云端二次登录取 token，并用返回 LoginInfo 覆盖响应 Data | src/controllers/login_controller.go；src/controllers/exlogin_controller.go；src/service/remote_service.go | POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| muen-device-login-failed | TikTok 场景 MuenDeviceLogin 返回 nil 或 UpdateUserToken 失败 | 返回 retcode.InternalFailed(-1)，Message "login failed"，登录失败 | src/controllers/login_controller.go；src/controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| report-device-login-event | 登录成功后上报 Login 事件失败 | 仅记录错误日志，不影响登录结果（代码 #toDo 标注） | src/controllers/login_controller.go；src/controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| get-user-bind-no-rows | 按 sessionID 查询 UserBind 命中 orm.ErrNoRows | 返回 HTTP 404（NotFound） | src/controllers/login_controller.go | GET /user-bind/v1/:sessionID；PUT /user-bind/v1/:sessionId |
| update-user-bind | UpdateUserBindRequest 中非空字段 | 仅覆盖非空字段，空字段保留原值；Heartbeats 强制刷新为当前时间 | src/service/user_service.go | POST /user-bind/v1/update |
| desensitize | 日志输出用户 Key 且长度 > minDesensitizeLength(4) | 仅保留首尾各 2 字符，中间以 "******" 脱敏 | src/service/user_service.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |

## 补充说明

登录链路执行顺序固定：参数解析 → 白名单鉴权 → 用户建档 → 实例分配 → （TikTok 二次登录）→ 事件上报；白名单鉴权拒绝为终止性规则（返回 -2），实例分配失败与事件上报失败均为降级短路规则，不阻断登录成功返回，仅 TikTok 场景下 Muen 登录失败会终止整个登录。LoginController 与 ExLoginController 规则完全一致（内部 HTTP / 外部 HTTPS 双平面），同一套规则服务两个平面入口。UserBind 过期判定（3 分钟）与实例健康判定共同决定是否复用既有绑定，过期绑定重新分配后 Token 重新生成，旧 Token 自然失效。
