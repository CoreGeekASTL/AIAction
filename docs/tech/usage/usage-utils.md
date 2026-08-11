# 基础库（uuid + 自研 utils）使用现状（基础库）

## 用途定位
- **google/uuid**：`github.com/google/uuid`，全仓唯一直接调用点在 `src/service/browser_service.go`（`uuid.New()` 生成实例标识），AGENTS.md 约定新代码 UUID 一律用它。
- **自研 utils**（`src/utils/`）：`fileutil`（文件落盘、权限/属主管理，配合 auditlog 常量）、`monitorutil`（监控时间处理）、`flagutil`（命令行 flag 解析覆盖 conf 默认值）。
- **公共常量与错误码**：`src/common/constants`（含 `retcode` 返回码：Success/InternalFailed/AuthFailed/ClientFailed 等，Controller 响应统一使用）、`src/common/storage/error.go`（`ErrNotExist` 跨存储统一"不存在"错误）。
- **stubs**：`src/stubs/` 下为 CSP 平台 SDK 的本地空实现（go.mod replace 目标），是本地可编译运行的基础设施，非业务框架。

## 使用模式

```go
// 来源：src/service/browser_service.go
uid := uuid.New()
```

```go
// 来源：src/common/logger/auditlog.go（响应码与审计序列化配合 retcode 使用）
// Controller 响应统一：resp.BaseResponse{Code: retcode.Success, Message: "success"}
// 来源：src/controllers/controller.go
```
