# go-redis 使用指导（存储/缓存）

> 版本：github.com/redis/go-redis/v9 v9.0.5 ｜ 调用点：封装层 18 接口方法 ｜ 涉及文件：common/storage/redis/redis.go（当前业务代码暂无直接调用方，Init 亦未被 main 调用）｜ 基线：main (6c93561)

## 用途定位
Redis 缓存/共享状态访问的预留封装层。当前版本业务链路未启用（grep 无 `redis.Init`/`redis.Instance` 调用点），属于已就位的公共能力。

## 初始化与配置
- `redis.Init(conf.RedisConfig)`：校验 DB 号 0-15，创建 client 并 Ping（`src/common/storage/redis/redis.go:24-42`）。
- 配置：app.conf `[redis] endpoint=...`；`conf.RedisConfig{Endpoint, DB}`（`src/common/conf/config.go:66-69`）。
- 测试：`InitForTest(data)` 返回清理函数（`src/common/storage/redis/redis.go:44`）。

## 核心使用模式

对象存取约定（来源：`src/common/storage/redis/redis.go:56-64`）：

```go
// 业务对象实现 Object 接口即可整体存取
type Object interface {
	GetKey() string
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
// Hash 场景实现 HFieldObject（多一个 GetField()）
```

读写（来源：`src/common/storage/redis/redis.go:91-225`）：

```go
client := redis.Instance()
err := client.Set(ctx, obj)                            // 无过期
ok, err := client.SetNx(ctx, obj, expiration)          // 分布式锁语义
err := client.SetWithExpiration(ctx, obj, ttl)
err := client.Get(ctx, &obj)                           // 未命中返回 storage.ErrNotExist
fields, err := client.HKeys(ctx, key)
```

## 封装层与扩展点
- 入口：`GIDS/common/storage/redis`，接口 `Client` + 包级单例 `Instance()`。
- 隐藏：go-redis Cmd 错误处理、`redis.Nil` → `storage.ErrNotExist` 统一转换（`src/common/storage/redis/redis.go:227-232`、`src/common/storage/error.go`）。
- 扩展点：Object/HFieldObject 接口，新业务对象实现即可复用全部方法。

## 并发与线程模型
go-redis Client 并发安全，进程级单例复用。

## 错误处理与容错
- key 不存在统一返回 `storage.ErrNotExist`（不是 redis.Nil），调用方按此判断（`src/common/storage/redis/redis.go:114-121`）。
- 无重试/熔断配置，依赖 go-redis 默认。

## 约定与规范
- 所有访问经 `redis.Instance()`，禁止业务代码直接 `redis.NewClient`。
- key 由对象 `GetKey()` 自描述，不集中定义常量。

## 已知问题与反模式
- 当前未启用：无 Init 调用，配置 section 存在但链路断开，启用前需先在 main 接入 `redis.Init`。
- `Init` 中 Errorf 格式串用 `%s` 打印 int（`src/common/storage/redis/redis.go:26`），小瑕疵。

## AI 编码指南
- 新增缓存/分布式锁能力：业务对象实现 `Object`（+`HFieldObject`）接口，经 `redis.Instance()` 调用；启用前在 main 启动流程补 `redis.Init(conf.Instance().Redis)`（依据：上文「初始化与配置」）。
- **禁止**绕过封装直接 import go-redis（依据：上文「封装层与扩展点」）。
- 未命中判断用 `errors.Is(err, storage.ErrNotExist)`（依据：`src/common/storage/redis/redis.go:227`）。
