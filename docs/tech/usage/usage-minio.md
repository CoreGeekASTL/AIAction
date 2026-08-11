# minio-go/v7 使用现状（存储/对象存储）

## 用途定位
`github.com/minio/minio-go/v7` 仅通过 `src/common/storage/oss/minio.go` 的自研 `oss.Client` 接口使用（PutObject/GetObject/DeleteObject/EnsureBucket/IsOnline），作为文件/插件包的对象存储后端（与 `t_file` 表配合，见 `src/controllers/file_controller.go`、`src/service` 插件相关逻辑）。全局单例：`Init(conf.OSSConfig)` 创建并健康检查、`Instance()` 取用。默认凭据 `minioadmin/minioadmin` 来自 conf 包默认值，`Secure: false`（明文 HTTP）。

## 使用模式

```go
// 来源：src/common/storage/oss/minio.go
func Init(config conf.OSSConfig) error {
	c, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, config.Token),
		Secure: false,
	})
	// ...
	client = &ossClient{client: c}
	if !client.IsOnline() {
		return errors.New("minio healthcheck failed")
	}
	return nil
}

func (c *ossClient) EnsureBucket(ctx context.Context, bucket string) error {
	exists, err := c.client.BucketExists(ctx, bucket)
	// 不存在则 MakeBucket
}
```
