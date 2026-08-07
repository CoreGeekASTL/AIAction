# 对外接口总览

| 元信息 | 值 |
|--------|-----|
| 分支 | ready/27.0-终端鉴权 分支 (2026-08-07) |
| 更新日期 | 2026-08-07 |
| Skill | spec-interface-analyze |

## 1. 接口全景

| 功能域 | 接口数 | 子文档 |
|---|---|---|
| 终端登录鉴权与用户绑定 | 6 | [spec-interface-device-login.md](spec-interface-device-login.md) |
| 用户数据缓存删除 | 1 | [spec-interface-cache-clean.md](spec-interface-cache-clean.md) |
| 客户端事件上报 | 2 | [spec-interface-client-event.md](spec-interface-client-event.md) |
| 文件上传下载管理 | 6 | [spec-interface-file-transfer.md](spec-interface-file-transfer.md) |
| 浏览器配置查询与同步 | 2 | [spec-interface-browser-config.md](spec-interface-browser-config.md) |
| 插件包管理 | 5 | [spec-interface-plugin-mgmt.md](spec-interface-plugin-mgmt.md) |
| 流量会话统计与导出 | 4 | [spec-interface-traffic-stats.md](spec-interface-traffic-stats.md) |
| 配置中心读写 | 2 | [spec-interface-config-center.md](spec-interface-config-center.md) |
| 连通性测试 | 1（仅测试使用） | [spec-interface-test-echo.md](spec-interface-test-echo.md) |
| 证书更新订阅 | 3（异步） | [spec-interface-cert-subscribe.md](spec-interface-cert-subscribe.md) |
| 统计数据定时清理 | 1（定时任务） | [spec-interface-data-cleanup.md](spec-interface-data-cleanup.md) |
| 终端鉴权（IMEI+IMSI 白名单） | 3（仅内部 server） | [spec-interface-terminal-auth.md](spec-interface-terminal-auth.md) |

> 未归类接口：无。注：main.go 中的 `service.StartRefreshConfigTask()`、`monitorService.InitMonitorSchedule()`、`service.CleanAllActiveAlarm()` 为进程内部后台循环（配置缓存刷新/运营指标上报/告警清理），不接收外部请求、不属于对外服务能力，未计入接口统计（src/main.go:72、src/main.go:91、src/main.go:93）。

自检：扫描 36 个接口，已记录 36 个，未归类 0 个，差集已清零（2026-07-30；含 27.0 终端鉴权需求新增 3 接口）
