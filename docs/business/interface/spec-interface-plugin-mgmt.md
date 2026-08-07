# 插件包管理

> 功能域：plugin-mgmt　接口数：5　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

管理浏览器插件包的生命周期：上传/删除/查询插件包，以及向全部浏览器实例下发加载插件并查询当前激活插件。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| UploadPluginPackage | 上传插件包（multipart 表单） | src/controllers/plugin_controller.go | POST /plugin/v1/upload |
| DeletePluginPackage | 删除插件包 | src/controllers/plugin_controller.go | POST /plugin/v1/delete |
| GetPluginPackages | 查询全部插件包列表 | src/controllers/plugin_controller.go | POST /plugin/v1/getAll |
| LoadPlugin | 向所有浏览器实例加载指定插件 | src/controllers/plugin_controller.go | POST /plugin/v1/load |
| GetCurrentPlugins | 查询当前激活插件列表 | src/controllers/plugin_controller.go | POST /plugin/v1/current |

## 3. 数据结构说明

- **UploadPluginPackage**
  - 请求 `req.UploadPluginPackageReq`（src/models/req/plugin_entity.go）：multipart 表单——`filename`（为空时取文件 header 文件名）、`file`（文件内容）、Size 取自 header；Validate 要求文件名/文件非空且 Size ≤ `constants.MaxFileSize`（src/models/req/plugin_entity.go:20）；controller 带 panic recover 兜底（src/controllers/plugin_controller.go:49-58）
  - 响应 `resp.BaseResponse`
- **DeletePluginPackage / LoadPlugin**
  - 请求 `req.PluginPackageReq`（src/models/req/plugin_entity.go）：Name、Type、Version（三者必填，组合成 key `type:name:version`）
  - 响应 `resp.BaseResponse`；LoadPlugin 会向 `browserService.GetAllServiceInstances()` 返回的全部实例下发（src/controllers/plugin_controller.go:138-140）
- **GetPluginPackages**
  - 请求：无 body
  - 响应 `[]db.PluginPackage`（src/models/db/plugin_info.go）：Name、Version、PackageName、Type、PackageBucket、Status（Completed/Failed/Doing/NotStart）、IfActive、Progress、CreatedAt
- **GetCurrentPlugins**
  - 请求：无 body
  - 响应 `[]db.PluginPackage`，结构同上（service/plugin_service.go:101）
