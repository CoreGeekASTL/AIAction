# 定时/调度框架使用指导（定时/调度）

## 用途定位
项目无 Quartz/cron 库，定时能力两种形态：
1. 自研 `DataCleanupScheduler`：定点（每日凌晨 2 点）任务调度器（`src/scheduler/task_scheduler.go`）。
2. 散落 `go func()+Ticker` 周期任务：DB 健康检查（5s）、话统上报（5min）、配置中心刷新（5min）、事件文件清理（1h）。


## 使用模式

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
