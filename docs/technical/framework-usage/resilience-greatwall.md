# 容错/服务治理使用指导（容错/服务治理）

## 用途定位
- **限流**：greatwall `overloadcontroller` 对全部 HTTP 入口做并发限流（APIService 维度）。
- **重试**：无统一框架，三处自研——HTTP Builder 指数退避、GSF/告警启动期固定间隔重试、调度任务重试。
- **健康检查/熔断配置**：`src/conf/circuit_breaker.yaml`、`recovery.yaml`、`load_balancing.yaml` 由 go-chassis 消费，业务代码不直接读写。


## 使用模式

限流过滤器（来源：`src/controllers/filter.go:30-49`）：

```go
func OverLoadFilter(ctx *beecontext.Context) {
	dimNameValues := map[string]string{FilterConfKey: ctx.Request.URL.Path + "/" + ctx.Input.Method()}
	isGranted, err := overloadcontroller.Process(dimNameValues)
	if !isGranted {
		ctx.ResponseWriter.Header().Add("Retry-After", "3")
		ctx.ResponseWriter.WriteHeader(http.StatusTooManyRequests)
		return
	}
}
```

HTTP 重试（见 [rpc-http-client.md](rpc-http-client.md)）：`.WithRetry(n)`，退避 2s×2^n 封顶 60s，可重试条件 = 网络瞬断/连接拒绝/EOF + 429/502/503/504（`src/common/https/builder.go:37-44,253-300`）。

启动期重试范式（来源：`src/main.go:130-144`）：

```go
for i := 0; i < InitGSFRETRYTIMES; i++ { // 360 次 × 5s
	err = gsfapi.CspInit()
	if err != nil && i == last { logger.Fatalf(...) }
	if err != nil { time.Sleep(5s); continue }
	break
}
```
