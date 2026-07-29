# CSP 证书管理服务出站调用

证书能力通过 `CSPGSOMF/CertSDK` 接入平台证书管理服务，封装在 `common/cert/cert.go`。SDK 内部通道出站（订阅/回调模式，具体下游地址由 SDK/平台注入，待确认）。

| 接口名 | 协议 | 调用位置 | 业务场景 |
|--------|------|----------|----------|
| CspCertSDKInit | CertSDK | common/cert/cert.go | 证书 SDK 初始化 |
| UnsubscribeExCert | CertSDK | common/cert/cert.go | 清理旧订阅 |
| SubscribeExCert | CertSDK | common/cert/cert.go | 订阅证书更新 |
| GetExCertPathInfo / GetExCertPrivateKeyPwd | CertSDK | common/cert/cert.go | 获取证书文件与私钥口令 |

## CspCertSDKInit

- 协议：CSPGSOMF CertSDK `certapi.CspCertSDKInit()`
- 调用位置：common/cert/cert.go（InitCert 函数，main.go startExternalHttpsServer 中调用，失败直接 Fatalf 退出）
- 业务场景：启动外部 HTTPS 服务前初始化证书 SDK
- 接口功能：初始化证书管理通道并获取 `CSPExCertManager` 实例

## UnsubscribeExCert

- 协议：CertSDK `exCertMgr.UnsubscribeExCert(subscriber, scenes)`
- 调用位置：common/cert/cert.go（InitCertScene 函数）
- 业务场景：换包重启场景下，先清理 gids-muen / gids-muenCa / gids 三类订阅者可能残留的旧订阅，避免客户端证书更新时服务端联动复位
- 接口功能：按订阅者 + 场景（sbg_external_ca_certificate / sbg_server_ca_certificate / sbg_server_device_certificate 等）取消订阅

## SubscribeExCert

- 协议：CertSDK `exCertMgr.SubscribeExCert(subscriber, scenes, callback, certDir)`
- 调用位置：common/cert/cert.go（SubscribeCert 函数）
- 业务场景：订阅三类证书更新事件——gids-muen（外部设备证书，用于访问沐恩云端的客户端证书）、gids-muenCa（外部 CA 证书）、gids（服务端 CA + 设备证书，用于本服务外部 HTTPS 端口）
- 接口功能：证书更新时触发回调：exCertInfoHandler 更新沐恩 HTTPS 客户端证书（`https.MuenCertUpdate`），serverCertInfoHandler 更新外部 HTTPS 服务端证书（`BeegoHttpsServer.UpdateCert`，必要时重启进程生效）

## GetExCertPathInfo / GetExCertPrivateKeyPwd

- 协议：CertSDK `exCertMgr.GetExCertPathInfo(sceneName)` / `GetExCertPrivateKeyPwd(sceneName)`
- 调用位置：common/cert/cert.go（exCertInfoHandler、serverCertInfoHandler 函数）
- 业务场景：证书更新回调中读取新证书文件路径（CA/设备证书/私钥）与私钥口令
- 接口功能：返回 `CspExCertPathInfo`（ExCaFilePath/ExDeviceFilePath/ExPrivateKeyFilePath）与私钥口令，用于重建 TLS 配置
