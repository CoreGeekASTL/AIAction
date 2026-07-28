# 流量统计

> 功能域概述：接收实例侧上报的会话/媒体流/控制流统计数据并落库（GaussDB），支持按月导出 CSV 压缩包，并由每日定时任务清理 3 个月前的过期数据。
> 接口数：5（外部 0 / 内部 4 + 定时任务 1）　核心模块：controllers, service, dao, scheduler

## 1. 功能故事（多彩建模）

**实现逻辑速览**

1. 统计数据逐条校验后批量入库，重复上报只更新结束时间。
2. 按月前缀匹配当月记录，分批导出 CSV 打包返回。
3. 每日凌晨两点清理三个月前数据，失败重试至多三次。

```mermaid
flowchart LR
    classDef mi fill:#ffd1dc,stroke:#c2185b,color:#4a0e24
    classDef role fill:#fff3b0,stroke:#f9a825,color:#4d3800
    classDef ppt fill:#c8e6c9,stroke:#388e3c,color:#123a15
    classDef desc fill:#bbdefb,stroke:#1976d2,color:#0d2b4d

    subgraph 上报线
        A[浏览器实例<br/>外部系统]:::role
        E1[上报会话/媒体流/控制面流量统计]:::mi
        P1[会话统计记录]:::ppt
        P2[媒体流流量记录]:::ppt
        P3[控制面流量记录]:::ppt
        P4[按月区分的统计数据存储]:::ppt
        R1[规则: started_at 必须为<br/>YYYY-MM 前缀字符串]:::desc
        A -->|触发者| E1
        A -->|输入: 会话/媒体/控制面统计数据| E1
        E1 -->|输出: 逐条校验入库| P1
        E1 -->|输出: 逐条校验入库| P2
        E1 -->|输出: 逐条校验入库| P3
        P1 --> P4
        P2 --> P4
        P3 --> P4
        E1 -->|后继: 数据可供导出与清理| P4
        R1 -.约束.-> P4
    end

    subgraph 清理线
        B[调度器<br/>内部系统]:::role
        E2[每日凌晨 2 点定时清理过期统计数据]:::mi
        R2[规则: 清理窗口 CleanupMonths<br/>= 3 个月]:::desc
        R3[规则: 失败重试 3 次<br/>间隔 10 分钟]:::desc
        B -->|触发者| E2
        B -->|输入: 当前时间与清理窗口| E2
        E2 -->|输出: 删除三类表中过期记录| P4
        E2 -->|后继: 失败进入重试| R3
        R2 -.约束.-> E2
    end

    subgraph 导出线
        C[管理面调用方<br/>身份代码中未体现]:::role
        E3[按月导出静态统计数据]:::mi
        OUT[当月三张 CSV 打包的压缩文件]:::ppt
        C -->|触发者| E3
        C -->|输入: 月份参数 YYYY-MM| E3
        E3 -->|读取| P4
        E3 -->|输出| OUT
        E3 -->|后继: 压缩包返回给调用方| C
    end
```

### 术语表

| 术语 | 人话解释 | 出处 |
|---|---|---|
| 会话统计（SessionStats） | 一次云浏览器会话从接入到结束的过程记录，含接入方式、网络协议、起止时间 | src/models/db/traffic_stats.go |
| 媒体流流量统计（MediaTrafficStats） | 会话音视频通道产生了多少字节流量的记录 | src/models/db/traffic_stats.go |
| 控制面流量统计（ControlTrafficStats） | 会话控制信令通道产生了多少字节流量的记录 | src/models/db/traffic_stats.go |
| tcp_unique_id（幂等键） | 同一条 TCP 连接的唯一编号，重复上报时只更新结束时间、不重复插入 | src/dao/traffic_stats_dao.go |
| 批量上报（MultiTableRequest） | 一次请求携带多条流量记录，逐条校验后在事务里按 1000 条一批写入 | src/models/req/request_entity.go、src/service/traffic_stats_service.go |
| 按月导出（exportStaticData） | 传入"年-月"，把当月三张统计表导成 CSV 打成 zip 返回 | src/controllers/traffic_stats_controller.go |
| 清理窗口（CleanupMonths） | 统计数据保留几个月，当前固定为 3 个月，超期即被定时任务删除 | src/common/constants/base.go |
| 内部监听（127.0.0.1:9090） | 本功能域接口只对本机内部开放，不对外暴露 | src/routers/beego_router.go |
| started_at 前缀过滤 | 用字符串开头匹配判断记录属于哪个月，无物理按月分区 | src/service/traffic_stats_service.go |
| 上报方身份/鉴权方式 | 上报与导出的调用方身份、鉴权机制代码中未体现 | 代码中未体现 |

## 2. 模块划分

