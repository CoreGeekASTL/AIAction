# GaussDB 数据访问规范

| 元信息 | 值 |
|--------|-----|
| 分支 | 27.0 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-data-access-guidelines-analyze |
| 运行模式 | 起草模式 |

## 用途定位

GaussDB 实例承载 GIDS 全部关系数据（t_config、t_file、t_plugin_package、t_user、t_user_bind、t_media_traffic_stats、t_control_traffic_stats、t_session_stats、t_config_center、t_white_list 共 10 张表），由 src/dao/ 模块经 Beego ORM 统一访问；文件内容也直接落库（t_file.content bytea），不依赖外部对象存储。LOCAL_MODE 本地训战时不走本中间件，改走嵌入式 SQLite（见 data-access-guidelines-sqlite.md）。

## 访问点分布

| 访问点 | 所在文件 | 所在函数 | 访问方式 | 业务场景 |
| --- | --- | --- | --- | --- |
| 数据源注册与切换 | src/dao/db_init.go | switchToAnotherDB | 裸 orm（orm.RegisterDataBase） | 连接 GaussDB 主节点、主备切换重连 |
| 建表与列变更 DDL | src/dao/db_init.go | initTables | 裸 orm（ormer.Raw().Exec()） | 启动时执行 initSql 建表/建索引 |
| 数据源地址发现 | src/dao/db_init.go | getGaussDBIP / getDataSourceFromDBService | HTTP 调用 DB Service + CSE 服务发现 | 获取主库 IP、连接串 |
| 通用 CRUD 基类 | src/dao/base_dao.go | BaseDao.List/Get/Delete/Insert/Update/QueryOne/QueryMulti/Exec/InsertMulti | 封装层：src/dao/base_dao.go | 全部表的基础读写 |
| 用户表访问 | src/dao/user.go | NewUserDaoDao / NewUserBindDao | 封装层（BaseDao，EntityType=db.User/db.UserBind） | 终端用户信息、绑定关系读写 |
| 配置表访问 | src/dao/browser_config.go | NewConfigDao | 封装层（BaseDao，EntityType=db.BrowserConfig） | 浏览器配置读写 |
| 配置中心表访问 | src/dao/config_center.go | NewConfigCenterDao | 封装层（BaseDao，EntityType=db.ConfigCenter） | 配置中心读写（src/service/config_center_service.go 事务写入） |
| 插件包表访问 | src/dao/plugin.go | NewPluginPackageDao | 封装层（BaseDao，EntityType=db.PluginInfo） | 插件包元数据读写（src/service/plugin_service.go 事务写入） |
| 文件表访问 | src/dao/file.go | FileDao.Exist / NewFileDao | 封装层（BaseDao，EntityType=db.File） | 文件元数据与内容读写（src/service/file_service.go） |
| 流量统计表访问 | src/dao/traffic_stats_dao.go | NewMediaTrafficStatsDao / NewControlTrafficStatsDao / NewSessionLogDao / SessionStatsDao.Exist / GetIdByTcpUniqueId / UpdatebySession | 封装层（BaseDao） | 媒体/控制流量、会话统计读写 |
| 白名单表访问 | src/dao/white_list.go | NewWhiteListDao / Count / GetByIMEIAndIMSI / InsertMulti / ClearAndInsert（DoTxWithCtx 清表+批量插入） / ListAll | 封装层（BaseDao，EntityType=db.WhiteList） | 终端白名单读写（src/service/auth_service.go 鉴权只读、src/service/whitelist_manage_service.go 导入导出） |
| 话统 SQL 执行 | src/service/traffic_stats_service.go | executeQuery | BaseDao.QueryMulti（Raw SQL，SQL 文本来自 src/conf/sql.yaml） | CSP 话统定时上报（在线数、流量） |
| 流量批量入库 | src/service/traffic_stats_service.go | InsertMultiWithOrm（DoTxWithCtx 内） | 封装层（BaseDao.InsertMultiWithOrm） | 流量统计数据批量写入 |

## 连接与客户端管理

### 现状

