# CSPGoMonitorSDK 话统监控使用指导（监控/可观测）

> 版本：CSPGoMonitorSDK（stub）｜ 调用点：1 个 service ｜ 涉及文件：service/monitor_service.go、models/monitor/metric.go、utils/monitorutil/time_util.go ｜ 基线：main (6c93561)

## 用途定位
CSP 平台运营指标（话统）上报：在线用户数、各机型在线数、单 VM 支撑数、应用/站点流量。注册指标模型后按 5min 周期从 DB 统计并打点。项目无 Prometheus 埋点（client_golang 为间接依赖，未直接使用）。

## 初始化与配置
- 启动：`service.NewMonitorService()` + `go monitorService.InitMonitorSchedule()`（`src/main.go:90-91`）。
- 指标模型文件：monitor.json，路径取 `cspmonitor::monitorJsonFile`，默认 `/opt/csp/gids/module/conf/monitor.json`（`src/service/monitor_service.go:37,77`）。
- 统计 SQL 文件：sql.yaml，路径取 `cspmonitor::sqlYamlFile`（`src/service/monitor_service.go:38,119`）。
- 身份来源：GSF 全局配置 `GetGlobalDefinition().AppID / GetSelfServiceName / GetSelfInstanceID`（`src/service/monitor_service.go:92-94`）。

## 核心使用模式

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

## 封装层与扩展点
- 封装：`MonitorService` 接口（`InitMonitorSchedule`），SDK 调用全部集中在 `MonitorServiceImpl`。
- 扩展点：新增指标 = monitor.json 加模型 + `metricFunMap` 加 `MetricID→getMetricFunc` 映射 + TrafficStatsService 加统计方法（`src/service/monitor_service.go:145-153`）。

## 并发与线程模型
- 上报在独立 goroutine 的 Ticker 循环（`src/service/monitor_service.go:123-135`）。
- `mocIdMap` 用 `sync.RWMutex` 保护（`src/service/monitor_service.go:45,214`）。
- 指标函数 panic 由 `getMetricResults` 内 `defer recover()` + 堆栈日志兜底，不中断整轮（`src/service/monitor_service.go:181-193`）。

## 错误处理与容错
- 注册失败 10s 间隔无限重试直至成功（`src/service/monitor_service.go:64-73`）。
- 单指标失败仅 Errorf 跳过，不影响其他指标。

## 约定与规范
- 打点对账周期固定 5min，与监控模板配置保持一致（代码注释，`src/service/monitor_service.go:25`）。
- moiId（对象）为空字符串的指标数据丢弃（`src/service/monitor_service.go:198-201`）。

## 已知问题与反模式
- monitor.json/sql.yaml 路径为生产绝对路径，本地无文件时 InitMonitorSchedule 直接失败返回（`src/service/monitor_service.go:78-81`），本地训练环境属预期。

## AI 编码指南
- 新增运营指标：monitor.json 增模型 → `models/monitor/metric.go` 增 MetricID → `metricFunMap` 注册统计函数 → TrafficStatsService 实现统计（sql.yaml 配 SQL）（依据：上文「核心使用模式」「扩展点」）。
- **禁止**业务代码绕过 MonitorService 直接调 `monitorsdk.MonSdkInstance`（依据：`src/service/monitor_service.go` 是唯一调用点）。
- 指标统计函数必须接受 `(startTime, endTime string)` 签名并返回 `[]Res`，内部自行防 panic 由框架兜底但不得依赖（依据：`src/service/monitor_service.go:41,181`）。
