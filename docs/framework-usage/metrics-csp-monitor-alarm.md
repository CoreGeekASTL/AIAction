# CSP 话统监控与告警使用指导（监控/可观测）

> 版本：CSPGoMonitorSDK（stub）+ AlarmSDK_GO（stub）｜ 调用点：~10 ｜ 涉及文件：2 ｜ 基线：main @ 5e78a48

## 用途定位

对接 CSP 平台运营体系的两个 SDK：
- **话统监控**（`service/monitor_service.go`）：定时把业务指标（在线用户数、各机型在线、流量等）上报 CSP 运营面
- **告警**（`service/alarm_service.go`）：向 FM 上报告警/清除告警，启动时清理历史活动告警

## 初始化与配置

### 话统

- 启动：`main.go:90-91` `go monitorService.InitMonitorSchedule()`
- 配置：`conf/monitor.json`（指标模型，路径由 `cspmonitor::monitorJsonFile` 指定，`monitor_service.go:37,77`）；SQL 模板 `conf/sql.yaml`（`cspmonitor::sqlYamlFile`，`monitor_service.go:38,119`）
- 流程：`buildMonitorConfig`（读 monitor.json）→ `InitCspMonitor`（`MonSdkInstance.InitMonitor` + `RegisterBasicInfo`，失败 10s 重试至死）→ `startCspMonitor`（5min ticker 打点）（`monitor_service.go:55-136`）

### 告警

- `service/alarm_service.go:63-72` `init()`：初始化 `CSPAlarmManager`、注册 reset-clear、建 999 缓冲 channel、`go handleEvent()` 单消费者
- 启动时清历史告警：`main.go:93` `go service.CleanAllActiveAlarm()`（查 FM 活动告警，按 sourceip 匹配本节点后清除，`alarm_service.go:191-215`）

## 核心使用模式

```go
// 话统：指标函数映射 + 周期上报（monitor_service.go:145-177, 196-211）
m.metricFunMap = map[monitor.MetricID]getMetricFunc{
	monitor.MetricOnlineUsers: m.statsService.GetOnline,
	...
}
// 每 5min：对每个 metric 执行函数 → 逐对象 addMoiIdIfNotExists(ObjChange 注册) → SetMetric 上报
monitorsdk.MonSdkInstance.SetMetric(int(metric.ID), item.Obj, float64(item.Cnt))
```

```go
// 告警：发/清都是异步入队（alarm_service.go:108-129）
a.sendAlarmEvent(AlarmEvent{AlarmID: alarmID, EventMessage: msg, Type: base.GenerateAlarm})
// 队列满时 5s 超时丢弃（alarm_service.go:108-115）
// 消费侧：10min 内同 AlarmID 抑制重发（alarm_service.go:90-97），失败重试 2 次间隔 10s（alarm_service.go:153-162）
```

## 封装层与扩展点

- 话统扩展点：`monitor.json` 增指标 + `createMetricFunctionMap` 加映射 + `TrafficStatsService` 加统计函数 + `sql.yaml` 加 SQL——四件套联动（`monitor_service.go:145-153`、`traffic_stats_service.go:27-32,80-100`）。
- SQL 外置：`sql.yaml` 按名字管理 SQL 与参数（`traffic_stats_service.go:80-100 LoadSQLConfig`），改统计口径不改代码。
- 告警 ID 集中登记：`AlarmList = []string{AlarmId300010}`（`alarm_service.go:189`），新告警 ID 需加入以便启动清理。

## 并发与线程模型

- 话统单 goroutine 5min ticker（`monitor_service.go:123-135`）；`mocIdMap` 用 `sync.RWMutex`（实际是写锁）保护（`monitor_service.go:213-228`）。
- 告警"单生产者多入队、单消费者"模型；告警指标函数 recover 模板在 `monitor_service.go:181-193`。

## 错误处理与容错

- 话统 SDK 初始化失败 10s 无限重试（`monitor_service.go:64-73`）；单个指标函数 panic 被 recover 不影响其他指标（`monitor_service.go:181-193`）。
- 告警发送失败重试 2 次后放弃（`alarm_service.go:153-165`）；FM 查询重试 360 次×5s 应对升级场景（`alarm_service.go:195-202`）。
- 告警/话统失败均不影响主业务（全部异步）。

## 约定与规范

- 上报窗口：`monitorutil.GetLastFiveMinuteWindow`（`monitor_service.go:190`），与 5min 打点周期对齐。
- 告警参数固定集合：kind/namespace/sourceip/EventMessage/EventSource/OriginalEventTime（`alarm_service.go:146-152`）。

## 已知问题与反模式

- 两个 SDK 本地均为 stub，话统/告警逻辑本地不可验证。
- `alarm_service.go` 大量 Errorf 级调试日志、`log.Println` 混用（见 log md）。
- `GetAllActiveAlarmFromFMService` 用全局 `mutex` 锁 + 无缓冲 channel 超时控制（`alarm_service.go:228-269`），goroutine 泄漏风险（超时后子 goroutine 仍阻塞在 `ch <- "success"`，`alarm_service.go:253`——无缓冲 channel 无人接收）。
- 话统 mocIdMap 用 `metricMapLock.Lock()`（写锁）包全部读判（`monitor_service.go:213-228`），RWMutex 名不副实但正确。

## AI 编码指南

- 新增运营指标：① `conf/monitor.json` 加 metric 定义 ② `models/monitor/metric.go` 加 MetricID 常量 ③ `TrafficStatsService` 加统计函数 + `sql.yaml` 加 SQL ④ `createMetricFunctionMap` 加映射。依据：上文「封装层与扩展点」。
- 新增告警：复用 `service.NewAlarmService().SendAlarm/ClearAlarm(alarmID, msg)`（异步入队即可，勿直接调 `reportAlarm`）；新告警 ID 加入 `AlarmList`（`alarm_service.go:189`）。依据：`alarm_service.go:108-129`。
- **禁止**在主业务流程同步等待话统/告警结果；**禁止**用无缓冲 channel 做超时控制（`alarm_service.go:245-266` 为反面教材，改用缓冲 channel 或 context）。依据：上文「已知问题与反模式」。
