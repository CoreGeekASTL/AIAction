# scheduler 模块结构文档

> 生成时间：2026-08-11
> 所属仓：AIAction（GIDS）
> 模块路径：src/scheduler

## 模块职责

定时任务调度层：task_scheduler.go 基于 time 与 sync 周期性触发 Service 层的后台任务（如流量统计、监控上报），由 main.go 启动。

## 子模块关系图

本模块为扁平包结构，无子模块。

## 子模块说明

| 关键文件 | 说明 |
| --- | --- |
| task_scheduler.go | 定时任务调度器，周期调用 Service 层 |
