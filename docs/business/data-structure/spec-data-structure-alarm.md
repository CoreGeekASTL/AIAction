# 告警

> 用途：alarm　实例数：3　返回 [README.md](README.md)

## 1. 核心作用

承载告警异步处理与抑制计数——告警事件以 channel 队列异步投递，抑制计数 map 判定 10min 内重复告警，全局告警 ID 以枚举列表承载，是告警链路的并发编排与状态中枢。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| alarmService.alarms | 告警抑制计数（AlarmID → 上次发送时间戳） | service/alarm_service.go |
| alarmEventChanel | 告警事件队列（异步处理告警生成/清除） | service/alarm_service.go |
| AlarmList | 全局告警 ID 枚举列表 | service/alarm_service.go |

## 3. 实例详解

- **alarmService.alarms**
  - 结构：`type alarmServiceImpl struct { alarms map[string]int64; alarmManager base.CSPAlarmManager }`（service/alarm_service.go:48）
  - 关键字段：alarms（key=AlarmID，value=上次发送 UnixMilli 时间戳）
  - 典型操作：sendAlarm 读 lastSendTime 判断 10min 内抑制；发送成功更新时间戳
  - 使用点：service/alarm_service.go:91（读抑制判断）、:94（抑制比较）、:69（init 初始化 map+channel+启消费 goroutine）
  - 并发模型：由 `handleEvent` 单 goroutine 消费 alarmEventChanel 时调用 sendAlarm/clearAlarm，alarms 访问限该 goroutine 内，当前无锁安全；跨 goroutine 直调须加锁
- **alarmEventChanel**
  - 结构：`var alarmEventChanel chan AlarmEvent`（service/alarm_service.go:54），缓冲 `maxAlarmListLen=999`
  - 关键操作：init 阶段 `make(chan AlarmEvent, 999)` + 启 `handleEvent` 消费 goroutine；SendAlarm/ClearAlarm 向 channel 投递事件；handleEvent 按 Type 分发 sendAlarm/clearAlarm
  - 使用点：service/alarm_service.go:69（创建+启消费者）、:77（消费 select）、投递点散布于告警上报处
  - 并发模型：多生产者 + 单消费者 goroutine，缓冲 999 背压
- **AlarmList**
  - 结构：`var AlarmList = []string{AlarmId300010}`（service/alarm_service.go:189），包级字面量
  - 关键操作：定义需订阅的告警 ID 集合，告警初始化时遍历注册
  - 使用点：service/alarm_service.go:189（字面量定义）、告警注册处遍历
  - 并发模型：字面量初始化后只读，无锁

## 4. 使用模式与约定

- 鉴权/告警类缓存统一"struct 内嵌 RWMutex + map 字段"封装，get/set/clear 各自加对应锁（authCache 范式）
- 异步事件处理统一 buffered chan + 单消费者 goroutine（alarmEventChanel 范式）
- 全局枚举列表用包级 `var x = []T{...}` 字面量，运行期只读（AlarmList 范式）
