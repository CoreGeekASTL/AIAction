# 文件管理 领域词典

> 子文档 of [lexicon.md](lexicon.md)；词汇口径、来源说明与全仓待确认清单见主文档。子域锚点 `file`，功能域口径与 docs/biz/interface/ 一致。

## 实体与业务概念

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| File | 文件存储表（t_file）：bucket/name/content(bytea)/size/created_at | - | `db.File`，`src/models/db/file.go` |
| Bucket | 文件归属桶名，File 表字段 | plugin 域 PluginPackage.PackageBucket 亦用 bucket 概念（包存储桶） | `db.File.Bucket`，`src/models/db/file.go` |
| ByteArrayField | 文件内容字节数组自定义 ORM 字段类型（底层 TypeTextField） | - | `db.ByteArrayField`，`src/models/db/file.go` |
| FileUploadRequest | 文件上传请求，fileName 必填 | - | `req.FileUploadRequest`，`src/models/req/request_entity.go` |

## 常量与状态枚举

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| MaxFileSize | 上传文件大小上限 300MB（300*1024*1024） | - | `constants.MaxFileSize`，`src/common/constants/base.go` |
| UploadBucket | 上传文件固定桶名 "upload-bucket" | - | `constants.UploadBucket`，`src/common/constants/base.go` |
