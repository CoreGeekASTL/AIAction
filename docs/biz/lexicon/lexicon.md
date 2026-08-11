# GIDS 领域词典

| 元信息 | 值 |
| --- | --- |
| 分支 | new_skill_test 分支（2026-08-11） |
| 更新日期 | 2026-08-11 |
| Skill | biz-lexicon-analyze |

## 说明

- 词汇口径：本词典收录业务与代码共用的受控术语，按「子域 × 术语类型」组织并**按业务子域拆分多篇**——每功能域 1 篇 `lexicon-{子域锚点}.md`（锚点与 docs/biz/interface/ 子文档一致），子域内按四类术语类型分组（实体与业务概念 → 常量与状态枚举 → 错误码 → 事件，固定顺序）；本主文档承载说明、全仓待确认清单、子域导航与跨子域共用的「通用」节（同四类分组）。
- 子域划分复用 docs/biz/interface/ 功能域口径：login / auth / cache / event / file / management / plugin / stats / config-center / test。
- 术语来源与类型映射：对外接口文档/请求响应模型、DB 实体注释 → 实体与业务概念；常量定义 → 常量与状态枚举；错误码定义 → 错误码；事件模型 → 事件。五类来源均已排查：仓内无独立消息订阅/topic 机制，「事件」类术语取自 `src/models/events/base.go` 的埋点事件模型（EventType/EventData）。
- 释义依据：词条释义以代码注释/字段名/枚举值为据；标注"代码未体现，待确认"的释义已汇总进下方「待确认清单」，需业务方确认。
- 同义异名已合并为一条并互指；同名异义按语境逐条列出。

## 待确认清单

| 术语 | 所在子域 | 待确认内容 |
| --- | --- | --- |
| AppType | login / event | 字段为 string/int 混合（login 请求为 string，流量统计表为 int），TikTok 时取常量 "2" 走 Muen 云二次鉴权（见 interface-login），完整枚举取值代码未体现，待确认 |
| Muen 云二次鉴权 | login | interface-login.md 提及 TikTok AppType 走 Muen 云二次鉴权，代码中未见对应实现位置，待确认 |
| BrowserCap | login | `UserBind.BrowserCap` 字段未落库（orm:"-"），业务含义代码未体现，待确认 |
| Lac / CI / Rxlev | login | LoginAuthRequest 中小区与信号强度字段，取值口径与使用场景代码未体现，待确认 |
| VideoMode / PlayMode | login / event | 取值枚举代码未体现，待确认 |
| RecordMode / MachineType / FFCode | management | ChromeConfig 字段枚举取值代码未体现，待确认 |
| AccessType | stats | 流量统计接入类型枚举取值代码未体现，待确认 |
| NodeIdent | management | URLConfig 节点标识格式与取值口径代码未体现，待确认 |
| LoginError | event | `browser_user_http_login_error` 事件已定义常量但不在 eventTypeMap 中，使用场景待确认 |

## 子域导航

| 子域 | 词条数 | 子文档 |
| --- | --- | --- |
| 登录鉴权与用户绑定 | 16 | [lexicon-login.md](lexicon-login.md) |
| 终端鉴权与白名单管理 | 12 | [lexicon-auth.md](lexicon-auth.md) |
| 缓存管理 | 1 | [lexicon-cache.md](lexicon-cache.md) |
| 客户端事件上报 | 12 | [lexicon-event.md](lexicon-event.md) |
| 文件管理 | 6 | [lexicon-file.md](lexicon-file.md) |
| 浏览器配置同步 | 6 | [lexicon-management.md](lexicon-management.md) |
| 插件管理 | 11 | [lexicon-plugin.md](lexicon-plugin.md) |
| 流量统计 | 10 | [lexicon-stats.md](lexicon-stats.md) |
| 配置中心 | 5 | [lexicon-config-center.md](lexicon-config-center.md) |
| 测试联通 | 1 | [lexicon-test.md](lexicon-test.md) |
| 通用 | 8 | 见本文「通用」节 |

## 通用

### 实体与业务概念

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| BaseResponse | 统一响应基座，含 code/msg 两字段，所有对外响应内嵌 | - | `resp.BaseResponse`，`src/models/resp/base.go` |
| DataResponse | 带 data 载荷的统一响应结构 | - | `resp.DataResponse`，`src/models/resp/response_entity.go` |
| IRequest | 请求校验接口契约（Validate() error） | - | `req.IRequest`，`src/models/req/request_entity.go` |
| ITable | 可落库请求契约（IRequest + orm.TableNameI） | - | `req.ITable`，`src/models/req/request_entity.go` |
| MultiTableRequest | 上报多条数据的批量请求（items 数组） | - | `req.MultiTableRequest`，`src/models/req/request_entity.go` |
| 定时任务选主 | 记录当前持有定时任务的节点（ip/mac/id），选主 key 固定为 gids.timerElection | - | `db.ScheduleElection`，`src/models/db/schedule_election.go` |

### 常量与状态枚举

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| ComponentName | 组件名 "GIDS"，埋点事件 service 字段取值 | - | `constants.ComponentName`，`src/common/constants/base.go` |
| ServiceName | 服务名 "gids" | - | `constants.ServiceName`，`src/common/constants/base.go` |
| TikTokAppType | TikTok 应用类型取值 "2" | - | `constants.TikTokAppType`，`src/common/constants/base.go` |
| CleanupMonths | 数据清理保留月数（3 个月） | - | `constants.CleanupMonths`，`src/common/constants/base.go` |
| 环境变量键 | APPID / NODENAME / NAMESPACE / SERVICENAME / ENABLE_HTTP 五个环境变量键名 | - | `constants.EnvAppId` 等，`src/common/constants/base.go` |

### 错误码

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| Success / AuthPassed | 成功/鉴权通过，取值 200 | - | `retcode.Success`、`retcode.AuthPassed`，`src/common/constants/retcode/retcode.go` |
| InternalFailed | 服务端内部失败，取值 -1 | - | `retcode.InternalFailed`，`src/common/constants/retcode/retcode.go` |
| ClientFailed | 客户端请求失败（参数/鉴权拒绝），取值 -2；login 路径拒绝返回本码 | - | `retcode.ClientFailed`，`src/common/constants/retcode/retcode.go` |
| AuthFailed | 鉴权失败，取值 401；event 路径拒绝返回本码 | - | `retcode.AuthFailed`，`src/common/constants/retcode/retcode.go` |
