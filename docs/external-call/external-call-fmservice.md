# FMService 出站调用

FMService（故障管理服务）通过 CSE 服务发现访问，使用 go-chassis rest invoker 以 `cse://FMService/...` 形式发起调用。

## POST /fmOperation/v1/alarms/get_alarms

- 协议：CSE REST POST `cse://FMService/fmOperation/v1/alarms/get_alarms`（底层走 go-chassis `core.NewRestInvoker().ContextDo`），请求体为 `{"cmd":"GET_ACTIVE_ALARMS","language":"en-us","data":{"appId":"...","alarmIds":"..."}}` JSON
- 调用位置：service/alarm_service.go（GetAllActiveAlarmFromFMService / handlerActivityAlarmData / OSHttpsGetRequestByCSE 函数）
- 业务场景：服务启动/升级后执行 `CleanAllActiveAlarm`，从 FMService 查询本服务（按 appId + alarmIds，当前仅 300010）的历史活动告警，用于清理升级前遗留的本节点告警
- 接口功能：查询活动告警列表，返回 `AlarmResponse`（retcode/data，data 含 alarmId 与 location 参数串）；按 `sourceip` 匹配本节点 IP 后对命中告警执行清除。内置 3 秒超时与最多 360 次（间隔 5s）重试，以应对升级场景 FMService 未就绪
