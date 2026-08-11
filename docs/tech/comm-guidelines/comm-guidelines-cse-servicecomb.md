# cse-servicecomb 通信规范

## 接口清单

| 接口名 | 协议 | 调用位置 | 业务场景 |
|---|---|---|---|
| WatchMicroServiceV1（browser-gateway） | go-chassis GSF Registry API | src/common/cse/cse.go（Init / watchServiceCallBack） | Watch browser-gateway 实例变更 |
| GetAllMicroServiceInstanceInfo | go-chassis GSF Registry API | src/common/cse/cse.go；src/dao/db_init.go | 查询指定微服务实例列表 |
| UpdateMicroServiceInstanceProperties | go-chassis GSF Registry API | src/common/cse/cse.go（Report） | 上报本实例 chainEndpoints 属性 |

## RPC / 平台 API

### WatchMicroServiceV1（browser-gateway）

- 业务场景：服务启动初始化时（cse.Init），订阅 `browser-gateway` 微服务实例变更事件，维护本地实例缓存 browserGWInstances，供 preOpen/插件加载/缓存删除选取目标实例
- 接口功能：Watch 服务名 "browser-gateway"（AppId 取环境变量 `APPID`，Version "0+"），回调处理 CREATE/UPDATE/DELETE/LIST 事件，从实例 Properties["status"] 反序列化出 ServiceInstance
- 调用位置：src/common/cse/cse.go（Init；browserGWNotifier.WatchServiceCallBack）
- 协议信息：
  - 协议：go-chassis GSF `api.Registry.WatchMicroServiceV1`（ServiceComb 注册中心）
  - 封装方式：平台 SDK 直连（Go-chassis-extend/api/GSF/api），本仓封装层 src/common/cse/cse.go
  - 超时重试：未设置（SDK 框架默认）
  - 错误码处理：Watch 失败仅记录日志（不中断启动）；事件解析失败记录日志并跳过该实例

### GetAllMicroServiceInstanceInfo

- 业务场景：按服务名查询微服务实例列表，当前用于启动阶段发现 DB 服务（GaussDB 管理服务）实例以获取 host/port
- 接口功能：入参服务名，返回 []base.MicroServiceInstance
- 调用位置：src/common/cse/cse.go（GetAllMicroServiceInstanceInfo）；src/dao/db_init.go
- 协议信息：
  - 协议：go-chassis GSF `api.Registry.GetAllMicroServiceInstanceInfo`
  - 封装方式：平台 SDK 直连，本仓封装层 src/common/cse/cse.go
  - 超时重试：未设置（SDK 框架默认）
  - 错误码处理：error 直接上抛调用方

### UpdateMicroServiceInstanceProperties

- 业务场景：服务注册阶段，将本实例的 HTTP/HTTPS 监听地址（chainEndpoints）上报到注册中心实例属性
- 接口功能：`UpdateMicroServiceInstanceProperties(selfServiceID, selfInstanceID, {"chainEndpoints": ...})`
- 调用位置：src/common/cse/cse.go（Report）；src/main.go（cse.NewCse().Report(maxRetryTimes)）
- 协议信息：
  - 协议：go-chassis GSF `api.Registry.UpdateMicroServiceInstanceProperties`
  - 封装方式：平台 SDK 直连，本仓封装层 src/common/cse/cse.go
  - 超时重试：失败后 sleep 30s 递归重试（`retryIntervalSec`），次数由调用方 maxRetry 控制；重试耗尽时 logger.Fatalf 退出进程
  - 错误码处理：err 非空即重试；无错误码分类
