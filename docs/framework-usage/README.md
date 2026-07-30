# 基础框架使用指导

| 元信息 | 值 |
|--------|-----|
| 代码仓 | AIAction（GIDS，GlobalInstanceDeliverService） |
| 分析基准 | main 分支 (2026-07-30, commit 6c93561) |
| 更新时间 | 2026-07-30 |
| Skill | spec-framework-usage-analyze |
| 主要语言 | Go 1.25（go module: GIDS，源码在 src/） |
| 分析深度 | 深（资产级） |

> 用途：AI 编码与新人上手时，按部件查阅对应框架的使用指导；每框架一篇 md。

## 框架全景清单

| # | 类别 | 框架 | 版本 | 调用点数 | 涉及文件数 | 封装层 | 使用指导 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | RPC/通信（服务端） | Beego v2 Web（server/web） | v2.1.0 | ~30 | 15 | 有（common/https 服务器包装 + controllers.BaseController） | [rpc-beego-web.md](rpc-beego-web.md) |
| 2 | RPC/通信（客户端） | net/http + 自研 Builder | — | NewRequest 5+裸用 1 | 8 | 即封装层（common/https builder.go） | [rpc-http-client.md](rpc-http-client.md) |
| 3 | RPC/服务发现 | go-chassis（CSE 注册发现 + rest invoker） | Go-chassis-extend v0.0.0（stub） | 9 文件 | 5 | 有（common/cse） | [rpc-go-chassis-cse.md](rpc-go-chassis-cse.md) |
| 4 | 并发/线程池 | Go 协程原语（goroutine/sync/channel） | 标准库 | 17+ | 10+ | 无（直接约定） | [concurrency-goroutine-sync.md](concurrency-goroutine-sync.md) |
| 5 | 日志 | go-chassis lager + CSP auditlog | stub | ~40 文件 | 40 | 有（common/logger） | [log-lager.md](log-lager.md) |
| 6 | 序列化 | encoding/json + yaml.v2 | stdlib / v2.4.0 | 43+1 | 19 | 无 | [codec-json-yaml.md](codec-json-yaml.md) |
| 7 | 配置管理 | Beego AppConfig + 自研 flagutil + DB 配置中心 | v2.1.0 | ~25 | 12 | 有（common/conf） | [config-beego-appconfig.md](config-beego-appconfig.md) |
| 8 | 依赖注入/组件 | 自研单例模式（无 DI 框架） | — | 全部 service | 12 | 即约定本身 | [di-singleton.md](di-singleton.md) |
| 9 | 存储/ORM | Beego ORM + 自研 GaussDB 驱动 + modernc sqlite | v2.1.0 / v1.0.7 / v1.54.0 | 20 文件 | 15 | 有（dao.BaseDao） | [storage-beego-orm.md](storage-beego-orm.md) |
| 10 | 存储/缓存 | go-redis/v9 | v9.0.5 | 接口 18 方法 | 1 | 有（common/storage/redis） | [storage-redis.md](storage-redis.md) |
| 11 | 存储/对象 | minio-go/v7 | v7.0.69 | 接口 5 方法 | 1 | 有（common/storage/oss） | [storage-minio.md](storage-minio.md) |
| 12 | 定时/调度 | 自研调度器 + time.Ticker/Timer | 标准库 | 5 | 4 | 有（scheduler.DataCleanupScheduler） | [schedule-ticker.md](schedule-ticker.md) |
| 13 | 容错/服务治理 | greatwall 限流 + 自研重试 | v1.9.6（stub） | 1+多处 | 5 | 有（controllers.OverLoadFilter） | [resilience-greatwall.md](resilience-greatwall.md) |
| 14 | 监控/可观测 | CSPGoMonitorSDK 话统上报 | stub | 1 | 1 | 有（service.MonitorService） | [metrics-csp-monitor.md](metrics-csp-monitor.md) |
| 15 | 基础库 | CSP/GSF 平台套件 + 自研 utils | stub | 多 | 10+ | 有 | [base-csp-gsf.md](base-csp-gsf.md) / [base-utils.md](base-utils.md) |
| 16 | 测试框架 | testing + testify + goconvey + gomockit | v1.8.4 / v1.8.1 / v1.1.0 | 17 | 17 | 有（test/util.It） | [test-go-testing.md](test-go-testing.md) |
| - | Actor 模型 | （未发现） | | | | | |
| - | 消息队列 | （未发现） | | | | | |
| - | 网络/事件循环 | （未发现，Beego/net.http 内建，无自研 Reactor） | | | | | |
| - | 资源池 | （未发现显式资源池；http.Transport 内置连接池复用） | | | | | |

## 不一致与风险汇总

| 风险 | 涉及框架 | 位置 | 建议 |
| --- | --- | --- | --- |
| HTTP 出站调用两轨并存：封装 Builder vs 裸 `client.Do` | rpc-http-client | `src/service/cache_service.go:78` 裸用 | 新代码一律走 `https.NewRequest().WithRetry()`，裸用点标注为禁模仿区 |
| 数据库连接串、内部错误等敏感/调试信息用 Errorf 输出 | log-lager | `src/dao/db_init.go:194,225,335`（含工号前缀 "s00893267"） | 新代码禁止日志输出密码/连接串；调试日志随需求清理 |
| 定时任务两套风格：自研 Scheduler 类 vs 散落 `go func()+Ticker` | schedule-ticker | `src/service/monitor_service.go:123`、`src/service/config_center_service.go:106`、`src/common/event/local_storage.go:77` | 周期任务优先复用 scheduler 模式；散落 ticker 至少保留 stopChan 退出通道 |
| 日志双轨：common/logger 与标准库 `log.Println` | log-lager | `src/service/alarm_service.go:95` | 新代码只用 `GIDS/common/logger` |
| service 单例三种实现：包级 var+init / sync.Once / 每次 new | di-singleton | `service/config_center_service.go:91`、`service/event_service.go:37`、`service/browser_service.go:47` | 无状态 service 用包级单例；带不可共享状态的每次 new |
| 配置三来源并存：app.conf 静态 / 环境变量 / DB 配置中心 | config-beego-appconfig | `src/service/remote_service.go:54` | 读取顺序约定：DB 配置中心 > 环境变量 > app.conf 默认值 |

## 使用说明

- AI 编码时：先按需求涉及的部件查阅对应 md 文末「AI 编码指南」，再动手写代码。
- 新增框架或框架用法变化时：更新对应 md，并同步本索引的全景清单。
- stubs/ 目录为华为内部 SDK 的本地空实现（go.mod replace），分析以 src/ 业务代码的调用方式为准，生产环境链接真实 SDK。
