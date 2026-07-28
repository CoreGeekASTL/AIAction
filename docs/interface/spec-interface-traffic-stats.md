# 流量统计

> 功能域：流量统计　接口数：4　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

会话与流量统计上报 + 静态数据导出，innerServer（HTTP）暴露，TrafficStatsController 承载。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| SessionStats | 上报会话统计 | controllers/traffic_stats_controller.go | POST /stats/v1/session |
| MediaTrafficStats | 上报媒体流量统计 | controllers/traffic_stats_controller.go | POST /stats/v1/traffic/media |
| ControlTrafficStats | 上报控制流量统计 | controllers/traffic_stats_controller.go | POST /stats/v1/traffic/control |
| ExportStaticData | 导出月度静态数据 | controllers/traffic_stats_controller.go | GET /stats/v1/exportStaticData/:month |

## 3. 数据结构说明

- **SessionStats / MediaTrafficStats / ControlTrafficStats**
  - 请求 `req.*StatsRequest`（models/req）：session_id；app_type；started_at；finished_at；out_bytes；access_type（对应 t_session_stats/t_media_traffic_stats/t_control_traffic_stats 表结构）
  - 响应：retcode 标准结构
- **ExportStaticData**
  - 请求：path 参数 `:month`（格式 `2006-01`）
  - 响应：CSV 文件流（session_stats.csv / media_stats.csv / control_stats.csv）

## 4. 风险与注意点

- **ExportStaticData 无月份范围校验**：controllers/traffic_stats_controller.go:43（`:month` 仅格式校验，无范围限制，可遍历历史月份导出敏感流量数据）
