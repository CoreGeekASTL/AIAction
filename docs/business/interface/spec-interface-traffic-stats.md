# 流量会话统计与导出

> 功能域：traffic-stats　接口数：4　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

接收浏览器/媒体/控制面实例上报的会话与流量统计数据入库，并支持按月导出三张统计表的 CSV 压缩包。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| SessionStats | 上报会话统计（按 sessionID 插入或更新） | src/controllers/traffic_stats_controller.go | POST /stats/v1/session |
| MediaTrafficStats | 批量上报媒体流量统计 | src/controllers/traffic_stats_controller.go | POST /stats/v1/traffic/media |
| ControlTrafficStats | 批量上报控制面流量统计 | src/controllers/traffic_stats_controller.go | POST /stats/v1/traffic/control |
| ExportStaticData | 按月导出三表统计数据（zip 包响应） | src/controllers/traffic_stats_controller.go | GET /stats/v1/exportStaticData/:month |

## 3. 数据结构说明

- **SessionStats**
  - 请求 `db.SessionStats`（src/models/db/traffic_stats.go）：SessionID、AppType（int）、StartedAt、FinishedAt（字符串时间）、TcpUniqueId
  - 响应 `resp.BaseResponse`
- **MediaTrafficStats / ControlTrafficStats**
  - 请求 `req.MultiTableRequest`（src/models/req/request_entity.go）：Items（[]json.RawMessage，必填非空；每条按目标表结构解析——媒体表 `db.MediaTrafficStats`：SessionID、AppType、StartedAt、FinishedAt、OutBytes、AccessType；控制表 `db.ControlTrafficStats` 同构）
  - 响应 `resp.BaseResponse`；两接口共用 `insertMultiData`，按路径 tag 区分目标表（src/controllers/traffic_stats_controller.go:90-104）
- **ExportStaticData**
  - 请求：路径参数 `month`，格式 `2006-01`（src/controllers/traffic_stats_controller.go:26），格式非法返回 ClientFailed
  - 响应：非 JSON 壳——zip 字节流，header `Content-Disposition: attachment; filename=<month>.zip`、`Content-Type: application/zip`；zip 内含 session_stats.csv / media_stats.csv / control_stats.csv（src/controllers/traffic_stats_controller.go:158-175）

