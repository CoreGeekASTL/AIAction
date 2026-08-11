# BeegoHttpsServer 证书监听与服务拉起 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | channel + 裸并发(goroutine) |

## 用途定位

外部 HTTPS 服务的证书驱动拉起：`Run()` 不直接监听，而是起 `monitorCertificate` goroutine 等待证书下发（`restartChan`），证书就绪后 `go b.server.Run("")` 拉起 beego；证书更新时通过 `os.Exit(3)` 重启进程换证书。从 `monitorCertificate`/`needStartServer` 实现推断。

- 代码标识符：`BeegoHttpsServer.restartChan` / `stopChan` / `monitorCertificate` / `UpdateCert`
- 定义位置：src/common/https/https_server.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/main.go（`externalServer.Run()`）；证书订阅方 src/common/cert（`cert.SubscribeCert(externalServer)`，回调中调用 `UpdateCert`，具体文件待确认）；src/main.go（`startExternalHttpsServer`）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 池化方式 | 裸 goroutine（无池化）：每 server 实例 1 个 `monitorCertificate` goroutine + 证书就绪后 1 个 beego `server.Run` goroutine | src/common/https/https_server.go |
| 容量配置 | `restartChan` buffer=1；`stopChan` 无缓冲 | src/common/https/https_server.go |
| 任务队列 | 有界 channel（buffer 1） | src/common/https/https_server.go |
| 拒绝策略 | 未设置：`UpdateCert` 直接 `restartChan <- certInfo`，buffer 满时阻塞生产者（证书订阅回调协程） | src/common/https/https_server.go |
| 隔离范围 | 每 server 实例独立 channel 对；`certInfo`、`isServerReady` 仅 monitor goroutine 读写，无锁（依赖单 goroutine 语义） | src/common/https/https_server.go |
| 关闭与等待 | `Stop()` close stopChan 使 monitor 退出；beego server 本体由 `close()` 关，main 未调用 Stop | src/common/https/https_server.go; src/main.go |
| 消息缓冲与背压 | buffer 1 允许在途 1 条证书事件，多余事件阻塞投递方 | src/common/https/https_server.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 允许裸 goroutine | 固定数量（每 server 2 个），无放大 |
| 容量基线 | buffer 1 合理，保持 | 证书事件低频 |
| 拒绝策略 | `UpdateCert` 阻塞投递建议改 select+default 并打日志（重复证书事件可合并，只需保留最新） | 阻塞证书订阅方协程可能反压证书框架 |
| 隔离 | 现状依赖 monitor 单 goroutine 语义，建议注释固化该约定 | 防未来改动引入 race |
| 命名与观测 | 证书更新走 `os.Exit(3)` 重启进程属重语义操作，建议加告警/审计日志（现状有 logger.Infof）并确认进程编排（K8s restartPolicy）依赖（待确认） | 退出码契约需与部署侧对齐 |
