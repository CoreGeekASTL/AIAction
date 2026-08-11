# 基础规范（日志 / 配置 / 告警）

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-foundation-guidelines-analyze |
| 运行模式 | 起草模式 |

## 机制全景

| 机制 | 使用框架 / SDK | 调用点分布概要 | 现状要点 |
|---|---|---|---|
| 日志 | go-chassis lager（桩 `src/stubs/Go-chassis-extend/.../lager`）+ 统一封装 `src/common/logger/logger.go`；审计日志独立通道 `src/common/logger/auditlog.go` | 40 个 Go 文件，共约 378 处调用（Errorf 232 / Infof 125 / Warnf 12 / Fatalf 4 / Debugf 3 / TeeErrorf 2） | 业务代码全部走 `logger.Xxxf` 封装；Errorf 严重滥用（正常流程也打 ERROR）；无敏感信息脱敏处理，但日志中基本不打印密码 / 密钥；审计日志仅 1 处调用 |
| 配置 | beego.AppConfig 读 `src/conf/app.conf` + `src/common/conf/config.go` Config 结构体 + `src/utils/flagutil` 命令行参数绑定 + os.Getenv 环境变量 | `src/common/conf/config.go`、`src/main.go`、`src/service/monitor_service.go`、`src/service/alarm_service.go`、`src/dao/db_init.go` 等 | 三通道混用（配置文件 / 命令行 flag / 环境变量）；app.conf 中数据库密码明文入库；缺配多用空字符串默认值（静默零值） |
| 告警 | AlarmSDK_GO（桩 `src/stubs/AlarmSDK_GO`），封装层 `src/service/alarm_service.go` | 封装层 1 个文件；业务上报点仅 `src/controllers/management_controller.go` 2 处（300010 上报告警 / 恢复） | 告警 ID 集中常量化（仅 AlarmId300010）；上报与恢复成对；10 分钟抑制窗口 + 重试 2 次；告警内部日志滥用 Errorf |
| 监控埋点 metrics（增补） | CSPGoMonitorSDK（桩 `src/stubs/CSPGoMonitorSDK`），封装 `src/service/monitor_service.go`，指标模型 `src/models/monitor/metric.go` | `src/service/monitor_service.go`、`src/utils/monitorutil/time_util.go`，配置文件 `src/conf/monitor.json`、`src/conf/sql.yaml` | 5 分钟定时注册 + 打点上报，注册失败 10s 重试直到成功 |

## 日志

### 使用模式

业务代码统一经 `src/common/logger/logger.go` 的封装函数打日志，底层为 go-chassis lager（本仓为桩实现，控制台无输出）：

```go
// 来源：src/common/logger/logger.go
func Infof(format string, args ...interface{})  { lager.Logger.Infof(format, args...) }
func Warnf(format string, args ...interface{})  { lager.Logger.Warnf(nil, format, args...) }
func Debugf(format string, args ...interface{}) { lager.Logger.Debugf(format, args...) }
func Errorf(format string, args ...interface{}) { lager.Logger.Errorf(nil, format, args...) }
// TeeErrorf 打 ERROR 日志并返回 error，便于 "log and return"
func TeeErrorf(format string, args ...interface{}) error { ... }
```

审计日志为独立通道：`AuditsLog` 组装操作 / 安全审计 JSON，经 GSF RestInvoker POST 到 `cse://AuditLog/plat/audit/v1/logs`（操作日志）与 `.../seculogs`（安全日志），appName / appId 取自环境变量 APPNAME / APPID（来源：`src/common/logger/auditlog.go`）。

```go
// 来源：src/controllers/cache_controller.go
para := &logger.AuditsPara{OperationZH: "...", OperateType: logger.GET, Level: logger.MinorLevel, Username: "...", Result: 0}
logger.AuditsLog(para, logger.OpsLog)
```

### 调用点分布

| 关注维度 | 分布（文件路径，不带行号） | 现状说明 |
|---|---|---|
| 级别使用 | 全部 `src/controllers/`、`src/service/`、`src/dao/`、`src/common/` 业务文件；滥用重灾区 `src/service/alarm_service.go`、`src/service/monitor_service.go` | 全仓 378 处：Errorf 232 / Infof 125 / Warnf 12 / Fatalf 4 / Debugf 3 / TeeErrorf 2。Errorf 占比 61%，且 `src/service/alarm_service.go` 中大量正常流程日志（"Enter GetAllActiveAlarmFromFMService"、"get body from fmservice is..."）也打 Errorf，级别语义失真 |
| 敏感信息脱敏 | `src/controllers/exlogin_controller.go`、`src/controllers/login_controller.go`（仅打 "update user token failed"，未打印 token 明文） | 无脱敏机制（无 mask / 脱敏工具函数）；现状日志基本不输出密码 / token / IMEI 明文，属"未约定但侥幸安全" |
| 审计日志 | 定义 `src/common/logger/auditlog.go`；调用点仅 `src/controllers/cache_controller.go` | 有操作审计 + 安全审计双通道及操作类型 / 级别枚举，但登录、鉴权、配置变更等关键操作（login_controller / management_controller / config_center_controller）均未接审计日志，覆盖严重不足 |

