# 文件管理

> 功能域：文件管理　接口数：6　所属 server：外部 + 内部
> 子文档 of [README.md](README.md)

## 1. 定位

云浏览器控制面文件的上传/下载。`/app-api/control/file/*` 两个接口经 externalServer（`ExFileController`）与 innerServer（`FileController`）双暴露，固定写入 UploadBucket；`/file/v1/*` 四个 bucket 级接口仅注册在内部 server，按 bucket/name 管理文件。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| HandleUpload | 上传文件（固定 UploadBucket，文件名走 query） | controllers/exfile_controller.go、controllers/file_controller.go | POST /app-api/control/file/upload |
| HandleDownload | 下载文件（固定 UploadBucket） | controllers/exfile_controller.go、controllers/file_controller.go | GET /app-api/control/file/download/:fileName |
| Download | 按 bucket/name 下载文件（仅内部） | controllers/file_controller.go | GET /file/v1/:bucket/:name |
| Upload | 按 bucket/name 上传文件（multipart，仅内部） | controllers/file_controller.go | POST /file/v1/:bucketName/:name |
| Exist | 判断文件是否存在（仅内部） | controllers/file_controller.go | GET /file/v1/:bucket/:name/exist |
| HandleDelete | 删除指定文件（仅内部） | controllers/file_controller.go | DELETE /file/v1/:bucketName/:fileName |

## 3. 数据结构说明

- **HandleUpload**
  - 请求：query 参数 fileName（对应 `req.FileUploadRequest`，models/req/request_entity.go，必填）；请求体为文件二进制内容
  - 响应 `resp.DataResponse`（models/resp/response_entity.go）：code/msg + data 为落盘后的文件路径（string）
- **HandleDownload**
  - 请求：路径参数 fileName
  - 响应：文件二进制流（Content-Disposition: attachment; Content-Type: application/octet-stream）
- **Upload**
  - 请求：路径参数 bucketName、name；multipart 表单字段 `file`（文件内容）
  - 响应：data 为文件路径（string）
- **Download**
  - 请求：路径参数 bucket、name
  - 响应：文件二进制流
- **Exist**
  - 请求：路径参数 bucket、name
  - 响应 `resp.BaseResponse`；不存在返回 404
- **HandleDelete**
  - 请求：路径参数 bucketName、fileName
  - 响应 `resp.BaseResponse`
