# 输出模板

两套模板：README.md 索引 + 每功能一篇软件要素 md。所有文档归档到 `<repo>/doc/AR/`。

## 模板一：README.md 索引

```markdown
# 功能架构资产（AR）

> 由 interface-feature-analyzer 生成/更新，面向人与 AI 共同消费。
> 代码仓：<仓库名>　分析基准：<commit/分支>　更新时间：<YYYY-MM-DD>

## 功能全景

| 功能域 | 接口数 | 核心模块 | 文档 |
|---|---|---|---|
| 用户认证 | 5 | auth, session | [feature-user-auth.md](feature-user-auth.md) |
| 订单管理 | 8 | order, billing | [feature-order-manage.md](feature-order-manage.md) |

## 接口统计

- 外部接口：N 个（IDL 契约 X / 框架路由 Y / 消息订阅 Z）
- 内部接口：N 个（语言级契约，详见各功能文档"模块划分"节）
- 已下线：N 个（清单见各功能文档接口表"状态"列）

## 未归类接口

以下接口探测到但未纳入任何功能域，原因逐条说明：

- `POST /internal/debug/dump`（debug 工具，非业务功能）

## 使用说明

- **AI 编码时**：按需求所属功能查阅对应 md，先读「AI 编码指南」，再按"接口清单 → 调用关系"定位改动点。
- **新人上手**：从功能全景表选入口，按"模块划分 → 调用关系"建立功能级认知。
```

## 模板二：每功能一篇 feature-<功能名>.md

```markdown
# <功能名>

> 功能域概述：一两句话说明该功能解决什么业务问题。
> 接口数：N（外部 X / 内部 Y）　核心模块：a, b, c

## 1. 模块划分

| 模块/包 | 职责 | 依赖方向 |
|---|---|---|
| api/handler/order | 接口入口、参数校验（order.go:12） | → service |
| service/order | 业务逻辑编排（service.go:30） | → dao, 外部服务 X |
| dao/order | 存储访问（dao.go:8） | → DB |

（依赖方向用箭头表达，标注分层约束，如"禁止 dao 反向依赖 service"）

## 2. 接口清单

| 接口 | 协议/路径 | 入口位置 | 状态 |
|---|---|---|---|
| CreateOrder | POST /api/v1/order | handler/order.go:45 | 在用 |
| QueryOrder | GET /api/v1/order/{id} | handler/order.go:78 | 在用 |
| 旧版下单 | POST /api/v0/order | handler/legacy.go:20 | 已下线（路由已注释 main.go:33） |

（语言级内部接口单列一节，含 interface 定义位置与实现列表）

## 3. 关键数据结构

| 结构 | 定义位置 | 用途 | 关键字段 |
|---|---|---|---|
| CreateOrderReq | api/types.go:15 | 下单请求 | ItemID（商品ID，必填）、Count（>0） |
| Order | model/order.go:10 | 领域模型 | Status（状态机：init→paid→done） |
| t_order 表 | dao/schema.sql:22 | 存储 | uk(order_no)，软删字段 deleted_at |

（只列"理解该功能必须知道"的结构，字段含义与约束是重点，不搬运全部字段）

## 4. 调用关系

主链路（以 CreateOrder 为例）：

```
handler/order.go:45 CreateOrder
  → service/order.go:102 CheckStock（外部调用库存服务，超时 500ms）
  → dao/order.go:60 Insert（事务，失败回滚）
  → mq/producer.go:33 SendOrderEvent（异步，失败仅告警）
```

关键分支与异步环节：

- 库存不足走 service/order.go:118 的预占回退分支
- 下单事件异步投递，消费方为 billing 服务（见 feature-billing.md）

（有条件时用 mermaid 时序图表达跨服务调用；多实现接口在此说明路由选择机制）

## 5. AI 编码指南

1. 新增订单类接口必须挂在 /api/v1/order 前缀下并在 handler/order.go 注册（依据：全部现有接口集中注册于此，handler/order.go:33-90）。
2. 写库操作必须走 dao 层事务封装 dao/order.go:60，禁止在 service 层直接拼 SQL（依据：service 层无裸 SQL 先例）。
3. 下单后的后续动作必须通过 mq 事件扩展，禁止在 CreateOrder 主链路同步追加（依据：主链路性能约束注释 service/order.go:95）。
```

## 撰写硬性要求

- 所有位置证据使用 `文件:行号` 格式，相对代码仓根目录。
- 表格行数过多时（接口 > 20 个）按子功能分小节，不要一张表撑到底。
- "状态"列只取：在用 / 已下线 / 灰度中。
- AI 编码指南每条 = 规则 + 依据（代码证据），1-3 条，禁止"建议合理设计""注意代码质量"类空泛表述。
