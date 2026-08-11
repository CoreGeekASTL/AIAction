# 终端鉴权 develop-task

## 1. 任务概述

为 GIDS 补齐终端鉴权：白名单 CSV 导入/导出 + IMEI+IMSI 联合鉴权（含逃生态与内存缓存），并在登录/事件链路注入鉴权。设计细节见 [story 设计文档](../storys/终端鉴权-story.md)，本文档只管怎么改代码。

## 2. 修改文件清单

### Story-1 白名单数据基础设施

| 文件 | 操作 | 修改点与实现逻辑 |
|------|------|-----------------|
| src/models/db/white_list.go（规划） | 新增 | WhiteList 实体：`Imei string \`orm:"pk;column(imei);size(15)"\``、`Imsi string \`orm:"column(imsi);size(15)"\``、`CreatedAt string \`orm:"column(created_at)"\``；`TableName()` 返回 "t_white_list"；`init()` 中 `orm.RegisterModel(&WhiteList{})`（对照 src/models/db/user.go；依据需求 §7、story §5） |
| src/dao/white_list.go（规划） | 新增 | `WhiteListDao struct { BaseInterface }` + `NewWhiteListDao()`（EntityType=&db.WhiteList{}，对照 src/dao/user.go）；方法：`Count() (int64, error)`（QueryTable().Count()）、`GetByIMEIAndIMSI(imei, imsi string) error`（QueryTable Filter imei+imsi 联合精确匹配，命中与否以 orm.ErrNoRows 区分）、`InsertMulti(list *[]db.WhiteList) error`、`ClearAndInsert(list *[]db.WhiteList) error`（DoTxWithCtx 内 Raw `DELETE FROM t_white_list` + InsertMultiWithOrm，同事务）、`ListAll(list *[]db.WhiteList) error`（依据需求 §2.4 DAO 职责、story §11） |
| src/dao/db_init.go | 修改 | initSql 常量追加：`CREATE TABLE IF NOT EXISTS t_white_list (imei char(15) NOT NULL PRIMARY KEY, imsi char(15) NOT NULL, created_at varchar(255) DEFAULT '');` + `CREATE UNIQUE INDEX IF NOT EXISTS idx_white_list_imsi ON t_white_list(imsi);`（GaussDB DDL；IMEI pk、IMSI 靠 UNIQUE INDEX 兜底——Beego 不支持复合 pk；依据需求 §7、story §5） |
| src/dao/db_local_sqlite.go | 修改 | localSqliteInitSql 常量追加等价 SQLite 语句：`CREATE TABLE IF NOT EXISTS t_white_list (imei TEXT NOT NULL PRIMARY KEY, imsi TEXT NOT NULL, created_at TEXT DEFAULT '');` + `CREATE UNIQUE INDEX IF NOT EXISTS idx_white_list_imsi ON t_white_list(imsi);`（story §5 双 DDL） |

### Story-2 白名单管理接口（CSV 导入/导出）

| 文件 | 操作 | 修改点与实现逻辑 |
|------|------|-----------------|
| src/service/whitelist_manage_service.go（规划） | 新增 | `WhiteListManageService` 接口 + `whiteListManageServiceImpl` + 包级变量 + sync.Once 单例（对照 src/service/event_service.go）。方法：`ImportIMEIList(reader io.Reader, operation string) (int, error)`——encoding/csv 解析（FieldsPerRecord=2，空行自动跳过，兼容 CRLF/LF）；逐行正则 `^[0-9]{15}$` 强校验 IMEI/IMSI，条数 >200000 报错；任一行非法整批拒绝（先全量校验后入库）；校验通过后在导入写锁（authImportLock.Lock，定义于 auth_cache.go）内：firstImport 先 Count>0 则返回 "already exists" 错误，否则 InsertMulti 按 1000 条/批分片写入；update 调 ClearAndInsert 事务覆盖。`ExportIMEIList() (string, error)`——ListAll 后按无 header、IMEI/IMSI 两列生成 CSV 文本（依据需求 §2.2.1/§2.2.2、story §4/§8） |
| src/controllers/auth_controller.go（规划） | 新增 | `AuthController struct { BaseController; authService service.AuthService; manageService service.WhiteListManageService }`；`RouteInfo()` 三路由：`/auth/v1/authIMEI` POST:AuthIMEI、`/auth/v1/importIMEIList` POST:ImportIMEIList、`/auth/v1/exportIMEIList` GET:ExportIMEIList；`Prepare()` 初始化两 Service。`ImportIMEIList()`：`c.GetFile("file")` 取 multipart 文件（不校文件名），header.Size>3MB 或 operation 非 firstImport/update → `c.OK(resp.BaseResponse{Code: retcode.ClientFailed, Message: "invalid parameter"})`；调 manageService.ImportIMEIList，校验失败 → code=InternalFailed(-1) + err msg；成功 → code=200 + 导入条数（依据需求 §2.2.1，响应统一 HTTP 200+body code）。`ExportIMEIList()`：Content-Type 置 text/csv，经 ResponseWriter 直接写 CSV 文本（story §8 链路 2） |

