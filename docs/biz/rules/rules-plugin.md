# 插件包管理 业务规则

> 生成时间：2026-08-11
> 覆盖入口：POST /plugin/v1/upload → controllers.PluginController.UploadPluginPackage；POST /plugin/v1/delete → controllers.PluginController.DeletePluginPackage；POST /plugin/v1/getAll → controllers.PluginController.GetPluginPackages；POST /plugin/v1/load → controllers.PluginController.LoadPlugin；POST /plugin/v1/current → controllers.PluginController.GetCurrentPlugins

## 概述

插件包管理功能域负责 Chrome 扩展插件包的上传、删除、查询与向全部 BrowserGW 实例分发加载，含插件状态机（NotStart → Doing → Complete/Failed）。规则提取自五条插件管理入口链路，覆盖参数校验、条件分支、状态迁移、阈值常量与事务回滚五类规则点。

## 规则表

| 规则名 | 条件 | 动作 | 依据 | 来源入口 |
| --- | --- | --- | --- | --- |
| upload-plugin-package-req-validate | Filename 为空、File 为 nil、Size 为 0 或 Size > constants.MaxFileSize(300MB) | 返回 "invalid param" 错误，上层返回 retcode.ClientFailed(-2)，拒绝上传 | src/models/req/plugin_entity.go；src/controllers/plugin_controller.go | POST /plugin/v1/upload |
| plugin-package-req-validate | PluginPackageReq 的 Name/Version/Type 任一为空 | 返回 "invalid param" 错误 | src/models/req/plugin_entity.go | POST /plugin/v1/delete；POST /plugin/v1/load |
| read-package-meta | 上传文件名不以 ".zip" 结尾 | 返回 "file %s is not a zip file" 错误，拒绝上传 | src/service/plugin_service.go | POST /plugin/v1/upload |
| extension-describe-file | zip 包内找不到 constants.ExtensionDescribeFile（package.json） | 返回 "file package.json not found" 错误，拒绝上传 | src/service/plugin_service.go | POST /plugin/v1/upload |
| max-file-size | package.json 文件大小 > constants.MaxFileSize | 返回错误，拒绝上传 | src/service/plugin_service.go | POST /plugin/v1/upload |
| chrome-extend-type | package.json 中 Type != constants.ChromeExtendType（"ChromeExtend"）或 Name/Version 为空 | 返回 "load package.json error, content exception" 错误，拒绝上传 | src/service/plugin_service.go | POST /plugin/v1/upload |
| plugin-key-exist | 同（Type,Name,Version）插件记录已存在 | 返回 "key %s is exist, forbidden upload" 错误，禁止重复上传 | src/service/plugin_service.go | POST /plugin/v1/upload |
| upload-plugin-package-tx | 上传落库（文件表 + 插件表）任一步失败 | 事务整体回滚，文件与插件记录同生同灭 | src/service/plugin_service.go | POST /plugin/v1/upload |
| delete-plugin-no-rows | 删除目标插件不存在（orm.ErrNoRows） | 视为已删除，直接返回成功（幂等） | src/service/plugin_service.go | POST /plugin/v1/delete |
| if-active-delete | 删除目标插件 IfActive=true（使用中） | 返回 "plugin is in used, cannot delete it" 错误，拒绝删除 | src/service/plugin_service.go | POST /plugin/v1/delete |
| delete-plugin-package-tx | 删除插件记录或对应文件记录失败 | 事务整体回滚 | src/service/plugin_service.go | POST /plugin/v1/delete |
| switch-active-plugin | LoadPlugin 激活插件 | 事务内先将同 Type（ChromeExtend）全部插件 if_active 置 false，再将目标插件置 IfActive=true、Status=Doing、Progress=0 | src/service/plugin_service.go | POST /plugin/v1/load |
| load-plugin-status | 某实例加载返回 Code==http.StatusOK | 累计完成数，Progress = 完成数*100/实例总数；Progress==100 时 Status 置 Complete | src/service/plugin_service.go | POST /plugin/v1/load |
| load-plugin-failed-status | 全部实例加载结束后 Status != Complete | Status 置 Failed | src/service/plugin_service.go | POST /plugin/v1/load |
| get-current-plugins | 查询当前生效插件 | 仅返回 Type=ChromeExtend 且 IfActive=true 的插件 | src/service/plugin_service.go | POST /plugin/v1/current |

## 补充说明

同一 Type 插件激活互斥：switchActivePlugin 在单事务内保证全局仅一个 ChromeExtend 插件处于 IfActive=true，删除在使用中的插件被前置拦截。加载为异步分发：LoadPlugin 同步切换激活态后立即返回，实例加载进度经 channel 异步落库，代码注释标注未考虑重试、服务重启中断与失败节点记录，待确认。上传与删除均以事务包裹保证文件表与插件表一致性。
