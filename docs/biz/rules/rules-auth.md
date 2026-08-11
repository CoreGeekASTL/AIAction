# 终端鉴权与白名单管理 业务规则

> 生成时间：2026-08-11
> 覆盖入口：POST /auth/v1/authIMEI → controllers.AuthController.AuthIMEI；POST /auth/v1/importIMEIList → controllers.AuthController.ImportIMEIList；GET /auth/v1/exportIMEIList → controllers.AuthController.ExportIMEIList（均仅内部 server 注册）

## 概述

终端鉴权功能域解决云浏览器终端（IMEI+IMSI 标识）接入合法性判定与白名单全生命周期管理问题：AuthIMEI 提供联合鉴权查询，ImportIMEIList/ExportIMEIList 提供白名单 CSV 导入导出；登录链路与事件上报链路在各自入口注入同一 AuthService.AuthIMEI 判定。规则提取自三条鉴权管理入口链路与鉴权注入点，覆盖条件分支、参数校验、阈值常量与错误码返回四类规则点。

## 规则表

| 规则名 | 条件 | 动作 | 依据 | 来源入口 |
| --- | --- | --- | --- | --- |
| request-body-unmarshal-to | AuthIMEI 请求体 JSON 反序列化失败 | 返回 retcode.AuthFailed(401)，Message "format invalid" | src/controllers/auth_controller.go | POST /auth/v1/authIMEI |
| imei-imsi-format-check | IMEI 或 IMSI 不匹配正则 `^[0-9]{15}$`（15 位纯数字） | 判 formatValid=false，鉴权拒绝且结果不入缓存 | src/service/auth_service.go | POST /auth/v1/authIMEI |
| auth-cache-hit | authCache 命中（key=imei\|imsi，TTL 30min 内） | 直接返回缓存结果，不回源 DB | src/service/auth_service.go；src/service/auth_cache.go | POST /auth/v1/authIMEI |
| white-list-empty-escape | t_white_list 记录数为 0（逃生态） | 放行（allowed=true）并缓存 true | src/service/auth_service.go | POST /auth/v1/authIMEI |
| white-list-exact-match | GetByIMEIAndIMSI 命中 / 未命中（orm.ErrNoRows） | 命中放行缓存 true；未命中拒绝缓存 false | src/service/auth_service.go；src/dao/white_list.go | POST /auth/v1/authIMEI |
| auth-db-error-safe-reject | Count 或 GetByIMEIAndIMSI 返回非 ErrNoRows 错误 | 安全优先拒绝（allowed=false）且不缓存，记错误日志 | src/service/auth_service.go | POST /auth/v1/authIMEI |
| auth-cache-capacity-evict | 缓存写入后条数 > authCacheCapacity(1000) | 同一 Lock 内按 expireAt 升序惰性清理最旧 authCacheEvictCount(500) 条 | src/service/auth_cache.go | POST /auth/v1/authIMEI |
| import-file-size-limit | 上传文件大小 > 3MB（maxImportFileSize） | 返回 retcode.ClientFailed(-2)，拒绝导入 | src/controllers/auth_controller.go | POST /auth/v1/importIMEIList |
| import-operation-check | operation 非 firstImport 且非 update | 返回 retcode.ClientFailed(-2)，拒绝导入 | src/controllers/auth_controller.go | POST /auth/v1/importIMEIList |
| import-csv-strict-validate | CSV 任一行非 2 列或 IMEI/IMSI 非 15 位纯数字 | 整体拒绝导入（不落任何数据），报行号错误 | src/service/whitelist_manage_service.go | POST /auth/v1/importIMEIList |
| import-count-limit | CSV 行数 > maxImportCount(200000) | 整体拒绝导入 | src/service/whitelist_manage_service.go | POST /auth/v1/importIMEIList |
| first-import-empty-check | operation=firstImport 且 t_white_list 非空 | 返回错误 "white list is not empty, please use update"，拒绝导入 | src/service/whitelist_manage_service.go | POST /auth/v1/importIMEIList |
| import-batch-insert | firstImport 模式入库 | 按 importBatchSize(1000) 条/批分片 InsertMulti | src/service/whitelist_manage_service.go | POST /auth/v1/importIMEIList |
| update-clear-and-insert | operation=update | 事务内清表（DELETE FROM t_white_list）+ 批量插入，失败整批回滚 | src/dao/white_list.go | POST /auth/v1/importIMEIList |
| import-hold-write-lock | 导入全程（两种模式） | 持 authImportLock 写锁；鉴权回源段持读锁互斥，杜绝清表窗口期逃生态误放行 | src/service/whitelist_manage_service.go；src/service/auth_service.go | POST /auth/v1/importIMEIList |
| export-failed | ListAll 查询失败 | 返回 HTTP 500（InternalServiceError） | src/controllers/auth_controller.go | GET /auth/v1/exportIMEIList |
| login-inject-auth-reject | 登录链路 loginAuth 注入 AuthIMEI 判定不通过 | 返回 retcode.ClientFailed(-2)，Message "auth rejected"，终止登录 | src/controllers/login_controller.go；src/controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth；POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser；POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |
| event-inject-auth-reject | 事件上报链路注入 AuthIMEI 判定不通过 | 返回 retcode.AuthFailed(401)，Message "auth rejected"，事件不记录 | src/controllers/event_controller.go | POST /app-api/center/public/client/sendClientEvent；POST /app-api/center/public/client/sendAppUseTimesEvent |

## 补充说明

AuthIMEI 判定顺序固定：格式校验短路 → 缓存查询 → 逃生态判定（表空放行）→ 联合精确匹配；DB 异常一律安全优先拒绝且不写缓存，缓存只承载确定性结果。AuthController.AuthIMEI 统一 HTTP 200 + body code 表达结果，与 login/event 链路注入鉴权的错误码口径不同：login 路径拒绝用 ClientFailed(-2)、event 路径拒绝用 AuthFailed(401)、authIMEI 接口拒绝用 AuthFailed(401)，三者不可混用。CSV 为纯数据无 header 口径，空行由 csv.Reader 自动跳过；导出不经鉴权缓存与导入锁。
