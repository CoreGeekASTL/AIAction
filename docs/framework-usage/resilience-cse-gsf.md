# GSF/CSE 服务发现与 greatwall 过载控制使用指导（容错/服务治理）

> 版本：Go-chassis-extend（内部 GSF SDK，本地 stub `src/stubs/Go-chassis-extend`）+ greatwall-sdk-go v1.9.6（stub）｜ 调用点：~15 ｜ 涉及文件：4 ｜ 基线：origin/main @ ae0a8a6（2026-07-29 复核）

## 用途定位

承担 CSP 平台的服务治理三件套：
1. **GSF 框架生命周期**：`gsfapi.CspInit/CspStart`、优雅退出回调、健康检查（`main.go:126-150`）
2. **CSE 服务注册与发现**：watch browser-gateway 实例、查询 GaussDB 主库、上报自身 chain endpoints（`common/cse/cse.go`）
3. **greatwall 过载控制**：所有入口流量先过限流过滤器（`controllers/filter.go`）

## 初始化与配置

- GSF 初始化：`main.go:126-150 initGSF()`，失败每 5s 重试 360 次后 `Fatalf`；随后 `RegistExitHandler`（退出时停调度器，`main.go:216-219`）+ `HealthCheckStart`。
- GSF 启动：`main.go:107-123 gsfStartHandler`，`PODNAME` 环境变量非空时带 `WithLocation` 参数；**`CspStart` 阻塞当前协程，必须 `go` 启动**（`main.go:118` 注释）。
- CSE 初始化：`main.go:71 cse.Init()`（`cse.go:72-90`），`appid` 取自 `APPID` 环境变量。
- 配置文件：`src/conf/chassis.yaml`、`microservice.yaml`、`circuit_breaker.yaml`、`load_balancing.yaml`、`tls.yaml`（go-chassis 标准配置，本地 stub 下不生效）。
- greatwall：`filter.go:23-28` `init()` 中 `overloadcontroller.Init()`；限流策略在 `src/conf/policy.json`。

## 核心使用模式

### 服务发现（查询依赖服务实例）

```go
// 来源：src/common/cse/cse.go:56-64
msKey := base.MicroServiceKey{AppId: c.appid, ServiceName: serviceName, Version: "0+"}
return c.register.GetAllMicroServiceInstanceInfo(config.GetSelfServiceID(), msKey)
```

```go
// 来源：src/dao/db_init.go:382-409 —— 消费方筛选 UP 且 Properties["status"]=="M" 的主库实例
instances, err := cse.NewCse().GetAllMicroServiceInstanceInfo(dbServiceName)
for _, instance := range instances {
	if instance.Status != "UP" || instance.Properties["status"] != "M" { continue }
	...
}
```

### Watch 订阅（实例变更 → 本地缓存）

```go
// 来源：src/common/cse/cse.go:84-116 —— watch browser-gateway，事件经回调写 sync.Map
cseService.register.WatchMicroServiceV1(selfServiceID, []base.MicroServiceKey{msKey}, browserGWNotifier{})
// 回调按 Action(CREATE/UPDATE/DELETE/LIST) 更新 c.browserGWInstances (sync.Map)
```

### 过载控制过滤器

```go
// 来源：src/controllers/filter.go:30-49
func OverLoadFilter(ctx *beecontext.Context) {
	dimNameValues := map[string]string{FilterConfKey: ctx.Request.URL.Path + "/" + ctx.Input.Method()}
	isGranted, err := overloadcontroller.Process(dimNameValues)
	if !isGranted {
		ctx.ResponseWriter.Header().Add("Retry-After", "3")
		ctx.ResponseWriter.WriteHeader(http.StatusTooManyRequests)
	}
}
```

### 实例属性上报

```go
// 来源：src/common/cse/cse.go:170-186 —— 上报 chainEndpoints 给 CSE，失败递归重试 maxRetry 次
c.register.UpdateMicroServiceInstanceProperties(selfServiceID, selfInstanceID,
	map[string]string{"chainEndpoints": strings.Join(c.chainEndpoints, ",")})
```

## 封装层与扩展点

- `common/cse.Cse` 接口（`cse.go:23-29`）：服务发现、browser-gateway 实例缓存、endpoint 上报。包级单例 `cseService`，`NewCse()` 返回（`cse.go:66-70`）。
- browser-gateway 实例的 `Properties["status"]` 是 JSON 字符串，反序列化为 `browsergateway.ServiceInstance`（`cse.go:118-145`）。
- **查询/订阅 CSE 一律经 `cse.NewCse()`，禁止直接 import go-chassis registry。**

## 并发与线程模型

- Watch 回调在 SDK 自己的 goroutine 触发，写 `sync.Map`（`cse.go:34,108,144`），读侧 `Range`（`cse.go:39-53,148-158`）——天然并发安全。
- `chainEndpoints` 是普通 slice，仅在启动期 append（`main.go:187,207`），无并发问题。
- `Report` 递归重试（`cse.go:184`），maxRetry≤0 时 `Fatalf`（`cse.go:177-179`）——**调用方必须给正数重试上限**（`main.go:94-95` 给 5）。

## 错误处理与容错

- GSF/CSE 初始化全部"重试至死"风格（initGSF 360 次、Report 5 次后 Fatal），平台依赖不可绕过。
- 过载控制 `Process` 出错时放行（`filter.go:37-40` 只记日志）——限流失败不阻断业务。
- 熔断配置在 `circuit_breaker.yaml`（go-chassis 标准），代码内无显式熔断调用。

## 约定与规范

- 环境变量约定：`APPID`、`PODNAME`、`NODENAME`、`NAMESPACE`、`SERVICENAME` 等由 CSP 平台注入（`cse.go:73`、`main.go:120`）。
- 外部入口（internal + external 两组路由）都挂 `OverLoadFilter`（`beego_router.go:19,30`）——新 server 也必须挂。

## 已知问题与反模式

- 本地/训战环境所有 GSF/CSE API 均为 stub 空实现（`src/stubs/Go-chassis-extend/...`），服务发现返回空——**依赖 CSE 的逻辑（GaussDB 发现、browser-gateway 路由）本地不可验证**，需 `LOCAL_MODE` 旁路（`db_init.go:237-241`）。
- `Report` 用递归实现重试（`cse.go:184`），深度由 maxRetry 限制，可接受但勿放大。
- greatwall `Init` 失败仅记 Errorf 后继续（`filter.go:24-27`），此时 `Process` 行为依赖 stub——生产环境需确认策略文件 `policy.json` 下发。

## AI 编码指南

- 新增"调用/发现某个 CSP 微服务"的需求：用 `cse.NewCse().GetAllMicroServiceInstanceInfo(serviceName)`，筛选 `Status=="UP"`，必要时按 `Properties` 过滤；调用走 `cse://服务名/路径`（见 rpc-http-client-builder.md 模式二）。依据：`cse.go:56-64`、`db_init.go:382-409`。
- 新增 HTTP server 或路由组：在注册函数首行加 `server.InsertFilter("*", beego.BeforeRouter, controllers.OverLoadFilter)`。依据：`beego_router.go:19,30`。
- **禁止**在 watch 回调里做阻塞 IO/DB 操作（回调线程是 SDK 的）；只更新内存缓存，重活投递到自己的 goroutine。依据：上文「并发与线程模型」。
