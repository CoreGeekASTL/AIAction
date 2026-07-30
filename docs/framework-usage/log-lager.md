# 日志框架使用指导（日志）

> 版本：go-chassis lager（stub）+ CSP auditlog v1.9.7（stub）｜ 调用点：~40 文件引用 common/logger ｜ 涉及文件：40 ｜ 基线：main (6c93561)

## 用途定位
- `GIDS/common/logger`：全项目唯一日志入口，包装 go-chassis `lager.Logger`（`src/common/logger/logger.go:9-14`）。
- `GIDS/common/logger/auditlog.go`：操作/安全审计日志，发往 `cse://AuditLog` 服务（`src/common/logger/auditlog.go:18-20`）。
- `code.huawei.com/fusionstage/auditlog`：事件日志文件引擎，用于 events 本地落盘（`src/common/event/local_storage.go:54`）。

## 初始化与配置
- lager 由 GSF 框架初始化时按 `src/conf/lager.yaml` 装配，业务代码无显式初始化。
- `logapi.AddFilterFileName("logger/logger.go")` 过滤封装层文件名，保证日志打印真实调用点（`src/main.go:49`）。
- 审计 sink 选择：配置了 `Logger.EventFile` 则落文件并起 1h 清理协程，否则写 stdout（`src/common/event/local_storage.go:64-75`）。

## 核心使用模式

```go
// 来源：全项目统一写法
logger.Infof("beego register router %v", k)
logger.Errorf("failed to start internal server:%v,exit", err)
logger.Warnf(...)
logger.Debugf(...)
logger.Fatalf("%v", err)            // 打印并退出，仅用于启动期致命错误
err := logger.TeeErrorf("read rquest body failed, err: %v", err) // 打日志+返回 error 一步到位（src/controllers/controller.go:74）
```

审计日志（来源：`src/common/logger/auditlog.go:44-60`）：

```go
// 组装 AuditsInfo/AuditsPara（操作类型 GET/ADD/MOD/DELETE...，级别 Minor/Important/Auto/Manual）
// 经 gsfapi 发往 OpsLog/SecLog 端点
```

事件落盘（来源：`src/common/event/local_storage.go:88-98`）：

```go
storage.engine.Print(event.ToJSON()) // auditlog.Logger，自动按 20MB×5 个文件滚动
```

## 封装层与扩展点
- 入口：`GIDS/common/logger`（5 个级别函数 + TeeErrorf）。
- 隐藏：lager.Logger 获取、格式化参数透传、审计端点。
- 扩展点：审计操作类型 `OperateType` 枚举与级别 `AuditLogLevel` 可扩展（`src/common/logger/auditlog.go:25-42`）。
- **禁止**业务代码直接 import lager 或标准库 `log`。

## 并发与线程模型
lager 与 auditlog 均并发安全，任意 goroutine 直接调用（存量证据：alarm_service 的消费 goroutine、scheduler 协程均直接调用）。

## 错误处理与容错
- 日志失败不返回错误（Infof/Errorf 无返回值）；需要错误传播用 `TeeErrorf`。
- `Fatalf` 会终止进程，只允许启动期使用（存量：`src/main.go:135,205`、`src/common/cse/cse.go:178`）。

## 约定与规范
- 格式串 + args 风格（非结构化 kv），消息用英文，关键路径打 Infof、错误打 Errorf。
- 请求/响应体打印：controller 层入口打印 request body、出口打印 response（`src/controllers/controller.go:76,108`）。

## 已知问题与反模式
- 敏感信息泄露：`src/dao/db_init.go:194,335,376` 打印完整 DB 连接串（含密码）；`src/dao/db_init.go:225` 等多处 Errorf 带工号前缀 "s00893267" 的调试残留——新代码禁止模仿。
- 双轨：`src/service/alarm_service.go:11,95` 使用标准库 `log.Println`，属历史遗留。
- 本地 stubs 的 lager 是空实现，控制台无输出属正常（AGENTS.md 已注明）。

## AI 编码指南
- 新代码打日志只用 `GIDS/common/logger` 的 Infof/Warnf/Debugf/Errorf；**禁止** import 标准库 `log` 或直用 lager（依据：上文「封装层与扩展点」，反例 `src/service/alarm_service.go:95`）。
- 返回 error 且需记录日志时优先 `logger.TeeErrorf`，避免日志与 error 文案不一致（依据：`src/common/logger/logger.go:33`）。
- **禁止**在日志中输出密码、连接串、token 等敏感信息；`Fatalf` 仅限启动失败场景（依据：上文「已知问题」）。
