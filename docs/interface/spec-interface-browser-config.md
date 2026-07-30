# 浏览器配置查询与同步

> 功能域：browser-config　接口数：2　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

向内部组件提供浏览器路由/Chrome 参数/URL 配置查询；支持手动触发从云端 Moon 配置服务同步最新配置到本地 DB。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| ListConfig | 查询浏览器配置（超过 24h 未更新时先自动同步） | src/controllers/management_controller.go | GET /config/v1 |
| SyncBrowserConfig | 手动触发从 Moon 云端同步浏览器配置 | src/controllers/management_controller.go | POST /rpc-api/center/config/syncBrowserConfig |

## 3. 数据结构说明

- **ListConfig**
  - 请求：无参数
  - 响应 `BrowserConfig`（src/controllers/management_controller.go:29-33）：RouteAPPConfigList（[]db.RouterAPPConfig：Manufacturer、Model、Type、Mode、ExtendModel、Name、Description）、ChromeConfigList（[]db.ChromeConfig：帧率/码率/采样率/分辨率等 14 字段）、URLConfigs（[]db.URLConfig：NodeIdent、APPType、URL、AppID、UserAgent 等）；配置不存在返回 HTTP 404
  - 注：本地配置（t_config 表 type=moon）超过 24h 未更新时先触发同步再返回（src/controllers/management_controller.go:79-83）
- **SyncBrowserConfig**
  - 请求：无 body
  - 响应 `resp.BaseResponse`；同步失败上报告警 AlarmId300010 并返回 HTTP 500，成功则恢复告警（src/controllers/management_controller.go:181-192）

## 4. 风险与注意点

- **多实例并发同步**：src/controllers/management_controller.go:45（代码内 TODO：分布式多实例时需加分布式锁，当前无防护）
- **配置源地址双通道**：src/controllers/management_controller.go:111-135（moon::enableHttps=true 时走 httpsConfigEndpoint + MuenInstance 客户端，否则走 configEndpoint）
