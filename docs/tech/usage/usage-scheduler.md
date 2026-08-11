# Go 协程原语 + 自研调度器 使用现状（并发/线程池、定时/调度）

## 用途定位
无线程池框架，并发全部用 Go 原生原语（goroutine + channel + sync.Mutex/WaitGroup/sync.Map/sync.Once），无第三方 cron——定时调度用 `time.Timer`/`time.Ticker` 自研。两类形态：
1. **独立 goroutine 后台任务**：`main.go` 直接 `go dao.EnsureConnectGaussDB()`、`go monitorService.InitMonitorSchedule()`、`go service.CleanAllActiveAlarm()`；server 证书监听 `go b.monitorCertificate()`；事件落盘 `src/common/event/local_storage.go` 的 `go func()`。
2. **自研可停调度器**：`src/scheduler/task_scheduler.go` 的 `DataCleanupScheduler`（每日凌晨 2 点清理旧流量数据）是标准范式——`stopChan` + `sync.WaitGroup` + `sync.Mutex` 保护启停状态，`time.NewTimer` + `calculateNextRunTime` 实现定点触发，任务失败重试 3 次（间隔 10min，重试 sleep 也可被 stop 打断）。
3. **ticker 周期任务**：DB 健康检查 5s（`src/dao/db_init.go`）、监控上报 5min（`src/service/monitor_service.go`）、配置刷新（`src/service/config_center_service.go`）。

单例服务初始化用 `sync.Once`（`src/service/event_service.go`、`src/dao/db_local_sqlite.go`）。

## 使用模式

可停调度器骨架：

```go
// 来源：src/scheduler/task_scheduler.go
type DataCleanupScheduler struct {
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
	isRunning bool
	mu        sync.Mutex
}

func (s *DataCleanupScheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isRunning {
		return
	}
	s.isRunning = true
	s.waitGroup.Add(1)
	go s.run()
}

func (s *DataCleanupScheduler) run() {
	defer s.waitGroup.Done()
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		nextRun := s.calculateNextRunTime()
		timer.Reset(nextRun.Sub(time.Now()))
		select {
		case <-timer.C:
			s.executeCleanup(maxRetries, retryInterval)
		case <-s.stopChan:
			return
		}
	}
}
```

并发缓存用 `sync.Map`（`src/common/cse/cse.go` 的 browserGWInstances）；优雅退出由 `GracefulExitHandler.Exit()` 调 `StopDataCleanupScheduler`（`src/main.go`）。
