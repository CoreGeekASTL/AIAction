# dbConnection 数据源连接与健康检查 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 锁 + 定时任务 |

## 用途定位

GaussDB 数据源连接管理：`dbConnection`（内嵌 sync.Mutex）保护当前数据源、已连接数据源表与健康检查失败计数；`EnsureConnectGaussDB` 由 main 以 goroutine 启动负责建连；`checkDBStatus` goroutine 每 5s 做一次健康检查，失败时切换备用数据源。从实现推断。

- 代码标识符：`dbConnection` / `EnsureConnectGaussDB` / `checkDBStatus` / `interval` / `maxHealthCheckCount`
- 定义位置：src/dao/db_init.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/main.go（`go dao.EnsureConnectGaussDB()`）；src/dao/db_init.go（`go checkDBStatus()`、`healthCheck`、`getCurrentDataSource`、`switchToAnotherDB`）

## 线程模型（可选）

```mermaid
flowchart LR
    main["main goroutine"] --> ensure["EnsureConnectGaussDB goroutine"]
    ensure --> check["checkDBStatus goroutine"]
    check --> ticker["time.Ticker 5s"]
    ticker --> hc["healthCheck 持 db.Lock"]
    hc -->|失败达到阈值| sw["switchToAnotherDB 持 db.Lock"]
```

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 锁类型 | sync.Mutex（`dbConnection` 内嵌） | src/dao/db_init.go |
| 锁粒度 | `healthCheck` / `getCurrentDataSource` / `switchToAnotherDB` 均全程持锁；`healthCheck` 临界区含 `Ping`（IO） | src/dao/db_init.go |
| 锁顺序约定 | 未设置（单锁） | src/dao/db_init.go |
| 临界区说明 | healthCheck：GetDB + Ping + 失败计数更新；switchToAnotherDB：RegisterDataBase + 换 ormer（包级变量 `ormer` 的并发访问未见该锁保护，待确认） | src/dao/db_init.go |
| 池化方式 | 裸 goroutine（无池化）：`go checkDBStatus()` 单实例循环 | src/dao/db_init.go |
| 容量配置 | 检查间隔硬编码 `interval = 5 * time.Second`；失败阈值 `maxHealthCheckCount = 3` | src/dao/db_init.go |
| 任务队列 | 无队列 | src/dao/db_init.go |
| 拒绝策略 | 未设置 | src/dao/db_init.go |
| 隔离范围 | 保护对象：`dbConnection` 内部状态；goroutine 独占健康检查职责 | src/dao/db_init.go |
| 关闭与等待 | 未设置（goroutine 无退出通道） | src/dao/db_init.go |
| 调度并发语义 | 单 goroutine ticker 串行检查，不允许重入；任务体内未再开并发 | src/dao/db_init.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 允许裸 goroutine | 单实例低频任务 |
| 容量基线 | 间隔与阈值建议外置配置 | 不同局点 DB 抖动容忍度不同 |
| 拒绝策略 | 不适用 | — |
| 隔离 | `healthCheck` 持锁做 Ping（IO），会阻塞 `getCurrentDataSource`/`switchToAnotherDB`；建议 Ping 移出临界区（锁内只读写计数与状态） | 锁内 IO 放大阻塞面 |
| 命名与观测 | 包级 `ormer` 在 `switchToAnotherDB` 中被替换，读写路径未统一走锁，建议明确 ormer 的并发访问约定（待确认实际访问点） | 潜在 data race |
