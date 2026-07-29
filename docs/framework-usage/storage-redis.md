# go-redis v9 封装使用指导（存储/缓存）

> 版本：github.com/redis/go-redis/v9 v9.0.5 ｜ 调用点：0 业务调用点（仅封装层 1 文件 + 测试）｜ 涉及文件：1 ｜ 基线：origin/main @ ae0a8a6（2026-07-29 复核）

## 用途定位

`src/common/storage/redis/redis.go` 是对 go-redis 的完整封装，提供对象化（GetKey + BinaryMarshaler）的 Get/Set、Hash、Set、SetNx 等能力。**当前仓库中没有任何业务代码调用它**（grep `redis.Instance|redis.Init` 无业务命中），属预置能力/历史遗留；`conf/app.conf` 的 `[redis] endpoint` 被 `common/conf/config.go:13` 读取但同样无人消费。

## 初始化与配置

```go
// 来源：src/common/storage/redis/redis.go:24-42
func Init(config conf.RedisConfig) error {
	client = New(&redis.Options{Addr: config.Endpoint, DB: config.DB})
	err := client.Ping(context.Background())
	...
}
```

- DB 编号校验 0-15（`redis.go:19,25-28`）。
- 初始化后通过 `redis.Instance() Client` 取包级单例（`redis.go:51-53`）。
- 测试复位：`InitForTest` 返回清 client 的闭包（`redis.go:44-48`）。

## 核心使用模式

```go
// 对象必须实现 Object 接口（redis.go:56-60）：GetKey() + MarshalBinary/UnmarshalBinary
// 参考实现：src/models/db/schedule_election.go:19-27
func (s *ScheduleElection) GetKey() string { ... }
func (s *ScheduleElection) MarshalBinary() ([]byte, error) { return json.Marshal(s) }

// 读写骨架
err := redis.Instance().Set(ctx, obj)              // SET key obj
err := redis.Instance().Get(ctx, dst)              // GET + UnmarshalBinary；不存在返回 storage.ErrNotExist
ok, err := redis.Instance().SetNx(ctx, obj, ttl)   // 分布式锁/选主语义
```

## 封装层与扩展点

- `Client` 接口（`redis.go:67-85`）18 个方法，全部面向 `Object/HFieldObject`，屏蔽原生 `*redis.Client`。
- `redis.Nil` 统一翻译为 `storage.ErrNotExist`（`redis.go:117-119,227-232`）。
- **若启用 Redis，业务代码只允许经此 `Client` 接口，禁止直接 import go-redis。**

## 并发与线程模型

`client` 包级变量，`Init` 后只读；go-redis client 内建连接池、协程安全。

## 错误处理与容错

Init 时 Ping 一次，失败即返回 error（`redis.go:36-39`）；无重连、无降级逻辑——调用方需自行处理 `ErrNotExist` 与其他错误的区别。

## 约定与规范

- 缓存对象放 `models/`，实现 `GetKey/MarshalBinary/UnmarshalBinary`（+json tag），hash 对象再实现 `GetField`（`HFieldObject`，`redis.go:61-64`）。
- key 设计含业务前缀（参考 `schedule_election.go` 的选主 key）。

## 已知问题与反模式

- **最大问题：未被初始化**。main 流程从未调用 `redis.Init`（`main.go:46-98` 无此调用），`Instance()` 将返回 nil——直接使用会 panic。使用前必须自行接线 `redis.Init(conf.Instance().Redis)` 并处理失败。
- `redis.go:26` 日志格式串用 `%s` 打印 int DB 编号（类型错误，仅日志瑕疵）。

## AI 编码指南

- 新增缓存能力前，先确认真的需要 Redis（当前系统无 Redis 依赖即可运行）；确定需要时：在启动流程加 `redis.Init(conf.Instance().Redis)` + 失败降级，再经 `redis.Instance()` 使用。依据：上文「已知问题与反模式」。
- 缓存对象实现 `Object` 接口（`GetKey`+`MarshalBinary/UnmarshalBinary`），miss 判 `storage.ErrNotExist`。依据：`redis.go:56-60,117-119`。
- **禁止**绕过封装直接 `redis.NewClient`。依据：上文「封装层与扩展点」。
