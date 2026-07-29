# Beego ORM + GaussDB/SQLite 使用指导（存储/ORM）

> 版本：beego/v2 v2.1.0（client/orm）+ openGauss-connector-go-pq v1.0.7 + modernc.org/sqlite v1.54.0 ｜ 调用点：全部 DAO 经封装层 ｜ 涉及文件：models/db(7 文件/9 实体) + dao(10) + db/driver ｜ 基线：origin/main @ ae0a8a6（2026-07-29 复核）

## 用途定位

唯一的关系型存储通道。双数据源：
- **生产**：GaussDB（Postgres 协议），经 CSE 服务发现定位主库，支持主备切换（`dao/db_init.go`）
- **本地训战**：`LOCAL_MODE=true` 时嵌入式 SQLite（纯 Go，无 CGO），落盘 `src/data/gids.db`（`dao/db_local_sqlite.go:127-183`）

业务代码永远面向 `dao.BaseInterface` 编程，不感知底层是 GaussDB 还是 SQLite。

## 初始化与配置

- 入口：`main.go:67` `go dao.EnsureConnectGaussDB()`（`db_init.go:228-263`）
  1. `orm.BootStrap()` + `orm.RegisterDriver("gaussdb_1", orm.DRPostgres)`
  2. LOCAL_MODE → `initLocalSQLite()`；否则循环：CSE 查主 GaussDB 实例（`Properties["status"]=="M"`，`db_init.go:382-409`）→ `switchToAnotherDB` 注册数据源 → `initTables()` 执行 DDL
- 驱动装饰：`db/driver/driver.go:82-87` 用 `sql.Register("gaussdb_1", Decorator)` 包了一层，作用是把 ORM 生成 SQL 中带引号的表名去引号（GaussDB 大小写适配，`driver.go:33-39`）。
- 健康检查与主备切换：`checkDBStatus` 每 5s ticker，Ping 失败 3 次或发现新主库则 `refresh()` 切换（`db_init.go:265-304`）；切换本质是改全局 `ormer = orm.NewOrmUsingDB(alias)`（`db_init.go:188,203`）。
- 连接串来源优先级：DB 微服务接口（`getGaussdbInfor`，`db_init.go:351-380`）→ app.conf `[gaussdb]` 段兜底（`db_init.go:342-349`）。

## 核心使用模式

### 三步曲：Model → DAO → Service

```go
// 1) Model：orm tag + TableName + init 注册（来源：src/models/db/user.go:10-59）
type User struct {
	Key   string `orm:"pk;column(key)"`
	Model string `json:"model" orm:"column(model)"`
}
func (u *User) TableName() string { return "t_user" }
func init() { orm.RegisterModel(&User{}) }
```

```go
// 2) DAO：内嵌 BaseInterface + 构造函数注入 EntityType（来源：src/dao/user.go:7-17）
type UserDao struct { BaseInterface }
func NewUserDaoDao() *UserDao {
	dao := &UserDao{}
	dao.BaseInterface = &BaseDao{EntityType: &db.User{}}
	return dao
}
```

```go
// 3) Service 使用（来源：src/service/user_service.go:86-117）
ub := &db.UserBind{Key: sessionId}
err := u.ubd.Get(ub)                 // 按主键读
err = u.ud.Insert(user)              // 插入
err = u.ubd.Update(ub)               // 更新
```

### 复杂查询：QueryOption / Raw SQL / 事务

```go
// QueryOption 链式（来源：src/dao/base_dao.go:29-54, 92-111）
opts := dao.NewQueryOption().Filter("Key", key).OrderBy("-CreatedAt").Limit(100, 0, "-CreatedAt")
err := d.List(&list, *opts)

// Raw SQL（base_dao.go:142-165）
err := d.QueryOne(&md, "SELECT ... WHERE key=?", arg)
n, err := d.Exec(ctx, "DELETE FROM t_x WHERE created_at<?", arg)

// 事务（来源：src/service/config_center_service.go:71-86）
err := c.dao.DoTxWithCtx(context.Background(), func(ctx context.Context, txOrm orm.TxOrmer) error {
	...
	return c.dao.Insert(&center)
})
```

## 封装层与扩展点

