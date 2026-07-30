# go-chassis / CSE 注册发现使用指导（RPC/服务发现）

> 版本：Go-chassis-extend v0.0.0-20231017151318（本地 stub 替换）｜ 调用点：9 文件 ｜ 涉及文件：common/cse/cse.go、main.go、service/alarm_service.go、dao/db_init.go、service/monitor_service.go ｜ 基线：main (6c93561)

## 用途定位
华为 CSE 微服务平台接入层：服务注册、实例发现（找 GaussDB 主节点、BrowserGW 实例）、实例属性上报、以及 `cse://` 协议的内部 rest 调用。框架本体是 GSF（Go-chassis-extend）+ go-chassis core。

## 初始化与配置
- 框架引导：`gsfapi.CspInit()` 带 360 次×5s 重试，失败 `Fatalf` 退出（`src/main.go:126-144`）；`gsfapi.CspStart()` 阻塞注册，`PODNAME` 环境变量非空时带 Location 启动（`src/main.go:107-123`）。
- 优雅退出回调：`gsfapi.RegistExitHandler(&GracefulExitHandler{})`（`src/main.go:147-148`）；健康检查：`gsfapi.HealthCheckStart(RestProtocal)`（`src/main.go:149`）。
- 配置文件：chassis.yaml / microservice.yaml / tls.yaml / lager.yaml 等在 `src/conf/`。
- 封装层初始化：`cse.Init()` 于 `src/main.go:71`，创建 Registry 并 Watch browser-gateway 实例（`src/common/cse/cse.go:72-90`）。

## 核心使用模式

实例发现（来源：`src/dao/db_init.go:382-409`）：

```go
instances, err := cse.NewCse().GetAllMicroServiceInstanceInfo(dbServiceName)
// 过滤 Status=="UP" 且 Properties["status"]=="M"（主节点）的实例，取 endpoint
```

Watch 订阅（来源：`src/common/cse/cse.go:84,96-116`）：

```go
register.WatchMicroServiceV1(selfServiceID, []base.MicroServiceKey{msKey}, browserGWNotifier{})
// 回调按 event.Action(CREATE/UPDATE/DELETE/LIST) 维护本地 sync.Map 缓存
```

cse:// rest 出站调用（来源：`src/service/alarm_service.go:338-368`）：

```go
request, err := rest.NewRequest(method, "cse://"+microServiceName+"/"+path, body)
defer request.Close()
response, err := core.NewRestInvoker().ContextDo(context.TODO(), request)
defer response.Close()
bodyStr := string(response.ReadBody())
```

实例属性上报：`cse.NewCse().Report(maxRetry)` 把 chainEndpoints 写入注册中心 properties，失败递归重试（`src/common/cse/cse.go:170-186`）。

## 封装层与扩展点
- 入口：`GIDS/common/cse`，接口 `Cse`（`src/common/cse/cse.go:23-29`），包级单例 `NewCse()`。
- 隐藏：Registry 创建、selfServiceID 获取、BrowserGW 实例本地缓存（sync.Map）。
- 扩展点：`Cse` 接口可整体 mock（单测 `src/service/browser_service_test.go` 注入 fake Cse）。

## 并发与线程模型
- Watch 回调在 go-chassis 内部 goroutine 触发，写 `sync.Map` 保证并发安全（`src/common/cse/cse.go:34`）。
- `Report` 递归重试间隔 30s，阻塞 main goroutine 属启动期行为。

## 错误处理与容错
- `CspInit` 失败重试 360 次后 `Fatalf`（进程退出）；`Report` 重试耗尽 `Fatalf`（`src/common/cse/cse.go:178`）。
- 实例查询失败返回空列表，调用方判空后 Sleep 重试（`src/dao/db_init.go:245-249`）。

## 约定与规范
- 所有 CSE 交互必须经 `cse.NewCse()` 接口，禁止业务代码直接 `api.NewRegistry()`。
- `ContextDo` 一律传 `context.TODO()`，禁止 nil（项目代码质量基线）。
- 测试环境 SDK 为 stubs/ 空实现，真实行为以生产为准；写单测时 mock `cse.Cse` 接口。

## 已知问题与反模式
- `alarmEventChanel` 接收方 `if !ok { break }` 只跳出 select 不退出 for（`src/service/alarm_service.go:77-80`），channel 关闭后空转（实际不会关闭，风险低）。
- `chainEndpoints` 普通 slice 无锁 append（`src/common/cse/cse.go:160`），仅启动期使用故无并发问题。

## AI 编码指南
- 新增服务发现/属性上报：调 `cse.NewCse().GetAllMicroServiceInstanceInfo(name)` / 扩展 `Cse` 接口；**禁止**绕过封装直接用 `api.NewRegistry()`（依据：`src/common/cse/cse.go:66-70`）。
- 新增 cse:// 内部调用：仿 `OSHttpsGetRequestByCSE` 模式，`ContextDo(context.TODO(), req)`，request/response 都必须 `defer Close()`（依据：`src/service/alarm_service.go:338-368`）。
- 新增长驻实例订阅：用 `WatchMicroServiceV1` + 回调维护 `sync.Map` 本地缓存，**禁止**每次查询都直连注册中心（依据：`src/common/cse/cse.go:84`）。
