# 用户缓存删除 业务规则

> 生成时间：2026-08-11
> 覆盖入口：POST /app-api/devicetcp/cache/deleteCache → controllers.CacheController.DeleteCache（内外部路由均注册）

## 概述

用户缓存删除功能域按 IMEI+IMSI 标识，向所有就绪 BrowserGW 实例广播删除该用户的对象存储缓存数据。规则提取自 deleteCache 入口链路，覆盖参数校验、条件分支与错误码返回三类规则点。

## 规则表

| 规则名 | 条件 | 动作 | 依据 | 来源入口 |
| --- | --- | --- | --- | --- |
| request-body-unmarshal-to | 请求体 JSON 反序列化失败 | 返回 retcode.ClientFailed(-2)，拒绝删除 | src/controllers/cache_controller.go | POST /app-api/devicetcp/cache/deleteCache |
| delete-cache-impl-empty | IMEI 或 IMSI 为空 | 返回 "IMEI or IMSI cannot be empty" 错误，不发起任何删除调用 | src/service/cache_service.go | POST /app-api/devicetcp/cache/deleteCache |
| get-all-ready-service-instances | 无就绪 BrowserGW 实例（len(instances)==0） | 返回 "no BrowserGW instances available" 错误，删除失败 | src/service/cache_service.go | POST /app-api/devicetcp/cache/deleteCache |
| call-browser-gw-failed | 单个 BrowserGW 实例删除调用失败（非 200 / 网络错误 / 超时 defaultCacheTimeoutSeconds=5s） | 仅记录错误日志，继续遍历其余实例，整体删除仍视为成功 | src/service/cache_service.go | POST /app-api/devicetcp/cache/deleteCache |
| delete-cache-failed | DeleteCache 返回错误 | 记录审计日志（Operation "删除用户数据"，Result=1 失败）并返回 retcode.InternalFailed(-1) | src/controllers/cache_controller.go | POST /app-api/devicetcp/cache/deleteCache |

## 补充说明

删除采用尽力而为（best-effort）广播语义：任一实例失败不阻断其余实例，也不影响整体成功返回；整体失败仅在参数为空或无可用实例两种前置条件下发生。失败时除业务返回外还写操作审计日志，成功路径不写审计日志。
