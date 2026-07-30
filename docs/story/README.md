# 功能软件要素文档

| 元信息 | 值 |
|--------|-----|
| 代码仓 | AIAction（GlobalInstanceDeliverService，go module: GIDS） |
| 分析基准 | main 分支 (2026-07-30) |
| 更新时间 | 2026-07-30 |
| Skill | spec-feature-analyze |
| 主要语言 | Go |
| 分析范围 | 全部接口类型（IDL 契约 / 框架路由 / 消息订阅定时） |

> 由 spec-feature-analyze 生成/更新，面向人与 AI 共同消费。仓内无 IDL 契约文件，全部为 Beego 框架路由 + SDK 订阅 + 定时任务（自研隐式接口）。

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

## 接口统计

- 对外接口：36 个（IDL 契约 0 / 框架路由 32 唯一路径 / 消息订阅 3 + 定时任务 1）
- 设计中：0 个（终端鉴权 3 接口已于 27.0 落地为在用；另含 5 个既有链路注入点）
- 已下线：0 个（各功能文档接口表"状态"列均为在用）
- 说明：语言级内部接口（仓内模块间契约，如 service.XxxService interface）仅用于分析，不写入功能文档
- 框架文档目录说明：docs/framework-usage/ 下子文档由 spec-framework-usage-analyze 生成中，本批文档链接以其 README 索引清单为准

## 未归类接口

以下接口探测到但未纳入任何功能域，原因逐条说明：

- `GET /test/v1/get`（src/controllers/test_controller.go；连通性测试桩，仅测试使用，非业务功能）

## 使用说明

- **新人上手**：每篇先读第 1 节「功能故事」（多彩建模图+术语表），再按需深入 L2 结构地图。
- **AI 编码时**：L1 建立业务认知后重点读「AI 编码指南」，再按"接口清单 → 调用关系"定位改动点。