```mermaid
graph LR
    main["main.go"]
    router["routers/beego_router.go<br/>(路由注册)"]
    ctrl["controllers/traffic_stats_controller.go<br/>(TrafficStatsController)"]
    sched["scheduler/task_scheduler.go<br/>(DataCleanupScheduler)"]
    svc["service/traffic_stats_service.go<br/>(TrafficStatsServiceImpl)"]
    mock["service/traffic_stats_service_mock.go<br/>(Mock, 仅测试)"]
    dao["dao/traffic_stats_dao.go<br/>(三表 Dao)"]
    basedao["dao/base_dao.go<br/>(事务/批量插入)"]
    db[("models/db/traffic_stats.go<br/>GaussDB 三表")]
    req["models/req, models/resp"]
    consts["common/constants<br/>(CleanupMonths)"]

    main --> sched
    router --> ctrl
    ctrl --> svc
    ctrl --> req
    sched --> svc
    sched --> consts
    svc --> dao
    dao --> basedao
    dao --> db
    mock -.仅单测注入.-> svc
```

| 模块/包 | 承载功能 |
|---|---|
| controllers（TrafficStatsController） | 4 个 HTTP 入口：路由注册、解析校验、导出 zip 组装（src/controllers/traffic_stats_controller.go） |
| service（TrafficStatsService/Impl） | 业务编排：会话插入或更新、批量事务插入、CSV 分批导出、过期清理（src/service/traffic_stats_service.go）；接口另有 Mock 实现（src/service/traffic_stats_service_mock.go），无运行时选择机制，生产固定走真实实现，Mock 仅在单元测试中字段注入（src/service/monitor_service_test.go） |
| dao（三张统计表 Dao） | t_session_stats / t_media_traffic_stats / t_control_traffic_stats 的 CRUD 与原生 SQL 查询（src/dao/traffic_stats_dao.go） |
| scheduler（DataCleanupScheduler） | 每日凌晨 2:00 调度清理，失败重试 3 次（src/scheduler/task_scheduler.go） |
| models/db | 三张统计表 ORM 实体（src/models/db/traffic_stats.go） |
| common/constants | 清理窗口常量 CleanupMonths = 3（src/common/constants/base.go） |

## 3. 接口清单

| 接口 | 路径/入口（含注册处） | 请求结构 | 响应结构 | 状态 |
|---|---|---|---|---|
| 会话统计上报 | POST /stats/v1/session（src/controllers/traffic_stats_controller.go，注册 src/routers/beego_router.go） | db.SessionStats（src/models/db/traffic_stats.go） | resp.BaseResponse OK（src/controllers/controller.go） | 在用（内部 127.0.0.1:9090） |
| 媒体流流量统计上报 | POST /stats/v1/traffic/media（src/controllers/traffic_stats_controller.go） | req.MultiTableRequest{Items: []MediaTrafficStats}（src/models/req/request_entity.go、src/models/db/traffic_stats.go） | resp.BaseResponse OK / Failed（src/controllers/controller.go） | 在用（内部） |
| 控制流流量统计上报 | POST /stats/v1/traffic/control（src/controllers/traffic_stats_controller.go） | req.MultiTableRequest{Items: []ControlTrafficStats}（src/models/db/traffic_stats.go） | resp.BaseResponse OK / Failed（src/controllers/controller.go） | 在用（内部） |
| 按月导出统计数据 | GET /stats/v1/exportStaticData/:month（src/controllers/traffic_stats_controller.go） | 路径参数 month，格式 "2006-01"（src/controllers/traffic_stats_controller.go） | zip 附件流（3 个 CSV，src/controllers/traffic_stats_controller.go） | 在用（内部） |
| 过期统计数据清理（定时任务） | 非 HTTP，每日凌晨 2:00（src/scheduler/task_scheduler.go，启动 src/main.go） | constants.CleanupMonths=3（src/common/constants/base.go） | 删除三表过期行，日志记录（src/service/traffic_stats_service.go） | 在用 |

## 4. 关键数据结构

| 结构 | 定义位置 | 关键字段（含义+约束） |
|---|---|---|
| db.SessionStats | src/models/db/traffic_stats.go | SessionID 会话标识；TcpUniqueId 为幂等键——存在则仅更新 finished_at（src/dao/traffic_stats_dao.go）；StartedAt/FinishedAt 为 string，是按月过滤与清理的比较基准 |
| db.MediaTrafficStats | src/models/db/traffic_stats.go | SessionID、AppType、OutBytes（int64 字节数）；Validate() 空实现，入库无字段校验 |
| db.ControlTrafficStats | src/models/db/traffic_stats.go | 字段同媒体表；FinishedAt/OutBytes 带 omitempty；Validate() 空实现 |
| req.MultiTableRequest | src/models/req/request_entity.go | Items []json.RawMessage；Validate 要求非空，但 Controller 未调用该校验（src/controllers/traffic_stats_controller.go） |
| service.SQLConfig / Res | src/service/traffic_stats_service.go | 外部化 SQL 配置（YAML 文件加载）；configPath 为空时 sqlConfig=nil，监控查询静默返回 nil |
| constants.CleanupMonths | src/common/constants/base.go | =3（编译期常量，非配置项）；截止点 time.Now().AddDate(0,-3,0)，按 started_at__lt 字符串比较删除（src/service/traffic_stats_service.go）；表无物理按月分区，"按月"仅为字符串前缀过滤 |

