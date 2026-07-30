# external-call-browser-gateway

> 下游服务：browser-gateway（浏览器网关实例，每个云浏览器节点一个实例）。
> 实例地址来源：通过 CSE Watch `browser-gateway` 微服务（common/cse/cse.go），实例内网地址取 `browsergateway.ServiceInstance.BrowserInnerEndpoint`。
> 调用方式：`common/https` 封装的 HTTP builder（`https.NewRequest(https.Instance())`），目标 URL `http://<BrowserInnerEndpoint>/browsergw/...`；缓存删除接口直接用 `net/http` 裸客户端。

| 接口名 | 协议 | 调用位置 | 业务场景一句话 |
|--------|------|----------|----------------|
| POST /browsergw/browser/preOpen | HTTP | service/browser_service.go | 登录成功后异步预开浏览器 |
| POST /browsergw/extension/load | HTTP | service/plugin_service.go | 向全部 BrowserGW 实例加载插件包 |
| DELETE /browsergw/browser/userdata/delete | HTTP | service/cache_service.go | 按 IMEI+IMSI 删除用户数据缓存 |

## HTTP

## POST /browsergw/browser/preOpen

- 协议：HTTP POST `http://<BrowserInnerEndpoint>/browsergw/browser/preOpen`
- 调用位置：service/browser_service.go（`instancePreOpenBrowser` 函数，由 `PreOpenBrowser` 遍历全部就绪实例 goroutine 异步发起）
- 业务场景：终端登录链路——`login_controller.go` / `exlogin_controller.go` 中登录鉴权成功后调用 `browserService.PreOpenBrowser`，在全部就绪 BrowserGW 实例上为用户预开浏览器，缩短首屏打开时延
- 接口功能：请求体 `browsergateway.InitBrowserRequest`（厂商/机型/分辨率/IMSI/IMEI/设备类型/客户端语言等），通知 BrowserGW 预初始化该用户的浏览器实例；仅记日志，失败不影响登录主流程

## POST /browsergw/extension/load

- 协议：HTTP POST `http://<BrowserInnerEndpoint>/browsergw/extension/load`（带重试，defaultRetryCount=2）
- 调用位置：service/plugin_service.go（`loadPluginToBrowserGW` 函数，由 `LoadPlugin` → `loadPlugin` 遍历实例异步发起）
- 业务场景：插件管理——管理面下发插件加载任务后，GIDS 逐个通知所有 BrowserGW 实例加载指定插件包，并汇总进度（progress）回写 `t_plugin_package` 表
- 接口功能：请求体 `browsergateway.ExtensionLoadRequest`（BucketName / ExtensionFilePath / Name / Version / Type），指示 BrowserGW 从存储桶加载插件；响应 `browsergateway.ExtensionLoadResponse`，Code=200 记为成功并累计进度

## DELETE /browsergw/browser/userdata/delete

- 协议：HTTP DELETE `http://<BrowserInnerEndpoint>/browsergw/browser/userdata/delete`（5s 超时，无重试）
- 调用位置：service/cache_service.go（`callBrowserGW` 函数，由 `DeleteCacheImpl` 遍历全部就绪实例发起）
- 业务场景：缓存清理——`cache_controller.go` 的 `/app-api/devicetcp/cache/deleteCache` 接口触发，按 IMEI+IMSI 删除该用户在全部 BrowserGW 上的对象存储页面缓存
- 接口功能：请求体 JSON `{imei, imsi}`，指示 BrowserGW 删除对应用户数据缓存；要求全部实例返回 200，单个失败仅记日志
