# Beego v2 Web 使用指导（RPC/通信-服务端）

> 版本：github.com/beego/beego/v2 v2.1.0 ｜ 调用点：~40（全部经封装层，无裸 `beego.Run()`）｜ 涉及文件：14 controllers + routers + common/https ｜ 基线：origin/main @ ae0a8a6（2026-07-29 复核）

## 用途定位

系统的唯一 Web 框架。同时承担三个监听端口的 HTTP/HTTPS 服务：
- 内部 HTTP（默认 `httpport=9090`，`src/conf/app.conf:5`）：管理面接口
- 外部 HTTPS（`httpsport`，默认 40051，`src/main.go:191-199`）：终端接入接口
- 外部 HTTP（可选，`EnableHTTP` 环境变量 + `moon::httpport`，默认 40050，`src/main.go:178-188`）

不用 Beego 的默认全局 server（从不调用 `beego.Run()`），而是通过 `beego.NewHttpServerWithCfg` 创建多实例。

## 初始化与配置

- 初始化位置：`src/main.go:153-209`（`startInternalServer` / `startExternalHttpsServer`）
- Server 封装：`src/common/https/http_server.go:36-46`（`NewHttpServer`）、`src/common/https/https_server.go:48-56`（`newBeegoHttpsServer`）
- 配置来源：`src/conf/app.conf`（`beego.BConfig` 拷贝后改 `Cfg.Listen.*`）；监听 IP 由 `https.GetLocalIP(ethEnv, defaultEth)` 从网口环境变量获取，失败回退 `127.0.0.1`（`src/common/https/https_server.go:70-85`）
- 生命周期：`Run()` 内部 `go server.Run("")` 非阻塞（`http_server.go:26-28`）；HTTPS server 等证书就绪后才真正拉起端口（`https_server.go:118-138` `monitorCertificate`）

## 核心使用模式

### 新增一个 Controller（标准骨架）

```go
// 来源：src/controllers/login_controller.go:20-44
type LoginController struct {
	BaseController                                    // 必须内嵌
	userService    service.UserService                // 依赖的 service
}

func (c *LoginController) RouteInfo() RouteInfo {     // 声明路由表
	return RouteInfo{
		RouteMapping: map[string]string{
			"/user-bind/v1/:sessionID": "GET:GetUserBind",   // "METHOD:FuncName"
		},
	}
}

func (c *LoginController) Prepare() {                 // 每请求初始化 service
	c.userService = service.NewUserService()
}
```

```go
// 来源：src/routers/beego_router.go:28-39 —— 在内部或外部路由注册函数中加一行
registerController(server, &controllers.LoginController{})
```

### 请求处理与响应

```go
// 来源：src/controllers/login_controller.go:46-60, 78-92
func (c *LoginController) GetUserBind() {
	sessionID := c.PathParameter(":sessionID")         // 路径参数
	ub, err := c.userService.GetUserBind(sessionID)
	if err == orm.ErrNoRows { c.NotFound(); return }   // 404
	if err != nil { c.InternalServiceError(); return } // 500
	c.OK(ub)                                           // 200 + JSON
}

func (c *LoginController) UpdateUserBind() {
	var request req.UpdateUserBindRequest
	err := c.RequestBodyUnmarshalTo(&request)          // body→JSON→Validate()
	if err != nil {
		c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: err.Error()})
		return
	}
	...
	c.OK(nil)                                          // nil → 默认 success 响应
}
```

## 封装层与扩展点

