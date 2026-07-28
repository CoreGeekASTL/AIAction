# 插件管理

> 功能域：插件管理　接口数：5　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

浏览器插件包上传/删除/查询/加载，innerServer（HTTP）暴露，PluginController 承载。底层走 service.PluginService → dao + common/storage/oss。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| UploadPluginPackage | 上传插件包 | controllers/plugin_controller.go | POST /plugin/v1/upload |
| DeletePluginPackage | 删除插件包 | controllers/plugin_controller.go | POST /plugin/v1/delete |
| GetPluginPackages | 查询全部插件包 | controllers/plugin_controller.go | POST /plugin/v1/getAll |
| LoadPlugin | 加载插件 | controllers/plugin_controller.go | POST /plugin/v1/load |
| GetCurrentPlugins | 查询当前已加载插件 | controllers/plugin_controller.go | POST /plugin/v1/current |

## 3. 数据结构说明

- **UploadPluginPackage**
  - 请求：multipart/form-data（parseUploadPluginPackageParam），含插件包文件元数据
  - 响应：retcode 标准结构；panic 自动 recover
- **DeletePluginPackage / GetPluginPackages / LoadPlugin / GetCurrentPlugins**
  - 请求：plugin key/name 等参数（plugin_service 层处理）
  - 响应：retcode 标准结构；GetPluginPackages/GetCurrentPlugins 返回插件列表

## 4. 风险与注意点

- **上传无大小限制**：controllers/plugin_controller.go:47（插件包上传未显式限制大小，可能内存/磁盘耗尽）
- **panic 自动 recover**：controllers/plugin_controller.go:49-58（捕获异常仅记日志，不返回错误，调用方感知不到失败）
