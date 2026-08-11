# SQLite 数据访问规范

| 元信息 | 值 |
|--------|-----|
| 分支 | 27.0 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-data-access-guidelines-analyze |
| 运行模式 | 起草模式 |

## 用途定位

嵌入式 SQLite（纯 Go 驱动 `modernc.org/sqlite`）承载 LOCAL_MODE 本地训战模式下 GIDS 的全部关系数据，表结构与 GaussDB 链路一一对应（src/dao/db_local_sqlite.go localSqliteInitSql，含 27.0 新增的 t_white_list 及 idx_white_list_imsi 唯一索引），数据落本地文件（默认 `./data/gids.db`）；仅当环境变量 `LOCAL_MODE=true` 时启用，生产不受影响。访问仍经 src/dao/ 同一套 Beego ORM 封装层。

## 访问点分布

| 访问点 | 所在文件 | 所在函数 | 访问方式 | 业务场景 |
| --- | --- | --- | --- | --- |
| 驱动探测与注册 | src/dao/db_local_sqlite.go | initLocalSQLite | 裸 database/sql（sql.Open/sql.Register） | LOCAL_MODE 启动时注册 sqlite3 驱动别名 |
| 数据库注册 | src/dao/db_local_sqlite.go | initLocalSQLite | 裸 orm（orm.RegisterDataBase("default", "sqlite3", dbPath)） | 打开本地库文件并设为默认数据源 |
| 建表 DDL | src/dao/db_local_sqlite.go | initLocalSQLite | 裸 orm（ormer.Raw().Exec()） | 执行 localSqliteInitSql 建表/建索引 |
| 业务读写 | src/dao/*.go | 全部 dao 方法 | 封装层：src/dao/base_dao.go | 与 GaussDB 链路共用，见 data-access-guidelines-gaussdb.md |

## 连接与客户端管理

### 现状

- 触发点：`EnsureConnectGaussDB` 中判断 `os.Getenv("LOCAL_MODE") == "true"` 后调用 `initLocalSQLite`（src/dao/db_init.go）。
- 库文件路径取配置 key `local::sqlitepath`，默认 `./data/gids.db`，目录不存在时 `os.MkdirAll` 自动创建（src/dao/db_local_sqlite.go）。
- 驱动注册用 `sync.Once`（`sqliteDriverOnce`）保护：先探测 `sqlite3` 驱动名是否已可用，不可用则用 `sql.Open("sqlite", ":memory:")` 取到 modernc 驱动再以 `sqlite3` 名字注册（src/dao/db_local_sqlite.go）。
- 连接池参数代码未显式设置（database/sql 框架默认）。
- 无显式关闭点；进程退出即关闭。
- 数据源别名固定为 `default`，重复注册时识别 "already registered" 并复用（src/dao/db_local_sqlite.go）。

### 约定（起草待评审）

1. SQLite 链路只允许由 LOCAL_MODE 分支初始化，生产代码路径不得引用 `initLocalSQLite`。
2. 库路径只允许取 `local::sqlitepath` 配置，禁止散落硬编码路径。
3. 驱动注册必须保持 `sync.Once` 保护，禁止重复 `sql.Register` 同名驱动。
4. 嵌入式单文件库不考虑连接池调优；新增代码不得为 SQLite 路径引入跨进程并发写假设。

## 事务使用

### 现状

与 GaussDB 链路共用同一套 `BaseDao.DoTxWithCtx` 事务入口（src/dao/base_dao.go），SQLite 模式下事务由 Beego ORM 的 sqlite3 驱动承载，无本中间件特有的事务代码。

### 约定（起草待评审）

1. 事务规则与 GaussDB 链路保持一致（见 data-access-guidelines-gaussdb.md 事务使用约定）。
2. SQLite 为单写者模型，本地模式下不得起多 goroutine 并发写同一事务。

## 分页与批量

### 现状

与 GaussDB 链路完全复用 `QueryOption.Limit` 分页与 `InsertMulti` 批量写（src/dao/base_dao.go），无本中间件特有实现。

### 约定（起草待评审）

同 GaussDB 链路约定：列表查询显式分页、批量写走 InsertMulti。

## SQL 拼接与注入防护

### 现状

- 建表 DDL 为 `localSqliteInitSql` 静态常量，按 `;` 拆分逐条 `Raw().Exec()`（src/dao/db_local_sqlite.go），无外部输入。
- 业务 SQL 与 GaussDB 链路共用参数化查询；注意 SQLite 与 GaussDB 语法差异由两套独立 DDL 常量分别维护（src/dao/db_init.go initSql 与 src/dao/db_local_sqlite.go localSqliteInitSql），src/conf/sql.yaml 中话统 SQL 含 GaussDB 特有语法（`::TIMESTAMP`、`INTERVAL '24 HOURS'`），本地模式执行会失败。

### 约定（起草待评审）

1. 业务 SQL 保持参数化绑定，禁止拼值（同 GaussDB 链路）。
2. 两套 DDL 常量必须保持表结构语义同步演进，新增表/列须同时改 initSql 与 localSqliteInitSql。
3. 含 GaussDB 方言的统计 SQL（src/conf/sql.yaml）应标注仅 GaussDB 可用，或提供 SQLite 兼容变体；本地训战覆盖话统链路前须先解决方言差异。

## 缓存读写模式

### 现状

本中间件不涉及。

### 约定（起草待评审）

本中间件不涉及。

## 错误处理

### 现状

- 数据目录创建失败、驱动打开失败直接返回 error 上抛（src/dao/db_local_sqlite.go）。
- 建表逐条执行，记录首个错误 `firstErr`，全部执行完后统一返回；单条失败不中断后续语句（src/dao/db_local_sqlite.go）。
- 数据源重复注册识别 "already registered" 字符串匹配并降级为复用（src/dao/db_local_sqlite.go）。
- 初始化失败由调用方 `EnsureConnectGaussDB` 仅记日志（src/dao/db_init.go），进程继续运行（此时 ormer 仍为 DoNothingOrm，后续读写静默无操作）。

### 约定（起草待评审）

1. 本地库初始化失败应视为启动失败并显式上抛/退出，禁止静默降级为 DoNothingOrm 继续运行导致读写无感知丢失。
2. 错误判定避免依赖错误文本匹配（"already registered"），如驱动版本提供错误类型应改为类型断言。

<!-- 规范未建，本次为现状盘点与约定起草。 -->
