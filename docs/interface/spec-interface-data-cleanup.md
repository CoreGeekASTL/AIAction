# 统计数据定时清理

> 功能域：data-cleanup　接口数：1（定时任务）　所属 server：进程内调度
> 子文档 of [README.md](README.md)

## 1. 定位

每日定时清理超期的流量/会话统计历史数据，防止统计表无限增长。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| DataCleanupScheduler | 清理 `constants.CleanupMonths` 个月前的统计数据 | src/scheduler/task_scheduler.go | 定时任务：每日 02:00（本地时区），失败最多重试 3 次、间隔 10 分钟 |

## 3. 数据结构说明

- **DataCleanupScheduler**
  - 无外部请求；调度逻辑：每日凌晨 02:00 触发（src/scheduler/task_scheduler.go:114-125），调用 `TrafficStatsService.CleanOldStats(CleanupMonths)` 清理媒体/控制/会话三张统计表（src/scheduler/task_scheduler.go:128-148）
  - 生命周期：main.go:87 启动（StartDataCleanupScheduler），优雅退出回调中停止（src/main.go:216-219）

## 4. 风险与注意点

- **多实例重复清理**：src/scheduler/task_scheduler.go（调度器无分布式选主，多实例部署时每个实例都会执行清理；仓内虽有 schedule_election 模型，此调度器未接入）
- **清理窗口常量**：清理月数由 `constants.CleanupMonths` 固定配置，非运行时配置项（src/scheduler/task_scheduler.go:132）
