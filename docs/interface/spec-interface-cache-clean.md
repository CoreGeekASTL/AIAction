# 用户数据缓存删除

> 功能域：cache-clean　接口数：1　所属 server：外部(HTTPS) + 内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

按 IMEI/IMSI 删除终端用户的缓存数据（GDPR/用户数据清理场景），内外 server 双暴露。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| DeleteCache | 删除指定终端用户的缓存数据 | src/controllers/cache_controller.go | POST /app-api/devicetcp/cache/deleteCache |

## 3. 数据结构说明

- **DeleteCache**
  - 请求 `req.DeleteCacheRequest`（src/models/req/request_entity.go）：IMEI（必填，Validate 校验非空）、IMSI（必填，Validate 校验非空）
  - 响应 `resp.DataResponse`（src/models/resp/response_entity.go）：`BaseResponse{code,msg}` + `data: true`；失败时记录操作审计日志（src/controllers/cache_controller.go:41-52）

## 4. 风险与注意点

- **内外双暴露**：同一删除能力同时挂在 HTTPS 外部与 HTTP 内部 server（routers/beego_router.go:21、routers/beego_router.go:31），内部链路无鉴权需注意访问控制