- 全局单例 `ormer`（src/dao/base_dao.go），初始为 `orm.DoNothingOrm{}`，连接成功后由 `switchToAnotherDB` 重建为 `orm.NewOrmUsingDB(alias)`；初始化入口为 `EnsureConnectGaussDB`（src/dao/db_init.go）。
- 驱动为 `gaussdb_1`：src/db/driver/driver.go 用 `Decorator` 包装 openGauss-connector-go-pq（`gitee.com/opengauss/openGauss-connector-go-pq`），在 `Prepare` 阶段去掉表名双引号；`orm.RegisterDriver("gaussdb_1", orm.DRPostgres)`（src/dao/db_init.go）。
- 连接串来源两路：优先 HTTP 调 DB Service `/service/api/getGaussdbInfor`（src/dao/db_init.go getDataSourceFromDBService），失败且 servicename=GaussDB 时回退读 app.conf `gaussdb::gaussdbuser/gaussdbport/gaussdbdbname/databasepassword`（getDataSourceByConfig）；主库 IP 经 CSE 服务发现（getGaussDBIP，仅取 Status=UP 且 Properties["status"]="M" 的实例）。
- 连接池参数（maxOpenConns 等）代码未显式设置（Beego/database-sql 框架默认）。
- 主备切换：后台 goroutine `checkDBStatus` 每 5s 健康检查（Ping，连续失败 3 次触发 refresh），`refresh` 循环重连新主（src/dao/db_init.go）。
- 初始连接失败时 for 循环每 5s 重试直至成功（src/dao/db_init.go EnsureConnectGaussDB）。

### 约定（起草待评审）

1. 业务代码一律经 src/dao/ 封装层访问数据库，禁止在 service/controller 层直接持有 `orm.Ormer` 或调用 `orm.NewOrmUsingDB`。
2. 连接串、账号密码只允许从 DB Service 接口或 app.conf `gaussdb::` 配置组获取，禁止在业务代码中硬编码 DSN。
3. 连接池上限应显式配置（当前框架默认，建议补充 maxOpenConns/maxIdleConns 并纳入配置管理）。
4. `ormer` 的重建只发生在 `switchToAnotherDB` 一处，新增代码不得另行替换全局 ormer。
5. 日志中不得打印含密码的连接串（现状 db_init.go 存在 `connStr` 明文入日志的打印点，建议脱敏）。

## 事务使用

### 现状

- 事务统一走 `BaseDao.DoTxWithCtx`（src/dao/base_dao.go），底层为 Beego `ormer.DoTxWithCtx`；事务内写操作通过 `InsertWithOrm/UpdateWithOrm/InsertMultiWithOrm` 传入 `txOrm` 完成。
- 事务使用点：配置中心写（src/service/config_center_service.go InsertOrUpdateConfig）、插件包增删改（src/service/plugin_service.go 多处 DoTxWithCtx）、流量统计批量写入（src/service/traffic_stats_service.go 批量 InsertMultiWithOrm 及 sd/md/cd 三表的 DoTxWithCtx）、白名单覆盖导入（src/dao/white_list.go ClearAndInsert：DELETE 全表 + InsertMultiWithOrm 同事务）。
- 其他写路径（user、file、browser_config 等）走 `Insert/Update` 非事务直写。
- 跨表一致性未覆盖的场景（如 t_file 与 t_plugin_package 联合写入是否同事务）未在代码中统一约束。

### 约定（起草待评审）

1. 涉及多行/多表的写路径必须走 `DoTxWithCtx`，禁止循环内逐条直写。
2. 事务边界收在 dao 层 `BaseDao.DoTxWithCtx`，service 层只组装 task 闭包，不得自行 `BeginTx/Rollback/Commit`。
3. 事务内一律使用 `txOrm`（`*WithOrm` 系列方法），禁止在 task 闭包内再调非事务的 `Insert/Update` 直写方法。

## 分页与批量

### 现状

- 分页：`QueryOption.Limit(limit, offset, orderBy)`（src/dao/base_dao.go），limit/offset 方式，负值忽略；目前仅流量统计清理路径使用（src/service/traffic_stats_service.go `opt.Limit(batchSize, offset, "id")`）。
- 批量写：`InsertMulti`/`InsertMultiWithOrm`（src/dao/base_dao.go），bulk 大小硬编码 `defaultBulk = 100`；使用点在流量统计批量入库（src/service/traffic_stats_service.go）与白名单导入（src/service/whitelist_manage_service.go 按 importBatchSize=1000 分片调用 InsertMulti）。
- 批量读：`List` 无条件时全量拉取，无单页上限强制。

