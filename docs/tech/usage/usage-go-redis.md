# go-redis/v9 使用现状（存储/缓存）

## 用途定位
`github.com/redis/go-redis/v9` 仅通过 `src/common/storage/redis/redis.go` 的自研 `Client` 接口使用（Get/Set/HSet/SAdd/SetNx 等 15 个方法），业务代码不接触 go-redis 原生 API。缓存对象须实现 `Object` 接口（`GetKey()` + `encoding.BinaryMarshaler/BinaryUnmarshaler`），hash 对象另实现 `GetField()`。全局单例：`Init(conf.RedisConfig)` 初始化、`Instance()` 取用；`InitForTest` 供测试注入 mock。key 未命中时把 `redis.Nil` 翻译成 `storage.ErrNotExist`（`src/common/storage/error.go`）。

## 使用模式

```go
// 来源：src/common/storage/redis/redis.go
func Init(config conf.RedisConfig) error {
	if config.DB < 0 || config.DB > maxRedisDB {
		return errors.New("illegal redis db number]")
	}
	client = New(&redis.Options{
		Addr:     config.Endpoint,
		Username: "",
		Password: "",
		DB:       config.DB,
	})
	err := client.Ping(context.Background())
	// ...
}

func (c *innerClient) Get(ctx context.Context, dst Object) error {
	result, err := c.client.Get(ctx, dst.GetKey()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return storage.ErrNotExist
		}
		return err
	}
	return dst.UnmarshalBinary([]byte(result))
}
```

带过期写入骨架：

```go
// 来源：src/common/storage/redis/redis.go
cmd := c.client.Set(ctx, src.GetKey(), src, expiration)      // SetWithExpiration
c.client.SetNX(ctx, src.GetKey(), src, expiration).Result()  // SetNx
```