### Story-3 联合鉴权核心服务（缓存+逃生态）

| 文件 | 操作 | 修改点与实现逻辑 |
|------|------|-----------------|
| src/models/req/auth_request.go（规划） | 新增 | `AuthIMEIRequest struct { IMEI string \`json:"imei"\`; IMSI string \`json:"imsi"\` }`，实现 req.IRequest：`Validate()` 仅做非空检查（15 位格式校验归 Service，裁定 Q1 的 format invalid 语义由 Service 判定）；JSON 解码大小写不敏感，兼容 testsuit 大写 "IMEI"/"IMSI" 键（story §3） |
| src/service/auth_cache.go（规划） | 新增 | 缓存组件：`cacheEntry struct { result bool; expireAt time.Time }`；`authCache struct { sync.RWMutex; entries map[string]cacheEntry }`；常量 ttl=30min、capacity=1000、evictCount=500。方法：`get(key) (bool, bool)`（RLock，过期视为 miss）、`set(key, result)`（Lock，写入后 len>capacity 按 expireAt 升序惰性删最旧 500，清理在同一 Lock 内，无独立 goroutine）。包级 `authImportLock sync.RWMutex`：导入管理 Service 持 Lock、鉴权回源段持 RLock（依据需求 §2.2.4、story §10） |
| src/service/auth_service.go（规划） | 新增 | `AuthService` 接口 + 单例。方法：`AuthIMEI(imei, imsi string) (allowed bool, formatValid bool)`——① 正则 `^[0-9]{15}$` 校验，非法返回 (false, false) 短路；② authCache.get 命中返回缓存结果；③ authImportLock.RLock 内回源：WhiteListDao.Count()==0 → 逃生态放行并 set(true)；④ 表非空 GetByIMEIAndIMSI 联合匹配，命中 set(true)/未命中 set(false)（防穿透）；⑤ DB 异常（非 ErrNoRows）→ 安全优先返回 (false, true)（依据需求 §2.2.3/§2.2.4、story §8 链路 3） |
| src/controllers/auth_controller.go（规划，同 Story-2 文件） | 新增 | `AuthIMEI()`：`RequestBodyUnmarshalTo` 解析失败（含 Validate 非空失败）→ `c.OK(resp.BaseResponse{Code: retcode.AuthFailed, Message: "format invalid"})`；调 authService.AuthIMEI：!formatValid → code=401 msg 含 "format invalid"（裁定 Q1）；!allowed → code=401 msg 含 "auth rejected"（裁定 Q2）；allowed → code=200 msg "success"。统一 HTTP 200 由 c.OK 保证（story §3） |
| src/routers/beego_router.go | 修改 | `RegisterInternalRouter` 增加 `registerController(server, &controllers.AuthController{})`（仅内网注册，不进 RegisterExternalRouter；依据需求 §2.1、story §3） |

### Story-4 登录链路鉴权注入

| 文件 | 操作 | 修改点与实现逻辑 |
|------|------|-----------------|
| src/controllers/login_controller.go | 修改 | LoginController struct 增加 `authService service.AuthService` 字段；`Prepare()` 增加初始化；`loginAuth(preOpenBrowser bool)`（login_controller.go:146）在 `RequestBodyUnmarshalTo` 成功后、`CreateOrUpdateUser` 之前插入：`allowed, _ := c.authService.AuthIMEI(request.IMEI, request.IMSI); if !allowed { c.Failed(resp.BaseResponse{Code: retcode.ClientFailed, Message: "auth rejected"}); return nil, nil }`（裁定 Q2/Q3；覆盖 GridLoginAuth/GridLoginAuthOpenBrowser/DeviceLoginAuth 三入口） |
| src/controllers/exlogin_controller.go | 修改 | ExLoginController struct 增加 authService 字段 + Prepare 初始化；`loginAuth`（exlogin_controller.go:94）同位置同逻辑注入（裁定 Q2/Q3） |

### Story-5 事件链路鉴权注入

| 文件 | 操作 | 修改点与实现逻辑 |
|------|------|-----------------|
| src/controllers/event_controller.go | 修改 | EventController struct 增加 `authService service.AuthService` 字段 + Prepare 初始化；`SendClientEvent()`（event_controller.go:32）与 `SendAppUseTimesEvent()`（event_controller.go:65）在 `RequestBodyUnmarshalTo` 成功后、构造 event 之前插入：`allowed, _ := c.authService.AuthIMEI(request.IMEI, request.IMSI); if !allowed { c.Failed(resp.BaseResponse{Code: retcode.AuthFailed, Message: "auth rejected"}); return }`（裁定 Q2） |

