# Beego v2 Web 服务端使用指导（RPC/通信）

## 用途定位
对外/对内 HTTP(S) 服务端框架。所有 REST 接口经 Beego `HttpServer` 暴露：内部服务（127.0.0.1/网口:9090）与外部 HTTPS 服务（:40051）共用一套 Controller，按路由注册范围区分（`src/main.go:153`、`src/main.go:172`）。


## 使用模式

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
