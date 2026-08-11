# 终端鉴权与白名单管理 领域词典

> 子文档 of [lexicon.md](lexicon.md)；词汇口径、来源说明与全仓待确认清单见主文档。子域锚点 `auth`，功能域口径与 docs/biz/interface/ 一致。

## 实体与业务概念

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| WhiteList | 终端白名单表（t_white_list）：IMEI（pk）+IMSI（唯一索引）组合，标识合法终端 | 鉴权与白名单导入导出语境 | `db.WhiteList`，`src/models/db/white_list.go` |
| AuthIMEIRequest | 终端联合鉴权请求（imei/imsi 两字段），Validate 仅非空检查 | 仅 authIMEI 接口请求侧 | `req.AuthIMEIRequest`，`src/models/req/auth_request.go` |
| 逃生态 | 白名单表为空时的放行状态：未导入白名单期间终端一律放行，避免未配置即全量拒止 | AuthIMEI 判定链路 | src/service/auth_service.go（count == 0 分支） |
| authCache | 进程内鉴权结果缓存（RWMutex+map，key=imei\|imsi，TTL 30min/容量 1000/超限惰性清 500） | 鉴权链路 | src/service/auth_cache.go |
| authImportLock | 白名单导入写锁 / 鉴权回源读锁，杜绝清表窗口期逃生态误放行 | 导入与鉴权回源互斥 | src/service/auth_cache.go |

## 常量与状态枚举

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| firstImport | 首次导入模式，要求白名单表为空，分片批量插入 | ImportIMEIList operation 取值 | `operationFirstImport`，`src/service/whitelist_manage_service.go` |
| update | 覆盖更新模式，事务清表后批量插入 | ImportIMEIList operation 取值 | `operationUpdate`，`src/service/whitelist_manage_service.go` |
| maxImportFileSize | 导入文件大小硬限 3MB | AuthController.ImportIMEIList | `maxImportFileSize`，`src/controllers/auth_controller.go` |
| maxImportCount / importBatchSize | 单文件最大导入 200000 条 / 分片批量插入 1000 条每批 | 白名单导入链路 | `maxImportCount`、`importBatchSize`，`src/service/whitelist_manage_service.go` |
| authCacheTTL / authCacheCapacity / authCacheEvictCount | 鉴权缓存 TTL 30min / 容量 1000 / 超限惰性清理 500 条 | 鉴权缓存 | `src/service/auth_cache.go` |
| imeiPattern | IMEI/IMSI 格式正则 `^[0-9]{15}$`（15 位纯数字） | 鉴权与导入强校验共用 | `imeiPattern`，`src/service/auth_service.go` |

## 错误码

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| AuthFailed（auth 语境） | 鉴权失败 401：authIMEI 接口格式非法/鉴权拒绝返回本码（统一 HTTP 200 + body code） | authIMEI 接口与 event 链路；login 链路拒绝用 ClientFailed(-2)，注意区分 | `retcode.AuthFailed`，`src/common/constants/retcode/retcode.go` |

## 事件

无（本子域不定义埋点事件；登录/事件上报链路的鉴权注入归 login/event 子域）。
