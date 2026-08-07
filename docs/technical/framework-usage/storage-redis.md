# go-redis 使用指导（存储/缓存）

## 用途定位
Redis 缓存/共享状态访问的预留封装层。当前版本业务链路未启用（grep 无 `redis.Init`/`redis.Instance` 调用点），属于已就位的公共能力。


## 使用模式

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
