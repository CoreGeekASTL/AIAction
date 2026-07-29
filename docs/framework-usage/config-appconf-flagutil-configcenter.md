# 配置体系使用指导（配置管理）

> 版本：Beego AppConfig(ini) + 自研 flagutil + 自研 DB 配置中心 ｜ 调用点：~25 ｜ 涉及文件：6 ｜ 基线：origin/main @ ae0a8a6（2026-07-29 复核）

## 用途定位

三层配置来源，优先级从低到高：

| 层 | 机制 | 位置 | 适用 |
| --- | --- | --- | --- |
| 1. 静态文件 | `beego.AppConfig` 读 `src/conf/app.conf`（ini） | 全仓 ~20 处 `beego.AppConfig.DefaultString("section::key", default)` | 端口、DB、endpoint 等部署期配置 |
| 2. 命令行 | `flagutil.Parse(conf.Instance())` 反射注册 flag | `utils/flagutil/flags.go`、`common/conf/config.go` | 启动期覆盖（`-redis-endpoint` 等） |
| 3. 动态配置 | DB 表 `t_config_center` + 5min 缓存刷新 | `service/config_center_service.go` | 运行期可调（moon 系列 URL 等） |

## 初始化与配置

- app.conf：Beego 启动自动加载；读取方式统一 `beego.AppConfig.DefaultString("gaussdb::servicename", "")`（`db_init.go:217`）。环境变量优先于文件的个例：`DB_SERVICE_NAME`/`DB_NAME`（`db_init.go:215-224`）、`EnableHTTP`（`main.go:178`）。
- flagutil：`main.go:52-54` `conf.Instance()` → `flagutil.Parse(c)` → `conf.SetDefault(c)`。Config 结构体字段带 `flag:"name" desc:"..."` tag，嵌套结构体以 `-` 连接（如 `redis-endpoint`，`flags.go:34-87`）。
- 配置中心：`main.go:72 service.StartRefreshConfigTask()`，每 5min 全量刷新到内存 map（`config_center_service.go:104-120,48-61`）；写路径 `InsertOrUpdateConfig` 走事务 upsert（`config_center_service.go:70-87`）。

## 核心使用模式

```go
// 静态读取（来源：src/service/monitor_service.go:77）
monitorJsonFile := beego.AppConfig.DefaultString("cspmonitor::monitorJsonFile", defaultMonitorFile)
```

```go
// 动态配置优先、app.conf 兜底（来源：src/service/remote_service.go:54-71）
cfgUrl := beego.AppConfig.DefaultString("moon::titokEndpoint", "")
config, ok := configCenter.GetConfig("moon::titokEndpoint")
if ok && config != "" { cfgUrl = config }   // DB 值覆盖文件值
```

```go
// 新增启动 flag（来源：src/common/conf/config.go:66-69）
type RedisConfig struct {
	Endpoint string `flag:"endpoint"`
	DB       int    `flag:"db" desc:"redis db number: 0-15"`
}
```

## 封装层与扩展点

- `common/conf.Config` 聚合 Logger/Redis/OSS/Node 四组启动配置，包级单例 `conf.Instance()`（`config.go:52-81`）。
- 配置中心接口：`GetConfig(key) (string, bool)`（内存缓存）、`GetFromDB(key)`（直查 DB）、`InsertOrUpdateConfig`（`config_center_service.go:19-25`）。
- 配置中心管理接口：`controllers/config_center_controller.go` + `dao/config_center.go`。

## 并发与线程模型

`configCenterServiceImpl.configs` 在 `Refresh()` 中整 map 替换（`config_center_service.go:56-60`），读侧无锁——map 替换是原子指针语义下的事实做法，但字段未加原子/锁保护（`config_center_service.go:30-34`），依赖"整 map 替换而非原地修改"保证安全；**修改 Refresh 实现时必须保持整替换**。

## 错误处理与容错

- 配置读取全部给默认值（`DefaultString/DefaultInt`），缺配置不致命。
- 配置中心刷新失败仅记日志，保留旧缓存（`config_center_service.go:52-55`）。
- `GetConfig` 返回 `(value, ok)`，调用方必须判 ok 并回退 app.conf（`remote_service.go:56-59`）。

## 约定与规范

- ini key 统一 `section::key` 双冒号引用。
- 新增运行期可调配置：走 `t_config_center`，key 复用 app.conf 同款命名（如 `moon::titokEndpoint`），读取处写"DB 覆盖文件"两段式。
- 新增部署期配置：app.conf 加项 + `DefaultString` 读取；需要命令行覆盖的再加 flag tag。

## 已知问题与反模式

- `conf.SetDefault` 里针对尼日局点的 `">>"` 特判（`config.go:83-88`）——局点定制逻辑进了通用代码，新增局点逻辑不要效仿堆在这里。
- 配置中心 `GetConfig` 直接读包级 `configCenter` 而非接收者（`config_center_service.go:64-67`），绕过了接口抽象。
- app.conf 内含明文 DB 密码与内网 IP（`conf/app.conf:14,32-39`）——注意不要外泄此文件内容。

## AI 编码指南

- 读取部署期配置：`beego.AppConfig.DefaultString("section::key", "")`；不存在的新 section 先在 `conf/app.conf` 加默认值。依据：上文「核心使用模式」。
- 需要运行期可调的参数：走配置中心"DB 覆盖 app.conf"两段式读取，key 命名与 app.conf 保持一致；**禁止**发明第四种配置来源。依据：`remote_service.go:54-77`。
- 新增命令行参数：在 `common/conf/config.go` 的 Config 结构体加字段并打 `flag`/`desc` tag，**禁止**直接调 `flag.StringVar`。依据：`flags.go:34-87`。
