# 插件加载

> 用途：plugin-loading　实例数：1　返回 [README.md](README.md)

## 1. 核心作用

承载插件并发加载的进度同步——以 buffered channel 汇总各 worker 加载进度，实现"等全部完成"的并发编排。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| progressChan | 插件加载进度同步队列 | service/plugin_service.go |

## 3. 实例详解

- **progressChan**
  - 结构：`var progressChan = make(chan db.PluginPackage, len(browserGWs))`（service/plugin_service.go:76）
  - 关键操作：并发加载插件，每个 worker 完成后向 channel 投递进度，主流程 `recordLoadPluginProgress`（:298）消费汇总
  - 使用点：service/plugin_service.go:76（创建，容量=worker 数）、:298（消费）、:313（签名传递）
  - 并发模型：多 worker 生产 + 单主流程消费，容量=worker 数实现"等全部完成"

## 4. 使用模式与约定

- 等待多并发结果用 buffered chan，容量=worker 数（progressChan 范式）
