# 配置同步与查询

> 功能域：配置同步与查询　接口数：2（仅内部）　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

浏览器配置同步（从沐恩云服务拉取配置落库）与配置查询，仅 innerServer（HTTP）暴露，ManagementController 承载（routers/beego_router.go:35）。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| SyncBrowserConfig | 同步浏览器配置（拉沐恩配置落库） | controllers/management_controller.go | POST /rpc-api/center/config/syncBrowserConfig |
| ListConfig | 查询浏览器配置 | controllers/management_controller.go | GET /config/v1 |

## 3. 数据结构说明

- **SyncBrowserConfig**
  - 请求：无参数
  - 响应：retcode 标准结构 BaseResponse{code,msg}
  - 内部逻辑：按 `moon::enableHttps`/`moon::configEndpoint`/`moon::httpsConfigEndpoint` 配置（优先读配置中心）GET 沐恩配置，`https.NewRequest().WithRetry(2)`（controllers/management_controller.go:149）；结果落库 `db.Config{Type: "moon", Content: <json>}`；失败上报告警 AlarmId300010，成功恢复告警
- **ListConfig**
  - 请求：无参数
  - 响应 `BrowserConfig`（controllers/management_controller.go:29）：RouteAPPConfigList（[]db.RouterAPPConfig）；ChromeConfigList（[]db.ChromeConfig）；URLConfigs（[]db.URLConfig），从 db.Config.Content 反序列化
  - 读前触发 `updateConfigIfNeed`：配置超 24h 未更新则先同步（controllers/management_controller.go:59-85）

## 4. 风险与注意点

- **多实例并发同步**：controllers/management_controller.go:45（代码内 TODO 注明：分布式多实例时同步操作需加分布式锁，当前无锁，多实例并发 SyncBrowserConfig/ListConfig 可能重复拉取与写库）
