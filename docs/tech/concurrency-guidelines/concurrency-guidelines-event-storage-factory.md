# eventStorageFactory 事件存储注册表锁 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 锁 |

## 用途定位

保护事件存储注册表 `factory map[string]Storage`：`Get` 按 location 取存储实例，`Register` 注册实例；读多写少（启动时注册、运行期按 location 读取）。从实现推断。

- 代码标识符：`eventStorageFactory`（内嵌 `*sync.RWMutex`）
- 定义位置：src/common/event/event_storage.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/common/event/event_storage.go（`Get`、`Register`）；初始化方 src/service/event_service.go（`initEventStorageFactory`，经 `sync.Once` 保护）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 锁类型 | sync.RWMutex（内嵌指针 `*sync.RWMutex`），Get 用 `RLock/RUnlock`，Register 用 `Lock/Unlock` | src/common/event/event_storage.go |
| 锁粒度 | 保护整个 map，临界区仅 map 读写，无 IO，持锁时间极短 | src/common/event/event_storage.go |
| 锁顺序约定 | 未设置（单锁无锁序问题） | src/common/event/event_storage.go |
| 临界区说明 | Get：map 查找；Register：map 写入；均未用 defer 解锁（临界区无 panic 风险点，但风格上建议 defer） | src/common/event/event_storage.go |
| 池化方式 | 未设置 | src/common/event/event_storage.go |
| 容量配置 | 未设置 | src/common/event/event_storage.go |
| 任务队列 | 未设置 | src/common/event/event_storage.go |
| 拒绝策略 | 未设置 | src/common/event/event_storage.go |
| 隔离范围 | 保护对象：`factory` map | src/common/event/event_storage.go |
| 关闭与等待 | 未设置 | src/common/event/event_storage.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | 不适用 | — |
| 容量基线 | 不适用 | — |
| 拒绝策略 | 不适用 | — |
| 隔离 | 现状读写锁使用正确；内嵌 `*sync.RWMutex` 指针形式建议改为值内嵌 `sync.RWMutex`，避免 nil 指针风险（`InitFactory` 未调用时 Get 会 panic） | `factoryInstance` 依赖初始化顺序 |
| 命名与观测 | 建议 Get/Register 用 defer 解锁 | 防御未来临界区扩展引入 panic 导致死锁 |
