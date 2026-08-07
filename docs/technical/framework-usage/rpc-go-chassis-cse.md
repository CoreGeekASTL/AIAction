# go-chassis / CSE 注册发现使用指导（RPC/服务发现）

## 用途定位
华为 CSE 微服务平台接入层：服务注册、实例发现（找 GaussDB 主节点、BrowserGW 实例）、实例属性上报、以及 `cse://` 协议的内部 rest 调用。框架本体是 GSF（Go-chassis-extend）+ go-chassis core。


## 使用模式

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
