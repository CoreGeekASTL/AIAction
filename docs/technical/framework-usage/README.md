# 基础框架使用指导

| 元信息 | 值 |
|--------|-----|
| 分支 | ready/27.0-终端鉴权 分支 (2026-08-07) |
| 更新日期 | 2026-08-07 |
| Skill | spec-framework-usage-analyze |

## 框架全景清单

| # | 类别 | 框架 | 使用指导 |
| --- | --- | --- | --- |
| 1 | RPC/通信（服务端） | Beego v2 Web（server/web） | [rpc-beego-web.md](rpc-beego-web.md) |
| 2 | RPC/通信（客户端） | net/http + 自研 Builder | [rpc-http-client.md](rpc-http-client.md) |
| 3 | RPC/服务发现 | go-chassis（CSE 注册发现 + rest invoker） | [rpc-go-chassis-cse.md](rpc-go-chassis-cse.md) |
| 4 | 并发/线程池 | Go 协程原语（goroutine/sync/channel） | [concurrency-goroutine-sync.md](concurrency-goroutine-sync.md) |
| 5 | 日志 | go-chassis lager + CSP auditlog | [log-lager.md](log-lager.md) |
| 6 | 序列化 | encoding/json + yaml.v2 | [codec-json-yaml.md](codec-json-yaml.md) |
| 7 | 配置管理 | Beego AppConfig + 自研 flagutil + DB 配置中心 | [config-beego-appconfig.md](config-beego-appconfig.md) |
| 8 | 依赖注入/组件 | 自研单例模式（无 DI 框架） | [di-singleton.md](di-singleton.md) |
| 9 | 存储/ORM | Beego ORM + 自研 GaussDB 驱动 + modernc sqlite | [storage-beego-orm.md](storage-beego-orm.md) |
| 10 | 存储/缓存 | go-redis/v9 | [storage-redis.md](storage-redis.md) |
| 11 | 存储/对象 | minio-go/v7 | [storage-minio.md](storage-minio.md) |
| 12 | 定时/调度 | 自研调度器 + time.Ticker/Timer | [schedule-ticker.md](schedule-ticker.md) |
| 13 | 容错/服务治理 | greatwall 限流 + 自研重试 | [resilience-greatwall.md](resilience-greatwall.md) |
| 14 | 监控/可观测 | CSPGoMonitorSDK 话统上报 | [metrics-csp-monitor.md](metrics-csp-monitor.md) |
| 15 | 基础库 | CSP/GSF 平台套件 + 自研 utils | [base-csp-gsf.md](base-csp-gsf.md) / [base-utils.md](base-utils.md) |
| 16 | 测试框架 | testing + testify + goconvey + gomockit | [test-go-testing.md](test-go-testing.md) |
| - | Actor 模型 | （未发现） | |
| - | 消息队列 | （未发现） | |
| - | 网络/事件循环 | （未发现，Beego/net.http 内建，无自研 Reactor） | |
| - | 资源池 | （未发现显式资源池；http.Transport 内置连接池复用） | |
