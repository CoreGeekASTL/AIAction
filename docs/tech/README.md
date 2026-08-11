# AIAction 技术要素（tech）资产索引

> 生成时间：2026-08-11（终端鉴权变更刷新）
> 生成工具：all-index（自动聚合产物，同名覆盖更新；资产变更后重跑 all-index 刷新，请勿手改）

## usage/（框架使用现状）

基础框架清单与使用方式盘点

| 文件 | 说明 |
| --- | --- |
| [README.md](usage/README.md) | 框架使用现状 |
| [usage-beego.md](usage/usage-beego.md) | Beego v2 使用现状（网络/事件循环——HTTP/HTTPS Server） |
| [usage-beego-orm.md](usage/usage-beego-orm.md) | Beego ORM 使用现状（存储/ORM） |
| [usage-config.md](usage/usage-config.md) | 配置管理 使用现状（配置管理） |
| [usage-cspgo-monitor-sdk.md](usage/usage-cspgo-monitor-sdk.md) | CSPGoMonitorSDK 使用现状（监控/可观测） |
| [usage-encoding-json.md](usage/usage-encoding-json.md) | encoding/json + yaml.v2 使用现状（序列化/编解码） |
| [usage-go-chassis-extend.md](usage/usage-go-chassis-extend.md) | Go-chassis-extend（GSF + CSE）使用现状（服务注册发现/平台框架） |
| [usage-go-redis.md](usage/usage-go-redis.md) | go-redis/v9 使用现状（存储/缓存） |
| [usage-go-testing.md](usage/usage-go-testing.md) | Go testing + testify + goconvey 使用现状（测试框架） |
| [usage-https-client.md](usage/usage-https-client.md) | 自研 https 客户端封装 使用现状（网络/通信——出站 HTTP） |
| [usage-logger.md](usage/usage-logger.md) | 自研 logger 封装 + auditlog 使用现状（日志） |
| [usage-minio.md](usage/usage-minio.md) | minio-go/v7 使用现状（存储/对象存储） |
| [usage-scheduler.md](usage/usage-scheduler.md) | Go 协程原语 + 自研调度器 使用现状（并发/线程池、定时/调度） |
| [usage-utils.md](usage/usage-utils.md) | 基础库（uuid + 自研 utils）使用现状（基础库） |

## comm-guidelines/（通信规范）

RPC/HTTP/MQ 跨服务调用指导

| 文件 | 说明 |
| --- | --- |
| [README.md](comm-guidelines/README.md) | 通信规范（外部服务调用） |
| [comm-guidelines-alarm-sdk.md](comm-guidelines/comm-guidelines-alarm-sdk.md) | alarm-sdk 通信规范 |
| [comm-guidelines-browser-gateway.md](comm-guidelines/comm-guidelines-browser-gateway.md) | browser-gateway 通信规范 |
| [comm-guidelines-cse-servicecomb.md](comm-guidelines/comm-guidelines-cse-servicecomb.md) | cse-servicecomb 通信规范 |
| [comm-guidelines-csp-go-monitor-sdk.md](comm-guidelines/comm-guidelines-csp-go-monitor-sdk.md) | csp-go-monitor-sdk 通信规范 |
| [comm-guidelines-cspntp-sdk.md](comm-guidelines/comm-guidelines-cspntp-sdk.md) | cspntp-sdk 通信规范 |
| [comm-guidelines-cspsomf-sdk.md](comm-guidelines/comm-guidelines-cspsomf-sdk.md) | cspsomf-sdk 通信规范 |
| [comm-guidelines-fmservice.md](comm-guidelines/comm-guidelines-fmservice.md) | FMService 通信规范 |
| [comm-guidelines-gaussdb.md](comm-guidelines/comm-guidelines-gaussdb.md) | GaussDB 通信规范 |
| [comm-guidelines-moon.md](comm-guidelines/comm-guidelines-moon.md) | moon 通信规范 |

## concurrency-guidelines/（并发规范）

线程池选型、隔离、拒绝策略

