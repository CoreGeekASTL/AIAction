# external-call-csp-omf-sdk

> 下游服务：CSP OMF 平台基础 SDK 集合（源码桩在 src/stubs/，真实实现由平台提供）。均为进程内 SDK 调用，仅在 main.go 启动流程中初始化一次，无业务期接口级调用。
> 各 SDK 的目标平台组件（ModuleKeeper 保活服务、运行日志服务、传输组件、NTP 时间源）由平台部署决定，本仓不可见具体地址，目标服务归属依据包名推断。

| 接口名 | 协议 | 调用位置 | 业务场景一句话 |
|--------|------|----------|----------------|
| modulekeeperapi.Init / AddServiceName | 平台 SDK | main.go | 接入 ModuleKeeper 进程保活上报 |
| runlogapi.InitServer | 平台 SDK | main.go | 初始化运行日志服务 |
| transportapi.Init | 平台 SDK | main.go | 初始化 OMF 传输组件 |
| ntp.Init | 平台 SDK | main.go | 初始化 NTP 时间同步 |

## SDK

## modulekeeperapi.Init / AddServiceName

- 协议：平台 SDK（`CSPGSOMF/ModulekeeperSDK/modulekeeperapi`，import 别名 `report`）
- 调用位置：main.go（`runLogInit` 函数：`report.Init()` + `report.AddServiceName("gids")`）
- 业务场景：进程启动时接入 ModuleKeeper（进程保活/状态上报组件，chassis.yaml references 亦声明 `ModuleKeeper`），注册本服务名用于保活心跳与状态上报
- 接口功能：SDK 初始化 + 服务名登记；具体上报协议由平台 SDK 内部实现

## runlogapi.InitServer

- 协议：平台 SDK（`CSPGSOMF/RunlogSDK/runlogapi`，import 别名 `RunlogSDK`）
- 调用位置：main.go（`runLogInit` 函数）
- 业务场景：进程启动时初始化运行日志服务端组件，接入平台运行日志采集链路
- 接口功能：SDK 初始化；无参数无返回

## transportapi.Init

- 协议：平台 SDK（`CSPGSOMF/TransportSDK/transportapi`）
- 调用位置：main.go（`main` 函数，GSF 初始化完成之后）
- 业务场景：初始化 OMF 传输组件（日志/数据上行传输通道）；业务场景细节待确认（代码中仅初始化调用，无后续业务调用点）
- 接口功能：SDK 初始化；无参数无返回

## ntp.Init

- 协议：平台 SDK（`CSPNTP_SDK_GO/api`，import 别名 `ntp`）
- 调用位置：main.go（`main` 函数）
- 业务场景：初始化 NTP 时间同步客户端，保证本机时间与平台时间源一致（告警 OriginalEventTime、话统时间窗等依赖准确时间）
- 接口功能：SDK 初始化；无参数无返回
