# 输出模板

只产出一份 markdown：`<repo>/docs/interface/external-interfaces.md`。结构固定三段：接口全景 → 各功能域详述 → 风险与注意点。每个功能域下先接口表格（接口名/作用/所在文件/方法或路径），表格下方逐个说明该接口相关的请求与响应数据结构。强调人类阅读友好。

## 模板：docs/interface/external-interfaces.md

```markdown
# 对外接口文档

> 代码仓：<仓库名>　分析基准：<commit/分支>　更新时间：<YYYY-MM-DD>
> 由 spec-feature-analyze 生成，面向人类阅读。
> 范围：本仓对外提供的接口（HTTP 路由 / RPC service / 消息订阅 handler / IDL 契约），不含本仓调用别人的接口。

## 1. 接口全景

一张 mermaid 图：本仓为中心节点，指向各功能域；功能域节点按业务功能命名。

​```mermaid
flowchart LR
  classDef repo fill:#e1f5ff,stroke:#0277bd,color:#000
  classDef http fill:#bbdefb,stroke:#1565c0,color:#000
  classDef rpc fill:#c8e6c9,stroke:#2e7d32,color:#000
  classDef async fill:#e1bee7,stroke:#6a1b9a,color:#000

  Repo[(本仓<br/>GIDS)]:::repo
  Auth[终端鉴权]:::http
  Whitelist[白名单管理]:::http
  Config[配置同步]:::http
  Session[会话管理]:::http
  Sched[定时巡检]:::async

  Repo --> Auth
  Repo --> Whitelist
  Repo --> Config
  Repo --> Session
  Repo --> Sched
​```

统计：共 **4** 个功能域，**18** 个对外接口（HTTP 16 / 消息订阅 0 / 定时任务 2）。

功能域清单：

| 功能域 | 接口数 | 核心模块 | 详见 |
|---|---|---|---|
| 终端鉴权 | 5 | controllers/auth, service/auth | [§2.1](#21-终端鉴权) |
| 白名单管理 | 6 | controllers/whitelist, dao/whitelist | [§2.2](#22-白名单管理) |
| 配置同步 | 4 | controllers/management, service/remote | [§2.3](#23-配置同步) |
| 会话管理 | 3 | controllers/session, service/browser | [§2.4](#24-会话管理) |
| 定时巡检 | 2 | scheduler, dao | [§2.5](#25-定时巡检) |

> 未归类接口：无（若有探测到但无法归入任何功能域的接口，在此逐条列出并说明原因）

## 2. 各功能域详述

每个功能域一节，固定三要素：定位 → 接口表格 → 表格下方逐个接口的数据结构说明。

### 2.1 终端鉴权

**定位**：终端登录与事件上报的鉴权校验，区分 login 路径（拒绝返回 -2）与 event 路径（拒绝返回 401）。

**接口清单**：

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| Login | 终端登录鉴权 | routers/router.go:12 | POST /api/v1/login |
| ReportEvent | 事件上报鉴权 | routers/router.go:15 | POST /api/v1/event |
| QueryAuthState | 查询鉴权态 | routers/router.go:18 | GET /api/v1/auth/state |
| Logout | 终端登出 | routers/router.go:22 | POST /api/v1/logout |
| Heartbeat | 心跳保活 | routers/router.go:25 | POST /api/v1/heartbeat |

**数据结构说明**（对应接口表格逐个接口）：

- **Login**
  - 请求 `LoginAuthRequest`（models/req/login.go:10）：IMEI（15 位纯数字，必填）；IMSI（15 位纯数字，必填）；Manufacturer；Model
  - 响应 `LoginInfo`（models/resp/login.go:15）：Token；ExpireAt；BrowserEndpoint
- **ReportEvent**
  - 请求 `EventRequest`（models/req/event.go:8）：IMEI；IMSI；EventType（枚举：login/logout/heartbeat）；Timestamp
  - 响应：`retcode.AuthFailed(401)`（event 路径拒绝码，controllers/auth/event.go:30）
- **QueryAuthState**
  - 请求：path 参数 `imei`（15 位纯数字）
  - 响应 `AuthState`（models/resp/auth.go:5）：Authed（bool）；ExpireAt
- **Logout / Heartbeat**
  - 请求 `LoginAuthRequest`（同上）
  - 响应：retcode 标准结构（Code/Msg）

### 2.2 白名单管理
（同上结构，略）

### 2.3 配置同步
（同上结构，略）

### 2.4 会话管理
（同上结构，略）

### 2.5 定时巡检
（同上结构，定时任务入口，方法/路径列填 cron 表达式或 @every）

## 3. 风险与注意点

逐条列出接口治理风险，每条带证据：

- **无鉴权**：routers/router.go:30（`/api/v1/public/*` 路径组未挂鉴权中间件）
- **无参数校验**：controllers/whitelist/import.go:45（IMEI 未校验 15 位，可能注入）
- **明文凭据**：controllers/auth/login.go:12（Token 明文返回，未脱敏日志）
- **已下线接口残留**：routers/router.go:40（`/api/v1/legacy/xxx` 路由仍注册但 controllers 标注 deprecated）
- **无超时**：controllers/management/sync.go:20（同步拉取配置无超时，可能挂起）

（无风险点时此节可省略，但需在全景图后注明"未发现显著风险点"）
```

## 撰写硬性要求

- **单文件产出**：只产出 `docs/interface/external-interfaces.md` 一份，禁止拆分多篇。
- **接口全景图**：必须 mermaid `flowchart`，本仓为中心节点；功能域节点按**业务功能命名**（如"终端鉴权""白名单管理"），禁止按技术层命名（如"controller 层"）；用 classDef 统一配色按接口类型区分（HTTP 矩形 / RPC 圆角 / 消息订阅 平行四边形 / 定时任务 子图）；功能域数 ≤ 12 全列，> 12 按业务域聚合。
- **功能域清单表**：每行含功能域名、接口数、核心模块、章节锚点链接；锚点必须能跳到对应详情节，无死链。
- **每功能域一节三要素**：定位（一句话）/ **接口表格** / **表格下方逐个接口数据结构说明**，缺一不可。
- **接口表格列固定**：接口名 | 作用 | 所在文件 | 方法/路径或调用方式（HTTP 填 `POST /api/v1/xxx`，RPC 填 `UserService.GetUser`，消息订阅填 `Consume topic-xxx`，定时任务填 `@every 1m` 或 cron 表达式）。
- **表下数据结构说明格式**：对应接口表格里每个接口，用列表逐个说明其请求与响应数据结构——格式为「结构名（定义位置）：关键字段+约束」；同功能域内多接口共用的结构只在首次出现处详述，后续接口注明"同上"。
- **表格化**：接口清单必须用表格，禁止散文段落描述单个接口。
- **证据锚点**：所有结论附 `文件:行号` 格式证据，相对代码仓根目录。
- **关键字段**：数据结构只列"理解该接口必须知道"的字段与约束，超过 8 字段只列关键，不整段搬运 struct 定义。
- **接口清单精简**：单功能域接口数 > 10 时，列代表性 5~10 处，注明"全量见 xxx"，禁止全量罗列刷屏。
- **风险点**：每条带证据，禁止"建议合理设计""注意代码质量"类空泛表述；无风险时省略本节并在全景图后注明。
- **状态标注**：已下线/灰度/仅测试使用的接口，在接口表格"作用"列或数据结构说明首行标注。
