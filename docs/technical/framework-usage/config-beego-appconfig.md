# 配置管理使用指导（配置管理）

## 用途定位
三层配置体系：
1. **静态配置**：`src/conf/app.conf`（ini），经 `beego.AppConfig` 读取（httpport、gaussdb、moon、redis、oss、local 等 section）。
2. **启动参数**：`conf.Config` 结构体 + `flagutil.Parse` 反射注册命令行 flag（`src/utils/flagutil/flags.go:13`）。
3. **动态配置中心**：DB 表 t_config_center + 内存缓存，5min 刷新（`src/service/config_center_service.go`）。


## 使用模式

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
