# 文件上传下载 业务规则

> 生成时间：2026-08-11
> 覆盖入口：POST /app-api/control/file/upload → controllers.FileController.HandleUpload / controllers.ExFileController.HandleUpload；GET /app-api/control/file/download/:fileName → controllers.FileController.HandleDownload / controllers.ExFileController.HandleDownload；GET /file/v1/:bucket/:name → controllers.FileController.Download；POST /file/v1/:bucketName/:name → controllers.FileController.Upload；GET /file/v1/:bucket/:name/exist → controllers.FileController.Exist；DELETE /file/v1/:bucketName/:fileName → controllers.FileController.HandleDelete

## 概述

文件功能域提供文件的上传、下载、存在性查询与删除，文件内容持久化在 DB 文件表（Bucket+Name 为键）。规则提取自 /app-api/control/file 与 /file/v1 两组入口链路，覆盖参数校验、安全校验（路径遍历防护）与错误码返回三类规则点。

## 规则表

| 规则名 | 条件 | 动作 | 依据 | 来源入口 |
| --- | --- | --- | --- | --- |
| handle-upload-file-name | 查询参数 fileName 为空 | 返回 retcode.ClientFailed(-2)，Message "Missing fileName parameter"，拒绝上传 | src/controllers/file_controller.go；src/controllers/exfile_controller.go | POST /app-api/control/file/upload |
| clean-file-name | filepath.Clean 后文件名为空，或仍含路径分隔符（pathSeparatorRegex `[/\\]`） | 返回 "invalid file name" 错误，阻断上传/下载/删除（防路径遍历） | src/service/file_service.go | POST /app-api/control/file/upload；GET /app-api/control/file/download/:fileName；POST /file/v1/:bucketName/:name；GET /file/v1/:bucket/:name；DELETE /file/v1/:bucketName/:fileName |
| upload-bucket | /app-api/control/file 入口上传/下载 | 固定使用 constants.UploadBucket（"upload-bucket"），bucket 不可由请求指定 | src/controllers/file_controller.go；src/controllers/exfile_controller.go | POST /app-api/control/file/upload；GET /app-api/control/file/download/:fileName |
| insert-or-update | 上传文件（Bucket,Name）记录不存在（orm.ErrNoRows） | 插入新文件记录（CreatedAt 置当前时间）；已存在则覆盖 Content 与 Size | src/service/file_service.go | POST /app-api/control/file/upload；POST /file/v1/:bucketName/:name |
| upload-file-failed | insertOrUpdate 落库失败 | 返回 "files upload to OSS failed" 错误，上传失败 | src/service/file_service.go | POST /app-api/control/file/upload；POST /file/v1/:bucketName/:name |
| download-file-failed | 按（Bucket,Name）查询文件失败（含不存在） | 返回 "failed to download file from OSS" 错误，下载失败 | src/service/file_service.go | GET /app-api/control/file/download/:fileName；GET /file/v1/:bucket/:name |
| exist-not-found | Exist 查询文件不存在（err==nil 且 ok==false） | 返回 HTTP 404（NotFound） | src/controllers/file_controller.go | GET /file/v1/:bucket/:name/exist |

## 补充说明

所有写/读路径先经 cleanFileName 安全校验，属短路规则，校验不过不触库。/app-api/control/file 与 /file/v1 两组入口底层共用同一 FileService，差异在于 bucket 来源（固定 vs 路径参数）与响应封装（DataResponse vs 裸 body）。HandleDownload 入口 fileName 为空时返回码为 InternalFailed（与 HandleUpload 的 ClientFailed 不一致，代码现状如此）。
