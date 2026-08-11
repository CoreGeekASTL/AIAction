# main 启动 goroutine 组 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 裸并发(goroutine) |

## 用途定位

main 启动期拉起的常驻/一次性后台 goroutine 集合（无池化、固定数量）：

| goroutine | 用途 | 生命周期 |
|---|---|---|
| `gsfStartHandler` | CSP GSF 框架启动（podName 为空时阻塞，否则完成后写 `done`） | 常驻/一次性（视 podName） |
| `dao.EnsureConnectGaussDB` | 异步建连数据库（含 `go checkDBStatus` 健康检查） | 常驻 |
| `monitorService.InitMonitorSchedule` | 监控注册+上报（见 concurrency-guidelines-monitor-service-schedule.md） | 常驻 |
| `service.CleanAllActiveAlarm` | 启动时清理本节点历史活动告警（重试 360 次×5s） | 一次性 |

main 以 `<-done` 阻塞主 goroutine。

- 代码标识符：`done chan bool` / `gsfStartHandler` / `EnsureConnectGaussDB` / `CleanAllActiveAlarm`
- 定义位置：src/main.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/main.go（`registerInstance`、`main`）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 池化方式 | 裸 goroutine（无池化），固定 4 个启动点 | src/main.go |
| 容量配置 | 未设置（数量固定）；`CleanAllActiveAlarm` 重试参数 `RetryTimes=360`、`TimePeriodClean=5s` 硬编码于 service 包 | src/main.go; src/service/alarm_service.go |
| 任务队列 | `done` 无缓冲 channel，仅 `gsfStartHandler`（podName 非空分支）写入、main 读取一次 | src/main.go |
| 拒绝策略 | 未设置 | src/main.go |
| 隔离范围 | 各 goroutine 职责独立，无共享可变状态经 channel/锁同步（`done` 仅作完成信号） | src/main.go |
| 关闭与等待 | 未设置：main 依赖 `<-done` 永久阻塞；podName 为空时 `gsfapi.CspStart()` 阻塞于该 goroutine，`done` 永不写入 | src/main.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 启动期固定数量裸 goroutine 合理，保持 | 数量有界且职责清晰 |
| 容量基线 | `CleanAllActiveAlarm` 最长 30min 重试在启动路径执行，建议确认是否阻塞后续启动步骤（现状 main 中在其后仍有 Report 等同步调用，待确认影响） | 启动时延影响 |
| 拒绝策略 | 不适用 | — |
| 隔离 | `done` 语义双态（阻塞型/信号型）建议注释说明两种部署形态 | 防误改导致主 goroutine 提前退出 |
| 命名与观测 | 优雅退出目前只停 DataCleanupScheduler，建议其余常驻 goroutine（config 刷新、monitor、健康检查）纳入统一退出编排 | 退出语义完整化 |
