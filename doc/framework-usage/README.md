# 基础框架使用指导

> 生成时间：2026-07-27 ｜ 代码基线：main @ 5e78a48 ｜ 分析深度：深（资产级）
> 用途：AI 编码与新人上手时，按部件查阅对应框架的使用指导；每框架一篇 md。

项目：GIDS（GlobalInstanceDeliverService），Go 1.25 + Beego v2.1.0，go module `GIDS`，源码在 `src/`。
特点：华为 CSP 局点服务，大量内部 SDK（GSF/CSE、AlarmSDK、CSPGoMonitorSDK 等）在本地以 `src/stubs/` 空实现替换（见 `src/go.mod:66-83` replace 块）；`LOCAL_MODE=true` 时使用嵌入式 SQLite 替代 GaussDB。

## 框架全景清单

| # | 类别 | 框架 | 版本 | 调用点数 | 涉及文件数 | 封装层 | 使用指导 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | RPC/通信（服务端） | Beego v2 Web | v2.1.0 | ~40 | 14 controllers + routers | 有（`common/https` server、`BaseController`） | [rpc-beego-web.md](rpc-beego-web.md) |
| 2 | RPC/通信（客户端） | 自研 HTTP Builder + net/http + GSF RestInvoker | — | ~10 | 6 | 即封装层（`common/https` builder） | [rpc-http-client-builder.md](rpc-http-client-builder.md) |
| 3 | 存储/ORM | Beego ORM + GaussDB(openGauss-pq) / SQLite(modernc) | beego v2.1.0 / v1.0.7 / v1.54.0 | 全量 DAO | models/db 7 实体 + dao 9 文件 | 有（`dao.BaseDao`/`BaseInterface`） | [storage-beego-orm.md](storage-beego-orm.md) |
| 4 | 存储/缓存 | go-redis v9 封装 | v9.0.5 | 0 业务调用点 | 1（仅封装层） | 有（`common/storage/redis`） | [storage-redis.md](storage-redis.md) |
| 5 | 存储/对象 | minio-go v7 封装 | v7.0.69 | 0 业务调用点 | 1（仅封装层） | 有（`common/storage/oss`） | [storage-minio.md](storage-minio.md) |
| 6 | 容错/服务治理 | GSF/CSE 服务发现 + greatwall 过载控制 | 内部 SDK（stub） | ~15 | 4 | 有（`common/cse`） | [resilience-cse-gsf.md](resilience-cse-gsf.md) |
| 7 | 日志 | lager 业务日志 + fusionstage auditlog 审计 + 自研 event 本地事件存储 | 内部 SDK（stub） | 全量 | 全部模块 | 有（`common/logger`、`common/event`） | [log-lager-auditlog-event.md](log-lager-auditlog-event.md) |
| 8 | 配置管理 | Beego AppConfig(ini) + flagutil 命令行 + DB 配置中心 | — | ~25 | 6 | 有（`common/conf`、`service/config_center_service`） | [config-appconf-flagutil-configcenter.md](config-appconf-flagutil-configcenter.md) |
| 9 | 定时/调度 | 自研 time.Timer/Ticker 调度（无 cron 库） | 标准库 | 6 | 5 | 部分（`scheduler.DataCleanupScheduler`） | [schedule-timer.md](schedule-timer.md) |
| 10 | 并发/线程池 | Go 协程原语 + sync.Once 包级单例 | 标准库 | 17+ | 10+ | 无独立线程池 | [concurrency-goroutine.md](concurrency-goroutine.md) |
| 11 | 监控/可观测 | CSPGoMonitorSDK 话统 + AlarmSDK 告警 | 内部 SDK（stub） | ~10 | 2 | 有（`service/monitor_service`、`service/alarm_service`） | [metrics-csp-monitor-alarm.md](metrics-csp-monitor-alarm.md) |
| 12 | 序列化 | encoding/json + yaml.v2 + encoding/csv | 标准库 / v2.4.0 | 43+ | 18 | 无 | [codec-json-yaml.md](codec-json-yaml.md) |
| 13 | 测试框架 | testify + goconvey + gomockit + Python testsuit | v1.8.4 / v1.8.1 / v1.1.0 | 17 测试文件 | src 全模块 + testsuit | 无 | [test-testify-goconvey.md](test-testify-goconvey.md) |
| 14 | 基础库 | google/uuid + 自研 utils（fileutil/monitorutil/response） | v1.6.0 | ~10 | 5 | 无 | [base-uuid-utils.md](base-uuid-utils.md) |
| - | Actor 模型 | （未发现） | | | | | |
| - | 消息队列 | （未发现；事件通过本地文件 + HTTP 上报，无 MQ） | | | | | |
| - | 网络/事件循环 | （未发现独立框架；由 Beego/net-http 内建承担） | | | | | |
| - | 资源池 | （未发现显式池；连接池由 database/sql、go-redis 内建） | | | | | |
| - | 依赖注入 | （无 DI 框架；以包级变量 + sync.Once 单例代替，见 concurrency-goroutine.md） | | | | | |

## 不一致与风险汇总

| 风险 | 涉及框架 | 位置 | 建议 |
| --- | --- | --- | --- |
| HTTP 客户端双轨：自研 builder vs 裸 `http.Client` | rpc-http-client-builder | `src/service/cache_service.go:49` 裸 new http.Client | 新代码必须用 `https.NewRequest()` builder，详见该 md |
| 配置三轨：app.conf / 命令行 flag / DB 配置中心优先级易混淆 | config | `src/service/remote_service.go:54-77` 手工"DB 覆盖文件" | 沿用既有"DB 配置覆盖 app.conf"模式，勿新增第四来源 |
| Redis/OSS 封装层存在但零业务调用（死代码/预置能力） | storage-redis/minio | `src/common/storage/` | 使用前需自行接线 Init；不要假设它已被初始化 |
| 生产 GaussDB 与 LOCAL_MODE SQLite DDL 双份维护，易漂移 | storage-beego-orm | `src/dao/db_init.go:30` vs `src/dao/db_local_sqlite.go:23` | 改表结构必须两处同步 |
| 日志 Errorf 被当作调试输出滥用（大量 `s00893267` 工号日志含敏感连接串） | log | `src/dao/db_init.go:194,335,354,369,376` | 新代码禁止打印密码/连接串；按级别规范使用 |
| 内部 SDK 全部 stub，本地无法验证真实行为 | resilience/metrics/log | `src/stubs/` | 平台相关逻辑变更需在真实环境验证 |
| 主备切换用全局可变 `ormer` 变量，非线程安全写入 | storage-beego-orm | `src/dao/db_init.go:188,203` 与 `base_dao.go:18` | 改造需谨慎；新 DAO 一律走 `BaseDao`，勿自管 ormer |

## 使用说明

- AI 编码时：先按需求涉及的部件查阅对应 md 文末「AI 编码指南」，再动手写代码。
- 新增框架或框架用法变化时：更新对应 md，并同步本索引的全景清单。
- 项目级代码风格基线另见 `AGENTS.md`「Go代码质量基线」与 `.claudecode/skills/code-quality-check/`。
