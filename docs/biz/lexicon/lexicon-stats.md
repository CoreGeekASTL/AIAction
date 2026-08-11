# 流量统计 领域词典

> 子文档 of [lexicon.md](lexicon.md)；词汇口径、来源说明与全仓待确认清单见主文档。子域锚点 `stats`，功能域口径与 docs/biz/interface/ 一致。

## 实体与业务概念

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| MediaTrafficStats | 媒体面流量统计表（t_media_traffic_stats）：session_id/app_type/起止时间/out_bytes/access_type | - | `db.MediaTrafficStats`，`src/models/db/traffic_stats.go` |
| ControlTrafficStats | 控制面流量统计表（t_control_traffic_stats），字段同媒体面 | - | `db.ControlTrafficStats`，`src/models/db/traffic_stats.go` |
| SessionStats | 会话统计表（t_session_stats）：session_id/app_type/起止时间/tcp_unique_id | - | `db.SessionStats`，`src/models/db/traffic_stats.go` |
| SessionID | 会话标识，三张统计表共有字段 | - | `db.MediaTrafficStats.SessionID` 等，`src/models/db/traffic_stats.go` |
| OutBytes | 出站字节数 | - | `db.MediaTrafficStats.OutBytes`、`db.ControlTrafficStats.OutBytes`，`src/models/db/traffic_stats.go` |
| AccessType | 接入类型（int，枚举取值待确认，见主文档） | - | `db.MediaTrafficStats.AccessType` 等，`src/models/db/traffic_stats.go` |
| TcpUniqueId | TCP 连接唯一标识 | - | `db.SessionStats.TcpUniqueId`，`src/models/db/traffic_stats.go` |

## 常量与状态枚举

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| MetricOnlineUsers | 在线用户数指标 ID 320101 | - | `monitor.MetricOnlineUsers`，`src/models/monitor/metric.go` |
| MetricOnlineUsersPerModel | 分机型在线用户数指标 ID 320102 | - | `monitor.MetricOnlineUsersPerModel`，`src/models/monitor/metric.go` |
| MetricUsersSupportedByVm | 单 VM 支持用户数指标 ID 320103 | - | `monitor.MetricUsersSupportedByVm`，`src/models/monitor/metric.go` |
| MetricApplicationTraffic | 应用流量指标 ID 320201 | - | `monitor.MetricApplicationTraffic`，`src/models/monitor/metric.go` |
| MetricSiteTraffic | 站点流量指标 ID 320202 | - | `monitor.MetricSiteTraffic`，`src/models/monitor/metric.go` |
| MonitorConfig | 监控配置结构（与 monitor.json 一致，含 MetricGroups/RealServicesName） | - | `monitor.MonitorConfig`，`src/models/monitor/metric.go` |
