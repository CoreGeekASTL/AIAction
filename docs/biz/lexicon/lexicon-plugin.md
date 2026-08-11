# 插件管理 领域词典

> 子文档 of [lexicon.md](lexicon.md)；词汇口径、来源说明与全仓待确认清单见主文档。子域锚点 `plugin`，功能域口径与 docs/biz/interface/ 一致。

## 实体与业务概念

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| PluginPackage | 插件包表（t_plugin_package）：key(type:name:version)/name/version/packageName/plugin_type/bucket/active_status/if_active/progress | - | `db.PluginPackage`，`src/models/db/plugin_info.go` |
| 插件包 Key | 插件包唯一标识，格式 `{type}:{name}:{version}` | - | `db.PluginPackage.GetField`，`src/models/db/plugin_info.go`；`req.PluginPackageReq.GetKey`，`src/models/req/plugin_entity.go` |
| PluginActive | 插件激活状态缓存结构（name/version/type/status/progress），Redis 序列化存储 | - | `db.PluginActive`，`src/models/db/plugin_info.go` |
| PluginPackageReq | 插件包操作请求（name/type/version，均必填） | - | `req.PluginPackageReq`，`src/models/req/plugin_entity.go` |
| UploadPluginPackageReq | 插件包上传请求（filename/file/size，大小受 MaxFileSize 限制） | - | `req.UploadPluginPackageReq`，`src/models/req/plugin_entity.go` |
| ExtensionLoadRequest | 浏览器网关扩展加载请求（bucket_name/extension_file_path/name/version/type） | - | `browsergateway.ExtensionLoadRequest`，`src/models/browsergateway/req.go` |
| PluginInfo | 插件信息响应项（name/version/type/status/progress） | - | `resp.PluginInfo`，`src/models/resp/plugin_entity.go` |
| PluginPackageResponse / PluginInfoResponse | 插件包列表 / 插件信息列表响应（BaseResponse + data） | - | `resp.PluginPackageResponse`、`resp.PluginInfoResponse`，`src/models/resp/plugin_entity.go` |

## 常量与状态枚举

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| ActiveStatus | 插件激活状态枚举（string） | - | `db.ActiveStatus`，`src/models/db/plugin_info.go` |
| Complete | 激活完成 "Completed" | - | `db.Complete`，`src/models/db/plugin_info.go` |
| Failed | 激活失败 "Failed" | - | `db.Failed`，`src/models/db/plugin_info.go` |
| Doing | 激活中 "Doing" | - | `db.Doing`，`src/models/db/plugin_info.go` |
| NotStart | 未开始 "NotStart" | - | `db.NotStart`，`src/models/db/plugin_info.go` |
| PluginPackageBucket | 插件包固定存储桶 "extension" | - | `constants.PluginPackageBucket`，`src/common/constants/base.go` |
| ChromeExtendType | Chrome 扩展类型标识 "ChromeExtend" | - | `constants.ChromeExtendType`，`src/common/constants/base.go` |
| ExtensionDescribeFile | 扩展描述文件名 "package.json" | - | `constants.ExtensionDescribeFile`，`src/common/constants/base.go` |
