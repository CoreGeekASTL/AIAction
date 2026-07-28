# 文件管理

> 功能域：文件管理　接口数：8　所属 server：外部(HTTPS) 2 + 内部(HTTP) 6
> 子文档 of [README.md](README.md)

## 1. 定位

对象存储文件上传/下载/删除。外部 ExFileController（HTTPS）暴露 upload/download；内部 FileController（HTTP）额外暴露 /file/v1/* RESTful 接口。底层走 service.FileService → common/storage/oss（MinIO SDK）。

## 2. 接口清单

**外部接口**（ExFileController，externalServer HTTPS）：

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| HandleUpload | 上传文件 | controllers/exfile_controller.go | POST /app-api/control/file/upload |
| HandleDownload | 下载文件 | controllers/exfile_controller.go | GET /app-api/control/file/download/:fileName |

**内部接口**（FileController，innerServer HTTP）：

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| HandleUpload | 上传文件 | controllers/file_controller.go | POST /app-api/control/file/upload |
| HandleDownload | 下载文件 | controllers/file_controller.go | GET /app-api/control/file/download/:fileName |
| Upload | RESTful 上传 | controllers/file_controller.go | POST /file/v1/:bucketName/:name |
| Download | RESTful 下载 | controllers/file_controller.go | GET /file/v1/:bucket/:name |
| Exist | 判断文件存在 | controllers/file_controller.go | GET /file/v1/:bucket/:name/exist |
| HandleDelete | 删除文件 | controllers/file_controller.go | DELETE /file/v1/:bucketName/:fileName |

## 3. 数据结构说明

- **HandleUpload**
  - 请求：query 参数 `fileName`（string，必填）；body 为文件内容（io.ReadAll）
  - 响应：retcode 标准结构
- **HandleDownload**
  - 请求：path 参数 `:fileName`（string，必填）
  - 响应：文件流（io.Reader）
- **Upload / Download / Exist / HandleDelete**
  - 请求：path 参数 `:bucket`/`:bucketName`/`:name`/`:fileName`
  - 响应：Download 返回文件流；Exist 返回 OK(nil) 或 NotFound()；HandleDelete 返回 retcode 标准结构

## 4. 风险与注意点

- **文件上传无大小校验**：controllers/exfile_controller.go:47（`io.ReadAll(c.Body())` 直接读全 body，无大小限制，可能内存耗尽）
- **文件下载 path 参数未校验**：controllers/file_controller.go:25/30（`:fileName`/`:bucket` 未做路径穿越校验，可能读到任意路径文件）
