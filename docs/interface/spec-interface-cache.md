# 缓存清理

> 功能域：缓存清理　接口数：1　所属 server：外部(HTTPS) + 内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

删除用户缓存数据，经 externalServer（HTTPS）与 innerServer（HTTP）双暴露，CacheController 承载。内部转调 BrowserGW 删缓存。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| DeleteCache | 删除用户数据缓存 | controllers/cache_controller.go | POST /app-api/devicetcp/cache/deleteCache |

## 3. 数据结构说明

- **DeleteCache**
  - 请求 `req.DeleteCacheRequest`（models/req）：IMEI（15 位纯数字）；IMSI（15 位纯数字）
  - 响应：retcode 标准结构；失败时记审计日志（OperationZH="删除用户数据"）
  - 内部转 service.DeleteCache → callBrowserGW 调 BrowserGW 删缓存（service/cache_service.go:23）

## 4. 风险与注意点

- **删缓存无重试**：service/cache_service.go:49（自建 http.Client Timeout=5s，失败仅记日志，:40-41，可能残留脏缓存）
