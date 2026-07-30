# 终端鉴权（IMEI+IMSI 白名单）develop-task

## 1. 任务概述

27.0 商用补齐：新增白名单管理（CSV 导入/导出）与 IMEI+IMSI 联合鉴权（缓存+逃生态），注入登录与事件上报链路。设计细节见 [story 设计文档](../story/feature-terminal-auth.md)，本文档只管怎么改代码。

## 2. 修改文件清单

| 文件 | 操作 | 修改点与实现逻辑 |
|------|------|-----------------|
| src/models/db/white_list.go（规划） | 新增 | AuthWhitelist 实体：IMEI `orm:"pk;column(imei)"` char(15)、IMSI `orm:"column(imsi)"` 不加 pk、CreatedAt；TableName() 返回 t_white_list；init() 中 orm.RegisterModel |
| src/dao/white_list.go（规划） | 新增 | WhiteListDao 继承 BaseInterface、EntityType 设 AuthWhitelist；实现 Count()、GetByIMEI(imei) 返回整行、InsertMulti(records)、ClearAndInsert(records)（事务：清表+批量插入，失败回滚）、ListAll() |
| src/dao/db_init.go | 修改 | GaussDB DDL 加建表 t_white_list + IMEI+IMSI 联合 UNIQUE INDEX（复合约束由 DDL 兜底，Beego 不支持复合 pk） |
| src/dao/db_local_sqlite.go | 修改 | SQLite DDL 同步同一表结构与联合索引，保持双 DDL 一致 |
| src/service/auth_cache.go（规划） | 新增 | 缓存组件：sync.RWMutex + map[string]cacheEntry{result, expireAt}；Get（RLock，过期视为未命中）、Set（Lock）、Clear()；Set 后容量>1000 按 expireAt 升序惰性删最旧 500，无独立 goroutine |
| src/service/auth_service.go（规划） | 新增 | AuthService 接口 + authServiceImpl + 包级变量 + sync.Once；Check(imei, imsi) bool：格式校验（15位纯数字，非法即 false）→ 缓存 Get → DAO Count==0 逃生态放行并缓存 → GetByIMEI 命中且 IMSI 一致放行，否则缓存 false 拒绝；ClearCache() 委托缓存组件 |
| src/service/auth_manage_service.go（规划） | 新增 | AuthManageService：Import(reader, operation)——校验 ≤3MB、逐行 15 位纯数字、≤20W 条、文件内重复组合整批拒；firstImport 先 Count 要求为 0，update 走 ClearAndInsert；事务提交成功后调 AuthService.ClearCache()；Export() 全量 ListAll 生成无 header CSV |
| src/controllers/auth_controller.go（规划） | 新增 | AuthController 继承 BaseController；RouteInfo() 声明三接口：POST /auth/v1/authIMEI（JSON{imei,imsi}→{code:200/-1}）、POST /auth/v1/importIMEIList（multipart 上传 + query operation；200/-1/-2）、GET /auth/v1/exportIMEIList（text/csv） |
| src/routers/beego_router.go | 修改 | 注册 AuthController 到内部监听路由列表 |
| src/controllers/login_controller.go | 修改 | GridLoginAuth、GridLoginAuthOpenBrowser 函数开头注入 AuthService.Check(imei, imsi)，false 时返回 retcode.ClientFailed(-2)，不进后续建档/分配 |
| src/controllers/exlogin_controller.go | 修改 | 与 login_controller.go 镜像同步注入（内外两 controller 必须同时改） |
| src/controllers/event_controller.go | 修改 | SendClientEvent、SendAppUseTimesEvent 函数开头注入 Check，false 时返回 retcode.AuthFailed(401) |
| src/service/auth_service_test.go（规划） | 新增 | DT：格式非法短路、逃生态放行、命中/未命中、缓存命中零 DB、TTL 过期回源、容量清理、ClearCache |
| src/service/auth_manage_service_test.go（规划） | 新增 | DT：CSV 校验各失败分支、firstImport 表非空拒绝、update 覆盖、导入后缓存被清空 |
| src/dao/white_list_test.go（规划） | 新增 | DT：五方法 CRUD 与事务回滚（LOCAL_MODE SQLite） |
| src/controllers/auth_controller_test.go（规划） | 新增 | DT：三接口路由映射与响应码 |

## 3. 要用的框架

| 框架 | 文档链接 | 按哪条约定执行 |
|------|---------|---------------|
| Beego Web | [rpc-beego-web.md](../framework-usage/rpc-beego-web.md) | 新 Controller 继承 BaseController，只动 RouteInfo()，注册由路由层统一遍历 |
| Beego ORM | [storage-beego-orm.md](../framework-usage/storage-beego-orm.md) | 实体三步曲（orm 标签+TableName+init 注册），DAO 继承 BaseInterface，ContextDo 传 context.TODO() |
| 自研单例 | [di-singleton.md](../framework-usage/di-singleton.md) | 接口 + 小写实现 + 包级变量 + sync.Once |
| 协程原语 | [concurrency-goroutine-sync.md](../framework-usage/concurrency-goroutine-sync.md) | 缓存 sync.RWMutex，清理在写锁内完成，不起 goroutine |
| Go 测试 | [test-go-testing.md](../framework-usage/test-go-testing.md) | testing + testify 表驱动，接口注入 fake |

## 4. 要调用的外部接口

无。本功能无出向调用；/auth/v1/authIMEI 为被 BrowserGW 反调的入向接口。

## 5. 验证方式

- DT 测试：src/service/auth_service_test.go、src/service/auth_manage_service_test.go、src/dao/white_list_test.go、src/controllers/auth_controller_test.go（均规划）
- 集成测试：python testsuit/TC_SBG_Func_GIDS_Auth_001.py（导入导出）、002（鉴权正常/异常/逃生态）、003（登录链路）、004（事件链路）、005（边界与缓存）
- 验证命令：`cd src && go build -o gids.exe . && go vet ./... && go test -v ./service/... ./dao/... ./controllers/...`，随后 LOCAL_MODE=true 启动跑 testsuit 五个 TC

## 6. 编码工作流提示

- TDD：先写测试再写实现，红灯→绿灯→重构
- 宣称完成前必须运行验证命令，用证据支撑"已完成"
- 遇 bug / 测试失败：先定位根因再改，禁止试凑式修改
- 已踩坑提醒：IMEI/IMSI 必须 15 位纯数字正则 ^[0-9]{15}$；login 拒绝 -2、event 拒绝 401 不可混用；CSV 无 header；Beego 不支持复合 pk，IMSI 唯一性靠 DDL
