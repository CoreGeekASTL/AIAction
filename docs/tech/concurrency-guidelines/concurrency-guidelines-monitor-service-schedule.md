# MonitorServiceImpl 监控上报定时任务 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 定时任务 |

## 用途定位

CSP 话统监控上报：先重试注册 CSPGoMonitorSDK（失败每 10s 重试直至成功），成功后每 5 分钟执行 `monitorSchedule`，遍历监控模板中的指标组，调用对应统计函数取数并调用 SDK `SetMetric` 上报。从定义处注释与实现推断。

- 代码标识符：`MonitorServiceImpl.InitMonitorSchedule` / `startCspMonitor` / `DotPeriodFiveMin` / `InitMonitorPeriod`
- 定义位置：src/service/monitor_service.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/main.go（`go monitorService.InitMonitorSchedule()`）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 池化方式 | 裸 goroutine（无池化）：main 中 `go InitMonitorSchedule()`；内部 `startCspMonitor` 再 `go func()` 起 ticker 循环 | src/service/monitor_service.go; src/main.go |
| 容量配置 | 间隔硬编码 `DotPeriodFiveMin = 5 * time.Minute`、注册重试 `InitMonitorPeriod = 10 * time.Second`（无限重试） | src/service/monitor_service.go |
| 任务队列 | 无队列；`stopChan chan struct{}` 无缓冲用于停止信号 | src/service/monitor_service.go |
| 拒绝策略 | 未设置 | src/service/monitor_service.go |
| 隔离范围 | 独占 goroutine；指标函数逐组串行执行，单指标 panic 通过 `getMetricResults` 内 recover 兜底 | src/service/monitor_service.go |
| 关闭与等待 | `Stop()` 中 `close(stopChan)`；无 WaitGroup 等待；main 未调用 Stop | src/service/monitor_service.go; src/main.go |
| 调度并发语义 | 单 goroutine 串行，不允许重入；任务体内未再开并发 | src/service/monitor_service.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 允许裸 goroutine，但注册重试循环与打点循环建议在同一个 goroutine 生命周期内管理（现状已是） | 减少协程数量与状态分散 |
| 容量基线 | 注册无限重试建议加上限 + 退避，长期失败上报告警 | 避免永久空转无感知 |
| 拒绝策略 | 不适用 | — |
| 隔离 | 保持独占 goroutine；指标函数为 IO 密集（DB 查询），如需提速可按指标组并发，但需评估 SDK 线程安全性（待确认） | 串行在指标多时可能超 5min 周期 |
| 命名与观测 | 建议每轮上报耗时打点；Stop 未被调用，建议接入优雅退出回调 | 进程退出时 goroutine 直接消失，无泄漏风险但语义不完整 |
