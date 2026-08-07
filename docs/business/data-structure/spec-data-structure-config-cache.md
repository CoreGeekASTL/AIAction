# 配置缓存

> 用途：config-cache　实例数：1　返回 [README.md](README.md)

## 1. 核心作用

承载配置中心的内存缓存，配置定时整体刷新，是业务流程的数据中枢。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| configCenter.configs | 配置中心缓存（5min 定时整体刷新） | service/config_center_service.go |

## 3. 实例详解

- **configCenter.configs**
  - 结构：`type configCenterServiceImpl struct { configs map[string]string; dao *dao.ConfigCenterDao; stopChan chan struct{} }`（service/config_center_service.go:30）
  - 关键字段：configs（key=配置 Key，value=配置 Value）；stopChan（控制刷新 goroutine 退出）
  - 典型操作：Refresh() 重新 List 全表构造新 map 后整体赋值 `c.configs = configMap`；GetConfig() 裸读 `configCenter.configs[key]`
  - 使用点：service/config_center_service.go:60（Refresh 整体替换）、:65（GetConfig 裸读）、:92-95（init 单例+stopChan）
  - 并发模型：**无锁**，独立 goroutine 定时 Refresh 整体替换，主线程并发裸读——存在并发读写隐患

## 4. 使用模式与约定

- 配置缓存定时整体刷新：独立 goroutine 重新 List 全表构造新 map 后整体替换，调用方通过 GetConfig 读取
