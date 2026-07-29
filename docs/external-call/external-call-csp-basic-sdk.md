# CSP 平台基础组件出站调用

以下平台级 SDK 在 `main.go` / `controllers/filter.go` 初始化或调用，均属进程内 SDK 封装、SDK 内部通道与平台组件通信（具体下游地址由平台注入，待确认）。业务代码仅触发初始化/单点调用，不直接持有连接。

| 接口名 | 协议 | 调用位置 | 业务场景 |
|--------|------|----------|----------|
| TransportSDK Init | CSPGSOMF TransportSDK | main.go | 传输通道初始化 |
| RunlogSDK InitServer | CSPGSOMF RunlogSDK | main.go | 运行日志初始化 |
| ModulekeeperSDK Init / AddServiceName | CSPGSOMF ModulekeeperSDK | main.go | 进程保活注册 |
| NTP SDK Init | CSPNTP_SDK_GO | main.go | 时钟同步初始化 |
| overloadcontroller Init / Process | greatwall-sdk-go | controllers/filter.go | 接口过载流控 |

## TransportSDK Init

- 协议：CSPGSOMF TransportSDK `transportapi.Init()`
- 调用位置：main.go（main 函数）
- 业务场景：服务启动时初始化平台传输通道（GSF 初始化之后）
- 接口功能：完成 TransportSDK 初始化，供平台内部通信使用

## RunlogSDK InitServer

- 协议：CSPGSOMF RunlogSDK `RunlogSDK.InitServer()`
- 调用位置：main.go（runLogInit 函数）
- 业务场景：服务启动时初始化运行日志通道
- 接口功能：将运行日志接入平台日志体系

## ModulekeeperSDK Init / AddServiceName

- 协议：CSPGSOMF ModulekeeperSDK `report.Init()` / `report.AddServiceName("gids")`
- 调用位置：main.go（runLogInit 函数）
- 业务场景：服务启动时向 Modulekeeper 注册本进程，纳入平台保活/监控
- 接口功能：进程心跳与保活信息上报

## NTP SDK Init

- 协议：CSPNTP_SDK_GO `ntp.Init()`
- 调用位置：main.go（main 函数）
- 业务场景：服务启动时初始化 NTP 时钟同步
- 接口功能：与平台 NTP 源对时，保证节点时钟一致

## overloadcontroller Init / Process

- 协议：greatwall-sdk-go `overloadcontroller.Init()` / `overloadcontroller.Process(dimNameValues)`
- 调用位置：controllers/filter.go（init 函数与 OverLoadFilter 过滤器）
- 业务场景：每个入站 HTTP 请求经过 Beego 过滤器时，向长城过载控制组件申请配额（维度 `APIService` = 请求路径 + 方法）
- 接口功能：未获授权时返回 429（Retry-After: 3）；SDK 初始化失败仅记录日志（降级为不限流）
