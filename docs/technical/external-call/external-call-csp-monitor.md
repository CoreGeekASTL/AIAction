# external-call-csp-monitor

> 下游服务：CSP 话统监控组件（平台 SDK `CSPGoMonitorSDK`，包 `CSPGoMonitorSDK/api/monitor`，源码桩在 src/stubs/CSPGoMonitorSDK/，真实 SDK 由平台提供）。
> 调用方式：进程内 SDK 调用（`monitorsdk.MonSdkInstance` 单例）；监控模型文件 monitor.json、统计 SQL 模板 sql.yaml（路径见 app.conf `[cspmonitor]`）。

| 接口名 | 协议 | 调用位置 | 业务场景一句话 |
|--------|------|----------|----------------|
| MonSdkInstance.InitMonitor | 平台 SDK | service/monitor_service.go | 初始化话统 SDK |
| MonSdkInstance.RegisterBasicInfo | 平台 SDK | service/monitor_service.go | 注册话统模型（monitor.json） |
| MonSdkInstance.ObjChange | 平台 SDK | service/monitor_service.go | 动态注册测量对象（MOI） |
| MonSdkInstance.SetMetric | 平台 SDK | service/monitor_service.go | 定时上报指标值 |

## SDK

## MonSdkInstance.InitMonitor

- 协议：平台 SDK
- 调用位置：service/monitor_service.go（`InitCspMonitor` 函数，由 `InitMonitorSchedule` 循环重试调用，间隔 10s 直到成功）；启动入口 main.go `go monitorService.InitMonitorSchedule()`
- 业务场景：服务启动后初始化 CSP 话统 SDK，携带 appId、serviceName、instanceName（取自 go-chassis 全局配置）
- 接口功能：SDK 初始化，失败返回 error 由上层重试

## MonSdkInstance.RegisterBasicInfo

- 协议：平台 SDK
- 调用位置：service/monitor_service.go（`InitCspMonitor` 函数）
- 业务场景：SDK 初始化成功后，把 monitor.json 中的话统模型（MetricGroups：MocId + Metrics 列表）注册到平台话统，先注册才能上报数据
- 接口功能：入参为 monitor.json 原文；失败返回 error 由上层重试

## MonSdkInstance.ObjChange

- 协议：平台 SDK
- 调用位置：service/monitor_service.go（`addMoiIdIfNotExists` 函数）
- 业务场景：打点过程中出现新的测量对象（如某机型、某实例、某应用）时，按 MocId 动态注册 MOI（操作类型 1=新增），本地 `mocIdMap` 去重避免重复注册
- 接口功能：入参 mocId、变更类型、moiId；失败返回 error，该对象本次打点跳过

## MonSdkInstance.SetMetric

- 协议：平台 SDK
- 调用位置：service/monitor_service.go（`processMetricResults` 函数，由 `monitorSchedule` 每 5 分钟定时触发）
- 业务场景：按 5 分钟周期执行各类统计（在线人数、分机型在线、单 VM 支撑用户数、应用流量、站点流量，数据来源为 GaussDB 统计 SQL），把结果逐条上报到平台话统
- 接口功能：入参 metricId、测量对象名、指标值（float64）；失败仅记日志
