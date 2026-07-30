# Go 协程原语使用指导（并发/线程池）

> 版本：Go 标准库（goroutine/sync/channel/time）｜ 调用点：17+ ｜ 涉及文件：10+ ｜ 基线：main (6c93561)

## 用途定位
项目无第三方线程池/并发框架，全部并发直接用 goroutine + channel + sync 原语。典型场景：后台初始化、定时任务、并发批量调用、串行化事件处理。

## 初始化与配置
无集中初始化；并发结构分散在各组件内，生命周期由 stopChan + WaitGroup 管理（范式见 `src/scheduler/task_scheduler.go:18-74`）。

## 核心使用模式

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

## 封装层与扩展点
无封装层，直接用标准库。约定见下。

## 并发与线程模型
- HTTP handler 已在 Beego 请求 goroutine 中，同步处理即可；耗时操作自行 `go` 出去（如 PreOpenBrowser）。
- 共享 map 必须锁或 channel 串行化；`alarmServiceImpl.alarms` 选择后者（单消费 goroutine）。

## 错误处理与容错
- goroutine 内错误只能打日志，无法上抛（存量统一做法，如 `src/service/browser_service.go:94-103`）。
- 关键 goroutine 用 recover 兜底：监控指标计算 `defer recover()` + `debug.Stack()`（`src/service/monitor_service.go:183-191`）。

## 约定与规范
- 长驻 goroutine 必须可停止：stopChan + `select` + WaitGroup；启动期一次性 goroutine 可直接 `go`。
- 定时循环必须 `defer ticker.Stop()`（`src/service/monitor_service.go:125`）。

## 已知问题与反模式
- `src/common/event/local_storage.go:78` 用 `time.Tick`（不可 Stop）且无退出通道，goroutine 泄漏风险（进程生命周期内可接受）。
- `src/service/alarm_service.go:77` channel 关闭后 `break` 仅跳出 select，空转 for。

## AI 编码指南
- 新增长驻后台任务：照「模式一」写 stopChan+WaitGroup+mu 幂等 Start/Stop；**禁止**裸 `go func(){ for { ... } }()` 无退出通道（依据：上文「核心使用模式」与「已知问题」）。
- 新增共享状态：优先 channel 串行化或 `sync.Map`；普通 map 必须配 RWMutex（依据：`src/common/cse/cse.go:34`、`src/service/monitor_service.go:45`）。
- goroutine 内执行可能 panic 的插件式逻辑时，必须 `defer recover()`（依据：`src/service/monitor_service.go:183`）。
