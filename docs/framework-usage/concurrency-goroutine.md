# Go 协程与单例模式使用指导（并发/线程池）

> 版本：Go 标准库 sync 包 ｜ 调用点：17+ ｜ 涉及文件：10+ ｜ 基线：origin/main @ ae0a8a6（2026-07-29 复核）

## 用途定位

无任何线程池/协程池框架，直接用 goroutine + sync 原语。同时承担"依赖注入"角色——**包级变量单例**是本项目唯一的组件装配方式（无 DI 框架）。

## 初始化与配置

单例模式三种变体（全部有实例佐证）：

```go
// 变体一：init() 直接构造包级单例（最常用）
// 来源：src/service/config_center_service.go:89-101
var configCenter *configCenterServiceImpl
func init() {
	configCenter = &configCenterServiceImpl{dao: dao.NewConfigCenterDao()}
	configCenter.stopChan = make(chan struct{})
}
func NewConfigCenterService() ConfigCenterService { return configCenter }
```

```go
// 变体二：sync.Once 延迟初始化
// 来源：src/service/event_service.go:19-46
var once = sync.Once{}
func NewEventService() *EventServiceImpl {
	once.Do(initEventStorageFactory)
	...
}
```

```go
// 变体三：每次新建实例（无共享状态时）
// 来源：src/service/user_service.go:30-35
func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{ud: dao.NewUserDaoDao(), ubd: dao.NewUserBindDao()}
}
```

## 核心使用模式

- 接口 + 小写实现 + 编译期断言：`type XxxService interface` + `xxxServiceImpl` + `var _ XxxService = &XxxServiceImpl{}`（`traffic_stats_service.go:35-58`、`user_service.go:23-42`）。
- 异步任务：`go dao.EnsureConnectGaussDB()`、`go monitorService.InitMonitorSchedule()`（`main.go:67,91`）——启动期的长循环全部 `go` 出去。
- 共享缓存：`sync.Map`（`cse.go:34` browserGWInstances）、`sync.RWMutex`（`monitor_service.go:45`、`alarm_service.go:228`）。
- channel 队列：告警事件 `make(chan AlarmEvent, 999)` 缓冲队列 + 单消费者 goroutine（`alarm_service.go:54,69-88`）。

## 封装层与扩展点

无封装层——这就是本项目的并发"框架"。测试替换点：`service/traffic_stats_service_mock.go`、`service/master_election_service_stub_test.go` 用手写 mock/stub 实现同一接口。

## 并发与线程模型

- Beego 请求 goroutine → Controller `Prepare()` 新建 service → 包级单例共享数据需自备同步。
- 后台 goroutine 分两类：与进程同寿（DB 健康检查 `db_init.go:265`、config 刷新、watch 回调） vs 可停止（DataCleanupScheduler，见 schedule-timer.md）。
- 发送告警带 5s 超时防阻塞：`select { case ch<-event / case <-time.After(5s) }`（`alarm_service.go:108-115`）。

## 错误处理与容错

goroutine 内 panic 无统一 recover（仅 `monitor_service.go:181-193` 对指标函数做了 recover + `debug.Stack()`）。**新增长驻 goroutine 内有 panic 风险的操作时，参考 monitor_service 的 recover 模板。**

## 约定与规范

- AGENTS.md 基线：单例初始化用 `sync.Once` 保护（搜索 `once.Do`）；`context.TODO()` 不传 nil。
- 命名：`NewXxxService()` 返回接口或实现指针，小写 impl 结构体。

## 已知问题与反模式

- `alarmService` 包级单例的 `alarms map` 无锁（`alarm_service.go:48-53,90-106`）——靠"单消费者 goroutine"保证安全，但 `CleanAllActiveAlarm` 路径（`alarm_service.go:217-226`）从另一 goroutine 调 `reportAlarm`，绕过了队列，与 `handleEvent` 并发读写 map 存在竞态。
- `dao` 全局 `ormer` 切换无锁（见 storage-beego-orm.md）。
- `monitor_service.go:6-16` import 分组混乱（标准库夹在第三方中间）。

## AI 编码指南

- 新增 Service：`type XxxService interface` + 小写 `xxxServiceImpl` + `var _ XxxService = &XxxServiceImpl{}` + `NewXxxService()`；有共享状态 → init/sync.Once 包级单例，无共享状态 → 每次新建。依据：上文「核心使用模式」。
- 共享可变状态必须选一种同步手段：`sync.Mutex/RWMutex`、`sync.Map`、或 channel 单消费者；**禁止**裸 map 跨 goroutine 读写（`alarm_service.go` 为反面教材）。依据：上文「已知问题与反模式」。
- 长驻 goroutine 必须可停止（stopChan）或与进程同寿二选一，并在函数体内 recover（模板：`monitor_service.go:181-193`）。依据：上文「错误处理与容错」。
