# 事件上报

> 功能域概述：接收云浏览器客户端上报的客户端异常事件与应用使用时长事件，封装为统一事件信封后经事件存储工厂写入本地事件日志文件（auditlog 组件落盘，不对外转发）。
> 接口数：2（外部 2 / 内部 2，同一 Controller 双侧注册）　核心模块：controllers.EventController, service.EventService, common/event（Storage/StorageFactory + localEventStorage）

## 1. 功能故事（多彩建模）

实现逻辑速览：
- 收到上报后按事件类型装配成统一格式。
- 逐行写入本地滚动日志文件，超限自动切分存档。
- 入参不做校验，客户端字段原样落盘。

```mermaid
flowchart LR
    classDef mi fill:#ffd1dc,stroke:#c2185b,color:#880e4f
    classDef role fill:#fff3b0,stroke:#f9a825,color:#5d4037
    classDef ppt fill:#c8e6c9,stroke:#388e3c,color:#1b5e20
    classDef desc fill:#bbdefb,stroke:#1976d2,color:#0d47a1

    Client["云浏览器客户端（外部系统）"]:::role
    Gids["GIDS 本服务"]:::role
    Audit["auditlog 日志框架"]:::role

    E1["事件：客户端异常上报<br/>触发者：云浏览器客户端<br/>输入：终端厂商/型号、IMEI、异常类型<br/>输出：客户端异常事件<br/>后继：按类型装配统一事件"]:::mi
    E2["事件：应用使用时长上报<br/>触发者：云浏览器客户端<br/>输入：使用时长、应用标识、分辨率<br/>输出：使用时长事件<br/>后继：按类型装配统一事件"]:::mi
    E3["事件：按类型装配统一事件<br/>触发者：GIDS 本服务<br/>输入：上报参数<br/>输出：统一格式的事件记录<br/>后继：写入本地审计存储"]:::mi
    E4["事件：写入本地审计存储<br/>触发者：GIDS 本服务<br/>输入：统一格式事件的一行文本<br/>输出：滚动日志文件追加一行<br/>后继：超限时滚动切分并清理旧文件"]:::mi

    Rec["事件记录"]:::ppt
    LogFile["滚动日志文件"]:::ppt

    R1["规则：事件类型须预注册才能装配"]:::desc
    R2["规则：单文件超 20MB 滚动切分<br/>转储最多留 5 份、保留 90 天"]:::desc
    R3["规则：入参校验为空实现<br/>客户端字段原样透传落盘"]:::desc

    Client --> E1
    Client --> E2
    E1 --> E3
    E2 --> E3
    Gids --> E3
    E3 --> Rec
    Gids --> E4
    Rec --> E4
    E4 --> Audit
    Audit --> LogFile
    R1 -.约束.-> E3
    R2 -.约束.-> E4
    R3 -.约束.-> E1
    R3 -.约束.-> E2
```

### 术语表

| 术语 | 人话解释 | 出处 |
| --- | --- | --- |
| 事件信封（events.Info） | 所有事件共用的外层包装：事件名、时间、来源服务固定为 GIDS，载荷挂在 EventData 上 | src/models/events/base.go |
| 事件描述表（eventTypeMap） | 事件类型到名称/描述/触发方的登记表，新类型必须先登记才能被装配 | src/models/events/base.go |
| 事件存储工厂（StorageFactory） | 按字符串名字挑选存储实现的注册表，当前只注册了本地文件这一种 | src/common/event/event_storage.go |
| 本地审计存储（localEventStorage） | 把事件逐行追加写进本地日志文件的实现，文件不可用时降级打到控制台 | src/common/event/local_storage.go |
| 滚动切分（rollOver） | 单文件超过 20MB 时把旧内容压缩成 zip 存档、清空原文件继续写 | src/common/event/local_storage.go |
| 转储文件清理（FileDeleter） | 定期删掉过期存档：最多留 5 份、最多保留 90 天 | src/common/event/local_storage.go |
| 事件文件路径（EventFile） | 启动参数指定事件日志写到哪个文件 | src/common/conf/config.go |
| 双侧注册 | 同一组接口同时挂在对外端口和对内 127.0.0.1:9090 端口上 | src/routers/beego_router.go |
| 空校验（Validate） | 请求结构自带的校验方法是空函数，收到什么字段就原样落盘 | src/models/req/event_request.go |
| 事件是否转发到外部审计服务 | 代码中未体现：仅见本地文件/控制台两种去处 | src/common/event/local_storage.go |

