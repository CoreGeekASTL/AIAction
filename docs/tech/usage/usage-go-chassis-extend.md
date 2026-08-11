# Go-chassis-extend（GSF + CSE）使用现状（服务注册发现/平台框架）

## 用途定位
华为内部 CSP 平台框架，本地开发时整包 replace 到 `src/stubs/`（空实现）。承担两类职责：
1. **GSF 生命周期**：`gsfapi.CspInit()`（失败重试 360 次×5s，最终 Fatal）、`gsfapi.CspStart()` 阻塞启动、`RegistExitHandler` 注册优雅退出（停调度器）、`HealthCheckStart` 健康检查。入口全在 `src/main.go`。
2. **CSE 注册发现**：`src/common/cse/cse.go` 封装为 `Cse` 接口（包级单例 `cseService`，`NewCse()` 取用）——查询服务实例（`GetAllMicroServiceInstanceInfo`，用于发现主 GaussDB）、Watch browser-gateway 实例变更（`WatchMicroServiceV1` + `browserGWNotifier` 回调，实例表存 `sync.Map`）、`AddChainEndpoint`/`Report` 上报本服务级联 endpoint。

配套平台 SDK（同在 main.go 初始化）：`CSPGSOMF`（TransportSDK/RunlogSDK/modulekeeper 告警上报）、`CSPNTP_SDK_GO`（时钟同步），均为 CSP 平台组件，本地为 stub。

## 使用模式

```go
// 来源：src/main.go
for i := 0; i < InitGSFRETRYTIMES; i++ {
	err = gsfapi.CspInit()
	if err != nil {
		time.Sleep(InitGSFSleepTime * time.Second)
		logger.Errorf("gids gsfapi.CspInit fail, err: %v", err)
	} else {
		break
	}
}
gsfapi.RegistExitHandler(&GracefulExitHandler{})
gsfapi.HealthCheckStart(gsfapibase.RestProtocal)
```

```go
// 来源：src/common/cse/cse.go
func Init() {
	cseService = cse{appid: os.Getenv("APPID"), register: api.NewRegistry(), ...}
	err := cseService.register.WatchMicroServiceV1(selfServiceID,
		[]base.MicroServiceKey{msKey}, browserGWNotifier{})
}

func (c *cse) GetAllMicroServiceInstanceInfo(serviceName string) ([]base.MicroServiceInstance, error) {
	msKey := base.MicroServiceKey{AppId: c.appid, ServiceName: serviceName, Version: "0+"}
	return c.register.GetAllMicroServiceInstanceInfo(config.GetSelfServiceID(), msKey)
}
```
