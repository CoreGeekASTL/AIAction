# 插件管理

> 功能域：插件管理　接口数：5（仅内部）　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

浏览器插件包上传/删除/查询/加载，仅 innerServer（HTTP）暴露，PluginController 承载（routers/beego_router.go:36）。底层走 service.PluginService → dao（t_plugin_package 表）+ 对象存储。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| UploadPluginPackage | 上传插件包 | controllers/plugin_controller.go | POST /plugin/v1/upload |
| DeletePluginPackage | 删除插件包 | controllers/plugin_controller.go | POST /plugin/v1/delete |
| GetPluginPackages | 查询全部插件包 | controllers/plugin_controller.go | POST /plugin/v1/getAll |
| LoadPlugin | 加载插件到浏览器实例 | controllers/plugin_controller.go | POST /plugin/v1/load |
| GetCurrentPlugins | 查询当前已加载插件 | controllers/plugin_controller.go | POST /plugin/v1/current |

## 3. 数据结构说明

- **UploadPluginPackage**
  - 请求：multipart/form-data，表单字段 `filename` + 文件域 `file`，内部封装为 `req.UploadPluginPackageReq`（models/req/plugin_entity.go）：Filename（必填）；File；Size（必填且 ≤ MaxFileSize=300MB，common/constants/base.go:11）
  - 响应：retcode 标准结构
- **DeletePluginPackage / LoadPlugin**
  - 请求 `req.PluginPackageReq`（models/req/plugin_entity.go）：Name（必填）；Type（必填）；Version（必填），key 形如 `type:name:version`
  - 响应：retcode 标准结构；LoadPlugin 会向所有浏览器实例下发加载
- **GetPluginPackages**
  - 请求：无参数（POST）
  - 响应 `resp.PluginPackageResponse`（models/resp/plugin_entity.go）：baseResponse + data=[]db.PluginPackage{Name, Version, PackageName, Type, Status(Completed/Failed/Doing/NotStart), Progress}
- **GetCurrentPlugins**
  - 请求：无参数（POST）
  - 响应 `resp.PluginInfoResponse`（models/resp/plugin_entity.go）：baseResponse + data=[]PluginInfo{Name, Version, Type, Status, Progress}

## 4. 风险与注意点

- **panic recover 吞掉失败响应**：controllers/plugin_controller.go:49-58（UploadPluginPackage defer recover 仅记日志不写错误响应，panic 时调用方可能收不到任何响应体）