- **Server 封装**：`common/https.BeegoServer` 接口（`http_server.go:30-34`），统一 `Router/InsertFilter/Run`，HTTP/HTTPS 两个实现。
- **Controller 基类**：`BaseController`（`controllers/controller.go:39-146`），封装参数获取（`QueryParameter/PathParameter/Body`）、`RequestBodyUnmarshalTo`（body 解析 + `req.IRequest.Validate()` 校验）、统一响应（`OK/Failed/NotFound/InternalServiceError`）。
- **路由注册约定**：`IController.RouteInfo()` 声明式路由表 + 每 Controller 过滤器（`routers/beego_router.go:41-72`）；`Filters` map 支持 `controllers.Before/After` 两个位置，注册在 `BeforeExec/AfterExec`。
- **全局过滤器**：`OverLoadFilter`（greatwall 限流）在 `BeforeRouter` 位置注册到所有路由（`beego_router.go:19,30`，实现见 resilience-cse-gsf.md）。
- **HTTPS 证书热更新**：`cert.SubscribeCert(externalServer)`（`main.go:203`），证书更新经 `restartChan` 触发；已运行时通过 `os.Exit(3)` 退出由平台重启换新证书（`https_server.go:151-163`）。
- **新业务必须继承 `BaseController` 并走 `registerController`，禁止直接用 `beego.Router` 全局函数。**

## 并发与线程模型

Beego 每请求一个 goroutine（标准 net/http 模型）。Controller 实例每请求新建（`Prepare()` 里创建 service 是安全的，`login_controller.go:40-44`）。**禁止在 Controller 上挂可变成员状态**——跨请求共享状态一律放 service 包级单例（见 concurrency-goroutine.md）。

## 错误处理与容错

- 统一返回体 `resp.BaseResponse{Code, Message}`（`controller.go:92-118`）；错误码集中在 `common/constants/retcode`。
- 约定：客户端参数错误用 `Failed(retcode.ClientFailed)`（400）；`orm.ErrNoRows` → `NotFound()`（404）；其余内部错误 → `InternalServiceError()`（500）。证据：`login_controller.go:46-92`。
- 特殊约定（已踩坑）：login 路径拒绝返回 `retcode.ClientFailed(-2)`，event 路径拒绝返回 `retcode.AuthFailed(401)`，不可混用（见 AGENTS.md 已踩坑 #2）。

## 约定与规范

- 路由路径带版本号（`/v1/`），内部接口与外部接口分别注册到不同 server（`beego_router.go:17-39`）。
- 请求结构体放 `models/req/`，实现 `IRequest`（`Validate()`）；响应放 `models/resp/`，内嵌 `BaseResponse`。
- 所有写响应必须经 `BaseController.OK/Failed/...`，不要直接写 `c.Ctx.ResponseWriter`（`writeHeaderAndJSON` 统一了 Content-Type 与状态码，`controller.go:137-146`）。

## 已知问题与反模式

- `RouteInfo` 中同一路径的大小写变体要分别注册（`/user-bind/v1/:sessionID` 与 `:sessionId` 并存，`login_controller.go:33-34`）——新增时注意路径参数名与 `PathParameter(":xxx")` 严格一致。
- HTTPS server 在证书未上传前不监听端口（`https_server.go:118-138`），本地调试外部 HTTPS 接口需注意。
- `registerFilters` 的 `routePathPre` 参数当前恒为 `""`（`beego_router.go:49`），前缀过滤能力是死代码。

## AI 编码指南

- 新增接口：照「核心使用模式」骨架——`models/req` 加请求结构体（含 `Validate()`）→ controller 继承 `BaseController` + `RouteInfo()` 映射 + `Prepare()` 注入 service → `routers/beego_router.go` 对应注册函数加一行 `registerController`。依据：上文「封装层与扩展点」。
- 响应一律用 `c.OK/c.Failed/c.NotFound/c.InternalServiceError` + `retcode` 错误码；**禁止**直接操作 `c.Ctx.ResponseWriter`、禁止自创错误码。依据：上文「错误处理与容错」及 `controller.go:92-146`。
- **禁止**调用 `beego.Run()`、`beego.Router()` 等全局函数——多实例端口模型下全局 server 不存在，必须用注入的 `https.BeegoServer`。依据：`main.go:153-209` 三端口启动方式。
