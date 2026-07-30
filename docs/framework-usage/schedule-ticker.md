# 定时/调度框架使用指导（定时/调度）

> 版本：Go 标准库 time.Timer/Ticker（无 cron 框架）｜ 调用点：5 ｜ 涉及文件：scheduler/task_scheduler.go、dao/db_init.go、service/monitor_service.go、service/config_center_service.go、common/event/local_storage.go ｜ 基线：main (6c93561)

## 用途定位
项目无 Quartz/cron 库，定时能力两种形态：
1. 自研 `DataCleanupScheduler`：定点（每日凌晨 2 点）任务调度器（`src/scheduler/task_scheduler.go`）。
2. 散落 `go func()+Ticker` 周期任务：DB 健康检查（5s）、话统上报（5min）、配置中心刷新（5min）、事件文件清理（1h）。

## 初始化与配置
- `scheduler.StartDataCleanupScheduler()` 于 `src/main.go:87`；退出回调中 `StopDataCleanupScheduler()`（`src/main.go:218`）。
- 包 init() 建全局单例 `globalScheduler`（`src/scheduler/task_scheduler.go:25-29`）。
- 散落 ticker 任务随各自 service 初始化启动（`src/main.go:72,91`）。

## 核心使用模式

定点调度器骨架（来源：`src/scheduler/task_scheduler.go:77-125`）：

```go
func (s *DataCleanupScheduler) run() {
	defer s.waitGroup.Done()
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		nextRun := s.calculateNextRunTime() // 算下一个凌晨 2 点
		timer.Reset(nextRun.Sub(time.Now()))
		select {
		case <-timer.C:
			s.mu.Lock()
			success := s.executeCleanup(maxRetries, retryInterval) // 失败重试 3 次×10min
			s.mu.Unlock()
		case <-s.stopChan:
			return
		}
	}
}
```

周期 ticker 骨架（来源：`src/service/config_center_service.go:104-120`）：

```go
go func() {
	ticker := time.NewTicker(RefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			service.Refresh()
		case <-configCenter.stopChan:
			return
		}
	}
}()
```

## 封装层与扩展点
- 定点任务封装：`scheduler.DataCleanupScheduler`（Start/Stop 幂等、重试、优雅停止）。
- 周期任务无封装，直接各 service 内联。
- 优雅退出链路：GSF `RegistExitHandler` → `scheduler.StopDataCleanupScheduler()`（`src/main.go:216-219`）。

## 并发与线程模型
- 调度器 Start/Stop 用 mu + isRunning 保证幂等（`src/scheduler/task_scheduler.go:46-74`）。
- 任务执行持 `s.mu`，避免与 Stop 竞态；重试间隔用 `sleepWithStopCheck` 保持可中断（`src/scheduler/task_scheduler.go:151-161`）。

## 错误处理与容错
- 任务失败重试 3 次，间隔 10min，重试中可被 Stop 中断（`src/scheduler/task_scheduler.go:128-148`）。
- ticker 任务失败仅日志，下个周期再试（`src/service/config_center_service.go:52-55`）。

## 约定与规范
- 定时周期定义为常量（`RefreshInterval`、`DotPeriodFiveMin`、`interval`），禁止魔法数字散落。
- 所有定时 goroutine 必须 `defer ticker.Stop()` 且响应 stopChan。

## 已知问题与反模式
- `src/common/event/local_storage.go:78` 用 `time.Tick` 无法 Stop、无 stopChan，属泄漏型写法。
- `src/dao/db_init.go:266-280` 的 checkDBStatus ticker 无停止通道（进程级任务可接受）。

## AI 编码指南
- 新增定点任务：复用/仿照 `DataCleanupScheduler` 结构（stopChan+WaitGroup+mu+可中断重试），注册到 GSF 退出回调（依据：`src/scheduler/task_scheduler.go:18-74`）。
- 新增周期任务：照 ticker 骨架写，**必须**带 stopChan 退出分支与 `defer ticker.Stop()`；**禁止** `time.Tick`（依据：上文「已知问题」，`src/common/event/local_storage.go:78`）。
- 定时周期定义为包级常量并注明单位（依据：`src/service/config_center_service.go:28`）。