**方案选型约束落实**：鉴权逻辑唯一实现于 AuthService（DRY），Controller 只做码转换（login -2 / event 401 / authIMEI 401+format invalid）；导入写锁与鉴权读锁同一把 authImportLock（service 层包级变量，导入全程 Lock、回源段 RLock），避免清表窗口期逃生态误放行（功能设计 G2）；未命中 false 也缓存防穿透；DB 异常安全优先拒绝不放行（G1）。

## 3. 要用的框架

| 框架 | 文档链接 | 按哪条约定执行 |
|------|---------|---------------|
| Beego v2 HTTP Server | [usage-beego.md](../tech/usage/usage-beego.md) | Controller 实现 RouteInfo() 路由表，RegisterInternalRouter 统一 registerController 注册；响应用 BaseController OK/Failed |
| Beego ORM | [usage-beego-orm.md](../tech/usage/usage-beego-orm.md) | 实体 orm 标签+TableName()+init 注册；DAO 继承 BaseInterface/BaseDao；事务用 DoTxWithCtx；手写双 DDL 追加 initSql/localSqliteInitSql（不走 RunSyncdb）；Beego 不支持复合 pk，IMEI 设 pk、IMSI 用 UNIQUE INDEX 兜底 |
| Go 单测 | [usage-go-testing.md](../tech/usage/usage-go-testing.md) | DT 测试按存量同目录 _test.go 约定搭建（service/dao/controllers 各层） |
| encoding/csv + regexp（标准库） | 无框架文档，待核实 | CSV 解析与 15 位纯数字强校验仅用 Go 标准库，不引第三方（story §11 最小化依赖） |

## 4. 要调用的外部接口

无（本功能无出站调用，story §2 总览对外部服务调用判定为不涉及）。

## 5. 验证方式

- DT 测试：新增 src/service/auth_service_test.go、src/service/whitelist_manage_service_test.go、src/service/auth_cache_test.go、src/dao/white_list_test.go、src/controllers/auth_controller_test.go（按存量 _test.go 约定；覆盖：15 位校验边界、逃生态、缓存命中/过期/清理、firstImport 幂等、update 覆盖、导入码 -1/-2/200、authIMEI 401 format invalid）
- 集成测试（testsuit 全部 5 个 TC，LOCAL_MODE=true 启动 127.0.0.1:9090 后执行，前置全新 gids.db）：
  - `python testsuit/TC_SBG_Func_GIDS_Auth_001.py`（白名单导入导出全链路）
  - `python testsuit/TC_SBG_Func_GIDS_Auth_002.py`（联合鉴权正常/异常/逃生态）
  - `python testsuit/TC_SBG_Func_GIDS_Auth_003.py`（登录链路鉴权）
  - `python testsuit/TC_SBG_Func_GIDS_Auth_004.py`（事件链路鉴权）
  - `python testsuit/TC_SBG_Func_GIDS_Auth_005.py`（边界与缓存覆盖，含 deviceLoginAuth）
- 验证命令（src/ 目录下）：`go build -o gids.exe .`、`go vet ./...`、`go test -v ./service/... ./dao/... ./controllers/...`

## 6. 编码工作流提示

- TDD：先写测试再写实现，红灯→绿灯→重构
- 宣称完成前必须运行验证命令，用证据支撑"已完成"
- 遇 bug / 测试失败：先定位根因再改，禁止试凑式修改
- 集成测试不被动等待：DT 通过后主动运行 testsuit 全部 5 个 TC

## 7. 澄清问题列表

| 编号 | 疑问描述 | 涉及源码位置 | 用户澄清结论 |
|------|---------|-------------|-------------|
| Q1 | authIMEI 参数非法返回码冲突：testsuitcase.md 期望 code=-2+"invalid IMEI/IMSI"，TC_002.py 脚本期望 code=401+"format invalid" | controllers/auth_controller.go（规划）AuthIMEI | 以可执行脚本 TC_002.py 为准：code=401，msg 含 "format invalid" |
| Q2 | 注入点拒绝 msg 文案冲突：TC_003/004.py 期望 "auth rejected"，testsuitcase.md 期望 login "not allowed"、event "Unauthorized" | controllers/login_controller.go、exlogin_controller.go loginAuth()；controllers/event_controller.go SendClientEvent/SendAppUseTimesEvent | 以脚本为准：注入点拒绝 msg 统一含 "auth rejected"（login code=-2 / event code=401） |
| Q3 | deviceLoginAuth 是否注入鉴权：需求 §2.2.3 只列 4 个注入点，loginAuth() 为三入口共用函数，TC_005_003/004 要求覆盖 | controllers/login_controller.go:146、exlogin_controller.go:94 loginAuth() | 在三入口共用函数 loginAuth() 注入，deviceLoginAuth 自然一并覆盖 |
