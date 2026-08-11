# 流量统计 业务规则

> 生成时间：2026-08-11
> 覆盖入口：POST /stats/v1/session → controllers.TrafficStatsController.SessionStats；POST /stats/v1/traffic/media → controllers.TrafficStatsController.MediaTrafficStats；POST /stats/v1/traffic/control → controllers.TrafficStatsController.ControlTrafficStats；GET /stats/v1/exportStaticData/:month → controllers.TrafficStatsController.ExportStaticData；定时任务每日 02:00 数据清理 → scheduler.DataCleanupScheduler

## 概述

流量统计功能域接收会话/媒体流/控制流三类统计数据上报，支持按月导出 CSV 打包下载，并由定时任务清理过期数据。规则提取自四条统计入口链路与一条定时清理链路，覆盖参数校验、条件分支、阈值常量与事务回滚四类规则点。

## 规则表

| 规则名 | 条件 | 动作 | 依据 | 来源入口 |
| --- | --- | --- | --- | --- |
| handle-session-stats | 按 TcpUniqueId 查询会话记录不存在 | 插入新会话记录；已存在则按会话更新（UpdatebySession） | src/service/traffic_stats_service.go | POST /stats/v1/session |
| batch-insert-stats-tag | BatchInsertStats 的 tag 非 "media traffic stats"/"control traffic stats" | 返回 "unsupported tag" 错误，拒绝批量插入 | src/service/traffic_stats_service.go | POST /stats/v1/traffic/media；POST /stats/v1/traffic/control |
| batch-insert-item-invalid | items 中任一记录反序列化失败或 Validate 失败 | 整批拒绝，返回错误，无任何记录入库 | src/service/traffic_stats_service.go | POST /stats/v1/traffic/media；POST /stats/v1/traffic/control |
| batch-insert-tx | 批量插入按 batchSize(1000) 分批，任一批次插入失败 | 事务整体回滚，已插入批次一并撤销 | src/service/traffic_stats_service.go | POST /stats/v1/traffic/media；POST /stats/v1/traffic/control |
| validate-export-static-data-param | 路径参数 month 为空或不满足 MonthFormat("2006-01") | 返回 retcode.ClientFailed(-2)，拒绝导出 | src/controllers/traffic_stats_controller.go | GET /stats/v1/exportStaticData/:month |
| export-csv-files | 导出三个 CSV（session_stats.csv / media_stats.csv / control_stats.csv）任一失败 | 整体导出失败，不返回 zip 包 | src/controllers/traffic_stats_controller.go；src/service/traffic_stats_service.go | GET /stats/v1/exportStaticData/:month |
| query-stats-data-month | 导出数据过滤条件 | 仅导出 started_at 以 month 前缀开头的记录，按 batchSize(1000) 分批查询 | src/service/traffic_stats_service.go | GET /stats/v1/exportStaticData/:month |
| clean-old-stats | 定时任务每日 02:00 触发 | 删除三张统计表中 started_at 早于（当前时间 - constants.CleanupMonths(3) 个月）的记录，三表各自由事务包裹 | src/scheduler/task_scheduler.go；src/service/traffic_stats_service.go | 定时任务每日 02:00（calculateNextRunTime） |
| execute-cleanup-retry | 单次清理执行失败 | 间隔 10 分钟重试，最多 maxRetries(3) 次；仍失败则等下一调度周期 | src/scheduler/task_scheduler.go | 定时任务每日 02:00（calculateNextRunTime） |

## 补充说明

批量插入采用"先全量校验、后事务写入"的两段式规则：任一记录非法则整批拒绝，写入阶段任一分批失败则全量回滚，保证批次级原子性。导出与清理以月份为共同维度：导出按 month 前缀过滤，清理按 CleanupMonths(3) 滚动窗口删除。会话统计按 TcpUniqueId 幂等（存在即更新），流量统计则纯追加。定时清理在重试期间收到停止信号会立即中断返回。
