# 容错/服务治理使用指导（容错/服务治理）

> 版本：greatwall-sdk-go v1.9.6（stub）+ 自研重试 ｜ 调用点：限流 1 处全局过滤器 + 多处重试 ｜ 涉及文件：5 ｜ 基线：main (6c93561)

## 用途定位
- **限流**：greatwall `overloadcontroller` 对全部 HTTP 入口做并发限流（APIService 维度）。
- **重试**：无统一框架，三处自研——HTTP Builder 指数退避、GSF/告警启动期固定间隔重试、调度任务重试。
- **健康检查/熔断配置**：`src/conf/circuit_breaker.yaml`、`recovery.yaml`、`load_balancing.yaml` 由 go-chassis 消费，业务代码不直接读写。

## 初始化与配置
- 限流：`controllers` 包 init() 调 `overloadcontroller.Init()`（`src/controllers/filter.go:23-28`）；策略文件 `src/conf/policy.json`（concurrent_limit，按 `路径/方法` 维度配置 rate_limit，测试样例见 `src/controllers/filter_test.go:17-30`）。
- 过滤器挂载：内外两个 server 都在 BeforeRouter 插 `OverLoadFilter`（`src/routers/beego_router.go:19,30`）。

## 核心使用模式

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

## 封装层与扩展点
- 限流封装在 greatwall SDK 内，业务只接触 `Init/Process` 两个函数。
- 重试扩展点：HTTP Builder `WithRetry`；其余场景手写 for+Sleep。

## 并发与线程模型
- `overloadcontroller.Process` 在 Beego 过滤器链（请求 goroutine）同步执行，内部并发计数线程安全。
- 限流拒绝直接写 ResponseWriter 并 return，不再进入 handler。

## 错误处理与容错
- 限流组件自身出错（err != nil）时**放行**（只打日志不拒绝，`src/controllers/filter.go:37-41`），属 fail-open 策略。
- DB 主备故障容错见 [storage-beego-orm.md](storage-beego-orm.md)（healthCheck 3 次失败阈值，`src/dao/db_init.go:140,158-175`）。

## 约定与规范
- 新接口自动被全局限流覆盖（过滤器挂 `*`），无需逐接口接入；特殊维度需在 policy.json 增配。

## 已知问题与反模式
- `retryAfter` 是变量却从不修改（`src/controllers/filter.go:18`）。
- 重试逻辑三处三种写法，无统一封装；新增重试场景优先复用 HTTP Builder 的 WithRetry。

## AI 编码指南
- 新增 HTTP 接口：无需关心限流，已被 `OverLoadFilter` 全局覆盖；限流阈值调整改 `src/conf/policy.json`（依据：`src/routers/beego_router.go:19`）。
- 新增出站调用重试：用 `https.NewRequest(...).WithRetry(n)`，**禁止**手写 for+Sleep 重试 HTTP 调用（依据：`src/common/https/builder.go:204`）。
- 新增启动期依赖初始化：照"固定次数+固定间隔+失败 Fatalf"范式（依据：`src/main.go:130-144`）。
