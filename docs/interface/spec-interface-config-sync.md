# 配置同步与查询

> 功能域：配置同步与查询　接口数：2　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

浏览器配置同步（拉取沐恩配置）与配置列表查询，innerServer（HTTP）暴露，ManagementController 承载。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| SyncBrowserConfig | 同步浏览器配置 | controllers/management_controller.go | POST /rpc-api/center/config/syncBrowserConfig |
| ListConfig | 查询配置列表 | controllers/management_controller.go | GET /config/v1 |

## 3. 数据结构说明

- **SyncBrowserConfig**
  - 请求：无参数（直接调沐恩 GET）
  - 响应 `resp.DataResponse{Data: *BrowserConfig}`，`BrowserConfig` 含 `RouteAPPConfigList`、`ChromeConfigList`、`URLConfigs`
  - 落库为 `db.Config{Type: moonConfig, Content: content}`
- **ListConfig**
  - 请求：无参数
  - 响应：配置列表（dao.ConfigDao 查询结果）

## 4. 风险与注意点

- **同步无超时**：controllers/management_controller.go:149（拉取沐恩配置无显式超时，可能挂起）
