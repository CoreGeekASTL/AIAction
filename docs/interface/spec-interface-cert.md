# 证书订阅

> 功能域：证书订阅　接口数：3（异步入口）　所属 server：异步（非 HTTP）
> 子文档 of [README.md](README.md)

## 1. 定位

通过 CertSDK 订阅外部证书推送，证书更新时由外部回调本仓 handler。异步入口，非 HTTP 路由。入口由 `cert.SubscribeCert(server)`（main.go:203）统一发起。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 调用方式 |
|---|---|---|---|
| SubscribeExCert(gids-muen) | 订阅沐恩业务证书推送 | common/cert/cert.go | SubscribeExCert("gids-muen", externalInfos, exCertInfoHandler, "/opt/csp/gids/") |
| SubscribeExCert(gids-muenCa) | 订阅沐恩 CA 证书推送 | common/cert/cert.go | SubscribeExCert("gids-muenCa", externalCaInfos, exCertInfoHandler, "/opt/csp/gids/") |
| SubscribeExCert(gids) | 订阅本仓服务证书推送 | common/cert/cert.go | SubscribeExCert("gids", serverInfos, handler, "/opt/csp/gids/") |

## 3. 数据结构说明

- **SubscribeExCert（三个订阅）**
  - 订阅参数（stubs/CSPGSOMF/CertSDK/api/base/base.go）：`appName`（string，如 "gids-muen"）；`scenes []CspExSceneInfo`（证书场景）；`handler func(certInfo []*CspExCertInfo, notifyType int) error`（证书到达回调）；`path string`（证书落盘路径 "/opt/csp/gids/"）
  - 回调数据 `CspExCertInfo`（stubs/.../base）：证书信息结构，notifyType 区分推送类型
  - 三个订阅共享 `exCertInfoHandler` 回调（common/cert/cert.go:90-98）

## 4. 风险与注意点

- **证书落盘明文路径**：common/cert/cert.go:90/94/98（`"/opt/csp/gids/"` 硬编码，证书文件落本地目录，权限未约束）
- **handler 共享**：common/cert/cert.go:90/94（gids-muen 与 gids-muenCa 共用同一 exCertInfoHandler，业务证书与 CA 证书处理逻辑未区分）
