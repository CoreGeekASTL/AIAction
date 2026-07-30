# 文件上传下载管理

> 功能域：file-transfer　接口数：6　所属 server：外部(HTTPS) + 内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

面向终端的文件上传/下载（固定上传桶），以及面向内部组件的通用 bucket 文件 CRUD。`/app-api/control/file/*` 两条路径内外双暴露，`/file/v1/*` 仅内部暴露。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| HandleUpload | 上传文件到固定上传桶（body 原始字节 + query 传文件名） | src/controllers/exfile_controller.go；src/controllers/file_controller.go | POST /app-api/control/file/upload?fileName=xxx |
| HandleDownload | 从固定上传桶下载文件（文件流响应） | src/controllers/exfile_controller.go；src/controllers/file_controller.go | GET /app-api/control/file/download/:fileName |
| Upload | 上传文件到指定 bucket（multipart 表单） | src/controllers/file_controller.go | POST /file/v1/:bucketName/:name |
| Download | 从指定 bucket 下载文件（文件流响应） | src/controllers/file_controller.go | GET /file/v1/:bucket/:name |
| Exist | 判断指定 bucket 中文件是否存在 | src/controllers/file_controller.go | GET /file/v1/:bucket/:name/exist |
| HandleDelete | 删除指定 bucket 中文件 | src/controllers/file_controller.go | DELETE /file/v1/:bucketName/:fileName |

## 3. 数据结构说明

- **HandleUpload（/app-api/control/file/upload）**
  - 请求：query 参数 `fileName`（必填，缺失返回 ClientFailed）；body 为文件原始字节流（非 multipart）
  - 响应 `resp.DataResponse`：`BaseResponse{code,msg}` + `data` 为上传后的文件路径（string）；bucket 固定为 `constants.UploadBucket`（src/controllers/file_controller.go:126）
- **HandleDownload（/app-api/control/file/download/:fileName）**
  - 请求：路径参数 `fileName`
  - 响应：非 JSON 壳——文件字节流，header `Content-Disposition: attachment; filename=...`、`Content-Type: application/octet-stream`
- **Upload（/file/v1/:bucketName/:name）**
  - 请求：路径参数 `bucketName`、`name`；multipart 表单字段 `file`（文件内容）
  - 响应：上传后的文件路径（string，经 `c.OK` 直接序列化）
- **Download（/file/v1/:bucket/:name）**
  - 请求：路径参数 `bucket`、`name`
  - 响应：文件字节流，header 同 HandleDownload；注意 Content-Length 用 `string(rune(len))` 写入，值不正确（src/controllers/file_controller.go:91）
- **Exist**
  - 请求：路径参数 `bucket`、`name`；响应：存在返回 200 `BaseResponse`，不存在返回 HTTP 404
- **HandleDelete**
  - 请求：路径参数 `bucketName`、`fileName`；响应 `resp.BaseResponse`（附带多余的下载类响应头，src/controllers/file_controller.go:190-191）

## 4. 风险与注意点

- **响应非统一壳**：下载类接口直接写文件流，不走 `resp.BaseResponse` JSON 壳，调用方需按流处理（src/controllers/exfile_controller.go:99-103）
- **Content-Length 缺陷**：src/controllers/file_controller.go:91（`string(rune(len(content)))` 把长度当 rune 转字符，非数字字符串）
