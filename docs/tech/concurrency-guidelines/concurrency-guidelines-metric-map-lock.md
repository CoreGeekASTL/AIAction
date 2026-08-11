# metricMapLock（mocIdMap 保护锁）并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 锁 |

## 用途定位

保护 `MonitorServiceImpl.mocIdMap`（mocId → moiId 注册集合）：在 `addMoiIdIfNotExists` 中判重并向监控 SDK 注册新对象（`ObjChange`）。从临界区实现推断。

- 代码标识符：`MonitorServiceImpl.metricMapLock`
- 定义位置：src/service/monitor_service.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/service/monitor_service.go（`addMoiIdIfNotExists`）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 锁类型 | sync.RWMutex（实际仅使用写锁 `Lock()`，未使用 `RLock()`） | src/service/monitor_service.go |
| 锁粒度 | 保护整个 `mocIdMap`，每次 addMoiId 调用全程持锁，临界区内包含 SDK 同步调用 `ObjChange`（IO） | src/service/monitor_service.go |
| 锁顺序约定 | 未设置（仓内唯一持该锁处，无锁序问题） | src/service/monitor_service.go |
| 临界区说明 | map 判重 + `MonSdkInstance.ObjChange` 远程调用 + map 写入 | src/service/monitor_service.go |
| 池化方式 | 未设置 | src/service/monitor_service.go |
| 容量配置 | 未设置 | src/service/monitor_service.go |
| 任务队列 | 未设置 | src/service/monitor_service.go |
| 拒绝策略 | 未设置 | src/service/monitor_service.go |
| 隔离范围 | 保护对象：`mocIdMap`；当前唯一调用方是同一 goroutine 内的 `processMetricResults`，锁实际无竞争 | src/service/monitor_service.go |
| 关闭与等待 | 未设置 | src/service/monitor_service.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 不适用 | — |
| 容量基线 | 不适用 | — |
| 拒绝策略 | 不适用 | — |
| 隔离 | 建议将 SDK 调用 `ObjChange` 移出临界区（锁内只做 map 判重与占位，SDK 失败再回滚占位），缩短持锁时间 | 锁内做远程 IO 调用，一旦未来多 goroutine 调用将放大尾延迟 |
| 命名与观测 | 若未来指标并发化，该锁为必经点，建议加锁等待耗时观测 | 提前暴露锁竞争 |
