# 沐恩云端服务（Muen）出站调用

沐恩（Muen）为云端服务，地址通过配置中心/本地配置 `moon::titokEndpoint`（鉴权）、`moon::configEndpoint`（配置同步）获取；`moon::enableHttps=true` 时切换到 `moon::httpsTitokEndpoint` / `moon::httpsConfigEndpoint` 并启用双向证书客户端（`https.MuenInstance()`）。

| 接口名 | 协议 | 调用位置 | 业务场景 |
|--------|------|----------|----------|
| POST /app-api/devicetcp/app/login/v1/deviceLoginAuth | HTTP/HTTPS | service/remote_service.go | 终端登录鉴权 |
| GET {moon::configEndpoint} | HTTP/HTTPS | controllers/management_controller.go | 浏览器配置同步 |

## POST /app-api/devicetcp/app/login/v1/deviceLoginAuth

- 协议：HTTP POST（可选 HTTPS，走 `moon::enableHttps` 开关），请求体为 `req.LoginAuthRequest` JSON
- 调用位置：service/remote_service.go（MuenDeviceLogin 函数）；上游调用方 controllers/exlogin_controller.go、controllers/login_controller.go
- 业务场景：终端（云浏览器客户端）发起登录请求时，GIDS 将终端信息（IMEI/IMSI/机型等）转发至沐恩云端进行设备登录鉴权
- 接口功能：请求设备登录鉴权，返回 `resp.DeviceLoginAuthResponse`（含 Token、分配的接入地址等 LoginInfo），鉴权失败返回 nil 由上层按拒绝流程处理

## GET {moon::configEndpoint}

- 协议：HTTP GET（可选 HTTPS），带重试（defaultRetryCount）
- 调用位置：controllers/management_controller.go（syncBrowserConfig 函数）
- 业务场景：管理面接口 SyncBrowserConfig 触发，从云端拉取最新浏览器配置（BrowserConfig）
- 接口功能：请求云端浏览器配置，响应体解析为 `resp.DataResponse{Data: BrowserConfig}` 后落库到 `t_config` 表；失败时上报告警 300010，成功时清除该告警
