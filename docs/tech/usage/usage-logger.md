# 自研 logger 封装 + auditlog 使用现状（日志）

## 用途定位
运行日志统一走 `GIDS/common/logger` 包级函数（Infof/Warnf/Debugf/Errorf/Fatalf/TeeErrorf），底层是 Go-chassis-extend 的 `lager.Logger`（本地开发时该依赖被 replace 到 `src/stubs/Go-chassis-extend`，为空实现，故本地控制台无输出）。业务代码不直接引用 lager。`TeeErrorf` 同时打日志并返回 error，用于"记录并上抛"场景。

事件/审计日志走另一条链路：`src/common/event/local_storage.go` 用 `code.huawei.com/fusionstage/auditlog`（本地同样 replace 为 stub）创建 `auditlog.Logger`，sink 为文件或 stdout，按 0440 权限落盘（`src/utils/fileutil/fileutil.go` 使用 `auditlog.PERMISSION0440`/`DefaultUserID`）。

## 使用模式

```go
// 来源：src/common/logger/logger.go
func Infof(format string, args ...interface{})  { lager.Logger.Infof(format, args...) }
func Errorf(format string, args ...interface{}) { lager.Logger.Errorf(nil, format, args...) }

func TeeErrorf(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args)
	lager.Logger.Errorf(nil, format, args...)
	return err
}
```

```go
// 来源：src/common/event/local_storage.go
engine: auditlog.NewLoggerBase(constants.ComponentName)
sink := auditlog.NewWriterSink(file)
```

约定：全仓日志调用点均为 `logger.Xxxf`，日志格式串用 `%v/%s` 占位；事件落盘文件权限统一用 auditlog 包常量。
