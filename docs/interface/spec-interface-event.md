# 客户端事件上报

> 功能域：客户端事件上报　接口数：2（双实现）　所属 server：外部(HTTPS) + 内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

客户端使用事件上报，经 externalServer（HTTPS）与 innerServer（HTTP）双暴露，EventController 承载（routers/beego_router.go:22/32）。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| SendClientEvent | 上报客户端事件 | controllers/event_controller.go | POST /app-api/center/public/client/sendClientEvent |
| SendAppUseTimesEvent | 上报应用使用时长事件 | controllers/event_controller.go | POST /app-api/center/public/client/sendAppUseTimesEvent |

## 3. 数据结构说明

- **SendClientEvent**
  - 请求 `req.ClientEventRequest`（models/req/event_request.go）：HSMan（厂商）；HSType（机型）；AppType；IMEI；IMSI；Type（事件类型）
  - 响应 `resp.DataResponse`：BaseResponse{code,msg} + Data=true
  - 内部转为 `events.ClientEventData` 经 EventService.ReportEvent 落库
- **SendAppUseTimesEvent**
  - 请求 `req.AppUseTimesEvent`（models/req/event_request.go）：UseTimes（使用时长）；HSMan；HSType；AppType；AppId；EXTType；PlayMode；SCWidth；SCHeight；IMEI；IMSI
  - 响应：同上 `resp.DataResponse`
  - 内部转为 `events.AppUseTimesEvent` 经 EventService.ReportEvent 落库

## 4. 风险与注意点

- **公共路径无鉴权前缀**：controllers/event_controller.go:22-23（路径含 `/center/public/`，命名示意外部可访问，仅 OverLoadFilter 限流，无鉴权 filter）