### 应有约定建议

| 编号 | 约定条目 | 现状符合度 |
|---|---|---|
| LOG-01 | ERROR 级仅用于真实异常（需人工介入或已影响功能）；正常流程入口 / 出口用 Infof，业务校验失败用 Warnf，禁止 Errorf 当 Infof 用 | 现状缺失（alarm_service / monitor_service 大面积 Errorf 滥用） |
| LOG-02 | 日志中的敏感信息（密码 / 密钥 / token / AccessKey / IMEI / IMSI）必须脱敏后输出，禁止明文打印 | 现状部分遵守（无脱敏工具、无约定，仅靠现状未打印） |
| LOG-03 | 关键业务操作（登录 / 鉴权 / 白名单与配置变更 / 文件上传下载）必须输出审计日志，含操作者 / 动作 / 结果 | 现状缺失（仅缓存操作 1 处接入） |
| LOG-04 | 业务代码统一走 `common/logger` 封装，禁止绕过封装直接调 lager 或标准库 log | 现状部分遵守（`src/service/alarm_service.go` 内混用标准库 `log.Println` 与 `log` 包） |

## 配置

### 使用模式

三通道混用：

1. 配置文件：`src/conf/app.conf` 由 beego.AppConfig 加载，`src/common/conf/config.go` 在包级 var 初始化时用 `DefaultString` 读取 redis / node / oss 的 endpoint；
2. 命令行参数：`src/main.go` 调 `flagutil.Parse(c)`，按 `src/common/conf/config.go` 中 Config 结构体的 `flag` tag 注册命令行参数覆盖配置；
3. 环境变量：部署差异项走 `os.Getenv`（APPID / APPNAME / NODENAME / NAMESPACE / SERVICENAME / PODNAME / LOCAL_MODE / DB_SERVICE_NAME / DB_NAME）。

```go
// 来源：src/common/conf/config.go
var defaultRedisEndpoint = beego.AppConfig.DefaultString("redis::endpoint", "")
// init() 中组装 Config{Logger, Redis, Node, OSS}，Instance() 全局获取

// 来源：src/main.go
c := conf.Instance()
flagutil.Parse(c)
conf.SetDefault(c)

// 来源：src/service/alarm_service.go
alarmManager: manager.CSPInitAlarmSDK(os.Getenv(constants.EnvAppId), constants.ServiceName, os.Getenv(constants.NODENAME), manager.GetNodeIP())
```

监控相关配置（monitor.json / sql.yaml 路径）在 `src/service/monitor_service.go` 中通过 `beego.AppConfig.DefaultString("cspmonitor::monitorJsonFile", defaultMonitorFile)` 读取并带硬编码兜底默认值。

### 调用点分布

| 关注维度 | 分布（文件路径，不带行号） | 现状说明 |
|---|---|---|
| 配置读取方式 | 配置文件入口 `src/common/conf/config.go`、`src/service/monitor_service.go`；命令行入口 `src/main.go` + `src/utils/flagutil/flags.go`；环境变量 `src/service/alarm_service.go`、`src/common/logger/auditlog.go`、`src/main.go`、`src/dao/db_init.go`、`src/dao/db_local_sqlite.go`、`src/common/cse/cse.go`、`src/common/https/https_server.go` | 三通道职责未文档化；存在硬编码默认值散落（如 OSS AccessKey/SecretKey 默认 "minioadmin" 写死在 `src/common/conf/config.go` init 中） |
| 默认值处理 | `src/common/conf/config.go`、`src/service/monitor_service.go` | endpoint 类缺配默认为空字符串（静默零值，运行时才发现连不上）；monitor 文件路径有兜底默认值；`src/common/conf/config.go` 的 `SetDefault` 针对尼日局点用 `strings.Contains(endpoint, ">>")` 魔数判定回填，语义晦涩 |
| 环境变量 | `src/service/alarm_service.go`（APPID/NODENAME/NAMESPACE/SERVICENAME）、`src/common/logger/auditlog.go`（APPNAME/APPID）、`src/main.go`（LOCAL_MODE）、`src/dao/db_init.go`（DB_SERVICE_NAME/DB_NAME） | 命名大写下划线，部分收敛在 `src/common/constants`（EnvAppId/NODENAME/NAMESPACE/SERVICENAME），但 auditlog 中直接 `os.Getenv("APPNAME")` / `os.Getenv("APPID")` 字面量与常量并存，口径不一 |

### 应有约定建议

| 编号 | 约定条目 | 现状符合度 |
|---|---|---|
| CONF-01 | 配置项必须有明确默认值策略：缺配按默认值运行或启动失败并报错，禁止静默取空串零值 | 现状部分遵守（endpoint 缺配为空串静默放行） |
| CONF-02 | 环境变量仅用于部署差异项（APPID / NODENAME / NAMESPACE / LOCAL_MODE / DB_*），业务参数走 app.conf / 命令行 flag；环境变量名统一收敛到 `common/constants` 常量 | 现状部分遵守（auditlog 中字面量与常量混用） |
| CONF-03 | 密钥类配置（数据库密码 / OSS SecretKey）禁止明文入库；`src/conf/app.conf` 中 `databasepassword=Cspdbg@2017` 及 config.go 中硬编码 `minioadmin` 应改为部署期注入 | 现状缺失（明文密码已入库） |
| CONF-04 | 新增配置项必须在 Config 结构体登记 flag tag 与 desc，禁止散落 `beego.AppConfig` 散读 | 现状部分遵守（monitor_service 中散读 AppConfig） |

