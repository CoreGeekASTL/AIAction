# Beego v2 Web 服务端使用指导（RPC/通信）

> 版本：github.com/beego/beego/v2 v2.1.0 ｜ 调用点：~30（封装层为主）｜ 涉及文件：controllers/*、routers/*、common/https/*_server.go ｜ 基线：main (6c93561)

## 用途定位
对外/对内 HTTP(S) 服务端框架。所有 REST 接口经 Beego `HttpServer` 暴露：内部服务（127.0.0.1/网口:9090）与外部 HTTPS 服务（:40051）共用一套 Controller，按路由注册范围区分（`src/main.go:153`、`src/main.go:172`）。

## 初始化与配置
- 不使用 Beego 默认全局 Server，而是 `beego.NewHttpServerWithCfg` 复制 `beego.BConfig` 建独立实例（`src/common/https/http_server.go:36`、`src/common/https/https_server.go:48`）。
- 监听地址：`https.GetLocalIP(env, defaultEth)` 从网口取 IP，失败回退 127.0.0.1（`src/common/https/https_server.go:70`）。
- 端口：`beego.AppConfig.Int("httpport"/"httpsport")`，conf/app.conf 中 httpport=9090、httpsport={tls_port}（`src/main.go:158,191`）。
- 生命周期：`Run()` 内部 `go server.Run("")` 异步启动（`src/common/https/http_server.go:26`）；HTTPS server 等证书上传后经 `restartChan` 触发启动，证书更新时 `os.Exit(3)` 重启进程（`src/common/https/https_server.go:151`）。

## 核心使用模式

新增一个接口的骨架（来源：`src/controllers/login_controller.go:19-47`、`src/routers/beego_router.go:41`）：

```go
// 1. Controller 嵌入 BaseController，实现 RouteInfo()
type LoginController struct {
	BaseController
	userService service.UserService
}

func (c *LoginController) RouteInfo() RouteInfo {
	return RouteInfo{
		RouteMapping: map[string]string{
			"/app-api/devicetcp/app/login/v1/gridLoginAuth": "POST:GridLoginAuth",
			"/user-bind/v1/:sessionID":                      "GET:GetUserBind",
		},
	}
}

// 2. Prepare() 中注入依赖（Beego 每请求新建 Controller）
func (c *LoginController) Prepare() {
	c.userService = service.NewUserService()
}

// 3. 处理方法：参数读取 → 业务调用 → c.OK()/c.Failed()/c.NotFound()
func (c *LoginController) GetUserBind() {
	sessionID := c.PathParameter(":sessionID")
	...
	c.OK(ub)
}

// 4. routers 注册：registerController(server, &controllers.LoginController{})
```

- 路由注册统一走 `registerController`（`src/routers/beego_router.go:41`）：读取 `RouteInfo().RouteMapping` 注册路由，再注册 `Filters`（Before→BeforeExec / After→AfterExec）。
- 响应统一走 BaseController 封装：`OK(data)`(200)、`Failed(resp)`(400)、`NotFound()`(404)、`InternalServiceError()`(500)（`src/controllers/controller.go:92-134`）。
- 请求体解析：`RequestBodyUnmarshalTo(param)` 完成读 body → json.Unmarshal → `param.Validate()`（`src/controllers/controller.go:71`）。

## 封装层与扩展点
- 服务器封装：`https.BeegoServer` 接口（Run/Router/InsertFilter 链式返回），实现 `BeegoHttpServer`/`BeegoHttpsServer`（`src/common/https/http_server.go:30`）。
- Controller 基类：`controllers.BaseController` 隐藏 `c.Ctx` 细节，暴露 `QueryParameter/PathParameter/Body/AddHeader/OK/Failed`（`src/controllers/controller.go:39`）。
- 扩展点：过滤器 `RouteInfo.Filters`（map[FilterAction]beego.FilterFunc）；全局限流过滤器 `OverLoadFilter` 在 BeforeRouter 插入（`src/routers/beego_router.go:19,30`）。
- **新业务代码必须继承 BaseController 并在 routers 注册，禁止直接用 beego.Router/裸 beego.Controller**。

## 并发与线程模型
- Beego 每请求独立 goroutine + 新建 Controller 实例，因此依赖在 `Prepare()` 注入而非结构体字段初始化（证据：`src/controllers/login_controller.go:41`）。
- Controller 方法内禁止写包级共享变量而不加锁；跨请求共享状态在 service 层单例内管理。

## 错误处理与容错
- 统一返回码：`retcode.Success=200 / InternalFailed=-1 / ClientFailed=-2 / AuthFailed=401`（`src/common/constants/retcode/retcode.go`）。
- 业务失败返回 200/400 + JSON body 内的 code，而非 HTTP 错误码（`src/controllers/controller.go:92,112`）。
- 登录链路拒绝用 `retcode.ClientFailed(-2)`、事件链路拒绝用 `retcode.AuthFailed(401)`，两条链路必须区分（项目坑记录 #2）。

## 约定与规范
- Controller 放 `src/controllers/`，命名 `XxxController`，嵌入 `BaseController`。
- 请求结构体放 `src/models/req/` 实现 `IRequest`（含 `Validate()`），响应放 `src/models/resp/`。
- 路由注册只改 `src/routers/beego_router.go` 的 `RegisterInternalRouter`/`RegisterExternalRouter`。
- 内部与外部接口靠注册列表区分，同一 Controller 可同时注册到两个 server（如 CacheController，`src/routers/beego_router.go:21,31`）。

## 已知问题与反模式
- HTTPS server 证书更新直接 `os.Exit(restartExitCode)` 依赖外部守护拉起（`src/common/https/https_server.go:159`），新代码不要模仿此重启方式。
- `GetLocalIP` 失败静默回退 127.0.0.1 且打 Errorf（`src/common/https/https_server.go:81`），生产排查时注意误判。

## AI 编码指南
- 新增 REST 接口：建 `XxxController{BaseController}` + `RouteInfo()` + `Prepare()` 注入 service + 处理方法用 `c.OK/c.Failed`，最后在 `routers/beego_router.go` 注册。**禁止**直接调 `beego.Router` 或裸 `c.Ctx.ResponseWriter.Write`（依据：上文「封装层与扩展点」，`src/routers/beego_router.go:41`）。
- 新请求结构体必须实现 `req.IRequest.Validate()`，解析一律走 `c.RequestBodyUnmarshalTo`（依据：`src/controllers/controller.go:71-90`）。
- 返回码必须来自 `retcode` 常量包，login/event 链路拒绝码分别为 -2/401，禁止自定义散落数字（依据：`src/common/constants/retcode/retcode.go`）。
