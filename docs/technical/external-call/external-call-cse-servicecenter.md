# external-call-cse-servicecenter

> 下游服务：CSE Service Center（ServiceComb 注册中心，地址 `https://cse-service-center.manage:30100`，见 src/conf/chassis.yaml `cse.service.registry`）。
> 调用方式：go-chassis / GSF 框架 registry 接口（`Go-chassis-extend/api/GSF/api`），封装在 common/cse/cse.go 与 main.go 初始化流程中。

| 接口名 | 协议 | 调用位置 | 业务场景一句话 |
|--------|------|----------|----------------|
| 服务注册与启动（CspInit/CspStart） | HTTPS（registry） | main.go | GSF 框架初始化并把本服务注册进 CSE |
| WatchMicroServiceV1(browser-gateway) | HTTPS（registry） | common/cse/cse.go | 订阅浏览器网关实例变更事件 |
| GetAllMicroServiceInstanceInfo | HTTPS（registry） | common/cse/cse.go、dao/db_init.go | 按服务名查询微服务实例列表 |
| UpdateMicroServiceInstanceProperties | HTTPS（registry） | common/cse/cse.go | 上报本实例 chainEndpoints 属性 |

## 服务注册与启动（gsfapi.CspInit / CspStart）

- 协议：HTTPS（registry 客户端，走 chassis.yaml 配置的 servicecenter 地址，内部 TLS 配置名 `registry`）
- 调用位置：main.go（`initGSF` → `gsfapi.CspInit` 重试最多 360 次；`registerInstance` → `gsfStartHandler` → `gsfapi.CspStart`，携带 podName location 信息）
- 业务场景：进程启动时初始化 CSP 微服务框架，完成自身服务注册、健康检查（`HealthCheckStart`）与优雅退出回调注册（`RegistExitHandler`）
- 接口功能：框架级出站交互，失败重试后仍失败则 `Fatalf` 退出进程

## WatchMicroServiceV1(browser-gateway)

- 协议：HTTPS（registry watch）
- 调用位置：common/cse/cse.go（`Init` → `register.WatchMicroServiceV1`；事件回调 `browserGWNotifier.WatchServiceCallBack`）
- 业务场景：持续订阅 `browser-gateway` 微服务实例的 CREATE/UPDATE/DELETE/LIST 事件，维护本地 `browserGWInstances` 缓存（从实例 properties.status 反序列化出 `browsergateway.ServiceInstance`），为登录分配实例、预开浏览器、插件加载、缓存删除提供目标实例列表
- 接口功能：事件驱动的实例发现；实例健康状态取 properties.isHealthy，异常解析默认视为健康

## GetAllMicroServiceInstanceInfo

- 协议：HTTPS（registry 查询）
- 调用位置：common/cse/cse.go（`GetAllMicroServiceInstanceInfo`）；业务调用方：dao/db_init.go（`getGaussDBIP`，查询 `SbgGaussDB` 服务实例）
- 业务场景：数据库连接初始化与主备切换——按服务名查询全部实例，过滤 Status=UP 且 properties.status="M"（主库）的实例，从 EndpointsList 解析出 GaussDB 主库 IP:port
- 接口功能：入参服务名+AppId+版本 `0+`；返回实例列表，查询失败返回错误（上层重试）

## UpdateMicroServiceInstanceProperties

- 协议：HTTPS（registry 属性更新）
- 调用位置：common/cse/cse.go（`Report` 函数，main.go 启动末尾 `cse.NewCse().Report(maxRetryTimes)`）
- 业务场景：外部 HTTP/HTTPS server 启动后，把本实例对外接入地址（`chainEndpoints`，如 `http://<ip>:40050`、`https://<ip>:<tls_port>`）写入本实例的 properties，供上游网关/客户端发现接入地址
- 接口功能：更新实例属性 `chainEndpoints`；失败间隔 30s 递归重试，超过 maxRetry 次 `Fatalf` 退出
