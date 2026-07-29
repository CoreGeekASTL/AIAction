# CSE 注册中心出站调用

CSE（ServiceComb 引擎）注册发现通过 `Go-chassis-extend`（GSF + go-chassis）SDK 完成，封装在 `common/cse/cse.go` 与 `main.go` 的 GSF 初始化流程中。

| 接口名 | 协议 | 调用位置 | 业务场景 |
|--------|------|----------|----------|
| GSF 框架初始化与实例注册 | GSF SDK | main.go | 服务接入 CSP 框架 |
| WatchMicroServiceV1（browser-gateway） | GSF Registry SDK | common/cse/cse.go | 发现 BrowserGW 实例 |
| GetAllMicroServiceInstanceInfo | GSF Registry SDK | common/cse/cse.go | 查询指定服务实例列表 |
| UpdateMicroServiceInstanceProperties | GSF Registry SDK | common/cse/cse.go | 上报 chainEndpoints |

## GSF 框架初始化与实例注册

- 协议：GSF SDK（`gsfapi.CspInit` / `gsfapi.CspStart` / `gsfapi.HealthCheckStart` / `gsfapi.RegistExitHandler`），SDK 内部与注册中心通信
- 调用位置：main.go（initGSF、registerInstance、gsfStartHandler 函数）
- 业务场景：服务启动时接入 CSP 框架，完成自身微服务实例注册、健康检查端点暴露与优雅退出回调注册
- 接口功能：CspInit 初始化框架（失败重试 360 次），CspStart 注册服务实例（支持携带 PodName Location），HealthCheckStart 开启 REST 健康检查

## WatchMicroServiceV1（browser-gateway）

- 协议：GSF Registry SDK `register.WatchMicroServiceV1`
- 调用位置：common/cse/cse.go（Init 函数，回调 browserGWNotifier.WatchServiceCallBack）
- 业务场景：持续监听 `browser-gateway` 微服务实例的增删改事件，维护本服务内存中的 BrowserGW 实例表（含健康状态、容量）
- 接口功能：接收 CREATE/UPDATE/DELETE/LIST 事件，解析实例 Properties["status"] 为 `browsergateway.ServiceInstance` 存入 sync.Map，供登录路由分配与插件下发使用

## GetAllMicroServiceInstanceInfo

- 协议：GSF Registry SDK `register.GetAllMicroServiceInstanceInfo(selfServiceID, MicroServiceKey{AppId, ServiceName, Version:"0+"})`
- 调用位置：common/cse/cse.go（GetAllMicroServiceInstanceInfo 函数）；业务调用方 dao/db_init.go（getGaussDBIP，查询 DB 服务主实例）
- 业务场景：按服务名查询某下游微服务的全部实例信息（当前用于发现 GaussDB DB 服务主节点）
- 接口功能：返回 `[]base.MicroServiceInstance`（含 Status、EndpointsList、Properties）

## UpdateMicroServiceInstanceProperties

- 协议：GSF Registry SDK `register.UpdateMicroServiceInstanceProperties`
- 调用位置：common/cse/cse.go（Report 函数，main.go 启动时调用，最多重试 5 次）
- 业务场景：服务启动后向注册中心上报本实例的对外接入地址（chainEndpoints：http/https 外部 endpoint 列表）
- 接口功能：将 `chainEndpoints` 属性写入本实例注册信息，供网关/路由层获取接入地址；失败间隔 30s 重试，超过次数直接退出进程
