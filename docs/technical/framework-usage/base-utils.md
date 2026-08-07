# 自研公共工具库使用指导（基础库）

## 用途定位
自研公共能力是"事实上的框架"，新代码优先复用而非新造。

| 包 | 能力 | 关键文件 |
| --- | --- | --- |
| `common/constants/retcode` | HTTP 返回码常量（Success=200/InternalFailed=-1/ClientFailed=-2/AuthFailed=401） | `src/common/constants/retcode/retcode.go` |
| `common/constants` | 全局常量（ServiceName、CleanupMonths、EnableHTTP 等） | `src/common/constants/base.go` |
| `utils/flagutil` | 结构体 flag tag 递归注册命令行参数 | `src/utils/flagutil/flags.go:13` |
| `utils/response` | 响应辅助 | `src/utils/response/response_util.go` |
| `utils/fileutil` | 文件/zip 操作（事件日志滚动删除依赖） | `src/utils/fileutil/fileutil.go`、`zip_util.go` |
| `utils/monitorutil` | 话统时间窗口计算（GetLastFiveMinuteWindow） | `src/utils/monitorutil/time_util.go` |
| `test/util` | 测试辅助 It()（goconvey 串行流程） | `src/test/util/utils.go` |
| `github.com/google/uuid` v1.6.0 | token 生成（`uuid.New().String()`） | `src/service/browser_service.go:204` |


## 使用模式

命令行参数（来源：`src/main.go:52`、`src/common/conf/config.go:52-57`）：

```go
type Config struct {
	Logger LoggerConfig `flag:"log"`   // → --log-file / --log-level / --log-event
	Redis  RedisConfig  `flag:"redis"` // → --redis-endpoint / --redis-db
}
c := conf.Instance()
flagutil.Parse(c) // 反射注册并 flag.Parse()
```

UUID token（来源：`src/service/browser_service.go:203-206`）：

```go
uid := uuid.New()
u.Token = uid.String()
```
