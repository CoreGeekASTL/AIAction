# browser-gateway（BrowserGW）出站调用

browser-gateway 实例通过 CSE Watch 机制发现（`common/cse/cse.go` Watch 服务名 `browser-gateway`），实例内网地址取 `ServiceInstance.BrowserInnerEndpoint`，调用均走 HTTP 直连实例。

| 接口名 | 协议 | 调用位置 | 业务场景 |
|--------|------|----------|----------|
| POST /browsergw/browser/preOpen | HTTP | service/browser_service.go | 登录后预开浏览器 |
| POST /browsergw/extension/load | HTTP | service/plugin_service.go | 插件包加载下发 |
| DELETE /browsergw/browser/userdata/delete | HTTP | service/cache_service.go | 删除用户页面缓存 |

## POST /browsergw/browser/preOpen

- 协议：HTTP POST `http://{BrowserInnerEndpoint}/browsergw/browser/preOpen`，请求体为 `browsergateway.InitBrowserRequest` JSON
- 调用位置：service/browser_service.go（instancePreOpenBrowser 函数，由 PreOpenBrowser 对每个就绪实例并发发起）
- 业务场景：终端登录鉴权通过后，提前通知所有就绪的 BrowserGW 实例预开浏览器，缩短用户首开时延
- 接口功能：请求 BrowserGW 按终端参数（机型/分辨率/IMEI/IMSI/语言等）预初始化浏览器实例；失败仅记录日志，不影响登录主流程

## POST /browsergw/extension/load

- 协议：HTTP POST `http://{BrowserInnerEndpoint}/browsergw/extension/load`，带重试（defaultRetryCount），请求体为 `browsergateway.ExtensionLoadRequest` JSON
- 调用位置：service/plugin_service.go（loadPluginToBrowserGW 函数）
- 业务场景：插件包激活时，将插件包（BucketName/ExtensionFilePath/Name/Version/Type）下发到所有 BrowserGW 实例加载
- 接口功能：请求 BrowserGW 加载指定插件扩展，返回 `browsergateway.ExtensionLoadResponse`（Code 200 为成功）；按成功实例数更新加载进度，全部成功置 Complete，否则置 Failed

## DELETE /browsergw/browser/userdata/delete

- 协议：HTTP DELETE `http://{BrowserInnerEndpoint}/browsergw/browser/userdata/delete`，请求体为 `{imei, imsi}` JSON，超时 5s
- 调用位置：service/cache_service.go（callBrowserGW 函数，由 DeleteCacheImpl 对所有就绪实例逐个调用）；上游入口 controllers/cache_controller.go（DeleteCache）
- 业务场景：管理面下发"删除页面缓存"请求时，通知各 BrowserGW 清除指定用户在对象存储侧的页面缓存数据
- 接口功能：请求删除指定 IMEI/IMSI 用户数据缓存，返回 200 视为成功；单实例失败仅记录日志，不中断其他实例处理
