# 配置管理 使用现状（配置管理）

## 用途定位
配置三轨并存：
1. **Beego AppConfig**：`src/conf/app.conf` 为唯一 ini 配置文件，通过 `beego.AppConfig.DefaultString("section::key", default)` / `.Int` 读取，调用点散布在 main、dao、service、controllers（如 `gaussdb::*`、`moon::*`、`local::sqlitepath`、`cspmonitor::*`）。
2. **自研 conf 包**：`src/common/conf/config.go` 定义 `Config`（Logger/Redis/OSS/Node 四节），`init()` 时从 AppConfig 读默认值生成单例 `conf.Instance()`；`src/utils/flagutil` 负责命令行 flag 解析覆盖。
3. **环境变量**：优先级最高——`LOCAL_MODE`、`PODNAME`、`APPID`、`DB_SERVICE_NAME`、`DB_NAME`、`FABRIC_ETH`/`SC_TRUNK_ETH`、`EnableHTTP`、`PORT`/`TLS_PORT` 等直接 `os.Getenv` 读取（`src/main.go`、`src/dao/db_init.go`）。

其余 yaml/json 配置文件（chassis.yaml、microservice.yaml、monitor.json、sql.yaml 等）属 Go-chassis/CSP SDK 自身配置，由对应框架读取。

## 使用模式

```go
// 来源：src/dao/db_init.go
dbServiceName = os.Getenv("DB_SERVICE_NAME")
if dbServiceName == "" {
	dbServiceName = beego.AppConfig.DefaultString("gaussdb::servicename", "")
}
```

```go
// 来源：src/common/conf/config.go
var config *Config

func init() {
	config = &Config{
		Logger: LoggerConfig{LogLevel: "INFO"},
		Redis:  RedisConfig{Endpoint: defaultRedisEndpoint},
		// ...
	}
}

func Instance() *Config { return config }
```

约定：配置键写法 `section::key`；读取顺序为 环境变量 → app.conf → 代码默认值；新增配置项需在 app.conf 登记并给默认值。