- **DAO 基类**：`dao.BaseDao`（`base_dao.go:78-174`）实现 `BaseInterface`（`base_dao.go:56-73`）：List/Get/Delete/Insert/Update/QueryOne/QueryMulti/Exec/InsertMulti/DoTxWithCtx，及配套 `*WithOrm` 变体（事务内复用 `orm.QueryExecutor`）。
- 全局 `ormer` 单例（`base_dao.go:18`），初始为 `orm.DoNothingOrm{}`——DB 未就绪时调用静默无操作，**注意这不是报错**。
- 批量插入默认 bulk=100（`base_dao.go:16,173`）。
- `InsertWithOrm` 吞掉 `orm.ErrLastInsertIdUnavailable`（`base_dao.go:131-133`）——GaussDB 无 LastInsertId，属正常。
- 测试替换：`dao/donothing_base_dao.go` + `//go:build !test` tag（`db_init.go:1`）实现 UT 隔离。
- **新表必须三步曲 + 双 DDL，禁止在业务代码中直接 `orm.NewOrm()`。**

## 并发与线程模型

`ormer` 是包级全局变量，主备切换时写、DAO 操作时读，无锁保护（`db_init.go:188` vs `base_dao.go:93`）——存量竞态，改造需评估。DAO 实例本身无状态（仅 `EntityType`），可每请求新建。

## 错误处理与容错

- `orm.ErrNoRows` 是"查无记录"的标准哨兵，调用方必须显式区分（`user_service.go:100-104`、`config_center_service.go:76`）。
- DB 连接失败：`EnsureConnectGaussDB` 每 5s 无限重试（`db_init.go:243-261`），服务可在 DB 就绪前先起端口。
- 双 DDL 维护：GaussDB 版 `db_init.go:30-138` 与 SQLite 版 `db_local_sqlite.go:23-123` 手工保持一致（SERIAL→AUTOINCREMENT、bytea→BLOB、boolean→INTEGER 等方言差异）。

## 约定与规范

- 表名 `t_` 前缀；主键统一名为 `key`（自增表用 `id`）；时间字段 `created_at/updated_at` 用字符串存 `time.Now().Format(time.DateTime)`（`user_service.go:106-107`）。
- Beego ORM 不支持复合主键：单字段 pk，其余唯一约束用 DDL `CREATE UNIQUE INDEX` 兜底（AGENTS.md 已踩坑 #6，证据 `db_init.go:126`）。
- DDL 变更必须幂等（`CREATE TABLE IF NOT EXISTS` / `ADD COLUMN IF NOT EXISTS`），启动时反复执行（`db_init.go:127-137`）。

## 已知问题与反模式

- 双 DDL 漂移风险（见上）；改表只改一处是历史高发错误。
- `db_init.go` 中大量 `logger.Errorf("s00893267 ...")` 调试日志打印连接串/密码（`db_init.go:194,335,376`）——敏感信息泄露，禁模仿。
- `dbConnection` 方法与包级函数混用 `dbConnections` 全局变量（`db_init.go:154,209-213`）。
- master election（`models/db/schedule_election.go`）目前仅有 stub 测试（`service/master_election_service_stub_test.go`），生产实现不在仓内——不要假设选主能力可用。

## AI 编码指南

- 新增表：① `models/db/` 加实体（`orm:"pk;column(...)"` + `TableName()` + `init(){orm.RegisterModel}`）② `dao/` 加 `XxxDao{BaseInterface}` + 构造函数 ③ **两处 DDL 同步加**（`db_init.go` initSql 与 `db_local_sqlite.go` localSqliteInitSql）。依据：上文「核心使用模式」「约定与规范」。
- 增删改查一律通过 `NewXxxDao()` 的 `BaseInterface`；查无记录判 `orm.ErrNoRows`；多写操作用 `DoTxWithCtx(context.Background(), ...)`。**禁止**业务层直接 `orm.NewOrm()`/`orm.Raw()`。依据：上文「封装层与扩展点」。
- **禁止**给实体加复合 pk tag；**禁止**在日志/SQL 拼接中打印密码与完整连接串。依据：上文「已知问题与反模式」及 AGENTS.md 已踩坑 #6。
