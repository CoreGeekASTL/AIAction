# FMService 通信规范

## HTTP

### POST /fmOperation/v1/alarms/get_alarms

- 业务场景：服务启动/升级后的告警清理流程（CleanAllActiveAlarm），从 FM 服务查询本节点历史活动告警，匹配 sourceip 后逐一清除，解决升级场景 FM 未启动导致告警不清除问题
- 接口功能：请求体 JSON `{"cmd":"GET_ACTIVE_ALARMS","language":"en-us","data":{"appId":...,"alarmIds":...}}`，返回 AlarmResponse（retcode "0" 为正常，data 含告警 location 参数列表）
- 调用位置：src/service/alarm_service.go（CleanAllActiveAlarm → GetAllActiveAlarmFromFMService → handlerActivityAlarmData → OSHttpsGetRequestByCSE）
- 协议信息：
  - 协议：HTTP POST，go-chassis rest 寻址 `cse://FMService/fmOperation/v1/alarms/get_alarms`（服务名显式硬编码 "FMService"）
  - 封装方式：go-chassis `rest.NewRequest` + `core.NewRestInvoker().ContextDo`（Go-chassis-extend，src/service/alarm_service.go 内 OSHttpsGetRequestByCSE 统一封装）
  - 超时重试：调用侧异步等待 3s 超时（`TimePeriodInit`）；CleanAllActiveAlarm 外层重试最多 360 次、间隔 5s（`RetryTimes`/`TimePeriodClean`，升级场景容忍 FM 晚启动约半小时）
  - 错误码处理：HTTP 状态码非 200 返回 error；响应 retcode 非 "0" 返回 error；error 触发外层重试，最终失败返回 false
