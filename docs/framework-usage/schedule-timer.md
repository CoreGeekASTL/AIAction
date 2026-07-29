# 定时调度使用指导（定时/调度）

> 版本：Go 标准库 time.Timer/Ticker（无 cron 库）｜ 调用点：6 ｜ 涉及文件：5 ｜ 基线：origin/main @ ae0a8a6（2026-07-29 复核）

## 用途定位

未引入任何 cron 框架，全部用标准库原语实现。仓内共 5 个调度场景：

| 场景 | 机制 | 周期 | 位置 |
| --- | --- | --- | --- |
| 数据清理（流量统计过期数据） | 自研 `DataCleanupScheduler`（Timer + 状态机） | 每天凌晨 2 点 | `scheduler/task_scheduler.go` |
| 配置中心缓存刷新 | Ticker | 5 min | `service/config_center_service.go:104-120` |
| CSP 话统上报 | Ticker | 5 min | `service/monitor_service.go:117-136` |
| GaussDB 健康检查/主备切换 | Ticker | 5 s | `dao/db_init.go:265-280` |
| 事件日志文件例行清理 | `time.Tick` | 1 h | `common/event/local_storage.go:77-84` |

## 初始化与配置

- 数据清理：`main.go:87 scheduler.StartDataCleanupScheduler()`；优雅退出经 GSF `GracefulExitHandler.Exit` → `StopDataCleanupScheduler()`（`main.go:216-219`）。
- 清理保留月数：`constants.CleanupMonths`（`task_scheduler.go:129-132`）。

## 核心使用模式

### 完整生命周期调度器（新调度任务的模板）

```go
// 来源：src/scheduler/task_scheduler.go:18-74 —— stopChan + WaitGroup + 幂等 Start/Stop
type DataCleanupScheduler struct {
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
	isRunning bool
	mu        sync.Mutex
}

func (s *DataCleanupScheduler) Start() {
	s.mu.Lock(); defer s.mu.Unlock()
	if s.isRunning { return }
	s.isRunning = true
	s.waitGroup.Add(1)
	go s.run()
}

func (s *DataCleanupScheduler) run() {
	defer s.waitGroup.Done()
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		timer.Reset(sleepDuration)
		select {
		case <-timer.C:  s.executeCleanup(maxRetries, retryInterval)
		case <-s.stopChan: return
		}
	}
}
```

### 轻量周期任务（无需停止语义的）

```go
// 来源：src/service/config_center_service.go:106-119
go func() {
	ticker := time.NewTicker(RefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C: service.Refresh()
		case <-configCenter.stopChan: return
		}
	}
}()
```

## 封装层与扩展点

- 唯一封装：`scheduler.DataCleanupScheduler` 包级单例 `globalScheduler`（`task_scheduler.go:25-29` init 创建），对外只有 `Start/StopDataCleanupScheduler` 两个函数。
- 其余调度都是"就地 `go func + ticker`"，无封装。

## 并发与线程模型

- 调度 goroutine 与业务执行同线程（`task_scheduler.go:100-102` 在 select 分支里直接执行，持 `s.mu` 锁）——执行期间 Stop 会等 `waitGroup.Wait()`（`task_scheduler.go:61-74`）。
- 重试间隔等待也响应停止：`sleepWithStopCheck`（`task_scheduler.go:151-161`）。
- **新调度任务若需优雅停止，必须复制"stopChan + WaitGroup + 幂等 Start/Stop"四件套。**

## 错误处理与容错

- 清理任务失败重试 3 次、间隔 10min，全失败等下一周期（`task_scheduler.go:80-83,128-148`）。
- `time.After` 用于单次超时（`alarm_service.go:112`、3s 收告警应答 `alarm_service.go:263`）。

## 约定与规范

- 计算下次定点执行用 `time.Date` 构造目标时刻（`task_scheduler.go:114-125`），不用 cron 表达式。
- `time.Tick`（不可 Stop）只在"与进程同寿"的任务用（`local_storage.go:78`）；其余一律 `NewTicker + defer Stop`。

## 已知问题与反模式

- 凌晨 2 点、5min、5s 等周期全部硬编码（`task_scheduler.go:117`、`config_center_service.go:28`），不可配置。
- `db_init.go:266` 与 `config_center_service.go:108` 的 ticker 循环没有退出通道（与进程同寿），单测中无法停止。
- master election（多实例只跑一个调度）仅有 stub 测试（`service/master_election_service_stub_test.go`），生产实现不在仓内——当前多实例部署下 `DataCleanupScheduler` 会在每个实例各跑一份。

## AI 编码指南

- 新增周期任务：需要可停止/可测试 → 复制 `DataCleanupScheduler` 四件套；与进程同寿的轻量任务 → `go func + NewTicker + defer ticker.Stop()`。依据：上文「核心使用模式」。
- 任务体内部允许中断的等待一律 `select { case <-timer.C / case <-stopChan }`，**禁止**裸 `time.Sleep`（`alarm_service.go:159,199` 为反面教材，停止时最多多睡 10s/5s）。依据：`task_scheduler.go:151-161`。
- **禁止**引入 robfig/cron 等第三方调度库——与存量风格保持一致。依据：本仓未使用任何 cron 库（go.mod 无此依赖）。
