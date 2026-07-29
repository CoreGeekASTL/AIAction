# 终端登录鉴权

> 功能域：终端登录鉴权　接口数：3（双实现）　所属 server：外部(HTTPS) + 内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

终端登录鉴权与浏览器预开。同一组 3 个路径经 externalServer（HTTPS，ExLoginController）与 innerServer（HTTP，LoginController）双暴露，两侧实现逻辑一致。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| GridLoginAuth | 网格登录鉴权 | controllers/exlogin_controller.go、controllers/login_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuth |
| GridLoginAuthOpenBrowser | 登录鉴权并预开浏览器 | controllers/exlogin_controller.go、controllers/login_controller.go | POST /app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser |
| DeviceLoginAuth | 设备登录鉴权（TikTok 类型经沐恩二次鉴权） | controllers/exlogin_controller.go、controllers/login_controller.go | POST /app-api/devicetcp/app/login/v1/deviceLoginAuth |

## 3. 数据结构说明

- **GridLoginAuth / GridLoginAuthOpenBrowser / DeviceLoginAuth**
  - 请求 `req.LoginAuthRequest`（models/req/request_entity.go）：IMEI；IMSI；AppType（"2"=TikTok，触发沐恩二次鉴权，constants/base.go:19）；Manufacturer；Model；ExtendModel；DeviceType；TotalKb；FreeKb。`Validate()` 为空实现（models/req/request_entity.go:50），controller 层不做格式校验
  - 响应 `resp.DeviceLoginAuthResponse`（models/resp/response_entity.go）：BaseResponse{code,msg} + Data=`resp.LoginInfo`（内嵌 AuthInfo{Token, ExpiresTime, TimeAxis} + AssignInfo{TcpAddr, TlsTcpAddr, VideoMode, ShortAddr, NodeGateWayURL, HttpsShortAddr, HttpsNodeGateWayUrl, NodeIntranetWayURL, NodeCapacity}）
  - GridLoginAuth / GridLoginAuthOpenBrowser 返回前将 TcpAddr/TlsTcpAddr/VideoMode/ShortAddr/HttpsShortAddr/NodeIntranetWayURL 置空；DeviceLoginAuth 仅置空 NodeIntranetWayURL（controllers/exlogin_controller.go:47-52/61-66/75）
  - DeviceLoginAuth 当 AppType=TikTok 时经沐恩云服务二次鉴权（service/remote_service.go:18 MuenDeviceLogin），失败返回 code=-1
  - 登录成功后上报 Login 事件（EventService.ReportEvent），事件上报失败不影响登录结果

## 4. 风险与注意点

- **同组接口双暴露**：routers/beego_router.go:20/34（3 个 login 路径在 HTTPS 外部与 HTTP 内部重复注册，鉴权策略若不一致成攻击面）
- **入参无格式校验**：models/req/request_entity.go:50（`LoginAuthRequest.Validate()` 返回 nil，IMEI/IMSI 位数/字符无校验，依赖下游 service 容错）
