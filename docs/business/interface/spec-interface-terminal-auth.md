# 终端鉴权（IMEI+IMSI 白名单）

> 功能域：terminal-auth　接口数：3　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)（27.0 终端鉴权需求新增，2026-07-30）

## 1. 定位

商用准入控制：运营批量导入/导出 IMEI+IMSI 白名单，BrowserGW 反调联合鉴权接口校验终端是否在授权范围内。仅注册到内部 server（127.0.0.1:9090/运维内网段），外部不可达，依赖网络隔离无接口级认证。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| AuthIMEI | 联合鉴权：IMEI+IMSI 组合精确匹配白名单同一行 | src/controllers/auth_controller.go | POST /auth/v1/authIMEI |
| ImportIMEIList | 白名单 CSV 批量导入（firstImport 首次/update 覆盖） | src/controllers/auth_controller.go | POST /auth/v1/importIMEIList |
| ExportIMEIList | 全量白名单导出为 CSV | src/controllers/auth_controller.go | GET /auth/v1/exportIMEIList |

## 3. 数据结构说明

- **AuthIMEI**
  - 请求 `req.AuthIMEIRequest`（src/models/req/auth_request.go）：IMEI、IMSI，Validate 校验非空；服务侧强校验 15 位纯数字（`^[0-9]{15}$`）
  - 响应 `resp.BaseResponse`：命中放行 `{code:200}`；格式非法 `{code:401, msg:"imei or imsi format invalid"}`；未命中 `{code:401, msg:"auth rejected"}`（HTTP 400 + body code，src/controllers/auth_controller.go）
  - 鉴权逻辑：格式短路 → 内存缓存（联合键，TTL 30min）→ 逃生态（白名单表 Count==0 一律放行）→ 按 IMEI 查行比对 IMSI；DB 异常 fail-open 放行（src/service/auth_service.go）
- **ImportIMEIList**
  - 请求：multipart/form-data，文件字段 `file`（CSV，无 header 纯数据，两列 IMEI,IMSI）+ form 字段 `operation`（firstImport/update）
  - 响应 `resp.DataResponse`：成功 `{code:200, data:导入条数}`；校验失败（文件>3MB、行数>20W、非 15 位纯数字、文件内重复组合、firstImport 时表非空）`{code:-1, msg}`，表非空 msg 含 "not empty"；参数错误（缺 file、operation 非法）`{code:-2, msg}`（src/service/auth_manage_service.go）
  - update 为事务清表+批量插入，失败整体回滚；导入成功后清空鉴权缓存立即生效
- **ExportIMEIList**
  - 请求：无参数
  - 响应：HTTP 200 `text/csv` 纯文本，首行表头 `IMEI,IMSI`，后续每行一条记录（src/controllers/auth_controller.go）

