# 对外接口总览

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | biz-interface-analyze |

## 1. 接口全景

本仓通过两个 HTTP server 对外暴露接口：externalServer（对外，注册于 `routers/beego_router.go` 的 `RegisterExternalRouter`）与 innerServer（内部，监听 127.0.0.1:9090，`RegisterInternalRouter`）。部分接口经双 server 同时暴露。

| 功能域 | 接口数 | 子文档 |
|---|---|---|
| 登录鉴权与用户绑定 | 6（3 个登录接口双 server 暴露） | [interface-login.md](interface-login.md) |
| 终端鉴权与白名单管理 | 3（仅内部） | [interface-auth.md](interface-auth.md) |
| 缓存管理 | 1（双 server 暴露） | [interface-cache.md](interface-cache.md) |
| 客户端事件上报 | 2（双 server 暴露） | [interface-event.md](interface-event.md) |
| 文件管理 | 6（2 个双 server 暴露，4 个仅内部） | [interface-file.md](interface-file.md) |
| 浏览器配置同步 | 2（仅内部） | [interface-management.md](interface-management.md) |
| 插件管理 | 5（仅内部） | [interface-plugin.md](interface-plugin.md) |
| 流量统计 | 4（仅内部）+ 1 个定时任务 | [interface-stats.md](interface-stats.md) |
| 配置中心 | 2（仅内部） | [interface-config-center.md](interface-config-center.md) |
| 测试联通 | 1（仅测试使用，externalServer） | [interface-test.md](interface-test.md) |

自检：扫描 32 个接口（去重后唯一路径）+ 1 个定时任务入口，已记录 32 + 1 个，未归类 0 个，差集已清零（2026-08-11）。