| 文件 | 说明 |
| --- | --- |
| [concurrency-guidelines-alarm-event-channel.md](concurrency-guidelines/concurrency-guidelines-alarm-event-channel.md) | alarmEventChanel 告警事件通道 并发规范 |
| [concurrency-guidelines-auth-cache.md](concurrency-guidelines/concurrency-guidelines-auth-cache.md) | authCache 鉴权缓存与 authImportLock 并发规范 |
| [concurrency-guidelines-beego-https-server-cert.md](concurrency-guidelines/concurrency-guidelines-beego-https-server-cert.md) | BeegoHttpsServer 证书监听与服务拉起 并发规范 |
| [concurrency-guidelines-browser-gw-instances.md](concurrency-guidelines/concurrency-guidelines-browser-gw-instances.md) | browserGWInstances（sync.Map 实例缓存）并发规范 |
| [concurrency-guidelines-data-cleanup-scheduler.md](concurrency-guidelines/concurrency-guidelines-data-cleanup-scheduler.md) | DataCleanupScheduler 并发规范 |
| [concurrency-guidelines-db-connection.md](concurrency-guidelines/concurrency-guidelines-db-connection.md) | dbConnection 数据源连接与健康检查 并发规范 |
| [concurrency-guidelines-event-storage-factory.md](concurrency-guidelines/concurrency-guidelines-event-storage-factory.md) | eventStorageFactory 事件存储注册表锁 并发规范 |
| [concurrency-guidelines-get-active-alarm-async.md](concurrency-guidelines/concurrency-guidelines-get-active-alarm-async.md) | GetAllActiveAlarmFromFMService 异步查询 并发规范 |
| [concurrency-guidelines-local-event-storage-deleter.md](concurrency-guidelines/concurrency-guidelines-local-event-storage-deleter.md) | localEventStorage 转储文件清理定时器 并发规范 |
| [concurrency-guidelines-main-startup-goroutines.md](concurrency-guidelines/concurrency-guidelines-main-startup-goroutines.md) | main 启动 goroutine 组 并发规范 |
| [concurrency-guidelines-metric-map-lock.md](concurrency-guidelines/concurrency-guidelines-metric-map-lock.md) | metricMapLock（mocIdMap 保护锁）并发规范 |
| [concurrency-guidelines-monitor-service-schedule.md](concurrency-guidelines/concurrency-guidelines-monitor-service-schedule.md) | MonitorServiceImpl 监控上报定时任务 并发规范 |
| [concurrency-guidelines-pre-open-browser.md](concurrency-guidelines/concurrency-guidelines-pre-open-browser.md) | PreOpenBrowser 预开浏览器扇出 并发规范 |
| [concurrency-guidelines-progress-chan-plugin-load.md](concurrency-guidelines/concurrency-guidelines-progress-chan-plugin-load.md) | progressChan 插件加载进度管道 并发规范 |
| [concurrency-guidelines-start-refresh-config-task.md](concurrency-guidelines/concurrency-guidelines-start-refresh-config-task.md) | StartRefreshConfigTask 并发规范 |
| [concurrency-guidelines-sync-once-singletons.md](concurrency-guidelines/concurrency-guidelines-sync-once-singletons.md) | sync.Once 单例初始化组 并发规范 |

## data-access-guidelines/（数据访问规范）

Redis/DB 等中间件访问指导

| 文件 | 说明 |
| --- | --- |
| [data-access-guidelines-gaussdb.md](data-access-guidelines/data-access-guidelines-gaussdb.md) | GaussDB 数据访问规范 |
| [data-access-guidelines-minio.md](data-access-guidelines/data-access-guidelines-minio.md) | MinIO 数据访问规范 |
| [data-access-guidelines-redis.md](data-access-guidelines/data-access-guidelines-redis.md) | Redis 数据访问规范 |
| [data-access-guidelines-sqlite.md](data-access-guidelines/data-access-guidelines-sqlite.md) | SQLite 数据访问规范 |

## resilience-guidelines/（韧性规范）

超时/重试/熔断/异常处理

| 文件 | 说明 |
| --- | --- |
| [resilience-guidelines.md](resilience-guidelines/resilience-guidelines.md) | 韧性规范（故障策略） |

## foundation-guidelines/（基础规范）

日志/配置/告警等编码指导

| 文件 | 说明 |
| --- | --- |
| [foundation-guidelines.md](foundation-guidelines/foundation-guidelines.md) | 基础规范（日志 / 配置 / 告警） |
