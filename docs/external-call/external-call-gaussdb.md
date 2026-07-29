# DB 服务 + GaussDB 出站调用

GIDS 业务数据持久化在 GaussDB（openGauss）。连接分两步：先通过 CSE 发现 DB 服务主实例，向 DB 服务请求 GaussDB 连接信息，再以 PostgreSQL 协议直连 GaussDB 主库。`LOCAL_MODE=true` 时跳过本组全部出站调用，改用嵌入式 SQLite。

| 接口名 | 协议 | 调用位置 | 业务场景 |
|--------|------|----------|----------|
| GET /service/api/getGaussdbInfor | HTTPS | dao/db_init.go | 获取 GaussDB 连接串 |
| GaussDB 主库连接 | PostgreSQL 协议 | dao/db_init.go + 全部 DAO | 业务数据读写 |

## GET /service/api/getGaussdbInfor

- 协议：HTTPS GET `https://{dbServiceIP}:{port}/service/api/getGaussdbInfor?serviceName={dbServiceName}&dbName={dbName}`，使用带 registry TLS 配置的内部客户端 `https.InnerInstance()`
- 调用位置：dao/db_init.go（getDataSourceFromDBService 函数）；DB 服务实例 IP 由 `cse.NewCse().GetAllMicroServiceInstanceInfo(dbServiceName)` 发现（仅取 Status=UP 且 Properties["status"]="M" 的主实例）
- 业务场景：服务启动建立数据库连接、运行期主备切换（refresh）时，向 DB 服务查询目标 GaussDB 的主库地址、端口、用户名、密码
- 接口功能：返回 `gaussDbInfo`（SERVICE_NAME/MASTER_IP_ADDR/PORT/NEW_DB_PWD/DB_NAME/DB_USER），拼装为连接串；失败且服务名为 "GaussDB" 时回退本地配置 `gaussdb::gaussdbuser` 等（getDataSourceByConfig）

## GaussDB 主库连接（数据读写）

- 协议：PostgreSQL 协议（驱动 `gaussdb_1` = orm.DRPostgres，连接器 `gitee.com/opengauss/openGauss-connector-go-pq`），通过 Beego ORM `orm.RegisterDataBase` 建立
- 调用位置：dao/db_init.go（switchToAnotherDB / EnsureConnectGaussDB / initTables / checkDBStatus / refresh）；业务读写在 dao/ 各 DAO（UserBindDao、FileDao、ConfigCenterDao、PluginPackageDao、流量统计相关等）
- 业务场景：全量业务数据持久化——用户绑定（t_user_bind）、终端信息（t_user）、文件（t_file）、插件包（t_plugin_package）、配置（t_config / t_config_center）、流量与会话统计（t_media_traffic_stats / t_control_traffic_stats / t_session_stats）
- 接口功能：建表初始化（initTables 执行 DDL）、每 5s 健康检查（Ping）、主备切换时自动重连新主库
