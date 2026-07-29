# 缓存清理

> 功能域：缓存清理　接口数：1（双实现）　所属 server：外部(HTTPS) + 内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

删除用户缓存数据，经 externalServer（HTTPS）与 innerServer（HTTP）双暴露，CacheController 承载（routers/beego_router.go:21/31）。内部遍历所有 Ready 态 BrowserGW 实例逐个删缓存。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| DeleteCache | 删除用户数据缓存 | controllers/cache_controller.go | POST /app-api/devicetcp/cache/deleteCache |

## 3. 数据结构说明

- **DeleteCache**
  - 请求 `req.DeleteCacheRequest`（models/req/request_entity.go）：IMEI（必填，仅非空校验）；IMSI（必填，仅非空校验）
  - 响应 `resp.DataResponse`：BaseResponse{code,msg} + Data=true；失败时记审计日志（OperationZH="删除用户数据"，controllers/cache_controller.go:41-52）
  - 内部转 `service.DeleteCache`（service/cache_service.go:23）→ 遍历 BrowserGW 实例 `callBrowserGW`（service/cache_service.go:48），HTTP DELETE `http://<browserGW>/browsergw/browser/userdata/delete`

## 4. 风险与注意点

- **部分实例删除失败被吞掉**：service/cache_service.go:39-41（单个 BrowserGW 调用失败仅记日志继续循环，函数仍返回 nil，调用方感知不到部分失败，可能残留脏缓存）
- **单实例调用无重试**：service/cache_service.go:49（自建 http.Client Timeout=5s，超时/失败不重试）