### 约定（起草待评审）

1. 列表查询必须显式 `Limit` 分页并给出排序字段，禁止无上限全量拉取大表。
2. 批量写一律走 `InsertMulti*`，禁止循环单条 Insert；bulk 大小沿用 `defaultBulk`，需要调整时应改常量而非散落硬编码。

## SQL 拼接与注入防护

### 现状

- ORM 条件查询走 `QueryOption.Filter(key, vals...)` → Beego `qs.Filter` 参数化绑定（src/dao/base_dao.go），无字符串拼值。
- 原生 SQL 出口：`QueryOne/QueryMulti/Exec`（`ormer.Raw`，src/dao/base_dao.go）；话统统计 SQL 文本来自配置文件 src/conf/sql.yaml（`?` 占位符参数化，参数仅取 startTime/endTime 两个固定变量，src/service/traffic_stats_service.go executeQuery）。
- DDL 拼接：`initSql` 常量按 `;` 拆分逐条 `Raw().Exec()`（src/dao/db_init.go initTables），为静态常量、无外部输入。
- 驱动层对表名做去引号字符串替换（src/db/driver/driver.go extractTableName），仅处理 SQL 文本中的表名标识符，不涉及值。
- 未发现 `fmt.Sprintf` 拼接查询值的用法。

### 约定（起草待评审）

1. 查询条件一律参数化绑定（`?` 占位符 / `qs.Filter`），禁止 `fmt.Sprintf` 或 `+` 拼接值进 SQL。
2. 新增统计类 SQL 统一收敛到 src/conf/sql.yaml 声明式配置，禁止在代码内散落裸 SQL 字符串。
3. 表名/列名等无法参数化的标识符必须来自代码常量或配置，禁止接收外部输入。

## 缓存读写模式

### 现状

GaussDB 为持久层；终端鉴权链路在 DB 前加进程内鉴权结果缓存 authCache（src/service/auth_cache.go，RWMutex+map，TTL 30min/容量 1000，cache-aside 模式：miss 回源 t_white_list 后写缓存），缓存不承载导入写路径（导入直写 DB，旧缓存条目靠 TTL 过期失效）；Redis 封装无业务调用方（见 data-access-guidelines-redis.md 附注）。

### 约定（起草待评审）

1. 鉴权缓存仅用于 AuthIMEI 只读判定路径，白名单写路径不得读写该缓存。
2. 白名单导入覆盖后旧缓存条目在 TTL 内可能仍旧放行/拒绝，业务须接受最长 30min 的生效延迟（代码未体现主动失效机制，待确认是否需导入后清缓存）。

## 错误处理

### 现状

- 初始化连接失败：for 循环 5s 间隔无限重试（src/dao/db_init.go EnsureConnectGaussDB）。
- 运行期连接失败：`checkDBStatus` 健康检查 + `refresh` 自动切换主库（src/dao/db_init.go）。
- `InsertWithOrm` 对 `orm.ErrLastInsertIdUnavailable` 特判返回 nil（src/dao/base_dao.go），其余错误原样上抛。
- 话统查询失败仅记录日志并返回空结果（src/service/traffic_stats_service.go executeQuery），无重试。
- 唯一键冲突、not found 等错误码未见统一分类处理，由各调用方自行判断；写失败无重试机制。

### 约定（起草待评审）

1. dao 层错误原样上抛，由 service 层决定转换业务错误码；禁止在 dao 层吞错只打日志。
2. not found（Beego `orm.ErrNoRows`）与系统错误须分开返回，调用方不得仅凭 err!=nil 一律按系统错误处理。
3. 写路径失败不重试（避免重复写）；只读统计查询失败允许降级返回空结果，但必须记录日志（与现状一致）。
4. 唯一键冲突应在 service 层显式识别并转为幂等或业务冲突错误码。

<!-- 规范未建，本次为现状盘点与约定起草。 -->
