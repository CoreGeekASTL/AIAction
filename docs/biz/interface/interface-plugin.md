# 插件管理

> 功能域：插件管理　接口数：5　所属 server：内部
> 子文档 of [README.md](README.md)

## 1. 定位

浏览器插件包的生命周期管理：上传、删除、查询、加载（激活）与当前生效插件查询。仅注册在内部 server。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| UploadPluginPackage | 上传插件包 | controllers/plugin_controller.go | POST /plugin/v1/upload |
| DeletePluginPackage | 删除插件包 | controllers/plugin_controller.go | POST /plugin/v1/delete |
| GetPluginPackages | 查询全部插件包 | controllers/plugin_controller.go | POST /plugin/v1/getAll |
| LoadPlugin | 加载（激活）指定插件 | controllers/plugin_controller.go | POST /plugin/v1/load |
| GetCurrentPlugins | 查询当前生效插件 | controllers/plugin_controller.go | POST /plugin/v1/current |

## 3. 数据结构说明

- **UploadPluginPackage**
  - 请求 `req.UploadPluginPackageReq`（models/req/plugin_entity.go）：Filename（必填）、multipart 文件 File、Size（>0 且不超过 MaxFileSize 上限）
  - 响应 `resp.BaseResponse`
- **DeletePluginPackage / LoadPlugin**
  - 请求 `req.PluginPackageReq`（models/req/plugin_entity.go）：Name、Type、Version 三字段均必填，组合键格式 `type:name:version`
  - 响应 `resp.BaseResponse`
- **GetPluginPackages**
  - 请求：无业务参数
  - 响应 `resp.PluginPackageResponse`（models/resp/plugin_entity.go）：BaseResponse + data 为 `db.PluginPackage` 列表（Name、Version、Type、PackageName）
- **GetCurrentPlugins**
  - 请求：无业务参数
  - 响应 `resp.PluginInfoResponse`（models/resp/plugin_entity.go）：BaseResponse + data 为 `PluginInfo` 列表（Name、Version、Type、Status（db.ActiveStatus 激活状态）、Progress）
