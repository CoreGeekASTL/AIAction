# CSP/GSF 平台套件使用指导（基础库）

## 用途定位
华为 CSP 平台强耦合的基础 SDK 集合，业务代码按其固定顺序初始化后使用：

| SDK | 用途 | 调用点 |
| --- | --- | --- |
| GSF（Go-chassis-extend api/GSF） | 框架引导 CspInit/CspStart、健康检查、退出回调、TLS 配置、审计日志发送 | main.go、common/logger/auditlog.go、common/https/client.go |
| CSPGSOMF TransportSDK/RunlogSDK/ModulekeeperSDK | 传输、运行日志、模块保活上报 | main.go:62,101-105 |
| CSPNTP_SDK_GO | NTP 对时 | main.go:64 |
| AlarmSDK_GO | 告警上报/清除 | service/alarm_service.go |
| fusionstage/auditlog | 事件日志文件引擎 | common/event/local_storage.go |
| greatwall-sdk-go | 限流（见 [resilience-greatwall.md](resilience-greatwall.md)） | controllers/filter.go |


## 使用模式
告警上报（来源：`src/service/alarm_service.go:118-166`）：

```go
// 业务接口：异步 channel 发送，单 goroutine 消费，10min 抑止重复
service.NewAlarmService().SendAlarm(alarmID, eventMessage)
service.NewAlarmService().ClearAlarm(alarmID, eventMessage)
// 消费端：manager.InitCSPAlarm(id, type) → AppendParameter(kind/namespace/sourceip/...) → SendAlarm，失败重试 2 次×10s
```

证书管理：`cert.InitCert()` + `cert.SubscribeCert(externalServer)`，证书到达后经 `UpdateCert` 触发 HTTPS server 启动（`src/main.go:202-206`、`src/common/https/https_server.go:118-138`）。

审计日志：见 [log-lager.md](log-lager.md)。
