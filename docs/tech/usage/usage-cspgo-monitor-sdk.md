# CSPGoMonitorSDK 使用现状（监控/可观测）

## 用途定位
CSP 话统监控 SDK（本地 replace 到 `src/stubs/CSPGoMonitorSDK`），用于向运营平台上报业务指标（流量统计等）。封装在 `src/service/monitor_service.go` 的 `MonitorService`：读取 `monitor.json`（指标模型，路径由 `cspmonitor::monitorJsonFile` 配置）→ `MonSdkInstance.InitMonitor` 初始化（appID/serviceName/instanceName 取自 go-chassis 全局配置）→ `RegisterBasicInfo` 注册模型 → 成功后启动 5 分钟 ticker 协程周期性采集并上报。注册失败按 10s 间隔无限重试。指标 SQL 模板来自 `sql.yaml`（`NewTrafficStatsService(sqlYamlFile)`）。无 Prometheus 埋点（client_golang 仅为间接依赖）。

## 使用模式

```go
// 来源：src/service/monitor_service.go
err = monitorsdk.MonSdkInstance.InitMonitor(id, serviceName, instanceName)
err = monitorsdk.MonSdkInstance.RegisterBasicInfo(m.monitorJson)

func (m *MonitorServiceImpl) startCspMonitor() {
	go func() {
		ticker := time.NewTicker(DotPeriodFiveMin)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.monitorSchedule()
			case <-m.stopChan:
				return
			}
		}
	}()
}
```
