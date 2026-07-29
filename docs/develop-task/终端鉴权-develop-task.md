# 终端鉴权 develop-task

## 1. 任务概述

落地 27.0 终端鉴权：白名单 CSV 导入/导出 3 个新接口 + IMEI/IMSI 联合鉴权核心服务（含缓存、逃生态、并发互斥），并注入登录与事件上报 4 条既有链路。设计细节见 [story 设计文档](../story/feature-terminal-auth.md)，本文档只管怎么改代码。内容依据《27.0终端鉴权需求设计文档-规范版.md》（含 8 处人工裁定）。

## 2. 修改文件清单

| 文件 | 操作 | 修改点与实现逻辑 |
|------|------|-----------------|
| src/models/db/white_list.go（规划） | 新增 | WhiteList 实体：Imei `orm:"pk;column(imei);size(15)"`、Imsi `orm:"column(imsi);size(15)"`（不加 pk，Beego 不支持复合 pk，DDL UNIQUE INDEX 兜底）、CreatedAt；TableName() 返回 t_white_list；init() 中 orm.RegisterModel |
| src/dao/white_list.go（规划） | 新增 | WhiteListDao 继承 BaseInterface，EntityType 设 WhiteList；实现 Count（逃生态判定）、GetByIMEI(imei, imsi) 联合精确查询、InsertMulti 事务批量插入、ClearAndInsert 事务清表+插入、ListAll 全量查询；ContextDo 传 context.TODO() |
| src/dao/db_init.go | 修改 | GaussDB DDL 加 t_white_list 建表语句（imei char(15)、imsi char(15)、created_at，imei+imsi 联合索引、唯一约束） |
| src/dao/db_local_sqlite.go | 修改 | SQLite DDL 同步加 t_white_list（与 GaussDB 双 DDL 一致，仅方言差异） |
| src/service/auth_cache.go（规划） | 新增 | AuthCache 包级单例（sync.Once）：sync.RWMutex + map[string]cacheEntry{result, expireAt}；Get/Set/惰性清理（超 1000 清最旧 500）；TTL 30min；不起独立 goroutine |
| src/service/auth_service.go（规划） | 新增 | AuthService 接口 + authServiceImpl 单例：CheckAuth(imei, imsi) —— ① `^[0-9]{15}$` 格式校验非法即拒 ② Count==0 逃生态放行 ③ 联合键查缓存 ④ 未命中 GetByIMEI 双匹配 ⑤ 写缓存（含 false 与放行标记） |
| src/service/whitelist_service.go（规划） | 新增 | WhiteListService：ImportCSV(file, operation) —— 互斥锁防并发（抢不到返回导入进行中）→ 文件 ≤3MB、逐行 15 位纯数字、≤100W 条、无表头逗号分隔、重复行整批拒 → firstImport 查 Count 非空报错 / update 调 ClearAndInsert；ExportCSV() 调 ListAll 生成无表头 CSV（空表返回空文本+200） |
| src/controllers/auth_controller.go（规划） | 新增 | AuthController 继承 BaseController：RouteInfo() 声明 3 条路由（POST /auth/v1/authIMEI、POST /auth/v1/importIMEIList、GET /auth/v1/exportIMEIList）；AuthIMEI 解析 JSON{imei,imsi} 调 AuthService 返回 {"code":200/-2,"msg"}；ImportIMEIList 接收 multipart CSV + operation 参数；ExportIMEIList 返回 CSV 文本；不做业务逻辑 |
| src/routers/beego_router.go | 修改 | RegisterInternalRouter 中将 AuthController 的 RouteMapping 注册进内部监听（端口 9090） |
| src/controllers/login_controller.go | 修改 | gridLoginAuth/gridLoginAuthOpenBrowser 内部入口：调用原业务 Service 之前注入 AuthService.CheckAuth，拒绝返回 code=-2 |
| src/controllers/exlogin_controller.go | 修改 | 同上两个登录接口外部入口的鉴权注入，拒绝返回 code=-2 |
| src/controllers/event_controller.go | 修改 | sendClientEvent/sendAppUseTimesEvent 注入 AuthService.CheckAuth，拒绝返回 code=401（与登录 -2 必须区分） |
| src/models/req/、src/models/resp/ | 修改 | 登录/事件请求结构确认 IMSI 字段可解析（协议新增字段，规范版断点2）；新增 AuthIMEIRequest{Imei, Imsi} |

## 3. 要用的框架

| 框架 | 文档链接 | 按哪条约定执行 |
|------|---------|---------------|
| Beego Web 路由/Controller | [rpc-beego-web.md](../framework-usage/rpc-beego-web.md) | RouteInfo().RouteMapping 集中声明，注册到内部监听；继承 BaseController |
| Beego ORM | [storage-beego-orm.md](../framework-usage/storage-beego-orm.md) | DAO 继承 BaseInterface + EntityType；Model 三步曲（orm tag+TableName+init 注册）；双 DDL 同步 |
| Go 协程与单例 | [concurrency-goroutine.md](../framework-usage/concurrency-goroutine.md) | Service/Cache 用 sync.Once 包级单例；缓存 sync.RWMutex；导入互斥锁 |
| 序列化 | [codec-json-yaml.md](../framework-usage/codec-json-yaml.md) | authIMEI 用 encoding/json；导入导出用 encoding/csv（无表头纯数据） |
| 日志 | [log-lager-auditlog-event.md](../framework-usage/log-lager-auditlog-event.md) | 鉴权拒绝/导入结果/逃生态触发记业务日志，不打印敏感信息 |
| 测试框架 | [test-testify-goconvey.md](../framework-usage/test-testify-goconvey.md) | DAO/Service 接口注入 mock 单测；LOCAL_MODE SQLite 集成验证 |

## 4. 要调用的外部接口

无。本功能全部为仓内闭环（白名单落本地 DB，无出站调用）。

## 5. 验证方式

- DT 测试：src/service/auth_service_test.go、auth_cache_test.go、whitelist_service_test.go、src/dao/white_list_test.go（规划），覆盖规范版 §3.6 的 TC-AUTH-01~16
- 集成测试：LOCAL_MODE=true 启动后手工验证 3 个接口 + 4 条注入链路（curl POST/GET 9090 端口）
- 验证命令：`cd src && go build -o gids.exe . && go vet ./... && go test -v ./service/... ./dao/... ./controllers/...`

## 6. 编码工作流提示

- TDD：先写测试再写实现，红灯→绿灯→重构
- 宣称完成前必须运行验证命令，用证据支撑"已完成"
- 遇 bug / 测试失败：先定位根因再改，禁止试凑式修改
- 登录链路拒绝 code=-2、事件链路 code=401，两处不可混用（历史踩坑）
- IMEI/IMSI 校验正则固定 `^[0-9]{15}$`，不要写成 16 位
