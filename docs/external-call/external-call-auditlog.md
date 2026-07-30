# external-call-auditlog

> 下游服务：AuditLog（审计日志微服务，CSE 注册名 `AuditLog`）。
> 调用方式：GSF 框架 CSE rest invoker——`gsfapi.NewCspRestInvoker().Invoke(POST, cse://AuditLog/...)`（common/logger/auditlog.go `AuditsLog`）。

| 接口名 | 协议 | 调用位置 | 业务场景一句话 |
|--------|------|----------|----------------|
| POST /plat/audit/v1/logs | HTTP（CSE） | common/logger/auditlog.go | 上报操作类审计日志 |
| POST /plat/audit/v1/seculogs | HTTP（CSE） | common/logger/auditlog.go | 上报安全类审计日志 |

## HTTP

## POST /plat/audit/v1/logs

- 协议：HTTP POST `cse://AuditLog/plat/audit/v1/logs`（常量 `OpsLog`，经 ServiceComb 服务发现路由）
- 调用位置：common/logger/auditlog.go（`AuditsLog` 函数）；业务调用方：controllers/cache_controller.go（DeleteCache 删除缓存操作审计）
- 业务场景：管理面/客户端发起敏感操作后，向平台审计服务上报操作日志，满足操作可审计要求
- 接口功能：请求体为二次 JSON 序列化的审计记录（operation 中英文、level、userName、dateTime、appName/appId、terminal、serviceName、result、detail/detail_zh、operateType）；HTTP 200 视为成功，失败仅记日志

## POST /plat/audit/v1/seculogs

- 协议：HTTP POST `cse://AuditLog/plat/audit/v1/seculogs`（常量 `SecLog`，经 ServiceComb 服务发现路由）
- 调用位置：common/logger/auditlog.go（`AuditsLog` / `AuditsSecAndOpsLog` 函数）；当前业务代码中仅 `AuditsSecAndOpsLog` 封装了同时上报入口，未见直接业务调用方（操作日志走 OpsLog）
- 业务场景：安全事件类审计上报（与操作日志并列的安全日志通道）；调用此 URL 时 level 强制置为 ImportantLevel
- 接口功能：请求体结构与操作日志一致（不含 operateType）；HTTP 200 视为成功，失败仅记日志
