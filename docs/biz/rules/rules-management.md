# 浏览器配置同步 业务规则

> 生成时间：2026-08-11
> 覆盖入口：POST /rpc-api/center/config/syncBrowserConfig → controllers.ManagementController.SyncBrowserConfig；GET /config/v1 → controllers.ManagementController.ListConfig

## 概述

浏览器配置功能域维护从云端 Moon 服务同步的浏览器配置（路由应用配置/Chrome 配置/URL 配置），本地落库缓存并支持被动过期刷新。规则提取自同步与查询两条入口链路，覆盖条件分支、阈值常量与告警上报三类规则点。

## 规则表

| 规则名 | 条件 | 动作 | 依据 | 来源入口 |
| --- | --- | --- | --- | --- |
| update-config-if-need | 查询 Type="moon" 配置失败、UpdatedAt 解析失败，或距 UpdatedAt 超过 interval（24h） | 触发 syncBrowserConfig 重新同步配置 | src/controllers/management_controller.go | GET /config/v1 |
| get-moon-config-url | 配置中心缓存存在 moon::configEndpoint / moon::httpsConfigEndpoint / moon::enableHttps 且非空 | 配置中心值覆盖 beego.AppConfig 文件默认值 | src/controllers/management_controller.go | POST /rpc-api/center/config/syncBrowserConfig；GET /config/v1 |
| enable-https | moon::enableHttps=="true" 且 httpsConfigUrl 非空 | 使用 HTTPS 地址与 https.MuenInstance 客户端；否则用 HTTP 地址与默认客户端 | src/controllers/management_controller.go | POST /rpc-api/center/config/syncBrowserConfig；GET /config/v1 |
| sync-browser-config-retry | 同步 Moon 配置 HTTP 请求 | 重试 defaultRetryCount(2) 次；响应非成功码或出错则同步失败 | src/controllers/management_controller.go | POST /rpc-api/center/config/syncBrowserConfig；GET /config/v1 |
| insert-or-update-config | Type="moon" 配置不存在（orm.ErrNoRows） | 插入新配置（CreatedAt/UpdatedAt 置当前时间）；已存在则覆盖 Content 并刷新 UpdatedAt | src/controllers/management_controller.go | POST /rpc-api/center/config/syncBrowserConfig；GET /config/v1 |
| alarm-id-300010 | SyncBrowserConfig 同步失败 | 上报告警 AlarmId300010；同步成功则恢复（ClearAlarm）该告警 | src/controllers/management_controller.go | POST /rpc-api/center/config/syncBrowserConfig |
| list-config-no-rows | 刷新后仍查不到 Type="moon" 配置（orm.ErrNoRows） | 返回 HTTP 404（NotFound）（随后仍走到 InternalServiceError 分支，代码现状如此） | src/controllers/management_controller.go | GET /config/v1 |

## 补充说明

ListConfig 是被动的惰性刷新入口：读配置时顺带判定是否超过 24h 未更新，过期才触发同步，属读路径上的短路刷新规则；SyncBrowserConfig 是主动同步入口并承担告警上报/恢复配对。配置来源优先级为配置中心缓存 > 配置文件默认值。代码内 TODO 标注：分布式多实例下同步操作需加分布式锁，现状未实现，待确认。
