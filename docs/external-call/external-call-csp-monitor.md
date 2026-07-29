# CSP 话统监控出站调用

运营指标（话统）通过 `CSPGoMonitorSDK` 上报到 CSP 监控平台，封装在 service/monitor_service.go。SDK 内部通道出站（具体下游地址由 SDK/平台注入，待确认）。

| 接口名 | 协议 | 调用位置 | 业务场景 |
|--------|------|----------|----------|
| InitMonitor | CSPGoMonitorSDK | service/monitor_service.go | 监控 SDK 初始化 |
| RegisterBasicInfo | CSPGoMonitorSDK | service/monitor_service.go | 监控模型注册 |
| ObjChange | CSPGoMonitorSDK | service/monitor_service.go | 监控对象（MOI）注册 |
| SetMetric | CSPGoMonitorSDK | service/monitor_service.go | 指标打点上报 |

## InitMonitor

- 协议：CSPGoMonitorSDK `MonSdkInstance.InitMonitor(appId, serviceName, instanceName)`
- 调用位置：service/monitor_service.go（InitCspMonitor 函数，main.go 中以协程启动，失败间隔 10s 重试直至成功）
- 业务场景：服务启动后初始化话统监控 SDK
- 接口功能：以应用 ID、服务名、实例名初始化监控上报通道

## RegisterBasicInfo

- 协议：CSPGoMonitorSDK `MonSdkInstance.RegisterBasicInfo(monitorJson)`
- 调用位置：service/monitor_service.go（InitCspMonitor 函数）
- 业务场景：初始化完成后注册监控模型（指标分组 MocId、指标 Metric 定义），模型文件默认 `/opt/csp/gids/module/conf/monitor.json`
- 接口功能：向监控平台注册指标元数据，注册成功后才开始打点

## ObjChange

- 协议：CSPGoMonitorSDK `MonSdkInstance.ObjChange(mocId, 1, moiId)`
- 调用位置：service/monitor_service.go（addMoiIdIfNotExists 函数）
- 业务场景：某指标出现新的监控对象（如某机型、某实例）时，先注册该对象再上报数值
- 接口功能：注册监控对象实例（MOI），已注册对象本地缓存去重

## SetMetric

- 协议：CSPGoMonitorSDK `MonSdkInstance.SetMetric(metricId, obj, value)`
- 调用位置：service/monitor_service.go（processMetricResults 函数，5 分钟定时器 monitorSchedule 驱动）
- 业务场景：每 5 分钟从流量统计服务（TrafficStatsService，数据源自 GaussDB 统计表）取最近 5 分钟窗口数据并上报运营指标
- 接口功能：上报指标数值。覆盖指标：在线人数、各机型在线人数、单 VM 支持用户数、应用流量、站点流量
