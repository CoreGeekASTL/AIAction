# external-call-gaussdb

> 下游服务：GaussDB（CSE 注册服务名 `SbgGaussDB`，见 app.conf `gaussdb::servicename`；数据库名 `sbggidsdb`，用户 `sbggidsdbuser`，默认端口 23521）。
> 主库地址来源：CSE 查询 `SbgGaussDB` 实例（dao/db_init.go `getGaussDBIP`，仅取主库 properties.status="M"）；连接串优先从 GaussDB 服务的 HTTP 接口获取，失败回落 app.conf `[gaussdb]` 配置拼接。
> 调用方式：Beego ORM（驱动 `gaussdb_1` = `orm.DRPostgres`，自研 driver 封装见 src/db/driver/）；`LOCAL_MODE=true` 时不连 GaussDB，改用嵌入式 SQLite。

| 接口名 | 协议 | 调用位置 | 业务场景一句话 |
|--------|------|----------|----------------|
| GaussDB SQL 连接（全量表 CRUD） | Postgres 协议 | dao/db_init.go + dao/*.go | 全部持久化数据的读写 |
| GET /service/api/getGaussdbInfor | HTTPS | dao/db_init.go | 从 DB 服务获取数据库连接串 |

## SQL

## GaussDB SQL 连接（全量表 CRUD 与统计查询）

- 协议：Postgres 协议（Beego ORM `RegisterDataBase(alias, "gaussdb_1", connStr)`；connStr 形如 `host=<ip> port=<port> user=<u> password=<p> dbname=<db>`）
- 调用位置：dao/db_init.go（`EnsureConnectGaussDB` 建连、`initTables` 建表、`checkDBStatus`/`refresh` 主备切换重连）；业务读写分散在 dao/ 全部 DAO（user.go、file.go、plugin.go、browser_config.go、config_center.go、traffic_stats_dao.go、base_dao.go）
- 业务场景：服务启动时连接 GaussDB 主库并初始化全部表（t_config / t_file / t_plugin_package / t_user / t_user_bind / t_media_traffic_stats / t_control_traffic_stats / t_session_stats / t_config_center）；运行期承载用户绑定、插件包、文件内容、浏览器配置、配置中心、流量/会话统计等全部持久化读写；每 5s 健康检查，主库 IP 变化或连续 3 次 ping 失败时自动切换重连
- 接口功能：DDL 建表 + DML CRUD + 统计类 SQL（traffic_stats_service.go 按 conf/sql.yaml 中的 SQL 模板做在线人数/流量聚合查询）

## HTTPS

## GET /service/api/getGaussdbInfor

- 协议：HTTPS GET `https://<gaussdbHost>:<port>/service/api/getGaussdbInfor?serviceName=<dbServiceName>&dbName=<dbName>`（客户端 `https.InnerInstance()`，内部 TLS 配置 `registry`）
- 调用位置：dao/db_init.go（`getDataSourceFromDBService` 函数，由 `getDataSourceUrl` 在每次建连/重连前调用）
- 业务场景：连库前从 GaussDB 服务侧动态获取主库连接信息（含轮换后的数据库密码），避免明文密码只依赖静态配置；获取失败且服务名为 `GaussDB` 时回落 `getDataSourceByConfig` 用 app.conf `[gaussdb]` 配置拼接
- 接口功能：响应体 JSON（`#` 替换为 `"` 后解析）`gaussDbInfo{SERVICE_NAME, MASTER_IP_ADDR, PORT, NEW_DB_PWD, DB_NAME, DB_USER}`，拼成 Postgres connStr 返回
