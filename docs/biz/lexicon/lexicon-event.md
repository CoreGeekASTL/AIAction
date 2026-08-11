# 客户端事件上报 领域词典

> 子文档 of [lexicon.md](lexicon.md)；词汇口径、来源说明与全仓待确认清单见主文档。子域锚点 `event`，功能域口径与 docs/biz/interface/ 一致。

## 实体与业务概念

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| ClientEventRequest | 客户端异常事件上报请求（hsman/hstype/appType/imei/imsi/type） | - | `req.ClientEventRequest`，`src/models/req/event_request.go` |
| AppUseTimesEvent（请求） | 客户端应用使用时长上报请求（useTimes/hsman/hstype/exttype/appType/appId/scheight/scwidth/imei/imsi/playMode） | - | `req.AppUseTimesEvent`，`src/models/req/event_request.go` |
| HSMan / HSType | 终端厂商与终端型号（handset man/type） | - | `req.ClientEventRequest.HSMan/HSType`，`src/models/req/event_request.go` |
| Info | 埋点事件上报体：EventDesc + service/eventTime/env/hostname/object/eventData | - | `events.Info`，`src/models/events/base.go` |
| LoginEventData | 登录埋点数据（imei/imsi/appType/extType/hsman/hstype/loginTime 及各分配地址/totalKb/freeKb） | - | `events.LoginEventData`，`src/models/events/base.go` |
| ClientEventData | 客户端异常埋点数据（hsman/hstype/appType/imei/imsi/type） | - | `events.ClientEventData`，`src/models/events/base.go` |
| AppUseTimesEvent（埋点） | 应用使用时长埋点数据，字段与上报请求一致 | 与 `req.AppUseTimesEvent` 同名不同包：events 包为埋点载荷，req 包为 HTTP 请求 | `events.AppUseTimesEvent`，`src/models/events/base.go` |

## 常量与状态枚举

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| Login 事件 | 云浏览器用户HTTP登录埋点（browser_user_http_login），触发方 client | - | `events.Login`，`src/models/events/base.go` |
| LoginError 事件 | 登录失败埋点（browser_user_http_login_error），已定义常量但未入 eventTypeMap（待确认，见主文档） | - | `events.LoginError`，`src/models/events/base.go` |
| Client 事件 | 客户端异常埋点（browser_client_error），触发方 client | - | `events.Client`，`src/models/events/base.go` |
| AppUseTimes 事件 | 客户端上传应用使用时长埋点（browser_client_app_use_times），触发方 client | - | `events.AppUseTimes`，`src/models/events/base.go` |
| CloudBrowser | 埋点事件 object 字段固定取值 "cloud-browser" | - | `events.CloudBrowser`，`src/models/events/base.go` |

## 事件

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| EventType | 埋点事件类型枚举（string），取值见上表四个常量 | - | `events.EventType`，`src/models/events/base.go` |
| EventDesc | 事件描述三元组（event/eventDesc/eventTrigger） | - | `events.EventDesc`，`src/models/events/base.go` |
