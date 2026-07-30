# CSP/GSF 平台套件使用指导（基础库）

> 版本：CSPGSOMF / CSPNTP_SDK_GO / AlarmSDK_GO / Go-chassis-extend / auditlog / greatwall（全部 stub 替换）｜ 涉及文件：main.go、service/alarm_service.go、common/cert/cert.go、common/logger/auditlog.go ｜ 基线：main (6c93561)

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

## 初始化与配置
main.go 中的固定顺序（`src/main.go:46-98`）：

```go
https.InitMuenClient()   // 1. 外部 client
initGSF()                // 2. CspInit 重试 + RegistExitHandler + HealthCheckStart
https.Init()             // 3. 默认/内部 client（innerClient 依赖 GSF TLS）
transportapi.Init()      // 4. TransportSDK
ntp.Init()               // 5. NTP
runLogInit()             // 6. RunlogSDK + Modulekeeper（AddServiceName("gids")）
go dao.EnsureConnectGaussDB() // 7. DB
registerInstance(done)   // 8. CspStart（阻塞，起 goroutine）
```

告警服务在包 init() 即初始化 SDK（import service 包即触发，`src/service/alarm_service.go:63-72`）：

```go
alarmManager: manager.CSPInitAlarmSDK(appId, ServiceName, nodeName, manager.GetNodeIP())
```

## 核心使用模式
告警上报（来源：`src/service/alarm_service.go:118-166`）：

```go
// 业务接口：异步 channel 发送，单 goroutine 消费，10min 抑止重复
service.NewAlarmService().SendAlarm(alarmID, eventMessage)
service.NewAlarmService().ClearAlarm(alarmID, eventMessage)
// 消费端：manager.InitCSPAlarm(id, type) → AppendParameter(kind/namespace/sourceip/...) → SendAlarm，失败重试 2 次×10s
```

证书管理：`cert.InitCert()` + `cert.SubscribeCert(externalServer)`，证书到达后经 `UpdateCert` 触发 HTTPS server 启动（`src/main.go:202-206`、`src/common/https/https_server.go:118-138`）。

审计日志：见 [log-lager.md](log-lager.md)。

## 封装层与扩展点
- 告警封装：`AlarmService` 接口（SendAlarm/ClearAlarm），隐藏 SDK 细节与抑止逻辑。
- 其余 SDK 无业务封装层，按 main.go 范式直接调用。
- stub 机制：go.mod `replace` 将全部内部 SDK 指向 `src/stubs/` 空实现，本地可编译运行；生产链接真实 SDK（`src/go.mod:66-83`）。

## 并发与线程模型
- `CspStart` 阻塞，必须 goroutine 启动（`src/main.go:118-122`）。
- 告警 channel 容量 999，发送 5s 超时丢弃（`src/service/alarm_service.go:108-115`）。

## 错误处理与容错
- 平台 SDK 初始化失败 → 重试 → 最终 `Fatalf`（进程退出由 CSP 守护重启）。
- 告警发送失败重试 2 次后放弃，仅日志。

## 约定与规范
- 新增平台 SDK 调用：初始化放 main.go 对应顺序位，或包 init()（仅无依赖时）。
- 测试环境行为以 stub 为准，涉及平台交互的逻辑必须用接口隔离便于 mock。

## 已知问题与反模式
- alarm_service.go 含大量 Errorf 级别的调试日志与标准库 log 混用（历史遗留）。
- `GetAllActiveAlarmFromFMService` 的 goroutine+channel+timeout 模式在超时后 goroutine 泄漏写 channel（缓冲 0，`src/service/alarm_service.go:245-266`）。

## AI 编码指南
- 新增告警：调 `service.NewAlarmService().SendAlarm/ClearAlarm`；新告警 ID 加入 `AlarmList` 以便升级时清理历史告警（依据：`src/service/alarm_service.go:189`）。
- 新增平台 SDK 初始化：插入 main.go 启动序列并保持"重试+Fatalf"范式；**禁止**在请求路径同步调用平台 SDK 初始化（依据：`src/main.go:126-150`）。
- 涉及 GSF/告警/审计的单测：依赖接口注入或 stub，**禁止**假定真实平台环境（依据：go.mod replace stub 机制）。
