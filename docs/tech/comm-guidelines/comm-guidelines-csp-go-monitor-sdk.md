# csp-go-monitor-sdk 通信规范

## 平台 SDK

### InitMonitor + RegisterBasicInfo（话统注册）

- 业务场景：服务启动后初始化 CSP 话统监控：先注册监控模型（monitor.json），注册成功后才能上报数据
- 接口功能：`MonSdkInstance.InitMonitor(appId, serviceName, instanceName)` 初始化 SDK；`MonSdkInstance.RegisterBasicInfo(monitorJson)` 注册指标模型（模型文件默认 /opt/csp/gids/module/conf/monitor.json，配置 key `cspmonitor::monitorJsonFile`）
- 调用位置：src/service/monitor_service.go（InitMonitorSchedule → InitCspMonitor）
- 协议信息：
  - 协议：CSPGoMonitorSDK（MonSdkInstance，SDK 内部通道；src/stubs/CSPGoMonitorSDK）
  - 封装方式：平台 SDK 直连（`CSPGoMonitorSDK/api/monitor`）
  - 超时重试：注册失败每 10s（InitMonitorPeriod）无限重试直到成功
  - 错误码处理：err 记日志后进入下一轮重试

### SetMetric（指标上报）

- 业务场景：注册成功后每 5 分钟（DotPeriodFiveMin）定时采集并上报各项指标：在线用户数、各机型在线人数、单 VM 支持人数、应用流量、站点流量（指标函数见 createMetricFunctionMap）
- 接口功能：`MonSdkInstance.SetMetric(metricId, moiId, value)` 上报单个指标值；新对象先 `MonSdkInstance.ObjChange(mocId, 1, moiId)` 注册 MOI
- 调用位置：src/service/monitor_service.go（startCspMonitor → monitorSchedule → processMetricResults / addMoiIdIfNotExists）
- 协议信息：
  - 协议：CSPGoMonitorSDK
  - 封装方式：平台 SDK 直连
  - 超时重试：未设置；单指标失败仅记日志跳过，不影响其余指标
  - 错误码处理：SetMetric/ObjChange err 记日志并 continue；指标函数 panic 由 recover 兜底记日志
