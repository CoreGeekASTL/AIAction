# external-call-csp-alarm

> 下游服务：CSP 告警组件（平台 SDK `AlarmSDK_GO`，包 `AlarmSDK_GO/api/alarmapi`，源码桩在 src/stubs/AlarmSDK_GO/，真实 SDK 由平台提供）。
> 调用方式：进程内 SDK 调用（`manager.CSPInitAlarmSDK` 初始化得到 `CSPAlarmManager`，`SendAlarm` 上报）；SDK 内部如何送达告警平台由平台实现，本仓不可见。

| 接口名 | 协议 | 调用位置 | 业务场景一句话 |
|--------|------|----------|----------------|
| CSPAlarmManager.SendAlarm | 平台 SDK | service/alarm_service.go | 上报告警产生/恢复 |

## SDK

## CSPAlarmManager.SendAlarm（告警上报与清除）

- 协议：平台 SDK（`AlarmSDK_GO/api/alarmapi`：`CSPInitAlarmSDK` 初始化，`InitCSPAlarm(alarmId, GenerateOrClearType)` 构造告警对象，`AppendParameter` 追加参数，`SendAlarm` 发送）
- 调用位置：service/alarm_service.go（`init` 初始化 SDK 与事件通道；`reportAlarm` 统一上报入口；业务调用方：controllers/management_controller.go 同步浏览器配置失败时 `SendAlarm(300010)`、成功时 `ClearAlarm(300010)`；`CleanAllActiveAlarm` → `clearHistoryAlarm` 清除历史告警）
- 业务场景：业务异常（浏览器配置同步失败）时向平台告警系统上报告警 300010，恢复后清除；告警 10 分钟内抑重，发送失败重试 1 次（间隔 10s）
- 接口功能：告警参数含 kind=App、namespace、sourceip、EventMessage、EventSource、OriginalEventTime；`RegisterRsetClear` 注册复位清除回调；返回 bool 表示发送成败
