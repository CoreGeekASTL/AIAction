# 缓存管理

> 功能域：缓存管理　接口数：1　所属 server：外部 + 内部
> 子文档 of [README.md](README.md)

## 1. 定位

按终端身份（IMEI/IMSI）清理服务端为该终端缓存的会话/实例数据，供终端或内部运维触发。同一接口经 externalServer 与 innerServer 双暴露。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| DeleteCache | 删除指定终端的缓存 | controllers/cache_controller.go | POST /app-api/devicetcp/cache/deleteCache |

## 3. 数据结构说明

- **DeleteCache**
  - 请求 `req.DeleteCacheRequest`（models/req/request_entity.go）：IMEI（必填）、IMSI（必填），任一为空校验失败返回 client 错误
  - 响应 `resp.BaseResponse`（models/resp/base.go）：code/msg
