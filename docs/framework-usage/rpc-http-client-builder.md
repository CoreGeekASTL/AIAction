# 自研 HTTP 客户端 Builder 使用指导（RPC/通信-客户端）

> 版本：自研封装（基于 net/http）+ GSF RestInvoker（内部 SDK stub）｜ 调用点：~10（builder 5 处 / GSF invoker 2 处 / 裸 http.Client 1 处）｜ 涉及文件：6 ｜ 基线：main @ 5e78a48

## 用途定位

所有出站 HTTP 调用的统一封装，位于 `src/common/https/`。分三个预置 client 实例，按对端选择：

| 实例 | 初始化 | 用途 | 证据 |
| --- | --- | --- | --- |
| `https.Instance()` | `https.Init()`（`client.go:26-37`） | 普通外部 HTTP，120s 超时 | `remote_service.go:21` |
| `https.MuenInstance()` | `InitMuenClient()`（`client.go:67-70`） | 沐恩（外部云）HTTPS，带外部 CA 证书 | `remote_service.go:25` |
| `https.InnerInstance()` | `initInnerClient()`（`client.go:49-60`） | CSP 内部通信，TLS 配置来自 GSF `tlsutil.GetTLSConfig("registry",...)` | `db_init.go:355-356` |

微服务间（CSE 内部）调用还有第二条路：GSF `RestInvoker` + `cse://` 协议（见下「模式二」）。

## 初始化与配置

- `https.InitMuenClient()` 在 GSF 初始化前调用（`main.go:56`）；`https.Init()` 在 `initGSF()` 之后（`main.go:59`）——顺序不能颠倒（`client.go:48` 注释"需要在csp框架初始化完成之后"）。
- 超时常量：`timeout=120s`、`idleTimeout=300s`、总超时 `2×120s`（`client.go:22-24`）。
- 证书更新：`https.MuenCertUpdate(info)` 重建 muenClient（`client.go:72-74`）。

## 核心使用模式

### 模式一：Builder（新代码标准用法）

```go
// 来源：src/service/remote_service.go:30-51
response := https.NewRequest(client).WithRetry(defaultRetryCount).
	Method("POST").
	URL(url).
	ParamFromInterface(request).   // 结构体→JSON body；另有 Param/Params(KeyValuePair)、ParamFromReader(IOReader)
	Complete().Do()
if response.Error() != nil { ... }
if !response.IsSuccessCode() { ... }        // 2xx 判定
err := response.ResponseToStruct(authResponse)  // 自动 defer CloseResponseBody
```

```go
// 来源：src/dao/db_init.go:355-362 —— GET + query 参数 + 手动读 body
response := https.NewRequest(instance).URL(dbUrl).Method("GET").Complete().Do()
body, err := io.ReadAll(response.ResponseBody())
```

### 模式二：CSE 微服务间调用（`cse://` 协议）

```go
// 来源：src/service/alarm_service.go:338-367
request, err := rest.NewRequest(method, "cse://"+microServiceName+"/"+path, body)
response, err := core.NewRestInvoker().ContextDo(context.TODO(), request)
if response.GetStatusCode() == RespOK {
	bodyStr = string(response.ReadBody())
}
```

```go
// 来源：src/common/logger/auditlog.go:111 —— GSF 版本 invoker
resp, err := gsfapi.NewCspRestInvoker().Invoke(http.MethodPost, requestURL, headers, bs2)
```

## 封装层与扩展点

- Builder 接口：`common/https/builder.go:53-66`，链式 `Method/URL/Header/Params/WithRetry/Complete/Do`。
- 注入点：`NewRequest(client HTTPDoer)` 接受任意 `Do(req)` 实现，测试可注入 mock（`builder.go:82-90`）。
- Response 封装：`StatusCode/Error/ResponseToStruct/ResponseToWriter/IsSuccessCode`（`builder.go:392-400`）。
- **新业务发 HTTP 必须用 `https.NewRequest(...)` builder，禁止裸 `http.Client`。**

## 并发与线程模型

三个 client 是包级变量，初始化后只读（`MuenCertUpdate` 除外），`http.Client` 本身协程安全。Builder 实例不可复用（内部累积 params），一次请求新建一个。

## 错误处理与容错

- 重试：`WithRetry(times)`，指数退避 `2s→4s→…→60s` 封顶（`builder.go:253-259`）。
- 可重试条件：网络瞬断（timeout/ECONNREFUSED/ECONNRESET/EOF，`builder.go:271-291`）或状态码 429/502/503/504（`builder.go:43-45`）。
- body 重试前经 `resetBody` 重置（`builder.go:244-251`），依赖 `http.Request.GetBody`，仅对 builder 构造的 body 有效。
- 响应 body 关闭：`ResponseToStruct/ResponseToWriter` 自动关闭；手动 `ResponseBody()` 读取时需自行关闭或调 `CloseResponseBody`（`builder.go:448-456`）。

## 约定与规范

- `AGENTS.md` 质量基线明文规定：HTTP 请求优先 `https.NewRequest().WithRetry()` builder（搜索 `NewRequest`）。
- GET 请求的 `Param` 值必须是 string，否则只记日志静默丢失参数（`builder.go:338-344`）——传参前自行格式化。
- 内部 CSP 服务间调用优先模式二（`cse://` + 服务名），由服务发现解析地址；外部固定地址用模式一。

## 已知问题与反模式

- **裸 http.Client 个例**：`src/service/cache_service.go:49-51` 自建 `http.Client{Timeout:5s}` 绕开封装层——禁模仿区，新代码不得照抄。
- `response.Error()` 非 nil 时 `StatusCode()` 返回 0（`builder.go:218,238-241` 返回空 `http.Response{}`），判错顺序必须先 `Error()` 再 `IsSuccessCode()`。
- `builder.go:379` 硬编码 `User-Agent: cos/v0.0.0`、`Content-Type: application/json`，发非 JSON 请求需 `Header()` 覆盖。

## AI 编码指南

- 新增出站 HTTP 调用：`https.NewRequest(https.Instance()|MuenInstance()|InnerInstance()).WithRetry(n).Method(...).URL(...).ParamFromInterface(req).Complete().Do()`，先判 `Error()` 再判 `IsSuccessCode()`，用 `ResponseToStruct` 解析。依据：上文「核心使用模式」「封装层与扩展点」。
- 调用 CSP 内部微服务：用 `cse://服务名/路径` + `core.NewRestInvoker().ContextDo(context.TODO(), req)`，`context.TODO()` 不传 nil（AGENTS.md 基线）。依据：`alarm_service.go:338-367`。
- **禁止** `&http.Client{}` 裸建客户端、**禁止**在 log 中打印带密码的完整 URL/连接串（`db_init.go:354,376` 为反面教材）。依据：上文「已知问题与反模式」。
