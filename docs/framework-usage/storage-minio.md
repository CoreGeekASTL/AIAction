# MinIO OSS 使用指导（存储/对象）

> 版本：github.com/minio/minio-go/v7 v7.0.69 ｜ 调用点：封装层 5 接口方法 ｜ 涉及文件：common/storage/oss/minio.go（业务调用方同 redis，当前未在 main 接入 Init）｜ 基线：main (6c93561)

## 用途定位
对象存储（文件/插件包等大对象）的封装层。当前业务链路文件内容存 DB bytea（t_file.content），OSS 属预留/按局点启用的能力。

## 初始化与配置
- `oss.Init(conf.OSSConfig)`：`minio.New(endpoint, creds, Secure:false)` + `IsOnline()` 健康检查（`src/common/storage/oss/minio.go:28-44`）。
- 配置：app.conf `[oss] endpoint`；默认 ak/sk "minioadmin"（`src/common/conf/config.go:36-41`）。
- 进程级单例 `oss.Instance()`（`src/common/storage/oss/minio.go:46`）。

## 核心使用模式

```go
// 来源：src/common/storage/oss/minio.go:54-80
client := oss.Instance()
err := client.EnsureBucket(ctx, bucket)                 // 不存在则 MakeBucket
err := client.PutObject(ctx, oss.PutObject{             // 上传
	BucketName: bucket, FileName: name, File: reader, Size: size,
})
rc, err := client.GetObject(ctx, oss.GetObjectOptions{BucketName: bucket, FileName: name}) // 返回 io.ReadCloser
err := client.DeleteObject(ctx, bucket, filename)
```

## 封装层与扩展点
- 入口：`GIDS/common/storage/oss`，接口 `Client`（PutObject/EnsureBucket/IsOnline/DeleteObject/GetObject）。
- 隐藏：minio 凭证、bucket 存在性检查、healthcheck。
- 扩展点：PutObject/GetObjectOptions 参数结构体，新增选项在此扩展。

## 并发与线程模型
minio.Client 并发安全，单例复用。

## 错误处理与容错
- Init 时健康检查失败返回 error 拒绝启动使用该能力（`src/common/storage/oss/minio.go:38-42`）。
- 运行期错误直接上抛 + Errorf 日志，无重试。

## 约定与规范
- 所有 OSS 访问经 `oss.Instance()`，禁止业务直接 import minio-go。
- bucket 由调用方命名并 Ensure，无集中常量。

## 已知问题与反模式
- 与 redis 同样未接入 main 启动流程，启用前需补 `oss.Init(conf.Instance().OSS)`。
- `Secure:false` 明文 HTTP 连接（`src/common/storage/oss/minio.go:31`），生产跨网段需评估。

## AI 编码指南
- 新增大对象存取：经 `oss.Instance()` + `EnsureBucket` 后 Put/Get/Delete；启用前在 main 补 `oss.Init`（依据：上文「初始化与配置」）。
- **禁止**业务代码直接 import minio-go（依据：上文「封装层与扩展点」）。
- `GetObject` 返回的 `io.ReadCloser` 调用方必须关闭（依据：`src/common/storage/oss/minio.go:25` 接口签名）。
