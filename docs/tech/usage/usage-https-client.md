# 自研 https 客户端封装 使用现状（网络/通信——出站 HTTP）

## 用途定位
`src/common/https`（builder.go、client.go）是全部出站 HTTP 调用的唯一入口，基于标准库 `net/http` 自研 builder 链式封装，内置重试（指数退避 2s 起步、上限 60s，对 429/502/503/504 及瞬时网络错误重试）。业务代码禁止直接用 `http.Client`，一律 `https.NewRequest(client)`。客户端实例分内外两套：`https.InnerInstance()`（内部调用，如取 GaussDB 连接信息）、`InitMuenClient` 初始化的外部实例（证书/TLS 在 `src/common/https/tls.go`）。

调用点示例：`src/dao/db_init.go`（getGaussdbInfor）、`src/service/remote_service.go`、`src/controllers/management_controller.go`。

## 使用模式

```go
// 来源：src/dao/db_init.go
instance := https.InnerInstance()
response := https.NewRequest(instance).URL(dbUrl).Method("GET").Complete().Do()
if !response.IsSuccessCode() || response.Error() != nil {
	logger.Errorf("failed to get db info, url is:%s, status is %d, err is %v", dbUrl,
		response.StatusCode(), response.Error())
	return "", errors.New("failed to get gaussDbInfo")
}
body, err := io.ReadAll(response.ResponseBody())
```

带重试与结构体参数的 POST 骨架（builder 接口定义）：

```go
// 来源：src/common/https/builder.go
type Builder interface {
	Method(method string) Builder
	Context(ctx context.Context) Builder
	URL(url string) Builder
	Header(key, value string) Builder
	Param(key string, value interface{}) Builder
	ParamFromInterface(interface{}) Builder // 结构体转 JSON body
	ParamFromReader(io.Reader) Builder      // 流式 body
	WithRetry(times int) Builder
	Complete() BuilderCompleter
}
```

约定：GET 请求 params 拼 query；非 GET 默认把 params/struct 序列化为 JSON body；响应体必须经 `CloseResponseBody` 或使用 `ResponseToStruct`/`ResponseToWriter`（内部 defer 关闭）。`Context` 不传时默认 `context.TODO()`（builder.go buildHTTPRequest）。

限流容错：入口侧用 `controllers.OverLoadFilter`（greatwall-sdk-go overloadcontroller，见 `src/routers/beego_router.go`）；出站侧容错即本 builder 的重试机制，无独立熔断框架。
