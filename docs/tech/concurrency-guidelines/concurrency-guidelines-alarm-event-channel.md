# alarmEventChanel 告警事件通道 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | channel |

## 用途定位

告警事件异步上报管道：业务侧通过 `SendAlarm/ClearAlarm` 把 `AlarmEvent` 投递到有缓冲 channel，init 时启动的单 goroutine `handleEvent` 消费并调用告警 SDK 上报；`alarms` map 用于同 ID 告警 10 分钟抑制。从 init 与 `handleEvent` 实现推断。

- 代码标识符：`alarmEventChanel` / `maxAlarmListLen` / `alarmServiceImpl.handleEvent`
- 定义位置：src/service/alarm_service.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：生产：src/service/alarm_service.go（`sendAlarmEvent`、`SendAlarm`、`ClearAlarm`、`clearHistoryAlarm`）；消费：src/service/alarm_service.go（`handleEvent`，init 中 `go alarmService.handleEvent()`）；src/service/alarm_service.go（`CleanAllActiveAlarm`，main 中 `go service.CleanAllActiveAlarm()`）

## 线程模型（可选）

```mermaid
flowchart LR
    biz["业务调用方 SendAlarm/ClearAlarm"] --> ch["alarmEventChanel buffer=999"]
    ch --> consumer["handleEvent 单 goroutine"]
    consumer --> sdk["AlarmSDK 上报"]
```

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 池化方式 | 裸 goroutine（无池化）：init 中 `go alarmService.handleEvent()` 单消费者 | src/service/alarm_service.go |
| 容量配置 | channel buffer 硬编码 `maxAlarmListLen = 999` | src/service/alarm_service.go |
| 任务队列 | 有界 channel（buffer 999） | src/service/alarm_service.go |
| 拒绝策略 | 生产者侧 `sendAlarmEvent`：投递失败（满）时等待 5s 超时后记录 error 日志丢弃事件 | src/service/alarm_service.go |
| 隔离范围 | 独占通道，仅承载告警事件 | src/service/alarm_service.go |
| 关闭与等待 | 未设置（channel 不关闭，消费者 `for` 循环中的 `break` 仅跳出 select 不跳出 for，属死循环结构） | src/service/alarm_service.go |
| 消息缓冲与背压 | buffer 999，满时生产者阻塞最多 5s 后丢弃；`alarms` 抑制 map 仅消费者单 goroutine 访问，无锁 | src/service/alarm_service.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 单消费者通道模型合理，保持；告警 SDK 调用若为慢 IO 可考虑多消费者，但需先给 `alarms` map 加锁 | 当前无锁依赖单消费者语义 |
| 容量基线 | buffer 999 偏大且无监控，建议保留有界并上报水位指标 | 有界防内存膨胀，水位可观测 |
| 拒绝策略 | 超时丢弃可接受（告警允许丢失），建议丢弃时计数打点 | 丢弃目前只有日志，无量化 |
| 隔离 | 保持独占通道；注意 `SendAlarm`/`ClearAlarm` 公开接口与 `sendAlarmEvent` 路径需确认是否都走超时丢弃（待确认 `SendAlarm` 是否直接投递） | 统一背压语义 |
| 命名与观测 | 消费者 `handleEvent` 中 `break` 不退出 for 属疑似 bug，建议改为带标签 break 或 return；channel 若永不关闭应删除 ok 判断 | 避免误读与空转 |
