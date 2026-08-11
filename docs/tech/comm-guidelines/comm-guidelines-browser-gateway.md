# browser-gateway 通信规范

## 接口清单

| 接口名 | 协议 | 调用位置 | 业务场景 |
|---|---|---|---|
| POST /browsergw/browser/preOpen | HTTP | src/service/browser_service.go（instancePreOpenBrowser） | 登录鉴权通过后预热浏览器实例 |
| POST /browsergw/extension/load | HTTP | src/service/plugin_service.go（loadPluginToBrowserGW） | 插件包加载到浏览器实例 |
| DELETE /browsergw/browser/userdata/delete | HTTP | src/service/cache_service.go（callBrowserGW） | 删除终端页面/用户数据缓存 |

## HTTP

### POST /browsergw/browser/preOpen

- 业务场景：终端登录鉴权流程中，鉴权通过后向所有 Ready 状态的 BrowserGW 实例并发发起浏览器预热
- 接口功能：请求体为 InitBrowserRequest（厂商/机型/屏幕宽高/IMEI/IMSI/设备类型/客户端语言），通知 BrowserGW 预创建浏览器实例
- 调用位置：src/service/browser_service.go（PreOpenBrowser → instancePreOpenBrowser，每实例一个 goroutine）
- 协议信息：
  - 协议：HTTP POST `http://{BrowserInnerEndpoint}/browsergw/browser/preOpen`，目标地址来自 CSE Watch 到的 browser-gateway 实例（src/common/cse/cse.go）
  - 封装方式：统一封装层 src/common/https/builder.go
  - 超时重试：未调用 WithRetry，未设置显式重试；超时为封装层框架默认（总超时 240s，ResponseHeader/TLS 握手 120s，src/common/https/client.go）
  - 错误码处理：响应为 nil 或非 2xx 仅记录日志，不向上传递失败（预热失败不影响登录主流程）

### POST /browsergw/extension/load

- 业务场景：插件管理流程，将插件包加载到各 BrowserGW 实例，并按实例完成数更新加载进度
- 接口功能：请求体为 ExtensionLoadRequest（BucketName/ExtensionFilePath/Name/Version/Type），返回 ExtensionLoadResponse（Code==200 记该实例加载完成）
- 调用位置：src/service/plugin_service.go（loadPlugin → loadPluginToBrowserGW）
- 协议信息：
  - 协议：HTTP POST `http://{BrowserInnerEndpoint}/browsergw/extension/load`
  - 封装方式：统一封装层 src/common/https/builder.go
  - 超时重试：重试 2 次（`defaultRetryCount`）；超时为封装层框架默认
  - 错误码处理：非 2xx 或出错返回 error 并跳过该实例；ExtensionLoadResponse.Code 非 200 记录日志；全部实例未完成时插件状态置为 Failed

### DELETE /browsergw/browser/userdata/delete

- 业务场景：删除页面缓存处理流程，向所有 BrowserGW 实例下发删除指定终端（IMEI+IMSI）对象存储缓存的请求
- 接口功能：请求体 JSON {imei, imsi}，要求响应 200 OK
- 调用位置：src/service/cache_service.go（DeleteCacheImpl → callBrowserGW）
- 协议信息：
  - 协议：HTTP DELETE `http://{BrowserInnerEndpoint}/browsergw/browser/userdata/delete`
  - 封装方式：直连裸 client（`net/http` 自建 `http.Client`，未走 src/common/https 统一封装层）
  - 超时重试：显式超时 5s（`defaultCacheTimeoutSeconds`）；未设置重试
  - 错误码处理：非 200 返回 error；调用方逐实例仅记录日志、不中断其余实例，整体仍返回 nil
