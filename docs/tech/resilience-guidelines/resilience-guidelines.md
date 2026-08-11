# 韧性规范（故障策略）

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-resilience-guidelines-analyze |
| 运行模式 | 起草模式 |

## 扫描范围总览

| 点位类型 | 数量 | 说明 |
|---|---|---|
| 出站调用点 | 9 | HTTP 7 处（moon 2 处、BrowserGW 3 处、DBService 1 处、FMService 1 处）、CSE 注册中心上报 1 处、外部存储（Redis/GaussDB）2 类封装点 |
| 后台任务 | 7 | 定时任务 3 个（数据清理、配置刷新、监控上报）、常驻 goroutine 4 个（告警事件循环、事件文件清理、GaussDB 连接保活、GSF 启动） |

| 策略 | 已设置点位数 | 未设置/框架默认点位数 | 现状一句话 |
|---|---|---|---|
| 超时 | 5 | 4 | 封装层 client 统一 120s/240s，个别调用点自设 5s/3s，InnerInstance 与 CSE rest 无显式超时 |
| 重试 | 7 | 2 | https builder 提供指数退避重试（2s 起、上限 60s），多数出站调用 WithRetry(2)；告警/CSE/DB 连接另有固定间隔重试 |
| 熔断与降级 | 1 | 8 | 全仓无熔断器；仅入口侧有 greatwall 过载过滤器（429 + Retry-After），个别点位失败后返回 nil/默认值兜底 |
| panic/recover 兜底 | 3 | 4 | 监控指标采集、插件控制器、File model init 有 recover；后台 goroutine 与 preOpen goroutine 均无兜底 |
| 错误 swallowing | — | 7 | 主要形态为"记日志后继续/返回 nil"，另有 1 处 `_` 丢弃 error、1 处 `_` 丢弃鉴权格式标志位 |

## 一、超时策略

### 现状调用点分布

| 调用点（文件 · 函数） | 被调目标 | 超时值与位置 | 证据（文件路径） |
|---|---|---|---|
| remote_service.go · MuenDeviceLogin | moon（沐恩云端）deviceLoginAuth | 封装层统一：ResponseHeader/TLSHandshake/ExpectContinue 120s、总超时 240s（https/client.go `Init`/`newHttpsClient`） | src/service/remote_service.go、src/common/https/client.go |
| management_controller.go · syncBrowserConfig | moon 配置同步 GET | 同上封装层 120s/240s | src/controllers/management_controller.go、src/common/https/client.go |
| plugin_service.go · loadPluginToBrowserGW | BrowserGW /browsergw/extension/load | 同上封装层 120s/240s（https.Instance） | src/service/plugin_service.go、src/common/https/client.go |
| browser_service.go · instancePreOpenBrowser | BrowserGW /browsergw/browser/preOpen | 同上封装层 120s/240s（https.Instance） | src/service/browser_service.go、src/common/https/client.go |
| cache_service.go · callBrowserGW | BrowserGW /browsergw/browser/userdata/delete | 5s，调用点自建 http.Client 显式设置（defaultCacheTimeoutSeconds） | src/service/cache_service.go |
| db_init.go · getDataSourceFromDBService | DBService getGaussdbInfor | 未设置（InnerInstance 仅配 TLSClientConfig，无 Timeout/Transport 超时项） | src/dao/db_init.go、src/common/https/client.go |
| alarm_service.go · OSHttpsGetRequestByCSE | FMService（经 CSE rest invoker） | 未设置（go-chassis rest 框架默认） | src/service/alarm_service.go |
| alarm_service.go · GetAllActiveAlarmFromFMService | FMService（异步调用 + select） | 3s 整体超时（time.After(TimePeriodInit)），调用点显式设置 | src/service/alarm_service.go |
| alarm_service.go · sendAlarmEvent | 告警事件 channel | 5s 写入超时（time.After(threshold)），超时后丢弃并记日志 | src/service/alarm_service.go |
| redis/redis.go · Init | Redis | 未设置（go-redis Options 未配 DialTimeout/ReadTimeout/WriteTimeout，框架默认） | src/common/storage/redis/redis.go |

### 应有约定建议

（以下为建议，非现状）出站 HTTP 调用应统一走 `common/https` 封装层 client，禁止在调用点裸建无超时 http.Client（cache_service 模式仅限特殊短超时场景并注释理由）；`InnerInstance` 建议补齐显式超时（建议 10s~30s），避免 DBService 不可达时启动流程无限挂起；按被调目标区分超时区间：内部微服务（BrowserGW/FMService）建议 3s~10s，外部云侧（moon）建议 30s~60s，超时值建议配置化（beego AppConfig key，如 `resilience::xxxTimeout`）；Redis 建议在 Init 时显式配置 DialTimeout/ReadTimeout/WriteTimeout。

## 二、重试策略

### 现状调用点分布

