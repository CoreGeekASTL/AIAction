# 浏览器配置同步

> 功能域：浏览器配置同步　接口数：2　所属 server：内部
> 子文档 of [README.md](README.md)

## 1. 定位

从 Muen 云服务同步浏览器运行配置（路由应用配置/Chrome 参数配置/URL 配置）到本地 DB，并向内部调用方提供查询。仅注册在内部 server。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| SyncBrowserConfig | 触发从云服务同步浏览器配置入库（失败上报告警 300010） | controllers/management_controller.go | POST /rpc-api/center/config/syncBrowserConfig |
| ListConfig | 查询本地浏览器配置（超过 24h 未更新则先触发同步） | controllers/management_controller.go | GET /config/v1 |

## 3. 数据结构说明

- **SyncBrowserConfig**
  - 请求：无业务参数
  - 响应 `resp.BaseResponse`（code/msg）
  - 同步源地址取配置 `moon::configEndpoint` / `moon::httpsConfigEndpoint`（ConfigCenter 优先，其次 app 配置），HTTPS 开关为 `moon::enableHttps`
- **ListConfig**
  - 请求：无业务参数
  - 响应 `controllers.BrowserConfig`（controllers/management_controller.go）：
    - RouteAPPConfigList：`db.RouterAPPConfig` 列表（Manufacturer、Model、Type、Mode、ExtendModel 等路由匹配字段）
    - ChromeConfigList：`db.ChromeConfig` 列表（帧率/码率/采样率/分辨率等 Chrome 参数）
    - URLConfigs：`db.URLConfig` 列表（NodeIdent、AppType、URL、UserAgent、IsVideoType/IsWebType/IsShortType 等）
    - DB 中无记录返回 404
