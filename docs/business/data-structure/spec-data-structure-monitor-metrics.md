# 监控指标

> 用途：monitor-metrics　实例数：2　返回 [README.md](README.md)

## 1. 核心作用

承载监控打点的配置/指标索引——用嵌套 map（外层 MocID→内层 set）建指标索引，并用 map 承载"指标 ID→采集函数"的函数注册表，是监控打点路由的数据中枢。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| mocIdMap | 监控指标索引（MocID → 指标名 set） | service/monitor_service.go |
| metricFunMap | 指标函数注册表（MetricID → 采集函数） | service/monitor_service.go |

## 3. 实例详解

- **mocIdMap**
  - 结构：`mocIdMap map[monitor.MocID]map[string]struct{}`（service/monitor_service.go:46），嵌套 map，内层 `map[string]struct{}` 当 set 用
  - 关键字段：外层 key=monitor.MocID（监控对象类），value=指标名集合（去重）；由 `metricMapLock sync.RWMutex`（:45）保护
  - 典型操作：buildMonitorConfig 按 MetricGroups 初始化外层 map（:87）；mocIdMap[mocId] 填内层 set（:217）；读取走 RLock
  - 使用点：service/monitor_service.go:87（初始化）、:217（填充内层 set）、:46（字段定义）
  - 并发模型：RWMutex 保护，读多写少；内层 set 复用 map 当 set 的去重语义
- **metricFunMap**
  - 结构：`metricFunMap map[monitor.MetricID]getMetricFunc`（service/monitor_service.go:50），key=指标 ID，value=采集函数（func(startTime,endTime string) []Res）
  - 关键操作：createMetricFunctionMap（:145-149）在 startCspMonitor 阶段一次性填充字面量（MetricOnlineUsers→GetOnline 等），运行期只读
  - 使用点：service/monitor_service.go:146（字面量填充）、:50（字段定义）
  - 并发模型：init 阶段填充后运行期只读，无需锁

## 4. 使用模式与约定

- 嵌套 map 须配 RWMutex 保护，外层建索引、内层用 `map[K]struct{}` 做 set（mocIdMap 范式）
- 函数注册表用 map 字面量在初始化阶段一次性填充，运行期只读，无需锁（metricFunMap 范式）
- map 当 set 用时 value 一律 `struct{}`（零内存），见 mocIdMap 内层
