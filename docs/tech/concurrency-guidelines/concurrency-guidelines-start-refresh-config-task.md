# StartRefreshConfigTask 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 定时任务 |

## 用途定位

配置中心缓存定时刷新：每 5 分钟从 DB 全量加载 `t_config_center` 到内存 map `configs`，供 `GetConfig` 快速读取。从定义处注释 `StartRefreshConfigTask` 与 `Refresh` 实现推断。

- 代码标识符：`StartRefreshConfigTask` / `configCenter.stopChan` / `RefreshInterval`
- 定义位置：src/service/config_center_service.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/main.go（`service.StartRefreshConfigTask()`）；`GetConfig` 读取点分散在各 service（src/service/config_center_service.go 内 `GetConfig`）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 池化方式 | 裸 goroutine（无池化），`go func()` 内 `time.NewTicker(RefreshInterval)` | src/service/config_center_service.go |
| 容量配置 | 刷新间隔硬编码 `RefreshInterval = 5 * time.Minute` | src/service/config_center_service.go |
| 任务队列 | 无队列；`stopChan chan struct{}` 无缓冲用于停止信号 | src/service/config_center_service.go |
| 拒绝策略 | 未设置（ticker 到点即执行，无拒绝概念） | src/service/config_center_service.go |
| 隔离范围 | 独占 goroutine；共享资源 `configCenter.configs`（map）被读写双方访问 | src/service/config_center_service.go |
| 关闭与等待 | `StopRefreshConfigTask()` 中 `close(stopChan)`；无 WaitGroup 等待退出 | src/service/config_center_service.go |
| 调度并发语义 | 单 goroutine 串行刷新，不允许重入；任务体内未再开并发 | src/service/config_center_service.go |

补充（共享状态事实）：`Refresh()` 整体替换 `c.configs = configMap`，`GetConfig()` 直接读 map，**未加任何锁**，读写存在 data race 风险（Go map 并发读写会 panic）。

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 单实例定时任务允许裸 goroutine | 低频任务池化无收益 |
| 容量基线 | 刷新间隔建议外置配置 key（如 `configcenter::refreshInterval`），默认值 5min | 不同局点刷新诉求不同 |
| 拒绝策略 | 不适用（定时拉模型） | — |
| 隔离 | `configs` map 应改用 `sync.RWMutex` 保护读写，或改用 `atomic.Value` 存整体快照 | 当前无锁读写 map，高并发下会 panic |
| 命名与观测 | 建议刷新失败次数打点上报；Stop 建议加 WaitGroup 保证优雅退出 | 避免 goroutine 泄漏与静默失败 |
