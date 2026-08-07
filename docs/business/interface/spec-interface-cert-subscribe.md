# 证书更新订阅

> 功能域：cert-subscribe　接口数：3（异步订阅）　所属 server：外部(HTTPS 证书支撑)
> 子文档 of [README.md](README.md)

## 1. 定位

通过 CSP CertSDK 订阅证书平台的证书变更通知，收到更新后热更新 Muen 客户端证书与外部 HTTPS server 的服务端证书，保障 TLS 链路证书不过期。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| 订阅 gids-muen | 订阅外部设备证书场景，更新 Muen 客户端证书 | src/common/cert/cert.go | Subscribe scene=sbg_external_device_certificate（CertSDK SubscribeExCert） |
| 订阅 gids-muenCa | 订阅外部 CA 证书场景，更新 Muen 客户端 CA | src/common/cert/cert.go | Subscribe scene=sbg_external_ca_certificate（CertSDK SubscribeExCert） |
| 订阅 gids | 订阅服务端 CA+设备证书场景，热更新 HTTPS server 证书 | src/common/cert/cert.go | Subscribe scene=sbg_server_ca_certificate / sbg_server_device_certificate（CertSDK SubscribeExCert） |

## 3. 数据结构说明

- **三个订阅**（订阅入口统一为 `base.CSPExCertManager.SubscribeExCert`，src/stubs/CSPGSOMF/CertSDK/api/base/base.go）
  - 入参 `base.CspExSceneInfo` 列表：SceneName、SceneDescCN/EN、SceneType（1=CA 证书，2=设备证书）、Feature；回调签名为 `func([]*base.CspExCertInfo, notifyType int) error`；落盘路径固定 `/opt/csp/gids/`（src/common/cert/cert.go:89-99）
  - 回调处理：`exCertInfoHandler`（外部证书，更新 `https.MuenCertUpdate` 客户端证书，src/common/cert/cert.go:106）；`serverCertInfoHandler`（服务端证书，热更新 HTTPS server，src/common/cert/cert.go:132）
  - 订阅前会先 UnsubscribeExCert 清理换包重启残留的旧订阅（src/common/cert/cert.go:67-79）

