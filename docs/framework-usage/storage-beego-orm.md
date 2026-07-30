# Beego ORM 存储层使用指导（存储/ORM）

> 版本：beego/v2 client/orm v2.1.0 + openGauss-connector-go-pq v1.0.7 + modernc.org/sqlite v1.54.0 ｜ 调用点：20 文件 ｜ 涉及文件：dao/*、models/db/*、db/driver/* ｜ 基线：main (6c93561)

## 用途定位
唯一关系库访问层。双数据源：生产 GaussDB（openGauss pq 驱动，注册名 "gaussdb_1"，DRPostgres）；本地训战 LOCAL_MODE=true 时嵌入式 SQLite（`src/dao/db_init.go:228-242`、`src/dao/db_local_sqlite.go`）。

## 初始化与配置
- 入口 `dao.EnsureConnectGaussDB()`（`src/main.go:67` goroutine 启动）：`orm.BootStrap()` → `orm.RegisterDriver("gaussdb_1", orm.DRPostgres)` → LOCAL_MODE 走 SQLite 否则 CSE 发现 GaussDB 主节点循环重连（`src/dao/db_init.go:228-263`）。
- 连接信息：优先 DB Service HTTP 接口（getGaussdbInfor），失败回退 app.conf `[gaussdb]`（`src/dao/db_init.go:329-349`）。
- 主备切换：`checkDBStatus` 每 5s Ping + 比对 CSE 主节点 IP，异常则 `switchToAnotherDB` 并重设全局 `ormer`（`src/dao/db_init.go:183-206,265-280`）。
- 建表：DDL 字符串按 `;` 拆分逐条 Raw Exec（GaussDB：`src/dao/db_init.go:30-138`；SQLite：`src/dao/db_local_sqlite.go:23-123`），**DDL 以代码内 initSql 为准**，`src/db/` 下脚本为辅。
- 全局 `ormer`：默认 `&orm.DoNothingOrm{}`（DB 未就绪时安全空转，`src/dao/base_dao.go:18`），连接成功后 `orm.NewOrmUsingDB(alias)` 替换。

## 核心使用模式

实体定义（来源：`src/models/db/user.go:10-59`）：

```go
type User struct {
	Key  string `orm:"pk;column(key)"`
	Model string `json:"model" orm:"column(model)"`
	// ...
}
func (u *User) TableName() string { return "t_user" }
func init() { orm.RegisterModel(&User{}, &UserBind{}) }
```

DAO 定义（来源：`src/dao/user.go:7-17`）：

```go
type UserDao struct { BaseInterface }
func NewUserDaoDao() *UserDao {
	dao := &UserDao{}
	dao.BaseInterface = &BaseDao{EntityType: &db.User{}}
	return dao
}
```

CRUD 与事务（来源：`src/dao/base_dao.go`、`src/service/config_center_service.go:71-85`）：

```go
err := ubd.Get(old, "Key")                 // Read by cols
err := ubd.Insert(new)                     // 插入（ErrLastInsertIdUnavailable 已吞）
list := &[]db.ConfigCenter{}
dao.List(list, dao.NewQueryOption().Filter("Key", k).Limit(10, 0, "-Id"))
dao.QueryMulti(&rows, "SELECT ... WHERE x=?", arg)  // 原生 SQL
dao.DoTxWithCtx(ctx, func(ctx goctx.Context, txOrm orm.TxOrmer) error {
	// 事务内操作
})
```

## 封装层与扩展点
- 封装层：`dao.BaseDao` + `BaseInterface`（List/Get/Delete/Insert/Update/QueryOne/QueryMulti/Exec/InsertMulti/DoTxWithCtx，`src/dao/base_dao.go:56-73`）。业务代码**只用 BaseInterface 方法，禁止直接碰全局 ormer**。
- 扩展点：`QueryOption`（Filter/OrderBy/Limit 链式，`src/dao/base_dao.go:22-54`）；`DoTxWithCtx` 暴露 `orm.TxOrmer` 做事务。
- 测试替身：`dao.DoNothingBase`（`src/dao/donothing_base_dao.go`）用于单测替换 DAO。

## 并发与线程模型
- `ormer` 全局变量在主备切换时被替换（`src/dao/db_init.go:188,203`），切换加 `dbConnection.Lock`，读端无锁，属可接受竞态窗口。
- orm.Ormer 本身并发安全（连接池由 database/sql 管理）。

## 错误处理与容错
- `orm.ErrNoRows` 表示无记录，调用方显式判断（`src/service/browser_service.go:113-117`、`src/controllers/login_controller.go:49`）。
- `ErrLastInsertIdUnavailable` 在 BaseDao.Insert 内吞掉返回 nil（`src/dao/base_dao.go:131-133`）。
- 所有 error 必须检查处理（质量基线，`go vet ./...` 兜底）。

## 约定与规范
- 表名 `t_xxx`，列 snake_case；实体放 `src/models/db/`，DAO 放 `src/dao/`。
- 单主键：pk tag 只标一个字段（Beego 不支持复合 pk）；唯一约束用 DDL UNIQUE INDEX 兜底（坑记录 #6）。
- DDL 双写：GaussDB initSql 与 SQLite localSqliteInitSql 必须同步演进（SQLite 用 INTEGER AUTOINCREMENT/TEXT 方言）。
- context 传 `goctx.Background()` 或 `context.TODO()`，禁止 nil。

## 已知问题与反模式
- 建表 DDL 散在 Go 字符串里，无迁移工具；改表需同时 ALTER 语句（`src/dao/db_init.go:127-137` 已有先例）。
- `initTables` 用 `ormer.Raw` 在切换数据源前执行可能打到旧连接（启动期可接受）。
- `NewUserDaoDao` 命名笔误（`src/dao/user.go:11`），新增 DAO 用 `NewXxxDao`。

## AI 编码指南
- 新增表：`models/db/` 建实体（orm tag + TableName + init 注册）→ `dao/` 建 `XxxDao{BaseInterface}` → 两份 initSql（GaussDB+SQLite）同步加 DDL（依据：上文「核心使用模式」「约定与规范」）。
- 业务读写一律走 `dao.BaseInterface` 方法；**禁止**业务代码 import `orm` 直接操作（orm.ErrNoRows 判断除外）（依据：`src/dao/base_dao.go:56`）。
- 多步写必须用 `DoTxWithCtx` 事务；单键查询用 `Get(md, "Col")` 并区分 `orm.ErrNoRows`（依据：`src/service/config_center_service.go:71`、`src/controllers/login_controller.go:49`）。
