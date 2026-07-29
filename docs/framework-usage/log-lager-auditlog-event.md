# 日志体系使用指导（日志）

> 版本：go-chassis lager（stub）+ fusionstage auditlog v1.9.7（stub）+ 自研 event 存储 ｜ 调用点：全模块 ｜ 涉及文件：common/logger(2) + common/event(2) + 全业务文件 ｜ 基线：origin/main @ ae0a8a6（2026-07-29 复核）

## 用途定位

三条相互独立的日志通道，**按用途选择，不可混用**：

| 通道 | 入口 | 用途 | 实现 |
| --- | --- | --- | --- |
| 业务运行日志 | `common/logger` | Infof/Warnf/Debugf/Errorf/Fatalf | 转发 `lager.Logger`（`logger.go:13-42`），配置 `conf/lager.yaml` |
| 操作/安全审计日志 | `logger.AuditsLog` / `AuditsSecAndOpsLog` | 管理面操作的审计上报 | HTTP POST 到 `cse://AuditLog/plat/audit/v1/logs|seculogs`（`auditlog.go:17-19,66-126`） |
| 业务事件日志 | `service.EventService.ReportEvent` | 登录/登出等业务事件留痕 | 本地文件滚动存储（`common/event/local_storage.go`） |

## 初始化与配置

- 业务日志：`main.go:49` `logapi.AddFilterFileName("logger/logger.go")` 用于正确显示调用点文件名；本地 stub 下 lager 为空实现，控制台无输出属正常（AGENTS.md 说明）。
- 事件存储：`service/event_service.go:37-46`，`sync.Once` 初始化 `event.InitFactory()` 并注册 `localAuditComponent` 存储；事件文件路径取 `conf.Instance().Logger.EventFile`（app.conf 无此项时为空 → 输出到 stdout，`local_storage.go:64-68`）。
- 审计日志：无需初始化，`gsfapi.NewCspRestInvoker()` 直接调 CSE。

## 核心使用模式

```go
// 业务日志（来源：遍布全仓，如 src/service/user_service.go:101）
logger.Infof("[CreateOrUpdateUser] get User[%+v] failed, err: %v", key, err)

// 打日志并返回 error 二合一（logger.go:33-37）
return logger.TeeErrorf("read rquest body failed, err: %v", err)
```

```go
// 审计日志（来源：src/common/logger/auditlog.go:66-126）
logger.AuditsLog(&logger.AuditsPara{
	OperationZH: "...", OperationEN: "...",
	OperateType: logger.ADD, Level: logger.ImportantLevel,
	Username: "...", Terminal: "...", Result: 0,
	Detail: "...", DetailZH: "...",
}, logger.OpsLog)   // 或 logger.SecLog；AuditsSecAndOpsLog 同时记两条
```

```go
// 业务事件（来源：src/controllers/login_controller.go:180-203）
event := events.NewInfo(events.Login)
event.SetEventData(events.LoginEventData{IMEI: ..., ...})
c.eventService.ReportEvent(event)   // 失败只记日志，不影响主流程（login_controller.go:200-203）
```

## 封装层与扩展点

- `common/logger` 6 个包级函数屏蔽 lager API（`logger.go:13-42`）。
- 事件存储是工厂模式：`event.StorageFactory` 按 location 注册/获取（`event_storage.go:16-72`），可新增 Storage 实现注册新 location。
- 本地事件文件滚动：单文件 >20MB 转储 zip，最多 5 份、保留 90 天，双 deleter 链式处理（`local_storage.go:25-38,157-222`）；每小时 ticker 例行清理（`local_storage.go:77-84`）。
- **业务代码一律用 `GIDS/common/logger`，禁止直接 import lager/log。**

## 并发与线程模型

`lager.Logger` 与 auditlog sink 协程安全；事件工厂用 `sync.RWMutex` 保护 map（`event_storage.go:28-72`）；`localEventStorage.Record` 无锁——并发写同一文件由 `auditlog.Logger` 内部保证（stub 下无从验证，生产 SDK 负责）。

## 错误处理与容错

- 审计日志发送失败只记 Errorf，不返回错误、不重试（`auditlog.go:116-124`）——审计丢失可接受。
- 事件 Record 失败由调用方决定：`login_controller.go:200-203` 明示"Event记录失败导致的异常暂不导致Login失败"。
- 事件文件创建失败降级 stdout（`local_storage.go:70-75`）。

## 约定与规范

- 日志消息带模块前缀 `[FuncName]`（惯例，`user_service.go:101`、`local_storage.go` 全文）。
- 注释风格 `// 函数名 功能说明`（AGENTS.md 基线）。
- 敏感字段脱敏：`desensitize()` 处理用户 key（`user_service.go:101,114`）。

## 已知问题与反模式

- **Errorf 当 DEBUG 用**：`alarm_service.go` 全文（如 232,250,262）、`db_init.go` 大量工号日志（194,335,354,369,376）用 Errorf 输出正常流程信息且含连接串/密码——严重反模式，禁模仿。
- `alarm_service.go:95` 混用标准库 `log.Println` 与 `logger`——不统一。
- 审计 body 双重 JSON 序列化（`auditlog.go:97-107`）是对端要求，改不得。
- `local_storage.go` 使用已废弃 `ioutil.ReadDir`（`local_storage.go:9,226`）。

## AI 编码指南

- 普通流程日志：`logger.Infof/Warnf`；异常且需要调用方处理：`logger.Errorf` + 返回 error，或直接用 `logger.TeeErrorf` 一步完成。级别按语义选，**禁止**用 Errorf 输出正常流程。依据：上文「已知问题与反模式」。
- 管理面增删改操作：调 `logger.AuditsSecAndOpsLog` 或 `AuditsLog` 记录审计；用户标识、操作结果必填。依据：`auditlog.go:44-63` 参数结构。
- 新增业务事件：`models/events` 加事件类型与 Data 结构 → `events.NewInfo(type)` + `SetEventData` → `NewEventService().ReportEvent(event)`；事件上报失败不得阻断主流程。依据：`login_controller.go:180-203`。
