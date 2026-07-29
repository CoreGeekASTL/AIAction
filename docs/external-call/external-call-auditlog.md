# AuditLog 服务出站调用

AuditLog 为平台审计日志服务，通过 GSF CSE Rest Invoker 以 `cse://AuditLog/...` 形式调用。封装在 `common/logger/auditlog.go`（AuditsLog / AuditsSecAndOpsLog）。

| 接口名 | 协议 | 调用位置 | 业务场景 |
|--------|------|----------|----------|
| POST /plat/audit/v1/logs | CSE REST | common/logger/auditlog.go | 操作日志上报 |
| POST /plat/audit/v1/seculogs | CSE REST | common/logger/auditlog.go | 安全日志上报 |

## POST /plat/audit/v1/logs

- 协议：CSE REST POST `cse://AuditLog/plat/audit/v1/logs`（`gsfapi.NewCspRestInvoker().Invoke`），请求体为二次 JSON 序列化的审计日志字符串
- 调用位置：common/logger/auditlog.go（AuditsLog 函数）；业务调用方 controllers/cache_controller.go（DeleteCache 操作审计）
- 业务场景：管理面敏感操作（如删除用户页面缓存）发生时，向平台上报操作审计日志
- 接口功能：上报操作日志（operation/level/userName/terminal/serviceName/result/operateType/detail 等字段），返回 200 视为成功

## POST /plat/audit/v1/seculogs

- 协议：CSE REST POST `cse://AuditLog/plat/audit/v1/seculogs`，请求体格式同上（level 固定为 ImportantLevel，不含 operateType）
- 调用位置：common/logger/auditlog.go（AuditsLog 函数，经 AuditsSecAndOpsLog 与操作日志成对调用）
- 业务场景：与操作日志配套的安全审计日志上报
- 接口功能：上报安全日志，返回 200 视为成功；失败仅记录错误日志，不影响业务流程
