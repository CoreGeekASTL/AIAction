# external-call-csp-cert

> 下游服务：CSP 证书管理组件（平台 SDK `CSPGSOMF/CertSDK`，包 `CSPGSOMF/CertSDK/api/certapi`，源码桩在 src/stubs/CSPGSOMF/CertSDK/，真实 SDK 由平台提供）。
> 调用方式：进程内 SDK 调用；订阅回调驱动证书热更新。

| 接口名 | 协议 | 调用位置 | 业务场景一句话 |
|--------|------|----------|----------------|
| certapi.CspCertSDKInit | 平台 SDK | common/cert/cert.go | 初始化证书 SDK |
| CSPExCertManager.SubscribeExCert / UnsubscribeExCert | 平台 SDK | common/cert/cert.go | 订阅/退订外部与服务端证书场景 |

## SDK

## certapi.CspCertSDKInit

- 协议：平台 SDK
- 调用位置：common/cert/cert.go（`InitCert` 函数，main.go `startExternalHttpsServer` 中外部 HTTPS server 启动后调用）
- 业务场景：外部 HTTPS 接入前初始化证书 SDK，获取外部证书管理器单例 `GetExCertManagerInstance`；初始化失败 `Fatalf` 退出进程
- 接口功能：SDK 初始化，返回 error

## CSPExCertManager.SubscribeExCert / UnsubscribeExCert

- 协议：平台 SDK（订阅回调 `func(certInfo []*base.CspExCertInfo, notifyType int) error`；配套 `GetExCertPathInfo` / `GetExCertPrivateKeyPwd` 取证书文件路径与私钥口令）
- 调用位置：common/cert/cert.go（`SubscribeCert` 订阅三个场景组：`gids-muen` 外部设备证书、`gids-muenCa` 外部 CA 证书、`gids` 服务端 CA+设备证书；`InitCertScene` 先对三组旧订阅做 UnsubscribeExCert 清理）
- 业务场景：证书生命周期管理——订阅沐恩通信所需的外部 CA/设备证书（用于 `https.MuenInstance()` 客户端 TLS，见 external-call-muen-cloud.md）和本服务外部 HTTPS server 的服务端证书；平台推送证书更新时回调 `exCertInfoHandler` 热更新沐恩客户端 TLS 配置、回调 `serverCertInfoHandler` 热更新 HTTPS server 证书（`server.UpdateCert`），免重启换证
- 接口功能：场景名 `sbg_external_ca_certificate` / `sbg_external_device_certificate` / `sbg_server_ca_certificate` / `sbg_server_device_certificate`；证书落盘目录 `/opt/csp/gids/`
