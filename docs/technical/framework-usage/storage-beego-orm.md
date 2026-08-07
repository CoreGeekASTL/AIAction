# Beego ORM 存储层使用指导（存储/ORM）

## 用途定位
唯一关系库访问层。双数据源：生产 GaussDB（openGauss pq 驱动，注册名 "gaussdb_1"，DRPostgres）；本地训战 LOCAL_MODE=true 时嵌入式 SQLite（`src/dao/db_init.go:228-242`、`src/dao/db_local_sqlite.go`）。


## 使用模式

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
