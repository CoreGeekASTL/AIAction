# 业务流程文档索引

| 元信息 | 值 |
|--------|-----|
| 分支 | ready/27.0-终端鉴权 分支 (2026-08-07) |
| 更新日期 | 2026-08-07 |
| Skill | spec-business-flow-analyze |

| 流程 | 一句话说明 | 入口位置 | 文档 |
|------|-----------|---------|------|
| 终端鉴权判定 | 终端携带 IMEI+IMSI 请求时联合鉴权，login 链路拒绝返回 -2、event 链路拒绝返回 401 | src/routers/beego_router.go RegisterInternalRouter/RegisterExternalRouter | [terminal-auth.md](terminal-auth.md) |
