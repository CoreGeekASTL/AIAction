# 终端登录鉴权

> 功能域：终端登录鉴权　接口数：3　所属 server：外部(HTTPS) + 内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

终端登录鉴权与浏览器预开。同一组 3 个路径经 externalServer（HTTPS，ExLoginController）与 innerServer（HTTP，LoginController）双暴露。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| GridLoginAuth | 网格登录鉴权 | controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth |
| GridLoginAuthOpenBrowser | 登录鉴权并预开浏览器 | controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser |
| DeviceLoginAuth | 设备登录鉴权 | controllers/exlogin_controller.go | POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |

## 3. 数据结构说明

- **GridLoginAuth / GridLoginAuthOpenBrowser / DeviceLoginAuth**
  - 请求 `req.LoginAuthRequest`（models/req）：IMEI（15 位纯数字，必填）；IMSI（15 位纯数字，必填）；Manufacturer；Model；ExtendModel；Platform；Width；Height；DeviceType；ClientLanguage
  - 响应 `resp.LoginInfo`（models/resp）：Token；ExpireAt；BrowserEndpoint；TcpAddr；TlsTcpAddr；VideoMode；ShortAddr；HttpsShortAddr；NodeIntranetWayURL（GridLoginAuth 将部分字段置空，见 controllers/exlogin_controller.go）
  - DeviceLoginAuth 经由沐恩云服务二次鉴权（service/remote_service.go:18 MuenDeviceLogin）

## 4. 风险与注意点

- **同组接口双暴露**：controllers/exlogin_controller.go:29（gridLoginAuth 等在 HTTPS 外部与 HTTP 内部重复注册，鉴权策略若不一致成攻击面）
- **IMEI/IMSI 校验依赖**：IMEI/IMSI 必须 15 位纯数字，校验在 service 层而非 controller 层，controller 仅做非空判断
