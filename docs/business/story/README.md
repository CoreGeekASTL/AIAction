# 功能软件要素文档

| 元信息 | 值 |
|--------|-----|
| 分支 | ready/27.0-终端鉴权 分支 (2026-08-07) |
| 更新日期 | 2026-08-07 |
| Skill | spec-feature-analyze |

## 功能全景

| 功能域 | 接口数 | 核心模块 | 文档 |
|---|---|---|---|
| 终端登录鉴权与实例分配 | 6（3 登录内外双注册 + 3 user-bind） | controllers/exlogin, controllers/login, service | [feature-device-login.md](feature-device-login.md) |
| 用户数据缓存删除 | 1（内外双注册） | controllers/cache, service | [feature-cache-clean.md](feature-cache-clean.md) |
| 客户端事件上报 | 2（内外双注册） | controllers/event, service, common/event | [feature-client-event.md](feature-client-event.md) |
| 文件上传下载管理 | 6（2 app-api 双注册 + 4 file/v1 内部） | controllers/exfile, controllers/file, service | [feature-file-transfer.md](feature-file-transfer.md) |
| 浏览器配置查询与同步 | 2 | controllers/management, service | [feature-browser-config.md](feature-browser-config.md) |
| 插件包管理 | 5 | controllers/plugin, service | [feature-plugin-mgmt.md](feature-plugin-mgmt.md) |
| 流量会话统计与数据治理 | 4 + 1 定时任务 | controllers/traffic_stats, service, scheduler | [feature-traffic-stats.md](feature-traffic-stats.md) |
| 配置中心读写 | 2 | controllers/config_center, service | [feature-config-center.md](feature-config-center.md) |
| 证书更新订阅 | 3（SDK 订阅，启动注册） | common/cert, common/https | [feature-cert-subscribe.md](feature-cert-subscribe.md) |
| 终端鉴权（IMEI+IMSI 白名单） | 3（在用）+ 5 注入点 | controllers/auth, service/auth, dao | [feature-terminal-auth.md](feature-terminal-auth.md) |

未归类接口：`GET /test/v1/get`（src/controllers/test_controller.go；连通性测试桩，仅测试使用，非业务功能）
