# GetAllActiveAlarmFromFMService 异步查询 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 锁 + 裸并发(goroutine) |

## 用途定位

从 FMService 异步拉取活动告警：起裸 goroutine 调 `handlerActivityAlarmData`（HTTP 调用），结果经无缓冲 channel `ch` 回传，主 goroutine select 等待结果或 3 秒超时。包级 `mutex` 串行化整个查询。从实现推断（防止并发查询 FMService）。

- 代码标识符：`mutex`（包级 var）/ `GetAllActiveAlarmFromFMService`
- 定义位置：src/service/alarm_service.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/service/alarm_service.go（`GetAllActiveAlarmFromFMService` 自身，以及调用方 `CleanAllActiveAlarm` 的 5s 重试循环）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 锁类型 | sync.RWMutex（包级 `var mutex sync.RWMutex`），实际只用写锁 `Lock()` | src/service/alarm_service.go |
| 锁粒度 | 整个函数临界区：从发起查询到拿到结果/超时才释放 | src/service/alarm_service.go |
| 锁顺序约定 | 未设置 | src/service/alarm_service.go |
| 临界区说明 | 临界区内起 goroutine 做 HTTP 调用，select 等待 `ch` 或 `time.After(3s)`；超时后 goroutine 仍持有 `mapAlarms/errRet` 闭包变量继续运行（写完 channel 后退出，无泄漏但结果丢弃） | src/service/alarm_service.go |
| 池化方式 | 裸 goroutine（无池化），每次调用起 1 个 | src/service/alarm_service.go |
| 容量配置 | 超时硬编码 `TimePeriodInit = 3`（秒）；调用方重试 `TimePeriodClean = 5`（秒）、`RetryTimes = 360` | src/service/alarm_service.go |
| 任务队列 | 无缓冲 channel `ch chan string` 做结果回传 | src/service/alarm_service.go |
| 拒绝策略 | 未设置（锁保证串行，无拒绝） | src/service/alarm_service.go |
| 隔离范围 | 保护对象：对 FMService 的并发查询；全局串行 | src/service/alarm_service.go |
| 关闭与等待 | 未设置 | src/service/alarm_service.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 一次性异步查询可用裸 goroutine，但建议改用 `context.WithTimeout` + 同步 HTTP 客户端超时，去掉 goroutine+channel 结构 | 现有模式超时后 goroutine 继续跑，语义绕 |
| 容量基线 | 超时应由 HTTP 客户端超时承担，而非 select+time.After 双保险 | 双超时语义易不一致 |
| 拒绝策略 | 串行化（锁）合理；调用方 `CleanAllActiveAlarm` 重试 360 次 * 5s ≈ 30min 在 main goroutine 启动阶段执行，建议评估是否阻塞启动（待确认调用时机） | 长时间持锁会挡住其他查询方 |
| 隔离 | 建议 `RWMutex` 改为 `Mutex`（未使用读锁，名实不符） | 减少误读 |
| 命名与观测 | 包级裸名 `mutex` 建议改名 `fmQueryMutex` 并加注释说明保护对象 | 包级变量语义不清 |
