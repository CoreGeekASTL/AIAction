# 关键数据结构总览

| 元信息 | 值 |
|--------|-----|
| 分支 | ready/27.0-终端鉴权 分支 (2026-08-07) |
| 更新日期 | 2026-08-07 |
| Skill | spec-data-structure-analyze |

## 1. 数据结构全景

| 用途 | 实例数 | 核心作用 | 子文档 |
|---|---|---|---|
| 鉴权缓存 | 1 | 鉴权结果内存缓存（TTL+容量惰性清理） | [spec-data-structure-auth-cache.md](spec-data-structure-auth-cache.md) |
| 配置缓存 | 1 | 配置中心缓存，定时整体刷新 | [spec-data-structure-config-cache.md](spec-data-structure-config-cache.md) |
| 告警 | 3 | 告警事件队列/抑制计数/全局告警 ID 枚举 | [spec-data-structure-alarm.md](spec-data-structure-alarm.md) |
| 服务发现 | 2 | 网关实例表/CSE 链路端点列表 | [spec-data-structure-service-discovery.md](spec-data-structure-service-discovery.md) |
| 监控指标 | 2 | 指标索引（嵌套 map）/指标函数注册表 | [spec-data-structure-monitor-metrics.md](spec-data-structure-monitor-metrics.md) |
| 白名单去重 | 1 | 白名单导入去重集合 | [spec-data-structure-whitelist-dedup.md](spec-data-structure-whitelist-dedup.md) |
| TLS 证书 | 2 | TLS 弱算法集合/证书重启信号 | [spec-data-structure-tls-cert.md](spec-data-structure-tls-cert.md) |
| 插件加载 | 1 | 插件加载进度同步队列 | [spec-data-structure-plugin-loading.md](spec-data-structure-plugin-loading.md) |
| 生命周期信号 | 1 | goroutine 停止信号 | [spec-data-structure-lifecycle-signal.md](spec-data-structure-lifecycle-signal.md) |

自检：扫描 14 个实例，已记录 14 个，未归类 0 个（见上表），差集已清零（2026-08-07）
