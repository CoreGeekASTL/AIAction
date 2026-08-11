# 流量统计

> 功能域：流量统计　接口数：4 + 1 个定时任务　所属 server：内部
> 子文档 of [README.md](README.md)

## 1. 定位

会话与媒体/控制面流量统计数据的上报入库、按月导出 CSV（zip 打包），以及每日凌晨 2 点定时清理过期统计数据。HTTP 接口仅注册在内部 server。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| SessionStats | 上报/更新会话统计 | controllers/traffic_stats_controller.go | POST /stats/v1/session |
| MediaTrafficStats | 批量上报媒体面流量统计 | controllers/traffic_stats_controller.go | POST /stats/v1/traffic/media |
| ControlTrafficStats | 批量上报控制面流量统计 | controllers/traffic_stats_controller.go | POST /stats/v1/traffic/control |
| ExportStaticData | 按月导出三类统计 CSV 并 zip 返回 | controllers/traffic_stats_controller.go | GET /stats/v1/exportStaticData/:month |
| DataCleanupScheduler | 每日 02:00 清理 CleanupMonths 个月前的统计数据（最多重试 3 次，间隔 10 分钟） | scheduler/task_scheduler.go | 定时任务：每天 02:00（timer 驱动，非 cron 表达式） |

## 3. 数据结构说明

- **SessionStats**
  - 请求 `db.SessionStats`（models/db/traffic_stats.go）：session_id、app_type、started_at、finished_at、tcp_unique_id
  - 响应 `resp.BaseResponse`
- **MediaTrafficStats / ControlTrafficStats**
  - 请求 `req.MultiTableRequest`（models/req/request_entity.go）：items 为 JSON RawMessage 数组（非空），逐条反序列化为 `db.MediaTrafficStats` / `db.ControlTrafficStats`（session_id、app_type、started_at、finished_at、out_bytes、access_type）
  - 响应 `resp.BaseResponse`
- **ExportStaticData**
  - 请求：路径参数 month，格式 `2006-01`（如 2026-08），非法格式返回 client 错误
  - 响应：zip 文件流（Content-Disposition: attachment; filename={month}.zip），内含 session_stats.csv、media_stats.csv、control_stats.csv
- **DataCleanupScheduler**
  - 无请求/响应；启动入口 `StartDataCleanupScheduler()`，调用 `TrafficStatsService.CleanOldStats`
