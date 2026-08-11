# moon 通信规范

## 接口清单

| 接口名 | 协议 | 调用位置 | 业务场景 |
|---|---|---|---|
| POST /app-api/devicetcp/app/login/v1/deviceLoginAuth | HTTP/HTTPS | src/service/remote_service.go（MuenDeviceLogin） | 终端登录鉴权转发云侧 |
| GET {moon::configEndpoint} | HTTP/HTTPS | src/controllers/management_controller.go（syncBrowserConfig） | 浏览器配置同步 |

## HTTP

### POST /app-api/devicetcp/app/login/v1/deviceLoginAuth

- 业务场景：终端登录鉴权流程中，GIDS 将终端登录鉴权请求转发给云侧 moon 服务进行云端鉴权
- 接口功能：请求体为 LoginAuthRequest（IMEI/IMSI/机型等终端信息），返回 DeviceLoginAuthResponse（含 LoginInfo 鉴权结果）
- 调用位置：src/service/remote_service.go（MuenDeviceLogin）
- 协议信息：
  - 协议：HTTP/HTTPS POST `{endpoint}/app-api/devicetcp/app/login/v1/deviceLoginAuth`，endpoint 来自配置 `moon::titokEndpoint` / `moon::httpsTitokEndpoint`（配置中心 ConfigCenter 可覆盖 app.conf 值），`moon::enableHttps=true` 时走 HTTPS
  - 封装方式：统一封装层 src/common/https/builder.go（`https.NewRequest(client)` builder），HTTPS 场景用 `https.MuenInstance()` 客户端
  - 超时重试：重试 2 次（`defaultRetryCount`，src/controllers/management_controller.go 与 src/service/plugin_service.go 中均为 2）；超时取封装层框架默认（src/common/https/client.go：总超时 240s，ResponseHeader/TLS 握手 120s）
  - 错误码处理：网络错误或响应非 2xx 时记录日志并返回 nil（鉴权失败按未通过处理），不做分类映射

### GET {moon::configEndpoint}

- 业务场景：浏览器配置管理流程，从云侧 moon 服务同步浏览器配置并落库（db.Config，Type=moonConfig）
- 接口功能：GET 请求配置端点，响应解析为 DataResponse{Data: BrowserConfig}，序列化后写库
- 调用位置：src/controllers/management_controller.go（syncBrowserConfig）
- 协议信息：
  - 协议：HTTP/HTTPS GET，URL 整体来自配置 `moon::configEndpoint` / `moon::httpsConfigEndpoint`，`moon::enableHttps=true` 时走 HTTPS
  - 封装方式：统一封装层 src/common/https/builder.go，HTTPS 场景用 `https.MuenInstance()`
  - 超时重试：重试 2 次（`defaultRetryCount`）；超时为封装层框架默认（同上 240s/120s）
  - 错误码处理：响应非 2xx 或出错时记录日志并返回 error；解析失败返回 error
