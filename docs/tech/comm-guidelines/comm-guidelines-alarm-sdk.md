# alarm-sdk 通信规范

## 平台 SDK

### SendAlarm（告警上报）

- 业务场景：业务运行中产生告警（如告警 ID 300010）时，向 CSP 告警平台上报告警事件
- 接口功能：`CSPAlarmManager.SendAlarm(alarm)`，告警对象由 `manager.InitCSPAlarm(alarmID, GenerateOrClearType)` 构造，附加 kind/namespace/sourceip/EventMessage/EventSource/OriginalEventTime 参数
- 调用位置：src/service/alarm_service.go（SendAlarm → sendAlarmEvent → sendAlarm → reportAlarm；CleanAllActiveAlarm → clearHistoryAlarm）
- 协议信息：
  - 协议：AlarmSDK_GO（CSPAlarmManager，SDK 内部通道；src/stubs/AlarmSDK_GO）
  - 封装方式：平台 SDK 直连（`AlarmSDK_GO/api/alarmapi`），初始化于 alarm_service.go init()
  - 超时重试：上报失败 sleep 10s 重试，最多 2 次（reportAlarm 内循环）；异步投递：告警事件先入 channel（容量 999），channel 满等待 5s 后放弃并记日志；同 alarmID 10 分钟内重复上报被抑制（alarmSuppressThresholdMs）
  - 错误码处理：SendAlarm 返回 bool，false 记日志；最终失败返回 false 由调用方记录

### ClearAlarm（告警清除）

- 业务场景：故障恢复或启动清理历史告警时，向告警平台清除指定告警
- 接口功能：同上 SDK 通道，事件类型为 ClearAlarm
- 调用位置：src/service/alarm_service.go（ClearAlarm → clearAlarm → reportAlarm）
- 协议信息：
  - 协议：AlarmSDK_GO（CSPAlarmManager）
  - 封装方式：平台 SDK 直连
  - 超时重试：同 SendAlarm（失败 sleep 10s 重试，最多 2 次）
  - 错误码处理：清除成功后从本地 alarms map 删除该告警；失败记日志