## 5. 调用关系

```mermaid
sequenceDiagram
    participant C as Client(实例侧)
    participant R as routers/beego_router.go
    participant CT as TrafficStatsController
    participant S as TrafficStatsServiceImpl
    participant D as Dao(traffic_stats_dao.go)
    participant DB as GaussDB 三表
    participant SCH as DataCleanupScheduler

    C->>R: POST /stats/v1/session
    R->>CT: SessionStats
    CT->>CT: Prepare 构造 Dao+Service
    CT->>S: HandleSessionStats
    S->>D: Exist(tcp_unique_id)
    D->>DB: SELECT EXISTS
    alt 不存在
        S->>D: Insert
    else 已存在(幂等)
        S->>D: UpdatebySession 仅更新 finished_at
    end
    D->>DB: 写入 t_session_stats

    C->>R: POST /stats/v1/traffic/media|control
    R->>CT: insertMultiData(tag)
    CT->>S: BatchInsertStats
    S->>S: 按 tag 分发+逐条 Validate
    S->>D: DoTxWithCtx 内 InsertMultiWithOrm 每批1000
    D->>DB: 批量写入 t_media/control_traffic_stats

    Note over SCH: main.go 启动调度循环
    loop 每日凌晨 2:00
        SCH->>SCH: calculateNextRunTime→timer 等待
        SCH->>S: NewTrafficStatsService("").CleanOldStats(3)
        S->>D: 三表事务内 DeleteByCond started_at__lt 截止点
        D->>DB: 删除 3 个月前数据
        alt 失败
            SCH->>SCH: 间隔 10min 重试，最多 3 次
        end
    end
```

- 导出分支：GET exportStaticData/:month 校验格式 → exportCSVFiles 导出三表 → QueryStatsDataAndWriteCSV 按 `started_at__istartswith month` 分批 1000 写 CSV → zip 打包流式写回（src/controllers/traffic_stats_controller.go、src/service/traffic_stats_service.go）。
- 幂等分支：会话上报按 tcp_unique_id 判重，重复上报只更新 finished_at 不重复插入（src/dao/traffic_stats_dao.go）。
- 清理重试分支：失败间隔 10 分钟重试、上限 3 次，重试期间可被 stopChan 中断（src/scheduler/task_scheduler.go）。

## 6. 框架引用

| 基础框架 | 框架文档 | 本功能中的用途（引用文件） |
|---|---|---|
| Beego Web（路由/Controller） | [rpc-beego-web.md](../framework-usage/rpc-beego-web.md) | 四个统计接口的路由注册与 Controller 请求处理（src/routers/beego_router.go、src/controllers/traffic_stats_controller.go） |
| Beego ORM（数据存储） | [storage-beego-orm.md](../framework-usage/storage-beego-orm.md) | 三张统计表的读写、批量事务插入与过期删除（src/dao/traffic_stats_dao.go、src/dao/base_dao.go、src/service/traffic_stats_service.go） |
| Timer 定时调度 | [schedule-timer.md](../framework-usage/schedule-timer.md) | 每日凌晨 2:00 清理调度循环与 10 分钟重试等待（src/scheduler/task_scheduler.go） |
| lager 业务日志 | [log-lager-auditlog-event.md](../framework-usage/log-lager-auditlog-event.md) | 上报/导出/清理全链路 Infof/Errorf/Warnf 日志（src/service/traffic_stats_service.go、src/scheduler/task_scheduler.go） |

说明：CleanupMonths 为编译期常量（src/common/constants/base.go），非 appconf/configcenter 配置项；监控查询 SQL 配置经 YAML 文件加载（src/service/traffic_stats_service.go），均未使用配置中心框架，故未列入。

## 7. AI 编码指南

- 批量上报 tag 与 Service case 字面量须一致（src/controllers/traffic_stats_controller.go vs src/service/traffic_stats_service.go）
- started_at 必须为 YYYY-MM 前缀字符串（src/service/traffic_stats_service.go）
- Mock 仅单测注入，新接口方法需自行扩展（src/service/monitor_service_test.go、src/service/traffic_stats_service.go）
- MultiTableRequest 校验未生效（src/controllers/traffic_stats_controller.go、src/service/traffic_stats_service.go）
- 改路由需同步 RouteMapping 及 routers（src/controllers/traffic_stats_controller.go、src/routers/beego_router.go）
