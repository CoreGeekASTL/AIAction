# queue（队列 / channel）

> 类型：queue　实例数：4　返回 [README.md](README.md)

## 1. 定位

Go channel 是 GIDS 的"队列"语义载体——告警异步处理、插件加载进度同步、goroutine 生命周期控制（停止/证书重启）全部走 channel，是并发编排的中枢。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| alarmEventChanel | 告警事件队列（异步处理告警生成/清除） | service/alarm_service.go |
| progressChan | 插件加载进度同步队列 | service/plugin_service.go |
| stopChan | goroutine 停止信号（监控/配置/调度/HTTPS） | service/monitor_service.go 等 |
| restartChan | HTTPS 证书重启信号 | common/https/https_server.go |

## 3. 实例详解

- **alarmEventChanel**
  - 结构：`var alarmEventChanel chan AlarmEvent`（service/alarm_service.go:54），缓冲 `maxAlarmListLen=999`
  - 关键操作：init 阶段 `make(chan AlarmEvent, 999)` + 启 `handleEvent` 消费 goroutine；SendAlarm/ClearAlarm 向 channel 投递事件；handleEvent 按 Type 分发 sendAlarm/clearAlarm
  - 使用点：service/alarm_service.go:69（创建+启消费者）、:77（消费 select）、投递点散布于告警上报处
  - 并发模型：多生产者 + 单消费者 goroutine，缓冲 999 背压
- **progressChan**
  - 结构：`var progressChan = make(chan db.PluginPackage, len(browserGWs))`（service/plugin_service.go:76）
  - 关键操作：并发加载插件，每个 worker 完成后向 channel 投递进度，主流程 `recordLoadPluginProgress`（:298）消费汇总
  - 使用点：service/plugin_service.go:76（创建，容量=worker 数）、:298（消费）、:313（签名传递）
  - 并发模型：多 worker 生产 + 单主流程消费，容量=worker 数实现"等全部完成"
- **stopChan**
  - 结构：`stopChan chan struct{}`，多处定义：service/monitor_service.go:51、service/config_center_service.go:33、scheduler/task_scheduler.go:19、common/https/https_server.go:37
  - 关键操作：`close(stopChan)` 触发消费 goroutine 的 `case <-stopChan: return` 退出
  - 使用点：service/monitor_service.go:122/130/140、service/config_center_service.go:95/114/123、scheduler/task_scheduler.go:34
  - 并发模型：close 广播退出信号，多消费 goroutine 同时感知
- **restartChan**
  - 结构：`restartChan chan CertInfo`（common/https/https_server.go:36），缓冲 1
  - 关键操作：证书变更时投递 CertInfo，HTTPS server goroutine select 命中后重启监听
  - 使用点：common/https/https_server.go:24（创建）、:36（字段）、消费 select
  - 并发模型：单生产者 + 单消费者，缓冲 1 防阻塞

## 4. 使用模式与约定

- 异步事件处理统一 buffered chan + 单消费者 goroutine（alarmEventChanel 范式）
- goroutine 生命周期统一 `chan struct{}` + close 广播退出（stopChan 范式，四处一致）
- 等待多并发结果用 buffered chan，容量=worker 数（progressChan 范式）

## 5. AI 编码指南

1. 新增异步事件处理用 buffered chan + 单消费者 goroutine，禁止裸 goroutine 无队列背压（依据：service/alarm_service.go:69，999 缓冲兜底）
2. 新增后台 goroutine 必须配 `chan struct{}` 停止信号并在 Stop() 中 close，禁止无退出机制的常驻 goroutine（依据：service/monitor_service.go:51/138，仓内统一范式）
3. 等待多并发结果用 buffered chan 容量=worker 数 + 单消费汇总，禁止 WaitGroup 之外的裸 sleep 轮询（依据：service/plugin_service.go:76）

## 6. 风险与注意点

- **alarmEventChanel 缓冲满则阻塞生产者**：service/alarm_service.go:69（缓冲 999，超过则 SendAlarm 阻塞调用方，告警风暴时可能卡住业务线程）
