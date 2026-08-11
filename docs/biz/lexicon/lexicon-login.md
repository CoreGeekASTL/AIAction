# 登录鉴权与用户绑定 领域词典

> 子文档 of [lexicon.md](lexicon.md)；词汇口径、来源说明与全仓待确认清单见主文档。子域锚点 `login`，功能域口径与 docs/biz/interface/ 一致。

## 实体与业务概念

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| IMEI | 国际移动设备识别码，终端身份之一（15 位纯数字） | 登录请求字段；cache 删除缓存参数；event 上报字段 | `req.UserIdentity.IMEI`，`src/models/req/request_entity.go`；`req.DeleteCacheRequest.IMEI`，`src/models/req/request_entity.go` |
| IMSI | 国际移动用户识别码，终端身份之一（15 位纯数字）（同义：与 IMEI 合称用户身份 UserIdentity） | 同上 | `req.UserIdentity.IMSI`，`src/models/req/request_entity.go` |
| UserIdentity | 用户身份组合（IMSI+IMEI），内嵌于登录请求 | - | `req.UserIdentity`，`src/models/req/request_entity.go` |
| LoginAuthRequest | 登录鉴权请求：用户身份 + 设备信息（Manufacturer/Model/AppType/ExtendModel/Country/Platform/Width/Height/MCC/MNC/Lac/CI/Rxlev/TotalKb/FreeKb/ClientLanguage/DeviceType） | - | `req.LoginAuthRequest`，`src/models/req/request_entity.go` |
| User | 终端设备档案表（t_user）：厂商/型号/扩展型号/国家/平台/分辨率/MCC/MNC/设备类型 | - | `db.User`，`src/models/db/user.go` |
| UserBind | 用户绑定关系表（t_user_bind）：session 到浏览器实例与各端点的绑定，含 Token 与心跳时间 | - | `db.UserBind`，`src/models/db/user.go` |
| BrowserInstance | 浏览器实例标识，UserBind 绑定目标 | - | `db.UserBind.BrowserInstance`，`src/models/db/user.go` |
| 媒体/控制端点 | 浏览器实例对外端点：MediaEndpoint/ControlEndpoint/MediaTlsEndpoint/ControlTlsEndpoint/InnerMediaEndpoint/InnerBrowserEndpoint | - | `db.UserBind` 各端点字段，`src/models/db/user.go` |
| Token | 登录鉴权通过后下发的访问令牌 | - | `db.UserBind.Token`，`src/models/db/user.go`；`resp.AuthInfo.Token`，`src/models/resp/response_entity.go` |
| AuthInfo | 鉴权信息（Token/ExpiresTime/TimeAxis） | - | `resp.AuthInfo`，`src/models/resp/response_entity.go` |
| AssignInfo | 分配信息（TcpAddr/TlsTcpAddr/VideoMode/ShortAddr/NodeGateWayURL/HttpsShortAddr/HttpsNodeGateWayUrl/NodeIntranetWayURL/NodeCapacity） | - | `resp.AssignInfo`，`src/models/resp/response_entity.go` |
| LoginInfo | 设备登录响应数据 = AuthInfo + AssignInfo 组合 | - | `resp.LoginInfo`，`src/models/resp/response_entity.go` |
| GridLoginAuthResponse | 网格登录响应，Data 含 Token/ExpiresTime/NodeGateWayURL/NodeIntranetWayURL/NodeCapacity/TimeAxis | - | `resp.GridLoginAuthResponse`，`src/models/resp/response_entity.go` |
| ServiceInstance | 浏览器服务实例（内外端点、容量 Cap/Used、插件状态、健康检查信息），按使用率排序 | - | `browsergateway.ServiceInstance`，`src/models/browsergateway/service_instance.go` |
| InitBrowserRequest | 预开浏览器请求（factory/dev_type/ext_type/plat_type/分辨率/app_type/appid/imsi/imei/device_type/client_language/play_mode） | - | `browsergateway.InitBrowserRequest`，`src/models/browsergateway/req.go` |
| Heartbeats | 用户绑定心跳时间（映射 updated_at 列） | - | `db.UserBind.Heartbeats`，`src/models/db/user.go` |
