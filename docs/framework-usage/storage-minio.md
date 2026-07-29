# MinIO（minio-go v7）封装使用指导（存储/对象）

> 版本：github.com/minio/minio-go/v7 v7.0.69 ｜ 调用点：0 业务调用点（仅封装层 1 文件）｜ 涉及文件：1 ｜ 基线：origin/main @ ae0a8a6（2026-07-29 复核）

## 用途定位

`src/common/storage/oss/minio.go` 是对 minio-go 的对象存储封装（PutObject/GetObject/DeleteObject/EnsureBucket）。**当前仓库中没有任何业务代码调用它**（grep `oss.Instance|oss.Init` 无业务命中），属预置能力；`conf/app.conf` 的 `[oss] endpoint` 被 `common/conf/config.go:16` 读取，默认凭据 `minioadmin/minioadmin`（`config.go:38-40`）。

## 初始化与配置

```go
// 来源：src/common/storage/oss/minio.go:28-44
func Init(config conf.OSSConfig) error {
	c, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, config.Token),
		Secure: false,   // 明文 HTTP，仅内网
	})
	client = &ossClient{client: c}
	if !client.IsOnline() { return errors.New("minio healthcheck failed") }
}
```

## 核心使用模式

```go
// 上传（minio.go:86-98）
err := oss.Instance().PutObject(ctx, oss.PutObject{
	BucketName: bucket, FileName: name, File: reader, Size: size,
})
// 下载（minio.go:111-117）——返回 io.ReadCloser，调用方负责 Close
rc, err := oss.Instance().GetObject(ctx, oss.GetObjectOptions{BucketName: bucket, FileName: name})
// 建桶（幂等，minio.go:66-80）
err := oss.Instance().EnsureBucket(ctx, bucket)
```

## 封装层与扩展点

- `Client` 接口（`minio.go:20-26`）5 个方法，参数全部用自定义 `PutObject/GetObjectOptions` 结构体，屏蔽 minio-go 原生类型。
- 删除是 `ForceDelete: true`（`minio.go:101-103`）。
- **若启用对象存储，业务代码只允许经此接口，禁止直接 import minio-go。**

## 并发与线程模型

`client` 包级变量，Init 后只读；minio.Client 协程安全。

## 错误处理与容错

Init 用 `IsOnline()` 做健康检查（`minio.go:38-42`）；运行时错误直接透传，无重试。

## 约定与规范

bucket 名与文件名由调用方规划；file 表（`models/db/file.go`）设计了 `bucket/name/content` 字段，说明文件内容历史上既可落 DB bytea 也可落 OSS——当前实现落 DB（`dao/file.go`），OSS 通路未启用。

## 已知问题与反模式

- **未接线**：main 流程无 `oss.Init` 调用，`oss.Instance()` 为 nil，直接用会 panic。
- `Secure: false` 明文传输（`minio.go:31`），只允许内网 endpoint。
- 默认凭据 `minioadmin/minioadmin` 硬编码为包级默认（`common/conf/config.go:38-39`），生产必须经命令行 flag 或配置覆盖。

## AI 编码指南

- 新增对象存储能力：先在启动流程接线 `oss.Init(conf.Instance().OSS)`（失败要降级/报错），再经 `oss.Instance()` 调用；上传前 `EnsureBucket`。依据：上文「已知问题与反模式」。
- `GetObject` 返回的 `io.ReadCloser` 必须 `defer Close()`。依据：`minio.go:111-117`。
- **禁止**绕过封装直接用 minio-go、**禁止**开启 `Secure:false` 对公网 endpoint。依据：`minio.go:28-32`。
