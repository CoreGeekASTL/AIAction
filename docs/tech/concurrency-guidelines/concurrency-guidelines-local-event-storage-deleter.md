# localEventStorage 转储文件清理定时器 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 定时任务 |

## 用途定位

本地事件日志转储文件（zip）定期清理：每小时检查一次备份目录，按保留天数（90 天）与最大个数（5 个）删除老旧 zip。从 `NewLocalEventStorage` 中 goroutine 实现推断。

- 代码标识符：`NewLocalEventStorage` 内匿名 goroutine / `deleteOldBackFiles`
- 定义位置：src/common/event/local_storage.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/common/event/local_storage.go（`NewLocalEventStorage`，经 event storage factory 初始化链路调用，见 src/service/event_service.go `initEventStorageFactory`）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 池化方式 | 裸 goroutine（无池化），每个 `localEventStorage` 实例各起 1 个 | src/common/event/local_storage.go |
| 容量配置 | 周期间隔硬编码 1 小时；清理策略参数常量 `FileRemainDay=90`、`FileMaxNum=5` | src/common/event/local_storage.go |
| 任务队列 | 无队列 | src/common/event/local_storage.go |
| 拒绝策略 | 未设置 | src/common/event/local_storage.go |
| 隔离范围 | 每实例独占 goroutine；清理与 `Record`/`rollOver` 之间无锁，`deleteOldBackFiles` 同时可被 ticker goroutine 与 `rollOver` 调用（并发调用时文件删除操作无互斥保护） | src/common/event/local_storage.go |
| 关闭与等待 | 未设置：使用 `time.Tick`（无法 Stop，ticker 永不释放），goroutine 无退出通道，随进程生命周期存在 | src/common/event/local_storage.go |
| 调度并发语义 | 单 goroutine 串行，不允许重入；任务体内未再开并发 | src/common/event/local_storage.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 允许裸 goroutine | 低频轻量任务 |
| 容量基线 | 间隔建议外置配置 | 不同局点磁盘水位不同 |
| 拒绝策略 | 不适用 | — |
| 隔离 | `deleteOldBackFiles` 存在 ticker goroutine 与 `rollOver`（Record 路径）并发调用的可能，建议加 `sync.Mutex` 串行化删除操作 | 并发删除同一批文件列表可能重复删除/报错 |
| 命名与观测 | `time.Tick` 应改为 `time.NewTicker` + 可 Stop；storage `Close()` 时应退出 goroutine | `time.Tick` 底层 ticker 无法回收，属已知泄漏模式 |
