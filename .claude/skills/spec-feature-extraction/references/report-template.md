# 输出模板

<<<<<<< HEAD:.claude/skills/spec-feature-extraction/references/report-template.md
两套模板：README.md 索引 + 每功能一篇软件要素 md。所有文档归档到 `<repo>/docs/story/`。
=======
只产出一份 markdown：`<repo>/docs/interface/external-interfaces.md`。结构固定三段：接口全景 → 各功能域详述 → 风险与注意点。每个功能域下先接口表格（接口名/作用/所在文件/方法或路径），表格下方逐个说明该接口相关的请求与响应数据结构。强调人类阅读友好。
>>>>>>> 512daa4 (修改接口分析 skill):.claude/skills/interface-feature-analyzer/references/report-template.md

## 模板：docs/interface/external-interfaces.md

```markdown
# 对外接口文档

> 代码仓：<仓库名>　分析基准：<commit/分支>　更新时间：<YYYY-MM-DD>
> 由 interface-feature-analyzer 生成，面向人类阅读。
> 范围：本仓对外提供的接口（HTTP 路由 / RPC service / 消息订阅 handler / IDL 契约），不含本仓调用别人的接口。

## 1. 接口全景

<<<<<<< HEAD:.claude/skills/spec-feature-extraction/references/report-template.md
| 功能域 | 接口数 | 核心模块 | 文档 |
|---|---|---|---|
| 用户认证 | 5 | auth, session | [feature-user-auth.md](feature-user-auth.md) |
| 订单管理 | 8 | order, billing | [feature-order-manage.md](feature-order-manage.md) |

## 接口统计

- 对外接口：N 个（IDL 契约 X / 框架路由 Y / 消息订阅 Z），按外部监听/内部监听分列
- 已下线：N 个（清单见各功能文档接口表"状态"列）
- 说明：语言级内部接口（仓内模块间契约）仅用于分析，不写入功能文档

## 未归类接口

以下接口探测到但未纳入任何功能域，原因逐条说明：

- `POST /internal/debug/dump`（debug 工具，非业务功能）

## 使用说明

- **AI 编码时**：按需求所属功能查阅对应 md，先看 mermaid 图建立结构认知，再按"接口清单 → 调用关系"定位改动点。
- **新人上手**：从功能全景表选入口，按"模块划分 → 调用关系"建立功能级认知。
```

## 模板二：每功能一篇 feature-<功能名>.md

```markdown
# <功能名>

> 功能域概述：一两句话说明该功能解决什么业务问题。
> 接口数：N（外部 X / 内部 Y）

## 1. 功能故事（多彩建模）

> L1 层，给新人读：用四色建模讲清"谁触发、经过什么业务环节、改变了什么、按什么规则"。方法论见 spec-logic-audit/references/color-modeling.md。
=======
一张 mermaid 图：本仓为中心节点，指向各功能域；功能域节点按业务功能命名。
>>>>>>> 512daa4 (修改接口分析 skill):.claude/skills/interface-feature-analyzer/references/report-template.md

实现逻辑速览（1~3 句，每句 ≤30 字）：

下单先锁库存，再写订单，支付超时自动释放。

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

<<<<<<< HEAD:.claude/skills/spec-feature-extraction/references/report-template.md
| 术语 | 人话解释 | 出处 |
|---|---|---|
| 预占 | 下单时先锁定库存、支付超时自动释放的机制 | service/order.go |
| 渠道单号 | 三方支付平台返回的交易流水号，用于对账 | dao/schema.sql |

## 2. 模块划分

用 mermaid 依赖图呈现模块分层与依赖方向，节点带入口文件锚点：

​```mermaid
graph LR
  R[routers<br/>beego_router.go] --> C[controllers<br/>order.go]
  C --> S[service<br/>order.go]
  S --> D[dao<br/>order.go]
  D --> DB[(GaussDB)]
  S --> EXT[外部服务 X]
