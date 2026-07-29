# CSP 告警平台出站调用

告警通过 `AlarmSDK_GO`（封装在 service/alarm_service.go 的 AlarmService）上报到 CSP 告警平台，SDK 内部通道出站（具体下游地址由 SDK/平台注入，待确认）。

| 接口名 | 协议 | 调用位置 | 业务场景 |
|--------|------|----------|----------|
| 告警上报（SendAlarm） | AlarmSDK_GO SDK | service/alarm_service.go | 故障告警上报 |
| 告警清除（ClearAlarm） | AlarmSDK_GO SDK | service/alarm_service.go | 告警恢复清除 |

## 告警上报 SendAlarm

- 协议：AlarmSDK_GO SDK `CSPAlarmManager.SendAlarm`（`manager.CSPInitAlarmSDK` 初始化，参数 appId/serviceName/nodeName/nodeIP）
- 调用位置：service/alarm_service.go（SendAlarm → reportAlarm）；业务调用方 controllers/management_controller.go（配置同步失败上报 300010）等
- 业务场景：业务流程出现需要运维介入的故障（如浏览器配置同步失败）时上报告警
- 接口功能：构造告警（alarmId + kind/namespace/sourceip/EventMessage/EventSource/OriginalEventTime 参数）经 channel 异步上报；10 分钟内同 ID 告警抑制，失败重试 1 次（间隔 10s）

## 告警清除 ClearAlarm

- 协议：AlarmSDK_GO SDK `CSPAlarmManager.SendAlarm`（Type=ClearAlarm）
- 调用位置：service/alarm_service.go（ClearAlarm → reportAlarm；clearHistoryAlarm 用于历史告警清除）；业务调用方 controllers/management_controller.go（配置同步成功恢复 300010）、main.go（CleanAllActiveAlarm）
- 业务场景：故障恢复或服务升级后清除对应告警
- 接口功能：按 alarmId 清除活动告警，成功后从本地告警台账删除
