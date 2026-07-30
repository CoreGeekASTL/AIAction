# external-call-fmservice

> 下游服务：FMService（故障管理微服务，CSE 注册名 `FMService`）。
> 调用方式：go-chassis CSE rest invoker——`rest.NewRequest` 构造 `cse://FMService/...` URL，`core.NewRestInvoker().ContextDo()` 发起（service/alarm_service.go `OSHttpsGetRequestByCSE`）。

| 接口名 | 协议 | 调用位置 | 业务场景一句话 |
|--------|------|----------|----------------|
| POST /fmOperation/v1/alarms/get_alarms | HTTP（CSE） | service/alarm_service.go | 启动时查询本服务活动告警以便清除历史告警 |

## HTTP

## POST /fmOperation/v1/alarms/get_alarms

- 协议：HTTP POST `cse://FMService/fmOperation/v1/alarms/get_alarms`（经 ServiceComb 服务发现路由）
- 调用位置：service/alarm_service.go（`GetAllActiveAlarmFromFMService` → `handlerActivityAlarmData` → `OSHttpsGetRequestByCSE`）
- 业务场景：服务启动时（main.go `go service.CleanAllActiveAlarm()`）查询本服务在 FM 上的活动告警，匹配本节点 IP（sourceip）的历史告警并逐条清除，解决升级重启后旧告警残留问题；FM 未就绪时最多重试 360 次（间隔 5s，约半小时）
- 接口功能：请求体 JSON `{"cmd":"GET_ACTIVE_ALARMS","language":"en-us","data":{"appId":<APPID>,"alarmIds":"300010"}}`；响应 `AlarmResponse`（retcode=0 为正常，data 为活动告警列表，含 alarmId/location/appendInfo），从 location 中解析 `key=value` 形式的告警参数
