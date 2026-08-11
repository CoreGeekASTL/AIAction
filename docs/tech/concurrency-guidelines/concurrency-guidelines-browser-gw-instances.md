# browserGWInstances（sync.Map 实例缓存）并发规范

| 元信息 | 值 |
|--------|-----|
| 分支 | new_skill_test 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-concurrency-guidelines-analyze |
| 运行模式 | 起草模式 |
| 实例类型 | 锁（sync.Map） |

## 用途定位

CSE 服务发现回调维护的 BrowserGW 就绪实例缓存：watch 回调中按事件 Store/Delete，业务侧（如 `GetAllReadyServiceInstances`、插件加载）Range 读取实例列表。读写并发高、读多写少，选型 sync.Map。从 `cse.go` 实现推断。

- 代码标识符：`cseService.browserGWInstances`（sync.Map）
- 定义位置：src/common/cse/cse.go
- 使用点（任务提交 / 加锁 / 消息投递位置）：src/common/cse/cse.go（watchServiceCallBack 中 Store/Delete/Range，实例列表读取处 Range）；读取方 src/service/browser_service.go（`GetAllReadyServiceInstances`）

## 容量 / 队列 / 拒绝策略现状（代码现状）

> 本节只写**代码事实**，逐条附证据文件路径（不带行号）；读不到写「未设置」或「框架默认」，禁止臆造取值。

| 维度 | 现状 | 证据（文件路径） |
|---|---|---|
| 锁类型 | sync.Map（无显式锁，内部读写分离） | src/common/cse/cse.go |
| 锁粒度 | 键级（实例 InstanceId 为 key）；Store/Delete/Range 均为 sync.Map 原子操作 | src/common/cse/cse.go |
| 锁顺序约定 | 未设置 | src/common/cse/cse.go |
| 临界区说明 | Range 遍历期间允许并发 Store/Delete（弱一致快照语义），读取方拿到的列表可能与最新事件存在短暂时差 | src/common/cse/cse.go |
| 池化方式 | 未设置 | src/common/cse/cse.go |
| 容量配置 | 未设置（随实例数增长） | src/common/cse/cse.go |
| 任务队列 | 未设置 | src/common/cse/cse.go |
| 拒绝策略 | 未设置 | src/common/cse/cse.go |
| 隔离范围 | 保护对象：BrowserGW 实例表；watch 回调 goroutine（CSE 框架）与业务 goroutine 并发访问 | src/common/cse/cse.go |
| 关闭与等待 | 未设置 | src/common/cse/cse.go |

## 应有约定建议（建议）

> 本节为**建议**（规范初稿，尚未在代码中落地），与上节「代码现状」严格区分；落地后相应内容转入现状节、从本节移除。

| 维度 | 建议约定 | 理由 |
|---|---|---|
| 池选型 | sync.Map 选型合理（读多写少、键集合稳定），保持 | 符合 sync.Map 适用场景 |
| 容量基线 | 不适用 | — |
| 拒绝策略 | 不适用 | — |
| 隔离 | Range 弱一致语义需在调用方注释中明示（实例列表允许短暂滞后） | 避免调用方误假设强一致 |
| 命名与观测 | 建议实例数变化时打 info 日志（现状 watchServiceCallBack 已有事件日志，保持） | 实例上下线可追踪 |
