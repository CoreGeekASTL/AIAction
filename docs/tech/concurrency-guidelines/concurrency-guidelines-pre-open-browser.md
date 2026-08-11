# PreOpenBrowser 预开浏览器扇出 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 裸并发(goroutine) |

## 用途定位

登录鉴权成功后向所有 Ready 状态的 BrowserGW 实例扇出预开浏览器请求（`POST /browsergw/browser/preOpen`），每个实例一个裸 goroutine，fire-and-forget（不等结果、不汇总）。从 `PreOpenBrowser` 实现推断：加速后续真实打开浏览器的时延。

- 代码标识符：`BrowserServiceImpl.PreOpenBrowser` / `instancePreOpenBrowser`
- 定义位置：src/service/browser_service.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/service/browser_service.go（`PreOpenBrowser` 中 `go instancePreOpenBrowser(...)`）；调用方在登录链路（src/controllers/login_controller.go 等，待确认具体调用点）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 池化方式 | 裸 goroutine（无池化），每次调用按 BrowserGW 实例数扇出 N 个 goroutine | src/service/browser_service.go |
| 容量配置 | 未设置（goroutine 数 = Ready 实例数，无上限控制） | src/service/browser_service.go |
| 任务队列 | 无队列 | src/service/browser_service.go |
| 拒绝策略 | 未设置（无并发上限，登录洪峰时 goroutine 数随请求量 * 实例数放大） | src/service/browser_service.go |
| 隔离范围 | 与登录链路共用 HTTP client `https.Instance()`；goroutine 内只做一次 HTTP POST，无共享可变状态 | src/service/browser_service.go |
| 关闭与等待 | 未设置（fire-and-forget，无 WaitGroup；依赖 HTTP client 自身超时，具体超时值见 https 包，此处未读取确认） | src/service/browser_service.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 建议引入有上限的 goroutine 池（如 ants）或信号量限流（如 `chan struct{}` 容量 N），限制并发扇出总数 | 登录洪峰下无界扇出可能打满 HTTP client 连接与调度器 |
| 容量基线 | 上限建议按「实例数 × 单实例并发预算」推导并外置配置 | 局点实例规模不同 |
| 拒绝策略 | 超限建议直接跳过预开（预开是优化而非必需），并计数打点 | 预开失败不影响登录主流程 |
| 隔离 | 预开请求与登录主请求共用 HTTP client，建议确认 client 连接池容量（待确认） | 避免预开挤占主链路连接 |
| 命名与观测 | 建议记录扇出 goroutine 数与失败率指标 | 无界并发的代价需要量化 |
