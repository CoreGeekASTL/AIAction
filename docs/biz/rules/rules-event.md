# 客户端事件上报 业务规则

> 生成时间：2026-08-11
> 覆盖入口：POST /app-api/center/public/client/sendClientEvent → controllers.EventController.SendClientEvent；POST /app-api/center/public/client/sendAppUseTimesEvent → controllers.EventController.SendAppUseTimesEvent

## 概述

客户端事件上报功能域负责接收终端上报的客户端事件与应用使用时长事件，并落本地审计事件存储。规则提取自两条事件上报入口链路，覆盖参数校验与错误码返回两类规则点。

## 规则表

| 规则名 | 条件 | 动作 | 依据 | 来源入口 |
| --- | --- | --- | --- | --- |
| event-inject-auth-reject | 事件上报前注入 AuthService.AuthIMEI 白名单联合鉴权判定不通过 | 返回 retcode.AuthFailed(401)，Message "auth rejected"，事件不记录（鉴权细节见 rules-auth.md） | src/controllers/event_controller.go | POST /app-api/center/public/client/sendClientEvent；POST /app-api/center/public/client/sendAppUseTimesEvent |
| request-body-unmarshal-to | 事件请求体 JSON 反序列化失败 | 返回 retcode.ClientFailed(-2)，拒绝本次上报 | src/controllers/event_controller.go | POST /app-api/center/public/client/sendClientEvent；POST /app-api/center/public/client/sendAppUseTimesEvent |
| report-event-failed | EventService.ReportEvent 落事件存储失败 | 返回 retcode.ClientFailed(-2)，事件不视为记录成功 | src/controllers/event_controller.go | POST /app-api/center/public/client/sendClientEvent；POST /app-api/center/public/client/sendAppUseTimesEvent |
| default-event-storage | 事件存储初始化时取到非法 eventStorage | 回退使用 DefaultEventStorage（"localAuditComponent"，本地文件事件存储） | src/service/event_service.go | POST /app-api/center/public/client/sendClientEvent；POST /app-api/center/public/client/sendAppUseTimesEvent |

## 补充说明

两条入口共用同一条 ReportEvent 链路，上报前先经 AuthService.AuthIMEI 白名单鉴权（拒绝返回 401，为终止性规则），事件类型分别为 events.Client 与 events.AppUseTimes，无字段级业务校验（ClientEventRequest/AppUseTimesEvent 的 Validate 均为空实现，代码未体现字段约束，待确认）。事件存储通过 sync.Once 初始化并注册到工厂，存储失败是唯一导致上报失败返回的业务分支。
