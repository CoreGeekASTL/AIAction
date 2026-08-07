# 客户端事件上报

> 功能域：client-event　接口数：2　所属 server：外部(HTTPS) + 内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

接收终端 App 上报的客户端行为事件与 App 使用时长事件，统一写入事件通道，内外 server 双暴露。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| SendClientEvent | 上报客户端事件（机型/应用类型/事件类型） | src/controllers/event_controller.go | POST /app-api/center/public/client/sendClientEvent |
| SendAppUseTimesEvent | 上报 App 使用时长事件 | src/controllers/event_controller.go | POST /app-api/center/public/client/sendAppUseTimesEvent |

## 3. 数据结构说明

- **SendClientEvent**
  - 请求 `req.ClientEventRequest`（src/models/req/event_request.go）：HSMan（厂商）、HSType（机型）、AppType、IMEI、IMSI、Type（事件类型）；Validate 为空实现
  - 响应 `resp.DataResponse`：`BaseResponse{code,msg}` + `data: true`；事件经 `events.ClientEventData` 写入 EventService（src/controllers/event_controller.go:40-53）
- **SendAppUseTimesEvent**
  - 请求 `req.AppUseTimesEvent`（src/models/req/event_request.go）：UseTimes（使用时长）、HSMan、HSType、AppType、AppId、EXTType、PlayMode、SCWidth/SCHeight（分辨率）、IMEI、IMSI；Validate 为空实现
  - 响应：同上

