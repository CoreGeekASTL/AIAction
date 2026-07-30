# 配置管理使用指导（配置管理）

> 版本：Beego AppConfig（ini，v2.1.0）+ 自研 flagutil + DB 配置中心 ｜ 调用点：~25 ｜ 涉及文件：12 ｜ 基线：main (6c93561)

## 用途定位
三层配置体系：
1. **静态配置**：`src/conf/app.conf`（ini），经 `beego.AppConfig` 读取（httpport、gaussdb、moon、redis、oss、local 等 section）。
2. **启动参数**：`conf.Config` 结构体 + `flagutil.Parse` 反射注册命令行 flag（`src/utils/flagutil/flags.go:13`）。
3. **动态配置中心**：DB 表 t_config_center + 内存缓存，5min 刷新（`src/service/config_center_service.go`）。

## 初始化与配置
- app.conf 由 Beego 启动时自动加载（工作目录下 conf/app.conf）。
- `conf.Instance()` 单例在包 init() 构建，默认值取自 app.conf（`src/common/conf/config.go:12-43`）；main 中 `flagutil.Parse(c)` → `conf.SetDefault(c)`（`src/main.go:52-54`）。
- 环境变量优先于 app.conf：如 `DB_SERVICE_NAME`/`DB_NAME`（`src/dao/db_init.go:215-223`）、`LOCAL_MODE`（`src/dao/db_init.go:237`）、`EnableHTTP`（`src/main.go:178`）。
- 配置中心：`StartRefreshConfigTask()` 于 `src/main.go:72` 启动，Ticker 5min 全量刷新内存 map（`src/service/config_center_service.go:104-120`）。

## 核心使用模式

静态读取（来源：`src/main.go:158,181`、`src/service/traffic_stats_service.go:119`）：

```go
port, err := beego.AppConfig.Int("httpport")
file := beego.AppConfig.DefaultString("cspmonitor::sqlYamlFile", defaultSqlYamlFile)
```

动态配置读取，DB 优先回退 app.conf（来源：`src/service/remote_service.go:54-76`）：

```go
cfgUrl := beego.AppConfig.DefaultString("moon::titokEndpoint", "")
config, ok := configCenter.GetConfig("moon::titokEndpoint")
if ok && config != "" {
	cfgUrl = config
}
```

动态配置写入（管理面接口，来源：`src/service/config_center_service.go:70-87`）：

```go
dao.DoTxWithCtx(ctx, func(ctx, txOrm) error { // 存在则 Update 否则 Insert
})
```

## 封装层与扩展点
- 启动配置封装：`GIDS/common/conf`（Config 单例 + flag tag 递归解析）。
- 动态配置封装：`service.ConfigCenterService` 接口（GetConfig/InsertOrUpdateConfig/Refresh/GetFromDB），包级单例 `NewConfigCenterService()`（`src/service/config_center_service.go:99`）。
- 扩展点：新配置项 = t_config_center 插一行 + `GetConfig(key)` 读取，无需改代码结构。

## 并发与线程模型
- `configCenterServiceImpl.configs` map 在 Refresh 整体替换（非逐 key 写），读端无锁（`src/service/config_center_service.go:56-66`），map 替换非原子存在理论竞态。
- beego.AppConfig 并发安全。

## 错误处理与容错
- 配置缺失统一用 `DefaultString/DefaultInt(key, default)` 给默认值，避免 err 分支。
- 配置中心 DB 读取失败保持旧缓存（Refresh 失败直接 return，`src/service/config_center_service.go:52-55`）。

## 约定与规范
- key 命名沿用 ini section 风格 `section::key`（如 `moon::titokEndpoint`、`cspmonitor::monitorJsonFile`）。
- 敏感配置（DB 密码）允许放 app.conf（存量如此，`src/conf/app.conf` [gaussdb]）。

## 已知问题与反模式
- 配置中心缓存替换无锁（上文），高频读写场景需注意。
- 环境变量/app.conf/DB 三来源优先级散在各调用点判断，无统一抽象（典型：`src/service/remote_service.go:54-76`）。

## AI 编码指南
- 新增静态配置：写 app.conf section + `beego.AppConfig.DefaultString("sec::key", def)` 读取；**禁止**硬编码默认值散落多处（依据：`src/common/conf/config.go:13-17`）。
- 新增可动态调整的配置：走 t_config_center + `NewConfigCenterService().GetConfig(key)`，读取顺序"DB 配置中心 > 环境变量 > app.conf"（依据：`src/service/remote_service.go:54-76`）。
- 新配置 key 必须采用 `section::key` 命名（依据：存量全部 key 风格）。
