# 流量统计

> 功能域：流量统计　接口数：4（仅内部）　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

会话与流量统计上报 + 月度静态数据导出，仅 innerServer（HTTP）暴露，TrafficStatsController 承载（routers/beego_router.go:37）。另有每日凌晨 2 点定时清理 3 个月前历史数据的调度任务（scheduler/task_scheduler.go:117，非 HTTP 入口）。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| SessionStats | 上报会话统计（插入或更新） | controllers/traffic_stats_controller.go | POST /stats/v1/session |
| MediaTrafficStats | 批量上报媒体流量统计 | controllers/traffic_stats_controller.go | POST /stats/v1/traffic/media |
| ControlTrafficStats | 批量上报控制流量统计 | controllers/traffic_stats_controller.go | POST /stats/v1/traffic/control |
| ExportStaticData | 导出月度静态数据 zip | controllers/traffic_stats_controller.go | GET /stats/v1/exportStaticData/:month |

## 3. 数据结构说明

- **SessionStats**
  - 请求 `db.SessionStats`（models/db/traffic_stats.go，t_session_stats 表）：session_id；app_type；started_at；finished_at；tcp_unique_id
  - 响应：retcode 标准结构
- **MediaTrafficStats / ControlTrafficStats**
  - 请求 `req.MultiTableRequest`（models/req/request_entity.go）：items []json.RawMessage（必填，非空），单条元素对应 `db.MediaTrafficStats` / `db.ControlTrafficStats`：session_id；app_type；started_at；finished_at；out_bytes；access_type
  - 响应：retcode 标准结构
- **ExportStaticData**
  - 请求：path 参数 `:month`（格式 `2006-01`，仅格式校验）
  - 响应：zip 文件流（`<month>.zip`，Content-Type: application/zip），内含 session_stats.csv / media_stats.csv / control_stats.csv（纯数据 CSV，无 header 行）

## 4. 风险与注意点

- **ExportStaticData 无月份范围校验**：controllers/traffic_stats_controller.go:177-188（`:month` 仅格式校验，无范围限制，可遍历历史月份导出流量数据）
- **批量上报条数无上限**：controllers/traffic_stats_controller.go:90-113（MultiTableRequest.items 仅校验非空，无条数/大小上限，超大批量可能打满 DB）