| 调用点（文件 · 函数） | 被调目标 | 重试次数/间隔/退避 | 重试触发条件 | 幂等保障 | 证据（文件路径） |
|---|---|---|---|---|---|
| https/builder.go · request.Do（封装层） | 全部走 builder 的 HTTP 调用 | 次数由 WithRetry 指定；指数退避 2s 起、上限 60s（InitialBackOff/MaxBackoff） | 429/502/503/504 + 网络层瞬时错误（timeout/ECONNREFUSED/ECONNRESET/io.EOF） | 无（POST 重试依赖被调方幂等） | src/common/https/builder.go |
| remote_service.go · MuenDeviceLogin | moon | WithRetry(defaultRetryCount=2)，走封装层退避 | 同上 | 无 | src/service/remote_service.go |
| plugin_service.go · loadPluginToBrowserGW | BrowserGW | WithRetry(2)，走封装层退避 | 同上 | 无 | src/service/plugin_service.go |
| management_controller.go · syncBrowserConfig | moon | WithRetry(2)，走封装层退避 | 同上 | GET 天然幂等 | src/controllers/management_controller.go |
| browser_service.go · instancePreOpenBrowser | BrowserGW | 无重试 | — | 不适用 | src/service/browser_service.go |
| cache_service.go · callBrowserGW | BrowserGW | 无重试 | — | DELETE 天然幂等 | src/service/cache_service.go |
| alarm_service.go · reportAlarm | 告警 SDK | 2 次，固定间隔 10s（alarmRetrySleepSeconds） | 发送失败（SDK 返回 false） | 同告警 ID 10 分钟内抑制重发（alarms map） | src/service/alarm_service.go |
| alarm_service.go · CleanAllActiveAlarm | FMService | 最多 360 次，固定间隔 5s（RetryTimes/TimePeriodClean），为升级场景 FM 未起设计 | err != nil | 查询天然幂等 | src/service/alarm_service.go |
| cse/cse.go · Report | CSE 注册中心 | 递归重试 maxRetry=5 次（main.go 传入），固定间隔 30s，耗尽后 logger.Fatalf 退出进程 | UpdateMicroServiceInstanceProperties err != nil | 属性覆盖写天然幂等 | src/common/cse/cse.go、src/main.go |
| main.go · initGSF | GSF/CSP 框架 | 最多 360 次，固定间隔 5s，耗尽后 logger.Fatalf | CspInit err != nil | 不适用 | src/main.go |
| dao/db_init.go · EnsureConnectGaussDB | GaussDB | 无限重试循环，固定间隔 5s（interval），无退出条件 | 连接/切换失败 | 不适用 | src/dao/db_init.go |
| scheduler/task_scheduler.go · executeCleanup | 数据清理（DB） | 3 次，固定间隔 10min，期间响应 stopChan 中止 | CleanOldStats err != nil | 清理操作按时间窗删除，重入安全 | src/scheduler/task_scheduler.go |

### 应有约定建议

（以下为建议，非现状）出站 HTTP 重试应统一使用 `https.Builder.WithRetry`，禁止调用点自造重试循环；重试次数建议不超过 3 次，POST 类非幂等调用重试前须确认被调方幂等性（如 deviceLoginAuth、extension/load）；无限重试仅限启动期必经依赖（GaussDB/CspInit），业务期禁止无界重试；指数退避建议增加随机抖动（jitter）防止多实例重试风暴；`cse.Report` 耗尽后 Fatalf 退出进程的策略建议改为告警 + 后台持续重试，避免注册中心短暂故障导致进程退出。

## 三、熔断与降级

### 现状调用点分布

| 调用点（文件 · 函数） | 被调目标 | 熔断器 | 降级逻辑 | 证据（文件路径） |
|---|---|---|---|---|
| filter.go · OverLoadFilter（入口侧） | 本服务全部对外接口 | 无（greatwall 过载控制器，非熔断器） | 过载时拒绝请求：429 + Retry-After: 3 | src/controllers/filter.go |
| remote_service.go · MuenDeviceLogin | moon | 无 | 失败返回 nil（由上层 login 流程按鉴权失败处理） | src/service/remote_service.go |
| cache_service.go · DeleteCacheImpl | BrowserGW（逐实例） | 无 | 单实例失败仅记日志、继续遍历其余实例，最终返回 nil | src/service/cache_service.go |
| db_init.go · getDataSource | DBService | 无 | DBService 查询失败且服务名为 GaussDB 时降级为本地配置文件连接串 | src/dao/db_init.go |
| plugin_service.go · loadPlugin | BrowserGW（逐实例） | 无 | 单实例失败 continue，全部失败置任务状态为 Failed | src/service/plugin_service.go |
| 其余出站调用（preOpenBrowser、syncBrowserConfig、alarm、CSE Report） | — | 无 | 直接失败（记日志上抛或 Fatalf） | src/service/browser_service.go、src/controllers/management_controller.go、src/service/alarm_service.go、src/common/cse/cse.go |

### 应有约定建议

（以下为建议，非现状）对 moon、BrowserGW 等可能级联故障的下游建议引入熔断器（如 go-chassis 自带熔断或自研计数式熔断），熔断打开期间按业务可接受口径降级：登录链路可快速失败并返回明确错误码，缓存删除/插件加载可采用队列暂存稍后补偿；降级返回值（如 MuenDeviceLogin 返回 nil）应有明确上层约定，避免 nil 被误读为"无此用户"；过载过滤（greatwall）建议保留并纳入统一韧性口径文档管理。

