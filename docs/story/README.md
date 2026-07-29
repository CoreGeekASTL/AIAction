# 功能软件要素文档

> 由 spec-feature-analyze 生成/更新，面向人与 AI 共同消费。
> 代码仓：GIDS（GlobalInstanceDeliverService，云浏览器全局实例交付服务，Go + Beego v2）　分析基准：personal/houle/test3 @ ae0a8a6（src/ 最后变更 b09632b）　更新时间：2026-07-29（全量复扫核对，与 2026-07-28 版无漂移）
> 分析范围：全部接口类型（IDL 契约 / 框架路由 / 消息订阅定时）

服务形态：双 HTTP 监听——外部 HTTPS（默认 40051，src/routers/beego_router.go `RegisterExternalRouter`）与内部 HTTP（默认 9090，同文件 `RegisterInternalRouter`）；路由由各 Controller 的 `RouteInfo().RouteMapping` 集中声明。无 IDL 契约文件（proto/thrift/OpenAPI 均未命中），接口以框架路由为准。

## 功能全景

| 功能域 | 接口数 | 核心模块 | 文档 |
|---|---|---|---|
| 登录鉴权 | 6（外部 3 / 内部 6，3 个登录接口双监听同名重复注册） | controllers(exlogin/login)、service(user/browser)、dao(user)、common/cse | [feature-login-auth.md](feature-login-auth.md) |
| 文件管理 | 8（外部 2 / 内部 6） | controllers(exfile/file)、service(file)、dao(file) | [feature-file-manage.md](feature-file-manage.md) |
| 插件管理 | 5（全部内部） | controllers(plugin)、service(plugin)、dao(plugin) | [feature-plugin-manage.md](feature-plugin-manage.md) |
| 流量统计 | 4 内部 HTTP + 1 定时清理任务 | controllers(traffic_stats)、service(traffic_stats)、dao(traffic_stats)、scheduler | [feature-traffic-stats.md](feature-traffic-stats.md) |
| 事件上报 | 2（外/内双监听同注册） | controllers(event)、service(event)、common/event | [feature-event-report.md](feature-event-report.md) |
| 缓存管理 | 1（外/内双监听同注册） | controllers(cache)、service(cache) | [feature-cache-manage.md](feature-cache-manage.md) |
| 配置管理 | 4（全部内部） | controllers(management/config_center)、service(config_center/alarm)、dao(browser_config/config_center) | [feature-config-manage.md](feature-config-manage.md) |
| 证书订阅 | 1 个订阅入口（3 个证书场景，消息订阅类，非 HTTP） | common/cert、common/https、CertSDK(stub) | [feature-cert-subscribe.md](feature-cert-subscribe.md) |
| 终端鉴权 | 3 个新增（设计中）+ 4 个既有链路注入点 | controllers(auth)、service(auth/auth_manage)、dao(white_list) | [feature-terminal-auth.md](feature-terminal-auth.md) |

## 接口统计

- 去重后 HTTP 业务路径 28 个：外部监听暴露 9 个（登录 3、文件 2、事件 2、缓存 1、测试桩 1），内部监听注册 28 个（含与外部同名重复的 8 个）。
- 非 HTTP 入口：证书订阅 3 个场景（src/common/cert/cert.go）、流量统计定时清理 1 个（src/scheduler/task_scheduler.go）。
- 语言级内部接口（仓内模块间契约）仅用于分析理解，不写入功能文档。
- 设计中：终端鉴权新增 3 个（/auth/v1/*，见 feature-terminal-auth.md），落地后计入上表。
- 已下线：0 个（探测到的接口均有路由注册，无注释残留路由）。

## 框架引用

各功能文档第 6 节「框架引用」逐功能列出使用的基础框架，框架使用指导见 [../framework-usage/README.md](../framework-usage/README.md)。

## 未归类接口

以下接口/组件探测到但未纳入任何功能域，原因逐条说明：

- `GET /test/v1/get`（TestController，src/controllers/test_controller.go，注册于外部监听）——测试连通性桩接口，非业务功能。
- `MonitorService`（src/service/monitor_service.go，main.go 启动）——指标采集上报内部支撑组件，无对外 HTTP 接口，属于平台监控框架范畴。
- `common/storage/redis` 包（src/common/storage/redis/redis.go）——已封装 Client/Object 接口但全仓无 importer、main 未初始化，属未接线代码；缓存管理实际走 HTTP 到 BrowserGW（详见 [feature-cache-manage.md](feature-cache-manage.md)）。
- `common/storage/oss` minio 封装（src/common/storage/oss/minio.go）——同样无调用方，文件实际落 DB（详见 [feature-file-manage.md](feature-file-manage.md)）。

## 使用说明

- **新人上手**：从功能全景表选入口，每篇先读第 1 节「功能故事」——多彩建模图（粉=事件、黄=角色、绿=实体、蓝=规则）讲清"谁触发、经过什么业务环节、改变了什么、按什么规则"，术语表逐个翻译业务黑话；读懂故事后再按需深入 L2 结构地图（模块划分 → 调用关系）。建议先读登录鉴权（主业务链路），再读文件/插件/配置。
- **AI 编码时**：按需求所属功能查阅对应 md，L1 建立业务认知后，重点读「AI 编码指南」，再按"接口清单 → 调用关系"定位改动点；跨功能改动需同时阅读涉及的多篇。
- **更新策略**：功能接口增删时按功能逐篇更新对应 md 并同步本索引；功能下线在接口表"状态"列标注，不删历史。
