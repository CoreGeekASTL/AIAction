# 客户端事件上报

> 功能域：客户端事件上报　接口数：2　所属 server：外部 + 内部
> 子文档 of [README.md](README.md)

## 1. 定位

接收云浏览器终端上报的客户端事件与应用使用时长事件，用于行为审计与运营统计。两个接口均经 externalServer 与 innerServer 双暴露。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| SendClientEvent | 上报终端客户端事件 | controllers/event_controller.go | POST /app-api/center/public/client/sendClientEvent |
| SendAppUseTimesEvent | 上报应用使用时长事件 | controllers/event_controller.go | POST /app-api/center/public/client/sendAppUseTimesEvent |

## 3. 数据结构说明

- **SendClientEvent**
  - 请求 `req.ClientEventRequest`（models/req/event_request.go）：HSMan（厂商）、HSType（机型）、AppType、IMEI、IMSI、Type（事件类型）
  - 响应 `resp.BaseResponse`（code/msg）
- **SendAppUseTimesEvent**
  - 请求 `req.AppUseTimesEvent`（models/req/event_request.go）：UseTimes（使用时长）、HSMan、HSType、EXTType、AppType、AppId、SCHeight/SCWidth（分辨率）、IMEI、IMSI、PlayMode
  - 响应 `resp.BaseResponse`
