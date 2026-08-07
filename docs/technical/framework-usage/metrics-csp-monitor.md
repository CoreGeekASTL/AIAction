# CSPGoMonitorSDK 话统监控使用指导（监控/可观测）

## 用途定位
CSP 平台运营指标（话统）上报：在线用户数、各机型在线数、单 VM 支撑数、应用/站点流量。注册指标模型后按 5min 周期从 DB 统计并打点。项目无 Prometheus 埋点（client_golang 为间接依赖，未直接使用）。


## 使用模式

注册→定时上报序列（来源：`src/service/monitor_service.go:55-136`）：

```go
// 1. 读 monitor.json 构建指标模型
// 2. 每 10s 重试直至注册成功
err = monitorsdk.MonSdkInstance.InitMonitor(appId, serviceName, instanceName)
err = monitorsdk.MonSdkInstance.RegisterBasicInfo(m.monitorJson)
// 3. 起 5min Ticker，每周期遍历指标组
```

指标计算与上报（来源：`src/service/monitor_service.go:155-228`）：

```go
metricFunc := m.metricFunMap[metric.ID]        // MetricID → statsService.GetXxx 映射
results := m.getMetricResults(metricFunc, id)  // defer recover 兜底
// 新对象先注册 moi
monitorsdk.MonSdkInstance.ObjChange(uint32(mocId), 1, moiId)
// 打点
monitorsdk.MonSdkInstance.SetMetric(int(metric.ID), item.Obj, float64(item.Cnt))
```

时间窗口：`monitorutil.GetLastFiveMinuteWindow(nil)`（`src/service/monitor_service.go:190`）。