## 2. 模块划分

```mermaid
graph LR
    R["routers/beego_router.go"] --> C["controllers/event_controller.go"]
    C --> S["service/event_service.go"]
    C --> M["models/events/base.go"]
    C --> Q["models/req/event_request.go"]
    S --> E["common/event/event_storage.go"]
    E --> L["common/event/local_storage.go"]
    L --> A["auditlog.WriterSink 外部库"]
    L --> M
    S --> F["common/conf/config.go"]
```

| 模块/包 | 承载功能 | 证据 |
| --- | --- | --- |
| routers | 双侧（外部+内部 127.0.0.1:9090）注册 EventController | src/routers/beego_router.go |
| controllers | HTTP 入口：解析请求体、组装事件、调用 EventService | src/controllers/event_controller.go |
| service | 业务门面：初始化存储工厂、选取存储实现、透传 Record；StorageFactory 按 location 键选实现，仅注册 "localAuditComponent" 且取用键写死为 DefaultEventStorage，无配置化，恒为本地文件存储 | src/service/event_service.go |
| common/event | 存储抽象与工厂 + 本地文件存储与转储清理 | src/common/event/event_storage.go；src/common/event/local_storage.go |
| models/events | 事件类型枚举、描述表、统一信封 Info 与载荷结构 | src/models/events/base.go |
| models/req | 上报请求体结构与校验（空实现） | src/models/req/event_request.go |

## 3. 接口清单

### HTTP 接口

| 接口 | 路径/入口 | 请求结构 | 响应结构 | 状态 |
| --- | --- | --- | --- | --- |
| SendClientEvent | POST /app-api/center/public/client/sendClientEvent；路由 src/controllers/event_controller.go；注册 src/routers/beego_router.go；入口 src/controllers/event_controller.go | req.ClientEventRequest（src/models/req/event_request.go） | resp.DataResponse{Code, Message:"record success", Data:true}；失败 400 + retcode.ClientFailed（src/controllers/event_controller.go） | 在用 |
| SendAppUseTimesEvent | POST /app-api/center/public/client/sendAppUseTimesEvent；路由 src/controllers/event_controller.go；注册 src/routers/beego_router.go；入口 src/controllers/event_controller.go | req.AppUseTimesEvent（src/models/req/event_request.go） | 同上（src/controllers/event_controller.go） | 在用 |

## 4. 关键数据结构

| 结构 | 定义位置 | 关键字段（含义+约束） |
| --- | --- | --- |
| events.Info | src/models/events/base.go | EventDesc（内嵌事件名/描述/触发方）；Service 固定 "GIDS"（src/common/constants/base.go）；EventTime 格式 "2006-01-02 15:04:05"；Object 固定 "cloud-browser"；Env/Hostname 恒为空串；EventData interface{} 承载载荷 |
| events.EventType / eventTypeMap | src/models/events/base.go | 四类事件；本域用 Client="browser_client_error"、AppUseTimes="browser_client_app_use_times" |
| events.ClientEventData | src/models/events/base.go | HSMan/HSType、AppType、IMEI/IMSI、Type（异常类型）；全字符串无必填约束 |
| events.AppUseTimesEvent | src/models/events/base.go | UseTimes（时长，字符串）、AppId、EXTType（json tag "exttype"）、SCWidth/SCHeight、PlayMode、IMEI/IMSI；无约束 |
| req.ClientEventRequest | src/models/req/event_request.go | 同 ClientEventData；Validate 空实现恒 nil |
| req.AppUseTimesEvent | src/models/req/event_request.go | 同 AppUseTimesEvent；Validate 空实现 |
| event.FileDeleter | src/common/event/local_storage.go | DayBasedDeleter 保留 90 天（FileRemainDay）；CountBasedDeleter 最多 5 个（FileMaxNum） |