## 四、panic/recover 与异常兜底

### 现状点位分布

| 点位（文件 · 函数） | 点位类型 | 兜底方式 | 证据（文件路径） |
|---|---|---|---|
| monitor_service.go · getMetricResults | 定时任务内指标函数调用 | 有 recover 兜底（defer + recover + 记日志与堆栈） | src/service/monitor_service.go |
| plugin_controller.go · UploadPluginPackage | HTTP 请求处理 | 有 recover 兜底（defer + recover 记日志与堆栈） | src/controllers/plugin_controller.go |
| models/db/file.go · init | 包初始化 | 有 recover 兜底（记日志与堆栈） | src/models/db/file.go |
| scheduler/task_scheduler.go · run | 定时任务 goroutine | 无 recover 兜底（panic 将致进程崩溃） | src/scheduler/task_scheduler.go |
| config_center_service.go · StartRefreshConfigTask | 定时任务 goroutine | 无 recover 兜底 | src/service/config_center_service.go |
| monitor_service.go · startCspMonitor | 定时任务 goroutine | 无 recover 兜底（仅指标函数级有 recover） | src/service/monitor_service.go |
| alarm_service.go · handleEvent | 常驻消费循环 goroutine | 无 recover 兜底 | src/service/alarm_service.go |
| event/local_storage.go · NewLocalEventStorage 内清理 goroutine | 定时清理 goroutine | 无 recover 兜底 | src/common/event/local_storage.go |
| browser_service.go · instancePreOpenBrowser | 一次性 goroutine（每实例一个） | 无 recover 兜底 | src/service/browser_service.go |
| dao/db_init.go · EnsureConnectGaussDB | 常驻保活 goroutine | 无 recover 兜底（for 循环内错误靠 sleep 重试，不含 panic 防护） | src/dao/db_init.go |
| 入站 HTTP 请求处理 | Beego 框架 | 顶层框架兜底（beego RecoverPanic 默认行为） | src/main.go、src/routers |

### 应有约定建议

（以下为建议，非现状）每个独立执行单元（goroutine/定时任务/消费循环）入口必须有 defer recover 兜底，兜底后应记日志（含堆栈）并上报告警，禁止 recover 后静默；一次性 fan-out goroutine（instancePreOpenBrowser）建议统一封装为带 recover 的 `safeGo` 工具函数；定时任务 panic 恢复后应保证调度循环可继续（recover 放在循环体内每轮执行单元外层）。

## 五、错误 swallowing

### 现状点位分布

| 点位（文件 · 函数） | swallowing 形态 | 影响说明 | 证据（文件路径） |
|---|---|---|---|
| event_service.go · NewEventService | 赋值给 `_`（第二次 Get 的 error） | 兜底获取默认 EventStorage 再失败时被忽略，eventStorage 可能为 nil 导致后续 ReportEvent panic | src/service/event_service.go |
| cache_service.go · DeleteCacheImpl | 只记日志不上抛（遍历中 err 丢弃，函数恒返回 nil） | 部分 BrowserGW 缓存删除失败被掩盖，调用方无法感知部分失败 | src/service/cache_service.go |
| alarm_service.go · sendAlarmEvent | 只记日志不处理（channel 写超时后事件丢弃） | 告警事件丢失，运维无法感知未上报的告警 | src/service/alarm_service.go |
| plugin_service.go · loadPlugin / updatePluginProgress 循环 | 只记日志后 continue | 单实例加载失败/进度写库失败被跳过，进度可能与实际不符 | src/service/plugin_service.go |
| event/local_storage.go · rollOver / deleteOldBackFiles / getBackFiles | 只记日志后 return/continue | 事件文件转储或清理失败被吞，可能导致磁盘占满 | src/common/event/local_storage.go |
| https/builder.go · resetBody | err 判断后静默忽略（err == nil 才处理，否则不重置 body） | 重试时 body 未重置可能发出空 body 请求，故障表现为下游 4xx，难定位 | src/common/https/builder.go |
| login_controller.go / exlogin_controller.go / event_controller.go · AuthIMEI 调用点 | 第二返回值 formatValid 赋值给 `_`（`allowed, _ := c.authService.AuthIMEI(...)`） | 格式非法与鉴权拒绝在 login/event 注入点被合并为同一拒绝响应（-2/401），调用方无法区分；AuthIMEI 内部 DB 异常已按安全优先拒绝处理并记日志，风险可控 | src/controllers/login_controller.go、src/controllers/exlogin_controller.go、src/controllers/event_controller.go |

### 应有约定建议

（以下为建议，非现状）error 必须上抛或显式处理，记日志不等于处理；确需忽略的错误应注释说明理由（如 best-effort 清理）；fan-out 场景（DeleteCacheImpl、loadPlugin）建议汇总各实例失败结果并向上抛聚合错误或在响应中携带部分失败明细；channel 丢弃类（sendAlarmEvent）建议增加丢弃计数指标与告警；`event_service.go` 中 `_` 丢弃 error 应消除，失败时 panic 或返回明确错误。
