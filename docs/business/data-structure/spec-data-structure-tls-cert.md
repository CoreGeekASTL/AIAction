# TLS 证书

> 用途：tls-cert　实例数：2　返回 [README.md](README.md)

## 1. 核心作用

承载 HTTPS/TLS 证书校验与重启——弱签名算法集合做证书校验判定的"存在性查询"，证书重启信号 channel 触发 HTTPS server 重启监听，是 TLS 安全与证书热更新的支撑。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| weakAlgorithms | TLS 弱签名算法集合（证书校验判定） | common/https/tls.go |
| restartChan | HTTPS 证书重启信号 | common/https/https_server.go |

## 3. 实例详解

- **weakAlgorithms**
  - 结构：`weakAlgorithms := map[x509.SignatureAlgorithm]bool{...}`（common/https/tls.go:205），包内字面量集合
  - 关键操作：证书校验时 `if weakAlgorithms[cert.SignatureAlgorithm]` 判定弱算法并告警
  - 使用点：common/https/tls.go:205（字面量定义）、后续判定查询
  - 并发模型：字面量初始化后只读，无并发写
- **restartChan**
  - 结构：`restartChan chan CertInfo`（common/https/https_server.go:36），缓冲 1
  - 关键操作：证书变更时投递 CertInfo，HTTPS server goroutine select 命中后重启监听
  - 使用点：common/https/https_server.go:24（创建）、:36（字段）、消费 select
  - 并发模型：单生产者 + 单消费者，缓冲 1 防阻塞

## 4. 使用模式与约定

- 包级常量集合可用 `map[K]bool` 字面量（weakAlgorithms 历史写法），新代码建议 `struct{}` value