​```

| 模块 | 承载功能 |
|---|---|
| controllers | 接口入口、参数校验（order.go） |
| service | 业务逻辑编排、库存校验（service/order.go） |
| dao | 订单表读写（dao/order.go） |

## 3. 接口清单

只列对外接口（对本仓外部暴露的 HTTP 路由/消息订阅），只记接口名与进出的数据结构，不写散文；语言级内部接口不写入本节：

| 接口 | 路径/入口 | 请求结构 | 响应结构 | 状态 |
|---|---|---|---|---|
| CreateOrder | POST /api/v1/order（handler/order.go） | CreateOrderReq | OrderResp | 在用 |
| QueryOrder | GET /api/v1/order/{id}（handler/order.go） | path id | OrderResp | 在用 |
=======
功能域清单：

| 功能域 | 接口数 | 核心模块 | 详见 |
|---|---|---|---|
| 终端鉴权 | 5 | controllers/auth, service/auth | [§2.1](#21-终端鉴权) |
| 白名单管理 | 6 | controllers/whitelist, dao/whitelist | [§2.2](#22-白名单管理) |
| 配置同步 | 4 | controllers/management, service/remote | [§2.3](#23-配置同步) |
| 会话管理 | 3 | controllers/session, service/browser | [§2.4](#24-会话管理) |
| 定时巡检 | 2 | scheduler, dao | [§2.5](#25-定时巡检) |
>>>>>>> 512daa4 (修改接口分析 skill):.claude/skills/interface-feature-analyzer/references/report-template.md

> 未归类接口：无（若有探测到但无法归入任何功能域的接口，在此逐条列出并说明原因）

## 2. 各功能域详述

<<<<<<< HEAD:.claude/skills/spec-feature-extraction/references/report-template.md
| 结构 | 定义位置 | 关键字段 |
|---|---|---|
| CreateOrderReq | api/types.go | ItemID（商品ID，必填）；Count（>0） |
| Order | model/order.go | Status（状态机：init→paid→done） |
| t_order 表 | dao/schema.sql | uk(order_no)；软删 deleted_at |
=======
每个功能域一节，固定三要素：定位 → 接口表格 → 表格下方逐个接口的数据结构说明。
>>>>>>> 512daa4 (修改接口分析 skill):.claude/skills/interface-feature-analyzer/references/report-template.md

### 2.1 终端鉴权

**定位**：终端登录与事件上报的鉴权校验，区分 login 路径（拒绝返回 -2）与 event 路径（拒绝返回 401）。

**接口清单**：

<<<<<<< HEAD:.claude/skills/spec-feature-extraction/references/report-template.md
关键分支与异步环节（各一句，带证据文件）：

- 库存不足走预占回退分支（service/order.go）
- 下单事件异步投递，消费方为 billing 服务

## 6. 框架引用

本功能使用的基础框架及对应框架使用文档（docs/framework-usage/，由 framework-usage-analyzer 产出；仓内无该目录时省略本节）：

| 基础框架 | 框架文档 | 本功能中的用途（引用文件） |
|---|---|---|
| Beego Web 路由/Controller | [rpc-beego-web.md](../framework-usage/rpc-beego-web.md) | 路由注册与请求处理（routers/beego_router.go、controllers/order.go） |
| beego ORM | [storage-beego-orm.md](../framework-usage/storage-beego-orm.md) | 订单表读写（dao/order.go） |
| lager 日志 | [log-lager-auditlog-event.md](../framework-usage/log-lager-auditlog-event.md) | 业务日志（service/order.go） |

## 7. AI 编码指南

只列实现要点，每条 ≤30 字，附证据文件锚点：

- 新接口挂 /api/v1/order 前缀并注册 RouteMapping（handler/order.go）
- 写库走 dao 事务封装，禁拼裸 SQL（dao/order.go）
- 后续动作走 MQ 事件，禁主链路同步追加（service/order.go）
=======
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
>>>>>>> 512daa4 (修改接口分析 skill):.claude/skills/interface-feature-analyzer/references/report-template.md
```

## 撰写硬性要求

<<<<<<< HEAD:.claude/skills/spec-feature-extraction/references/report-template.md
- 所有位置证据使用**文件级**路径（相对代码仓根目录），禁止写行号；术语表"出处"列同样只写文件。
- **L1 功能故事**：开头 1~3 句实现逻辑速览，每句 ≤30 字，用业务语言概括代码实现逻辑，禁止文件名/函数名/行号；mermaid 四色用 classDef 固定配色（粉 #ffd1dc / 黄 #fff3b0 / 绿 #c8e6c9 / 蓝 #bbdefb）；粉色事件按时序连成主干，每个事件必须具备 触发者(黄)/输入(绿/蓝)/输出(绿)/后继(粉) 四要素；事件名用人话业务动作，禁止出现文件名/函数名/行号；实体标注状态变更（X→Y）；规则用虚线挂到实体；代码推断不出的业务背景标注"代码中未体现"，禁止脑补。
- **术语表**：覆盖功能故事与接口清单中全部业务缩写/黑话/外部系统名；每条=一句话人话解释+出处文件，禁止只写英文全称。
- 模块划分必须先 mermaid 图后表格，图节点带入口文件名；依赖方向在图中用箭头表达，禁止用纯文字罗列。
- 接口清单只列对外接口（对本仓外部暴露的 HTTP 路由/消息订阅），语言级内部接口一律不写入文档；只填表格五列，禁止为每个接口写段落描述；接口 > 20 个时按子功能拆多张表。
- 数据结构只列"理解该功能必须知道"的结构与关键字段；超过 8 个字段的结构只写关键字段。
- 调用关系必须先 mermaid 时序图，再用短句补关键分支/异步；同接口多实现的选择机制在模块划分表或关键分支中说明。
- **框架引用**：逐行必须有代码事实依据（import/调用点所在文件），禁止按框架文档目录全量罗列；框架文档列必须链接到仓内存在的文档，相对路径，无死链；仓内无框架使用文档目录时整节省略并在索引 README 注明。
- AI 编码指南每条 ≤30 字（不含证据锚点），1-5 条，禁止"建议合理设计""注意代码质量"类空泛表述。
- "状态"列只取：在用 / 已下线 / 灰度中。
=======
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
>>>>>>> 512daa4 (修改接口分析 skill):.claude/skills/interface-feature-analyzer/references/report-template.md
