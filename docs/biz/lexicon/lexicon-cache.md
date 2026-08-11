# 缓存管理 领域词典

> 子文档 of [lexicon.md](lexicon.md)；词汇口径、来源说明与全仓待确认清单见主文档。子域锚点 `cache`，功能域口径与 docs/biz/interface/ 一致。

## 实体与业务概念

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| DeleteCacheRequest | 删除缓存请求，IMEI/IMSI 均必填（为空校验报错） | IMEI/IMSI 在本域为 form 表单参数；在 login 域为 JSON 字段 | `req.DeleteCacheRequest`，`src/models/req/request_entity.go` |
