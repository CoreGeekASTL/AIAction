# authCache 鉴权缓存与 authImportLock 并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | 27.0 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 锁（sync.RWMutex）+ 进程内缓存（map） |

## 用途定位

终端鉴权的两个并发原语：

| 原语实例 | 保护内容 | 定义位置 |
|---|---|---|
| `authCache`（内嵌 sync.RWMutex） | 进程内鉴权结果缓存 entries map（读 RLock、写/清理 Lock） | src/service/auth_cache.go |
| `authImportLock`（包级 sync.RWMutex） | 白名单导入（持写锁）与鉴权回源 DB 段（持读锁）互斥 | src/service/auth_cache.go |

- 代码标识符：`authCache.get/set`（RLock/Lock）；`authImportLock.Lock/RLock`
- 定义位置：src/service/auth_cache.go
- 使用点：src/service/auth_service.go（AuthIMEI 缓存读写、回源段 RLock）；src/service/whitelist_manage_service.go（ImportIMEIList 全程 Lock）

## 容量 / 队列 / 拒绝策略现状（代码现状）

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 锁类型 | sync.RWMutex 读写锁 ×2（缓存自身一把；导入/回源互斥一把包级） | src/service/auth_cache.go |
| 锁粒度 | authCache：每次 get/set 整表级读写锁；authImportLock：导入全程持写锁、回源段持读锁 | src/service/auth_cache.go；src/service/whitelist_manage_service.go；src/service/auth_service.go |
| 锁顺序约定 | 固定顺序：先查缓存（cache 锁）→ 释放后再取 authImportLock 回源；两锁不嵌套持有 | src/service/auth_service.go |
| 临界区说明 | authCache.set 内完成写入 + 超限惰性清理（同一 Lock 内）；authImportLock 写锁覆盖清表+批量插入事务全程 | src/service/auth_cache.go；src/service/whitelist_manage_service.go |
| 池化方式 | 未设置 | — |
| 容量配置 | 缓存容量上限 authCacheCapacity=1000，超限按 expireAt 升序惰性清理最旧 authCacheEvictCount=500 条；条目 TTL authCacheTTL=30min | src/service/auth_cache.go |
| 任务队列 | 未设置 | — |
| 拒绝策略 | 缓存过期视为 miss 回源 DB；DB 异常安全优先拒绝且不缓存 | src/service/auth_service.go |
| 隔离范围 | 进程内单例（getAuthCache 经 sync.Once），跨 HTTP server 双平面共享 | src/service/auth_cache.go |
| 关闭与等待 | 未设置（进程级缓存，随进程生命周期） | — |

## 应有约定建议（建议）

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | RWMutex+map 进程内缓存符合当前规模（白名单上限 200000、缓存容量 1000），保持 | 读多写少、单进程口径，无需引入外部缓存组件 |
| 容量基线 | 容量/TTL/清理阈值集中常量管理，调整时改常量而非散落硬编码 | 与现状常量定义一致 |
| 拒绝策略 | DB 异常安全优先拒绝且不写缓存的口径须保持，禁止改为"异常放行" | 鉴权场景可用性让步于安全性 |
| 隔离 | authImportLock 为包级共享锁，新增白名单写路径必须同持写锁，禁止绕开锁直写 t_white_list | 防止清表窗口期逃生态误放行 |
| 命名与观测 | 建议补充缓存命中率/清理次数观测（代码未体现，待确认） | 容量 1000 偏小场景需可观测性支撑调优 |
