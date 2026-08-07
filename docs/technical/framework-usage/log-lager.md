# 日志框架使用指导（日志）

## 用途定位
- `GIDS/common/logger`：全项目唯一日志入口，包装 go-chassis `lager.Logger`（`src/common/logger/logger.go:9-14`）。
- `GIDS/common/logger/auditlog.go`：操作/安全审计日志，发往 `cse://AuditLog` 服务（`src/common/logger/auditlog.go:18-20`）。
- `code.huawei.com/fusionstage/auditlog`：事件日志文件引擎，用于 events 本地落盘（`src/common/event/local_storage.go:54`）。


## 使用模式

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
