# Redis 数据访问规范

| 元信息 | 值 |
|--------|-----|
| 分支 | 27.0 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-data-access-guidelines-analyze |
| 运行模式 | 起草模式 |

## 用途定位

Redis（`github.com/redis/go-redis/v9`）封装为通用 KV/Hash/Set 对象缓存客户端（src/common/storage/redis/redis.go），规划上用于缓存终端绑定等热点数据；连接配置取 app.conf `redis::endpoint` 与启动参数 `redis.db`（src/common/conf/config.go）。**当前仓内未发现任何业务调用方**（仅 src/service/browser_service.go、src/service/plugin_service.go 的日志文案提及 redis，实际走的是 DB 读写），属于预留能力，详见附注。

## 访问点分布

| 访问点 | 所在文件 | 所在函数 | 访问方式 | 业务场景 |
| --- | --- | --- | --- | --- |
| 客户端初始化 | src/common/storage/redis/redis.go | Init | 封装层自身（go-redis NewClient + Ping 健康检查） | 服务启动初始化（当前仓内未发现 Init 调用方） |
| 对象读写接口 | src/common/storage/redis/redis.go | innerClient.Get/Set/HSet/HGet/SetNx/SetWithExpiration 等 | 封装层：src/common/storage/redis/redis.go | 预留，无业务调用方 |

## 连接与客户端管理

### 现状

- 包级单例 `var client Client`，经 `Init(config conf.RedisConfig)` 创建（go-redis `NewClient`，内部连接池由 go-redis 默认管理，池参数未显式配置），`Instance()` 获取（src/common/storage/redis/redis.go）。
- 配置来源：`redis::endpoint`（app.conf）+ `RedisConfig.DB`（启动 flag，校验范围 0-15，src/common/conf/config.go）。
- 初始化后立即 `Ping` 健康检查，失败返回 error（src/common/storage/redis/redis.go）。
- 无用户名/密码配置（`Username:""`, `Password:""` 硬编码为空，src/common/storage/redis/redis.go）。
- 未发现 Init 的启动调用方（src/main.go 等启动路径未调用），当前进程内 client 保持 nil。

### 约定（起草待评审）

1. 启用前必须在启动路径（main/init）统一调用 `redis.Init` 并处理失败，业务代码一律经 `redis.Instance()` 获取单例，禁止自行 `New`。
2. 认证信息（Username/Password）不得硬编码为空串上线，应从配置注入。
3. 连接池参数（PoolSize/MinIdleConns）启用时应显式配置。

## 事务使用

### 现状

代码未使用事务（封装未暴露 MULTI/EXEC 或 pipeline，src/common/storage/redis/redis.go）。

### 约定（起草待评审）

1. 多命令需原子性时优先用 pipeline/事务封装，启用前应在封装层补齐对应接口，业务层不得绕过封装直连 go-redis。

## 分页与批量

### 现状

未提供批量/pipeline 接口；Set 成员读取 `SMembers`、Hash 读取 `HGetAll`/`HKeys` 为一次性全量返回（src/common/storage/redis/redis.go），大 key 场景无游标式（SSCAN/HSCAN）接口。

### 约定（起草待评审）

1. 大集合读取应使用 SCAN 系列游标接口，启用前需在封装层补齐；禁止对大 key 使用 SMembers/HGetAll 全量拉取。

## SQL 拼接与注入防护

### 现状

本中间件不涉及。

### 约定（起草待评审）

本中间件不涉及。

## 缓存读写模式

### 现状

- 对象模型：`Object` 接口要求 `GetKey() + BinaryMarshaler/BinaryUnmarshaler`，key 由业务对象自带，封装层无统一 key 命名约束（src/common/storage/redis/redis.go）。
- TTL：`Set` 默认不过期（过期时间 0）；`SetWithExpiration`/`SetNx` 由调用方传入 `time.Duration`（src/common/storage/redis/redis.go）。
- 穿透/击穿/雪崩处理：代码未实现（无空值缓存、互斥重建、TTL 抖动逻辑）。
- 读写策略：无业务调用方，读写策略未落地（框架上是裸 KV 读写，未内置 Cache-Aside 回源逻辑）。

### 约定（起草待评审）

1. key 统一 `{业务域}:{实体}:{id}` 格式，由业务对象的 `GetKey` 实现，禁止随意拼接。
2. 所有写入必须显式 TTL（优先 `SetWithExpiration`），并加随机抖动防雪崩；禁止用 `Set` 写永久 key，除非有显式淘汰设计。
3. 启用缓存的业务路径须采用 Cache-Aside（先查缓存，miss 回源后回写），并对空结果短 TTL 回写防穿透。

## 错误处理

### 现状

- `redis.Nil` 统一翻译为 `storage.ErrNotExist`（wrapErr，src/common/storage/redis/redis.go、src/common/storage/error.go），区分"不存在"与系统错误。
- 其余错误原样上抛；无重试、无降级逻辑。
- client 为 nil 时调用会 panic，无 nil 防护（src/common/storage/redis/redis.go Instance 直接返回包级变量）。

### 约定（起草待评审）

1. 保留 `ErrNotExist` 语义，调用方必须将"缓存 miss"与"Redis 故障"分开处理：miss 回源，故障按业务决定降级（通常直接回源 DB 并记日志）。
2. `Instance()` 使用前须有 nil 防护或在 Init 失败时拒绝启动，禁止裸调用。
3. 读失败可降级回源；写缓存失败不影响主流程（记日志即可）。

## 附注

- **预留死代码**：src/common/storage/redis/redis.go 整套 Client 封装（含 `InitForTest`）在当前仓内无任何业务调用方，`redis.Init` 亦无启动调用点；src/service/browser_service.go、src/service/plugin_service.go 中的 "redis" 字样仅为日志文案，实际数据读写走 src/dao/ 链路。启用前需先接启动初始化并接入业务路径。

<!-- 规范未建，本次为现状盘点与约定起草。 -->
