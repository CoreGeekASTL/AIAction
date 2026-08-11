# Beego v2 使用现状（网络/事件循环——HTTP/HTTPS Server）

## 用途定位
Beego v2（`github.com/beego/beego/v2/server/web`）承担本服务全部 HTTP/HTTPS 服务端能力：内部服务（默认 9090 端口）与外部服务（HTTPS 默认 40051，可选 HTTP 40050）各创建一个独立的 `beego.HttpServer` 实例，通过 `NewHttpServerWithCfg` 复制全局配置实现多 server 并存。路由注册不用注解，而是各 Controller 实现 `IController.RouteInfo()` 返回路由表与过滤器，由 `routers.RegisterInternalRouter`/`RegisterExternalRouter` 统一注册。

调用点分布：入口 `src/main.go`；Server 封装 `src/common/https/http_server.go`、`src/common/https/https_server.go`；路由注册 `src/routers/beego_router.go`；Controller 基类 `src/controllers/controller.go`（全部 Controller 通过 `BaseController` 嵌入 `beego.Controller`）。

## 使用模式

创建 HTTPS Server（复制 BConfig，避免多实例互相污染）：

```go
// 来源：src/common/https/https_server.go
func newBeegoHttpsServer(ip string, port int) *beego.HttpServer {
	config := *beego.BConfig
	server := beego.NewHttpServerWithCfg(&config)
	server.Cfg.Listen.EnableHTTPS = true
	server.Cfg.Listen.HTTPSPort = port
	server.Cfg.Listen.EnableHTTP = false
	server.Cfg.Listen.HTTPSAddr = ip
	return server
}
```

路由与过滤器注册（Controller 自带路由表，统一遍历注册）：

```go
// 来源：src/routers/beego_router.go
func RegisterInternalRouter(server https.BeegoServer) {
	server.InsertFilter("*", beego.BeforeRouter, controllers.OverLoadFilter)
	registerController(server, &controllers.LoginController{})
	// ...
}

func registerController(server https.BeegoServer, controller controllers.IController) {
	routeInfo := controller.RouteInfo()
	for k, v := range routeInfo.RouteMapping {
		server.Router(k, controller, v)
	}
	registerFilters(server, routeInfo, "")
}
```

Controller 骨架（嵌入 BaseController，声明路由表，手写 JSON 响应）：

```go
// 来源：src/controllers/controller.go
type IController interface {
	beego.ControllerInterface
	RouteInfo() RouteInfo
}

type BaseController struct {
	beego.Controller
}

func (c *BaseController) OK(data interface{}) {
	if data == nil {
		data = resp.BaseResponse{Code: retcode.Success, Message: "success"}
	}
	err := c.writeHeaderAndJSON(http.StatusOK, data, "application/json")
	// ...
}
```

HTTPS Server 的启动被证书订阅驱动：证书未上传前只跑 `monitorCertificate` 协程，收到 `UpdateCert` 事件后才 `server.Run("")` 拉起端口；证书更新时以退出码 3 重启进程（`src/common/https/https_server.go`）。
