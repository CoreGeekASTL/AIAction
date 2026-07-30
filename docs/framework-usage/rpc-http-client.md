# 自研 HTTP Client（Builder 封装）使用指导（RPC/通信）

> 版本：net/http 标准库 + 自研封装 ｜ 调用点：封装层 5 处（service/*、dao/db_init.go）/ 裸 API 1 处 ｜ 涉及文件：8 ｜ 基线：main (6c93561)

## 用途定位
全部出站 HTTP 调用的统一封装：调用 BrowserGW、沐恩云端、DB Service 等下游。屏蔽重试、超时、TLS、body 序列化。

## 初始化与配置
- 三个进程级 client（`src/common/https/client.go:18-20`）：
  - `client`（`https.Instance()`）：普通明文/默认 TLS，`Init()` 于 `src/main.go:59`，超时 120s×2=240s、IdleConnTimeout 300s（`src/common/https/client.go:22-35`）。
  - `muenClient`（`https.MuenInstance()`）：对外沐恩通信，`InitMuenClient()` 于 `src/main.go:56`，证书更新走 `MuenCertUpdate`（`src/common/https/client.go:72`）。
  - `innerClient`（`https.InnerInstance()`）：CSE 内部通信，TLS 配置取自 GSF `tlsutil.GetTLSConfig("registry")`，须在 CspInit 之后（`src/common/https/client.go:49-60`）。
- TLS 构造：`GetTLS(info, ServerType/ClientType)`（`src/common/https/tls.go`）。

## 核心使用模式

标准用法骨架（来源：`src/service/remote_service.go:30-44`）：

```go
response := https.NewRequest(client).WithRetry(defaultRetryCount).
	Method("POST").
	URL(url).
	ParamFromInterface(request).   // struct → JSON body；Param(k,v) → kv JSON；ParamFromReader → 流
	Complete().Do()
if response.Error() != nil || !response.IsSuccessCode() {
	logger.Errorf("call failed, status %d, err %v", response.StatusCode(), response.Error())
	return nil
}
err := response.ResponseToStruct(&resp) // 内部 defer 关 body
```

- 内部 CSE HTTP 调用（绕过封装的另一套，go-chassis rest invoker）：见 [rpc-go-chassis-cse.md](rpc-go-chassis-cse.md)（`src/service/alarm_service.go:338`）。
- 响应体关闭：`ResponseToStruct`/`ResponseToWriter` 内部 `defer CloseResponseBody`（`src/common/https/builder.go:423`）；用 `ResponseBody()` 裸读时必须自行 `defer https.CloseResponseBody(resp)`。

## 封装层与扩展点
- 入口：`GIDS/common/https`，`NewRequest(client HTTPDoer) Builder`（`src/common/https/builder.go:88`）。
- 隐藏：重试（指数退避 InitialBackOff 2s→MaxBackoff 60s）、可重试错误判定（net timeout/ECONNREFUSED/ECONNRESET/EOF + 429/502/503/504，`src/common/https/builder.go:44,261-300`）、JSON 序列化、User-Agent。
- 扩展点：`HTTPDoer` 接口可注入 mock client 做测试（`src/common/https/builder_test.go`）。
- **新业务出站 HTTP 必须走 Builder，禁止裸 `http.NewRequest + client.Do`**。

## 并发与线程模型
- client 进程级单例可并发复用（http.Client 线程安全，Transport 内置连接池）。
- `Builder` 实例不可复用（每次 NewRequest 新建），`Do()` 同步阻塞。

## 错误处理与容错
- 错误双通道：`Response.Error()`（传输/构造错误）与 `StatusCode()`（业务码），两者都要判（存量代码均如此，如 `src/service/remote_service.go:34-42`）。
- 重试只对幂等场景开启；存量 `WithRetry` 用于登录鉴权等只读/幂等接口。

## 约定与规范
- GET 请求的 `Param` 值必须是 string（否则仅 Errorf 后丢弃，`src/common/https/builder.go:337-345`）。
- ctx 不传时内部补 `context.TODO()`（`src/common/https/builder.go:318-320`）。

## 已知问题与反模式
- 裸 API 个例：`src/service/cache_service.go:70-90` 直接 `http.NewRequest + client.Do`，无重试，属禁模仿区。
- `Do()` 失败后返回空 `http.Response{}`，`StatusCode()` 为 0，必须先判 `Error()`（`src/common/https/builder.go:218,238`）。
- `MuenCertUpdate` 直接替换全局 client 变量，无锁（`src/common/https/client.go:72`），并发读取存在理论竞态。

## AI 编码指南
- 新增出站 HTTP 调用：`https.NewRequest(https.Instance()).Method(...).URL(...).ParamFromInterface(req).Complete().Do()`；需重试加 `.WithRetry(n)`（仅限幂等接口）；外部 HTTPS 用 `https.MuenInstance()`，CSE 内网 TLS 用 `https.InnerInstance()`（依据：上文「初始化与配置」）。
- **禁止** `http.NewRequest + client.Do` 裸调用（依据：上文「已知问题」，`src/service/cache_service.go:78`）。
- 判定失败必须同时检查 `response.Error()` 与 `response.IsSuccessCode()`；裸读 body 必须 `defer https.CloseResponseBody(resp)`（依据：`src/common/https/builder.go:447`）。