## 告警

### 使用模式

告警统一经 `src/service/alarm_service.go` 封装：业务方调 `SendAlarm(alarmID, msg)` / `ClearAlarm(alarmID, msg)`，事件进 channel 由 goroutine `handleEvent` 异步消费；上报前按告警 ID 做 10 分钟抑制窗口（`alarmSuppressThresholdMs`），失败重试 2 次（间隔 10s）；启动时 `CleanAllActiveAlarm` 从 FMService 拉活动告警并按本节点 sourceip 清历史告警。

```go
// 来源：src/service/alarm_service.go
const AlarmId300010 = "300010" // 告警 ID 集中常量
var AlarmList = []string{AlarmId300010}

// 来源：src/controllers/management_controller.go
c.alarm.SendAlarm(service.AlarmId300010, "Failed to sync browser configuration")
// 同步成功后恢复：
c.alarm.ClearAlarm(service.AlarmId300010, "Browser configuration synced successfully")
```

### 调用点分布

| 关注维度 | 分布（文件路径，不带行号） | 现状说明 |
|---|---|---|
| 告警 ID 使用 | 定义 `src/service/alarm_service.go`（AlarmId300010）；引用 `src/controllers/management_controller.go` | 集中常量化、引用走常量，无散落魔法数；但全仓仅 1 个告警 ID，覆盖面窄 |
| 上报与恢复配对 | `src/controllers/management_controller.go`（同一函数内失败上报 / 成功恢复）；启动清理 `src/service/alarm_service.go` CleanAllActiveAlarm | 300010 上报与恢复成对且恢复条件对称；无只报不恢复项 |

### 应有约定建议

| 编号 | 约定条目 | 现状符合度 |
|---|---|---|
| ALM-01 | 告警 ID 集中定义为常量、禁止散落魔法数，新增告警 ID 同步登记到 AlarmList 供启动清理 | 现状已遵守 |
| ALM-02 | 可恢复故障的告警上报必须配对恢复上报，恢复条件与上报条件对称 | 现状已遵守（现有唯一告警） |
| ALM-03 | 告警上报携带定位上下文（kind / namespace / sourceip / EventMessage / EventSource / OriginalEventTime） | 现状已遵守（reportAlarm 统一附加） |
| ALM-04 | 告警封装内部诊断日志用 Infof / Errorf 分级，禁止把正常流程打进 Errorf，禁止混用标准库 log.Println | 现状缺失（alarm_service.go 内部 Errorf 滥用 + log.Println 混用） |

## 监控埋点 metrics

### 使用模式

`src/service/monitor_service.go` 启动时读 `monitor.json`（路径来自 app.conf `cspmonitor::monitorJsonFile`，兜底硬编码默认值），调 CSPGoMonitorSDK 注册话统（失败 10s 重试直到成功），随后按 5 分钟窗口（`src/utils/monitorutil/time_util.go`）定时采集指标并打点上报；指标模型定义在 `src/models/monitor/metric.go`。

### 调用点分布

| 关注维度 | 分布（文件路径，不带行号） | 现状说明 |
|---|---|---|
| 指标定义 | `src/models/monitor/metric.go`、`src/conf/monitor.json` | 指标 ID / 分组走 monitor.json 模板，代码中 MocID / MetricID 类型化，未散落魔法数 |
| 采集与上报 | `src/service/monitor_service.go`、`src/utils/monitorutil/time_util.go` | 5min 定时窗口统一由 monitorutil 计算；注册失败无限重试（10s 间隔），无熔断退出路径 |

### 应有约定建议

| 编号 | 约定条目 | 现状符合度 |
|---|---|---|
| MON-01 | 监控指标的新增以 monitor.json 模板为准，代码仅按 MocID / MetricID 引用，禁止在代码中硬编码指标定义 | 现状已遵守 |
| MON-02 | 注册 / 上报失败的重试应有日志可观测，长期失败需有告警联动 | 现状部分遵守（仅 Errorf 日志，无告警联动） |

## 附注

- `src/service/alarm_service.go` 中 `logger.Errorf("GetAllActiveAlarmFromFMService  parameter is : ", jsonParam)` 等多处 Errorf 误用且格式化参数个数与占位符不匹配，待确认是否为历史遗留调试日志。
- `src/common/conf/config.go` 的 `SetDefault` 中 `strings.Contains(c.Node.HttpsEndpoint, ">>")` 判定逻辑语义不明，标注「待确认」。
- `src/stubs/` 下 lager / AlarmSDK / MonitorSDK 均为空实现桩，本仓本地行为不代表生产真实 SDK 行为。
