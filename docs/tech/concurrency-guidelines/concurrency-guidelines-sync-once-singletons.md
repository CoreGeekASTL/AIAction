# sync.Once 单例初始化组 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 锁（sync.Once） |

## 用途定位

仓内 sync.Once 单例初始化点归集：

| Once 实例 | 初始化内容 | 定义位置 |
|---|---|---|
| `once`（service 包级） | `initEventStorageFactory`：初始化事件存储工厂并注册 local storage | src/service/event_service.go |
| `sqliteDriverOnce` | LOCAL_MODE 下 SQLite 驱动注册（单次） | src/dao/db_local_sqlite.go |
| `authServiceOnce` | AuthService 单例（含 whiteListDao 与 authCache 装配） | src/service/auth_service.go |
| `whiteListManageServiceOnce` | WhiteListManageService 单例 | src/service/whitelist_manage_service.go |
| `authCacheOnce` | authCache 进程内鉴权缓存单例 | src/service/auth_cache.go |

- 代码标识符：`once.Do(initEventStorageFactory)` / `sqliteDriverOnce.Do` / `authServiceOnce.Do` / `whiteListManageServiceOnce.Do` / `authCacheOnce.Do`
- 定义位置：src/service/event_service.go; src/dao/db_local_sqlite.go; src/service/auth_service.go; src/service/whitelist_manage_service.go; src/service/auth_cache.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/service/event_service.go（事件服务构造路径）；src/dao/db_local_sqlite.go（本地驱动注册路径）；src/service/auth_service.go / whitelist_manage_service.go / auth_cache.go（鉴权与白名单管理服务构造路径）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 锁类型 | sync.Once（并发安全的一次性初始化原语） | src/service/event_service.go; src/dao/db_local_sqlite.go; src/service/auth_service.go; src/service/whitelist_manage_service.go; src/service/auth_cache.go |
| 锁粒度 | 每 Once 实例保护一次初始化函数执行 | 同上 |
| 锁顺序约定 | 未设置 | — |
| 临界区说明 | 初始化函数内含 DB/文件/注册操作；Once 保证并发调用方只执行一次，其余阻塞至完成 | 同上 |
| 池化方式 | 未设置 | — |
| 容量配置 | 未设置 | — |
| 任务队列 | 未设置 | — |
| 拒绝策略 | 未设置（重复调用语义为跳过，非拒绝） | — |
| 隔离范围 | 保护对象：各自的单例初始化状态 | — |
| 关闭与等待 | 未设置 | — |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | sync.Once 用法符合 Go 惯例，保持 | 与项目代码风格基线一致（单例初始化用 sync.Once 保护） |
| 容量基线 | 不适用 | — |
| 拒绝策略 | 不适用 | — |
| 隔离 | 初始化失败时 Once 不会重试，建议初始化函数内部自行处理失败（记录日志 + 后续走兜底），避免「失败后永远未初始化」 | Once 语义只执行一次，失败无重入机会 |
| 命名与观测 | 不适用 | — |
