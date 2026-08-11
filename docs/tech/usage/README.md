# 框架使用现状

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-usage-analyze |

## 框架全景清单

| # | 类别 | 框架 | 使用现状文档 |
| --- | --- | --- | --- |
| 1 | RPC/通信 | 未发现 RPC 框架（跨服务调用均为 HTTP，见「网络/事件循环」自研 https 封装） | |
| 2 | 并发/线程池 | Go 协程原语（goroutine + channel + sync，含自研 DataCleanupScheduler 调度器） | [usage-scheduler.md](usage-scheduler.md) |
| 3 | Actor 模型 | （未发现） | |
| 4 | 日志 | 自研 logger 封装（基于 Go-chassis-extend lager）+ auditlog 事件日志 | [usage-logger.md](usage-logger.md) |
| 5 | 序列化/编解码 | Go encoding/json + gopkg.in/yaml.v2 | [usage-encoding-json.md](usage-encoding-json.md) |
| 6 | 配置管理 | Beego AppConfig（app.conf）+ 自研 conf 包 + 环境变量 | [usage-config.md](usage-config.md) |
| 7 | 依赖注入/组件管理 | 未发现 DI 框架（手动单例：包级变量 + init() + sync.Once） | |
| 8 | 存储/ORM | Beego ORM（GaussDB/Postgres + LOCAL_MODE 嵌入式 SQLite） | [usage-beego-orm.md](usage-beego-orm.md) |
| 9 | 存储/缓存 | go-redis/v9（自研 Client 封装） | [usage-go-redis.md](usage-go-redis.md) |
| 10 | 存储/对象存储 | minio-go/v7（自研 oss.Client 封装） | [usage-minio.md](usage-minio.md) |
| 11 | 消息队列 | （未发现） | |
| 12 | 定时/调度 | time.Timer/time.Ticker（无 cron 框架） | [usage-scheduler.md](usage-scheduler.md) |
| 13 | 网络/事件循环 | Beego v2 HTTP/HTTPS Server + 自研 https 客户端封装（net/http builder + 重试） | [usage-beego.md](usage-beego.md)、[usage-https-client.md](usage-https-client.md) |
| 14 | 资源池 | 未发现独立资源池框架（DB 连接池由 Beego ORM 内部管理） | |
| 15 | 容错/服务治理 | 无独立容错框架；重试/退避内建于 https builder，限流用 greatwall-sdk-go OverLoadFilter（见 usage-https-client.md） | |
| 16 | 监控/可观测 | CSPGoMonitorSDK（CSP 话统监控上报，stubs 本地桩） | [usage-cspgo-monitor-sdk.md](usage-cspgo-monitor-sdk.md) |
| 17 | 服务注册发现 | Go-chassis-extend（GSF 框架 + CSE 注册/发现/Watch） | [usage-go-chassis-extend.md](usage-go-chassis-extend.md) |
| 18 | 基础库 | google/uuid、自研 utils（fileutil/monitorutil/flagutil）、retcode | [usage-utils.md](usage-utils.md) |
| 19 | 测试框架 | Go testing + testify + goconvey | [usage-go-testing.md](usage-go-testing.md) |
