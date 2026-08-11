# Beego ORM 使用现状（存储/ORM）

## 用途定位
Beego ORM（`github.com/beego/beego/v2/client/orm`）承担全部关系库访问。生产连 GaussDB（`orm.DRPostgres` 驱动，openGauss-connector-go-pq），`LOCAL_MODE=true` 时切换为嵌入式 SQLite（`modernc.org/sqlite`，纯 Go 无 CGO）。全局只有一个包级 `ormer` 变量（`src/dao/base_dao.go`），所有 DAO 通过 `BaseDao` 封装访问，业务代码不直接使用 ORM 原生 API。连接管理（主 GaussDB 发现、健康检查、主备切换）由 `src/dao/db_init.go` 自研逻辑承担：通过 CSE 发现主库 endpoint、5s ticker 健康检查、失败自动 `switchToAnotherDB`。

实体注册集中在 `src/models/db/*.go`（orm tag + `TableName()` + `init(){ orm.RegisterModel() }`）；建表不走 `orm.RunSyncdb`，而是手写 DDL 常量（`src/dao/db_init.go` 的 `initSql`、`src/dao/db_local_sqlite.go` 的 `localSqliteInitSql`）用 `ormer.Raw(stmt).Exec()` 逐条执行。

## 使用模式

实体定义骨架（注意：Beego 不支持复合主键，唯一约束靠 DDL UNIQUE INDEX 兜底）：

```go
// 来源：src/models/db/user.go
type User struct {
	Key          string `orm:"pk;column(key)"`
	Manufacturer string `json:"manufacturer" orm:"column(manufacturer)"`
	// ...
}

func (u *User) TableName() string { return "t_user" }

func init() {
	orm.RegisterModel(&User{})
	orm.RegisterModel(&UserBind{})
}
```

DAO 骨架（继承 BaseDao，设置 EntityType；事务用 DoTxWithCtx）：

```go
// 来源：src/dao/base_dao.go
type BaseDao struct {
	EntityType interface{} // 数据表对应的entity类型
}

func (base *BaseDao) List(md interface{}, opts ...QueryOption) error {
	qs := ormer.QueryTable(base.EntityType)
	// Filter/OrderBy/Limit 链式拼接
	_, err := qs.All(md)
	return err
}

func (base *BaseDao) DoTxWithCtx(ctx goctx.Context, task func(ctx goctx.Context, txOrm orm.TxOrmer) error) error {
	return ormer.DoTxWithCtx(ctx, task)
}
```

数据源注册与切换（生产 GaussDB / 本地 SQLite 双轨）：

```go
// 来源：src/dao/db_init.go
orm.RegisterDriver("gaussdb_1", orm.DRPostgres)
orm.RegisterDataBase(aliasName, "gaussdb_1", dataSourceUrl)
ormer = orm.NewOrmUsingDB(aliasName)

// 来源：src/dao/db_local_sqlite.go
orm.RegisterDataBase("default", "sqlite3", dbPath)
ormer = orm.NewOrmUsingDB("default")
```

原生 SQL 走 `ormer.Raw(query, args...).QueryRow/QueryRows/Exec()`，事务内操作用 `InsertWithOrm/UpdateWithOrm/ExecWithOrm` 传入 `orm.QueryExecutor`（`src/dao/base_dao.go`）。
