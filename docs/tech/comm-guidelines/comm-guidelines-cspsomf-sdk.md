# cspsomf-sdk 通信规范

## 接口清单

| 接口名 | 协议 | 调用位置 | 业务场景 |
|---|---|---|---|
| CspCertSDKInit / 证书场景订阅 | CSPGSOMF CertSDK | src/common/cert/cert.go | 证书订阅与 TLS 配置 |
| transportapi.Init | CSPGSOMF TransportSDK | src/main.go | 传输通道初始化 |
| RunlogSDK.InitServer | CSPGSOMF RunlogSDK | src/main.go（runLogInit） | 运行日志初始化 |
| report.Init / AddServiceName | CSPGSOMF ModulekeeperSDK | src/main.go（runLogInit） | 进程保活上报 |

## 平台 SDK

### CspCertSDKInit / 证书场景订阅

- 业务场景：服务启动阶段订阅 SBG 证书场景（sbg_external_ca_certificate / sbg_external_device_certificate / sbg_server_ca_certificate），获取证书变更通知并更新 HTTPS 服务端 TLS 配置
- 接口功能：`certapi.CspCertSDKInit()` 初始化；`GetExCertManagerInstance()` 获取证书管理器；证书变更回调 exCertInfoHandler / serverCertInfoHandler 处理 CspExCertInfo
- 调用位置：src/common/cert/cert.go（InitCert / InitCertScene / SubscribeCert）
- 协议信息：
  - 协议：CSPGSOMF CertSDK（src/stubs/CSPGSOMF/CertSDK）
  - 封装方式：平台 SDK 直连
  - 超时重试：未设置
  - 错误码处理：初始化失败 logger.Fatalf 退出进程；回调处理错误记日志

### transportapi.Init

- 业务场景：服务启动阶段初始化 CSP 传输通道 SDK（CSPGSOMF TransportSDK）
- 接口功能：`transportapi.Init()`
- 调用位置：src/main.go
- 协议信息：
  - 协议：CSPGSOMF TransportSDK
  - 封装方式：平台 SDK 直连
  - 超时重试：未设置
  - 错误码处理：未识别（SDK 内部处理）

### RunlogSDK.InitServer

- 业务场景：服务启动阶段初始化运行日志 SDK
- 接口功能：`RunlogSDK.InitServer()`
- 调用位置：src/main.go（runLogInit）
- 协议信息：
  - 协议：CSPGSOMF RunlogSDK
  - 封装方式：平台 SDK 直连
  - 超时重试：未设置
  - 错误码处理：未识别（SDK 内部处理）

### report.Init / AddServiceName（进程保活）

- 业务场景：服务启动阶段初始化 Modulekeeper 保活上报并注册本服务名 "gids"
- 接口功能：`report.Init()`；`report.AddServiceName("gids")`
- 调用位置：src/main.go（runLogInit）
- 协议信息：
  - 协议：CSPGSOMF ModulekeeperSDK（modulekeeperapi）
  - 封装方式：平台 SDK 直连
  - 超时重试：未设置
  - 错误码处理：未识别（SDK 内部处理）
