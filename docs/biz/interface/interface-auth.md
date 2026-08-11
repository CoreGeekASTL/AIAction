# 终端鉴权与白名单管理

> 功能域：终端鉴权　接口数：3　所属 server：仅内部（innerServer 127.0.0.1:9090）
> 子文档 of [README.md](README.md)

## 1. 定位

为终端 IMEI+IMSI 联合鉴权与白名单管理提供内部管理面接口：鉴权查询、白名单 CSV 导入/导出。三个接口均仅注册于 innerServer（routers/beego_router.go RegisterInternalRouter），不对外暴露。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| AuthIMEI | 终端 IMEI+IMSI 联合鉴权 | controllers/auth_controller.go | POST /auth/v1/authIMEI |
| ImportIMEIList | 白名单 CSV 导入（firstImport/update） | controllers/auth_controller.go | POST /auth/v1/importIMEIList |
| ExportIMEIList | 白名单 CSV 导出 | controllers/auth_controller.go | GET /auth/v1/exportIMEIList |

## 3. 数据结构说明

- **AuthIMEI**
  - 请求 `req.AuthIMEIRequest`（models/req/auth_request.go）：IMEI、IMSI（json: imei/imsi）
  - 响应 `resp.BaseResponse`（code/msg）：统一 HTTP 200，code 取值 retcode.Success(200) 或 retcode.AuthFailed(401)（format invalid / auth rejected）
- **ImportIMEIList**
  - 请求 multipart 表单：file（CSV 文件，无 header、IMEI/IMSI 两列，≤3MB）、operation（firstImport / update）
  - 响应 `resp.DataResponse`（code/msg/data）：data 为导入条数；失败 code 为 retcode.ClientFailed(-2)（参数非法）或 retcode.InternalFailed(-1)（导入失败）
- **ExportIMEIList**
  - 请求 无参数
  - 响应 `text/csv` 文本（无 header、IMEI/IMSI 两列全量）；失败 HTTP 500（InternalServiceError）
