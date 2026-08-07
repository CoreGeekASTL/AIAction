# MinIO OSS 使用指导（存储/对象）

## 用途定位
对象存储（文件/插件包等大对象）的封装层。当前业务链路文件内容存 DB bytea（t_file.content），OSS 属预留/按局点启用的能力。


## 使用模式

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
