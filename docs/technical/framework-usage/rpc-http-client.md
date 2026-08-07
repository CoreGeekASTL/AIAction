# 自研 HTTP Client（Builder 封装）使用指导（RPC/通信）

## 用途定位
全部出站 HTTP 调用的统一封装：调用 BrowserGW、沐恩云端、DB Service 等下游。屏蔽重试、超时、TLS、body 序列化。


## 使用模式

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
