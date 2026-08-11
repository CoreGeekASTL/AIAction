# AIAction 架构要素（arch）资产索引

> 生成时间：2026-08-11（终端鉴权变更刷新）
> 生成工具：all-index（自动聚合产物，同名覆盖更新；资产变更后重跑 all-index 刷新，请勿手改）

## structure-model/（结构模型）

模块划分、分层、职责与依赖关系（UML 包图 + 依赖矩阵）

| 文件 | 说明 |
| --- | --- |
| [structure-model.md](structure-model/structure-model.md) | AIAction（GIDS）结构模型总览 |
| [structure-model-common.md](structure-model/structure-model-common.md) | common 模块结构文档 |
| [structure-model-controllers.md](structure-model/structure-model-controllers.md) | controllers 模块结构文档 |
| [structure-model-dao.md](structure-model/structure-model-dao.md) | dao 模块结构文档 |
| [structure-model-db.md](structure-model/structure-model-db.md) | db 模块结构文档 |
| [structure-model-models.md](structure-model/structure-model-models.md) | models 模块结构文档 |
| [structure-model-routers.md](structure-model/structure-model-routers.md) | routers 模块结构文档 |
| [structure-model-scheduler.md](structure-model/structure-model-scheduler.md) | scheduler 模块结构文档 |
| [structure-model-service.md](structure-model/structure-model-service.md) | service 模块结构文档 |
| [structure-model-utils.md](structure-model/structure-model-utils.md) | utils 模块结构文档 |

## interaction-model/（交互模型）

模块间主业务流程、消息走向（UML 时序图）

| 文件 | 说明 |
| --- | --- |
| [interaction-model-auth-imei.md](interaction-model/interaction-model-auth-imei.md) | 终端联合鉴权（AuthIMEI） 交互模型 |
| [interaction-model-control-traffic-stats.md](interaction-model/interaction-model-control-traffic-stats.md) | 控制流量统计批量上报（ControlTrafficStats） 交互模型 |
| [interaction-model-data-cleanup.md](interaction-model/interaction-model-data-cleanup.md) | 过期统计数据定时清理（DataCleanupScheduler） 交互模型 |
| [interaction-model-delete-cache.md](interaction-model/interaction-model-delete-cache.md) | 删除用户缓存（DeleteCache） 交互模型 |
| [interaction-model-delete-plugin-package.md](interaction-model/interaction-model-delete-plugin-package.md) | 删除插件包（DeletePluginPackage） 交互模型 |
| [interaction-model-device-login-auth.md](interaction-model/interaction-model-device-login-auth.md) | 设备登录鉴权（DeviceLoginAuth） 交互模型 |
| [interaction-model-download.md](interaction-model/interaction-model-download.md) | 通用文件下载（Download） 交互模型 |
| [interaction-model-exist.md](interaction-model/interaction-model-exist.md) | 文件存在性检查（Exist） 交互模型 |
| [interaction-model-expired-user-bind.md](interaction-model/interaction-model-expired-user-bind.md) | 用户绑定过期（ExpiredUserBind） 交互模型 |
| [interaction-model-export-imei-list.md](interaction-model/interaction-model-export-imei-list.md) | 白名单导出（ExportIMEIList） 交互模型 |
| [interaction-model-export-static-data.md](interaction-model/interaction-model-export-static-data.md) | 统计数据导出（ExportStaticData） 交互模型 |
| [interaction-model-get-current-plugins.md](interaction-model/interaction-model-get-current-plugins.md) | 查询当前激活插件（GetCurrentPlugins） 交互模型 |
| [interaction-model-get-from-db.md](interaction-model/interaction-model-get-from-db.md) | 配置中心查询（GetFromDB） 交互模型 |
| [interaction-model-get-plugin-packages.md](interaction-model/interaction-model-get-plugin-packages.md) | 查询全部插件包（GetPluginPackages） 交互模型 |
| [interaction-model-get-user-bind.md](interaction-model/interaction-model-get-user-bind.md) | 查询用户绑定（GetUserBind） 交互模型 |
| [interaction-model-grid-login-auth.md](interaction-model/interaction-model-grid-login-auth.md) | 网格登录鉴权（GridLoginAuth） 交互模型 |
| [interaction-model-grid-login-auth-open-browser.md](interaction-model/interaction-model-grid-login-auth-open-browser.md) | 网格登录鉴权并预开浏览器（GridLoginAuthOpenBrowser） 交互模型 |
| [interaction-model-handle-delete.md](interaction-model/interaction-model-handle-delete.md) | 文件删除（HandleDelete） 交互模型 |
| [interaction-model-handle-download.md](interaction-model/interaction-model-handle-download.md) | 应用文件下载（HandleDownload） 交互模型 |
| [interaction-model-handle-upload.md](interaction-model/interaction-model-handle-upload.md) | 应用文件上传（HandleUpload） 交互模型 |
| [interaction-model-import-imei-list.md](interaction-model/interaction-model-import-imei-list.md) | 白名单导入（ImportIMEIList） 交互模型 |
| [interaction-model-insert-or-update.md](interaction-model/interaction-model-insert-or-update.md) | 配置中心写入（InsertOrUpdate） 交互模型 |
| [interaction-model-list-config.md](interaction-model/interaction-model-list-config.md) | 查询浏览器配置（ListConfig） 交互模型 |
| [interaction-model-load-plugin.md](interaction-model/interaction-model-load-plugin.md) | 加载插件（LoadPlugin） 交互模型 |
| [interaction-model-media-traffic-stats.md](interaction-model/interaction-model-media-traffic-stats.md) | 媒体流量统计批量上报（MediaTrafficStats） 交互模型 |
| [interaction-model-monitor-schedule.md](interaction-model/interaction-model-monitor-schedule.md) | CSP 话统监控上报（MonitorSchedule） 交互模型 |
| [interaction-model-refresh-config-task.md](interaction-model/interaction-model-refresh-config-task.md) | 配置中心缓存定时刷新（RefreshConfigTask） 交互模型 |
| [interaction-model-send-app-use-times-event.md](interaction-model/interaction-model-send-app-use-times-event.md) | 上报应用使用时长事件（SendAppUseTimesEvent） 交互模型 |
| [interaction-model-send-client-event.md](interaction-model/interaction-model-send-client-event.md) | 上报终端事件（SendClientEvent） 交互模型 |
| [interaction-model-session-stats.md](interaction-model/interaction-model-session-stats.md) | 会话统计上报（SessionStats） 交互模型 |
| [interaction-model-sync-browser-config.md](interaction-model/interaction-model-sync-browser-config.md) | 同步浏览器配置（SyncBrowserConfig） 交互模型 |
| [interaction-model-update-user-bind.md](interaction-model/interaction-model-update-user-bind.md) | 更新用户绑定（UpdateUserBind） 交互模型 |
| [interaction-model-upload.md](interaction-model/interaction-model-upload.md) | 通用文件上传（Upload） 交互模型 |
| [interaction-model-upload-plugin-package.md](interaction-model/interaction-model-upload-plugin-package.md) | 上传插件包（UploadPluginPackage） 交互模型 |