## 5. 调用关系

以 SendClientEvent 为代表（SendAppUseTimesEvent 同构，src/controllers/event_controller.go）：

```mermaid
sequenceDiagram
    participant Client
    participant Router as routers/beego_router.go
    participant Ctrl as EventController<br/>event_controller.go
    participant Base as BaseController<br/>controller.go
    participant Svc as EventServiceImpl<br/>event_service.go
    participant Fac as StorageFactory<br/>event_storage.go
    participant Store as localEventStorage<br/>local_storage.go
    participant Sink as auditlog.WriterSink

    Client->>Router: POST /sendClientEvent
    Router->>Ctrl: 路由分发（beego_router.go）
    Ctrl->>Svc: Prepare 中 NewEventService（event_controller.go）
    Svc->>Fac: sync.Once 注册 localAuditComponent + Get(DefaultEventStorage)<br/>（event_service.go）
    Ctrl->>Base: RequestBodyUnmarshalTo（event_controller.go）
    Base->>Base: JSON 反序列化 + Validate（空实现，event_request.go）
    Ctrl->>Ctrl: events.NewInfo(Client) + SetEventData<br/>（event_controller.go、models/events/base.go）
    Ctrl->>Svc: ReportEvent（event_controller.go）
    Svc->>Store: Record（event_service.go）
    Store->>Store: needRollOver>20MB? rollOver 转储 zip<br/>（local_storage.go）
    Store->>Sink: engine.Print(event.ToJSON()) 追加写 EventFile<br/>（local_storage.go；配置 common/conf/config.go）
    Store-->>Ctrl: err?
    Ctrl-->>Client: 200 DataResponse / 400 ClientFailed（event_controller.go）
```

- 无外发：auditlog 仅作 WriterSink 写本地文件/stdout，无向事件服务转发的 HTTP 逻辑（src/common/event/local_storage.go）。
- 同步逐条写入，无批量/异步队列/重试；Record 失败仅记日志并返回 err（src/common/event/local_storage.go）。
- 仅清理异步：goroutine 每小时删过期转储 zip，rollOver 末尾同步再清一次（src/common/event/local_storage.go）。

## 6. 框架引用

| 基础框架 | 框架文档 | 本功能中的用途（引用文件） |
| --- | --- | --- |
| Beego Web | ../framework-usage/rpc-beego-web.md | 路由注册与 Controller 基座：EventController 内嵌 beego.Controller，经 beego.Router 双侧挂载（src/controllers/event_controller.go、src/routers/beego_router.go） |
| auditlog 审计日志 | ../framework-usage/log-lager-auditlog-event.md | 事件落盘：auditlog.NewLoggerBase + NewWriterSink 将事件 JSON 追加写入 EventFile 或 stdout（src/common/event/local_storage.go） |
| JSON 编解码 | ../framework-usage/codec-json-yaml.md | 请求体 JSON 反序列化与事件 ToJSON 序列化落盘（src/controllers/controller.go、src/models/events/base.go） |
| goroutine 并发 | ../framework-usage/concurrency-goroutine.md | 后台 goroutine 每小时触发过期转储文件清理（src/common/event/local_storage.go） |

## 7. AI 编码指南

- 新事件类型注册eventTypeMap并复用NewInfo（src/models/events/base.go）。
- 新存储实现须同步改注册键与取用键（src/service/event_service.go）。
- 请求Validate为空，校验在Service层补（src/models/req/event_request.go）。
- Record/rollOver无锁，勿加并发状态（src/common/event/local_storage.go）。
