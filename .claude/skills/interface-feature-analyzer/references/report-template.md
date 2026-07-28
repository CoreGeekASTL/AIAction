# 输出模板

两套模板：README.md 索引 + 每功能一篇软件要素 md。所有文档归档到 `<repo>/docs/interface/`。

## 模板一：README.md 索引

```markdown
# 功能软件要素文档

> 由 interface-feature-analyzer 生成/更新，面向人与 AI 共同消费。
> 代码仓：<仓库名>　分析基准：<commit/分支>　更新时间：<YYYY-MM-DD>

## 功能全景

| 功能域 | 接口数 | 核心模块 | 文档 |
|---|---|---|---|
| 用户认证 | 5 | auth, session | [feature-user-auth.md](feature-user-auth.md) |
| 订单管理 | 8 | order, billing | [feature-order-manage.md](feature-order-manage.md) |

## 接口统计

- 外部接口：N 个（IDL 契约 X / 框架路由 Y / 消息订阅 Z）
- 内部接口：N 个（语言级契约，详见各功能文档"接口清单"节）
- 已下线：N 个（清单见各功能文档接口表"状态"列）

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

​```mermaid
flowchart LR
  classDef mi fill:#ffd1dc,stroke:#c2185b,color:#000
  classDef role fill:#fff3b0,stroke:#f9a825,color:#000
  classDef ppt fill:#c8e6c9,stroke:#2e7d32,color:#000
  classDef desc fill:#bbdefb,stroke:#1565c0,color:#000

  Buyer[买家<br/>角色]:::role
  Pay[三方支付渠道<br/>外部系统]:::role
  E1[提交订单]:::mi
  E2[支付]:::mi
  E3[发货]:::mi
  Order[(订单<br/>待支付→已支付→已发货)]:::ppt
  Stock[(库存)]:::ppt
  Rule[库存不足禁止下单]:::desc

  Buyer --> E1
  E1 -->|校验| Stock
  Rule -.约束.-> Stock
  E1 --> E2 --> E3
  Pay --> E2
  E2 -->|状态变更| Order
  E3 -->|状态变更| Order
​```

术语表：

| 术语 | 人话解释 | 出处 |
|---|---|---|
| 预占 | 下单时先锁定库存、支付超时自动释放的机制（order.go:118） | service/order.go:118 |
| 渠道单号 | 三方支付平台返回的交易流水号，用于对账 | dao/schema.sql:30 |

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
| controllers | 接口入口、参数校验（order.go:45） |
| service | 业务逻辑编排、库存校验（service/order.go:30） |
| dao | 订单表读写（dao/order.go:8） |

## 3. 接口清单

只记接口名与进出的数据结构，不写散文：

| 接口 | 路径/入口 | 请求结构 | 响应结构 | 状态 |
|---|---|---|---|---|
| CreateOrder | POST /api/v1/order（handler/order.go:45） | CreateOrderReq | OrderResp | 在用 |
| QueryOrder | GET /api/v1/order/{id}（handler/order.go:78） | path id | OrderResp | 在用 |

语言级内部接口（模块间契约）：

| 接口 | 定义位置 | 实现 | 选择机制 |
|---|---|---|---|
| OrderService | service/order.go:15 | OrderServiceImpl | 唯一实现，构造函数注入 |

## 4. 关键数据结构

表格呈现；字段多时只列关键字段（含义+约束），不搬运全部字段：

| 结构 | 定义位置 | 关键字段 |
|---|---|---|
| CreateOrderReq | api/types.go:15 | ItemID（商品ID，必填）；Count（>0） |
| Order | model/order.go:10 | Status（状态机：init→paid→done） |
| t_order 表 | dao/schema.sql:22 | uk(order_no)；软删 deleted_at |

## 5. 调用关系

用 mermaid 时序图呈现主链路（数据从哪进、到哪出，含跨服务跳转与异步）：

​```mermaid
sequenceDiagram
  participant Cli as Client
  participant C as Controller
  participant S as Service
  participant D as DAO
  participant M as MQ
  Cli->>C: POST /api/v1/order
  C->>S: CreateOrder(req)
  S->>S: CheckStock（外部调用，超时500ms）
  S->>D: Insert（事务）
  S-->>M: SendOrderEvent（异步，失败仅告警）
  C-->>Cli: 200 OrderResp
​```

关键分支与异步环节（各一句，带 文件:行号）：

- 库存不足走预占回退分支（service/order.go:118）
- 下单事件异步投递，消费方为 billing 服务

## 6. AI 编码指南

只列实现要点，每条 ≤30 字，附证据锚点：

- 新接口挂 /api/v1/order 前缀并注册 RouteMapping（handler/order.go:33）
- 写库走 dao 事务封装，禁拼裸 SQL（dao/order.go:60）
- 后续动作走 MQ 事件，禁主链路同步追加（service/order.go:95）
```

## 撰写硬性要求

- 所有位置证据使用 `文件:行号` 格式，相对代码仓根目录。
- **L1 功能故事**：mermaid 四色用 classDef 固定配色（粉 #ffd1dc / 黄 #fff3b0 / 绿 #c8e6c9 / 蓝 #bbdefb）；粉色事件按时序连成主干，每个事件必须具备 触发者(黄)/输入(绿/蓝)/输出(绿)/后继(粉) 四要素；事件名用人话业务动作，禁止出现文件名/函数名/行号；实体标注状态变更（X→Y）；规则用虚线挂到实体；代码推断不出的业务背景标注"代码中未体现"，禁止脑补。
- **术语表**：覆盖功能故事与接口清单中全部业务缩写/黑话/外部系统名；每条=一句话人话解释+出处行号，禁止只写英文全称。
- 模块划分必须先 mermaid 图后表格，图节点带入口文件名；依赖方向在图中用箭头表达，禁止用纯文字罗列。
- 接口清单只填表格五列，禁止为每个接口写段落描述；接口 > 20 个时按子功能拆多张表。
- 数据结构只列"理解该功能必须知道"的结构与关键字段；超过 8 个字段的结构只写关键字段。
- 调用关系必须先 mermaid 时序图，再用短句补关键分支/异步；多实现选择机制填在"接口清单"的语言级接口表中。
- AI 编码指南每条 ≤30 字（不含证据锚点），1-5 条，禁止"建议合理设计""注意代码质量"类空泛表述。
- "状态"列只取：在用 / 已下线 / 灰度中。
