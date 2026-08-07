# external-call-muen-cloud

> 下游服务：沐恩云服务（代码与配置中命名 `moon` / `Muen`，云端管控面）。
> 地址来源：配置中心（t_config_center 表）优先，回落到 app.conf 配置项 `moon::titokEndpoint` / `moon::configEndpoint`（HTTP）或 `moon::httpsTitokEndpoint` / `moon::httpsConfigEndpoint`（HTTPS，`moon::enableHttps=true` 时启用）。
> 调用方式：HTTP 用 `https.Instance()` 客户端；HTTPS 用 `https.MuenInstance()` 客户端（TLS 证书来自 CertSDK 订阅的沐恩 CA 证书，见 external-call-csp-cert.md）。

| 接口名 | 协议 | 调用位置 | 业务场景一句话 |
|--------|------|----------|----------------|
| POST /app-api/devicetcp/app/login/v1/deviceLoginAuth | HTTP/HTTPS | service/remote_service.go | 终端登录时到云侧做设备鉴权 |
| GET {moon::configEndpoint} | HTTP/HTTPS | controllers/management_controller.go | 定期/手动同步浏览器配置 |

## HTTP

## POST /app-api/devicetcp/app/login/v1/deviceLoginAuth

- 协议：HTTP/HTTPS POST `<titokEndpoint>/app-api/devicetcp/app/login/v1/deviceLoginAuth`（带重试，defaultRetryCount=2）
- 调用位置：service/remote_service.go（`MuenDeviceLogin` 函数）；调用方：controllers/login_controller.go（登录链路）、controllers/exlogin_controller.go（外部登录链路）
- 业务场景：终端登录鉴权——设备发起登录请求后，GIDS 先将设备信息转发到沐恩云端做设备级登录鉴权，云侧返回成功后才继续本端的实例分配与 token 生成流程
- 接口功能：请求体 `req.LoginAuthRequest`（IMEI/IMSI/机型/平台/分辨率等设备信息）；响应 `resp.DeviceLoginAuthResponse`，返回云侧分配的登录信息 `resp.LoginInfo`（AuthInfo/AssignInfo）；失败或超时返回 nil，登录链路按鉴权失败处理

## GET {moon::configEndpoint}

- 协议：HTTP/HTTPS GET `<configEndpoint>`（完整 URL 由配置给出，带重试，defaultRetryCount=2）
- 调用位置：controllers/management_controller.go（`syncBrowserConfig` 函数）；触发方：`SyncBrowserConfig`（POST /rpc-api/center/config/syncBrowserConfig 手动触发）、`updateConfigIfNeed`（`ListConfig` 时配置超过 24h 未更新自动触发）
- 业务场景：浏览器配置同步——定期从沐恩云端拉取浏览器运行配置（路由 APP 配置、Chrome 配置、URL 配置），落库到 `t_config` 表供 `GET /config/v1` 查询；同步失败上报告警 300010
- 接口功能：无请求参数；响应 `resp.DataResponse{Data: BrowserConfig}`（routeAppConfigList / chromeConfigList / urlConfigList），解析后 JSON 序列化存库
