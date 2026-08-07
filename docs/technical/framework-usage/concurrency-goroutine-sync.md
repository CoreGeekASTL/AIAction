# Go 协程原语使用指导（并发/线程池）

## 用途定位
项目无第三方线程池/并发框架，全部并发直接用 goroutine + channel + sync 原语。典型场景：后台初始化、定时任务、并发批量调用、串行化事件处理。


## 使用模式

模式一：可停止的后台循环（来源：`src/scheduler/task_scheduler.go:46-110`）：

```go
type XxxScheduler struct {
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
	isRunning bool
	mu        sync.Mutex
}
func (s *XxxScheduler) Start() { // mu 保护 isRunning 幂等
	s.waitGroup.Add(1)
	go s.run()
}
func (s *XxxScheduler) run() {
	defer s.waitGroup.Done()
	for {
		select {
		case <-timer.C: /* 干活 */
		case <-s.stopChan:
			return
		}
	}
}
```

模式二：启动期 fire-and-forget 后台初始化（来源：`src/main.go:67,91,93`）：

```go
go dao.EnsureConnectGaussDB()
go monitorService.InitMonitorSchedule()
go service.CleanAllActiveAlarm()
```

模式三：扇出并发调用（来源：`src/service/browser_service.go:75-78`）：

```go
for _, instance := range instances {
	go instancePreOpenBrowser(instance, initBrowserRequest, b.httpClient) // 不等结果
}
```

模式四：channel 串行化事件队列（来源：`src/service/alarm_service.go:69-88`）：

```go
alarmEventChanel = make(chan AlarmEvent, maxAlarmListLen) // 容量 999
go alarmService.handleEvent() // 单 goroutine 消费，保护 alarms map 免锁
// 发送端带 5s 超时防阻塞：select { case ch <- e: case <-time.After(5s): }
```

模式五：并发安全缓存用 `sync.Map`（`src/common/cse/cse.go:34`）；普通 map 用 `sync.RWMutex`（`src/service/monitor_service.go:45`）。
