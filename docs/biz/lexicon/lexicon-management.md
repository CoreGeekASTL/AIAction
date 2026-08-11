# 浏览器配置同步 领域词典

> 子文档 of [lexicon.md](lexicon.md)；词汇口径、来源说明与全仓待确认清单见主文档。子域锚点 `management`，功能域口径与 docs/biz/interface/ 一致。

## 实体与业务概念

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| Config | 配置存储表（t_config）：type/content(text)/时间戳，浏览器配置持久化载体 | - | `db.Config`，`src/models/db/browser_config.go` |
| SyncBrowserConfigRequest | 浏览器配置同步请求，含 RouteAPPConfigList/ChromeConfigList/URLConfigs 三类配置 | - | `req.SyncBrowserConfigRequest`，`src/models/req/request_entity.go` |
| RouteAppConfig | 路由应用配置（manufacturer/model/type/mode/extendModel/name/description） | db 包另有同义 `db.RouterAPPConfig`（DB 侧结构，type 字段名同为 Type） | `req.RouteAppConfig`，`src/models/req/request_entity.go`；`db.RouterAPPConfig`，`src/models/db/browser_config.go` |
| ChromeConfig | Chrome 运行配置：帧率/码率/采样率/声道/机型/FFCode/分辨率/RecordMode 等 | req 与 db 两包各有一份同名字段结构（req 带 omitempty） | `req.ChromeConfig`，`src/models/req/request_entity.go`；`db.ChromeConfig`，`src/models/db/browser_config.go` |
| URLConfigs | URL 配置（nodeIdent/appType/url/appID/name/isVideoType/isWebType/isShortType/userAgent） | req 与 db 两包各有一份（db 包名为 URLConfig） | `req.URLConfigs`，`src/models/req/request_entity.go`；`db.URLConfig`，`src/models/db/browser_config.go` |
| UserAgent | URL 配置中的 UA 串 | - | `req.URLConfigs.UserAgent`，`src/models/req/request_entity.go` |
