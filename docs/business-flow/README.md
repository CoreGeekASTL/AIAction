# 业务流程文档索引

| 流程 | 一句话说明 | 入口位置 | 文档 |
|------|-----------|---------|------|
| 终端鉴权判定 | 终端携带 IMEI+IMSI 请求时联合鉴权，login 链路拒绝返回 -2、event 链路拒绝返回 401 | src/routers/beego_router.go RegisterInternalRouter/RegisterExternalRouter | [terminal-auth.md](terminal-auth.md) |
