# 生命周期信号

> 用途：lifecycle-signal　实例数：1　返回 [README.md](README.md)

## 1. 核心作用

承载 goroutine 生命周期控制——以 `chan struct{}` + close 广播停止信号，监控/配置/调度/HTTPS 多处后台 goroutine 统一感知退出。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| stopChan | goroutine 停止信号（监控/配置/调度/HTTPS） | service/monitor_service.go 等 |

## 3. 实例详解

- **stopChan**
  - 结构：`stopChan chan struct{}`，多处定义：service/monitor_service.go:51、service/config_center_service.go:33、scheduler/task_scheduler.go:19、common/https/https_server.go:37
  - 关键操作：`close(stopChan)` 触发消费 goroutine 的 `case <-stopChan: return` 退出
  - 使用点：service/monitor_service.go:122/130/140、service/config_center_service.go:95/114/123、scheduler/task_scheduler.go:34
  - 并发模型：close 广播退出信号，多消费 goroutine 同时感知

## 4. 使用模式与约定

- goroutine 生命周期统一 `chan struct{}` + close 广播退出（stopChan 范式，四处一致）
