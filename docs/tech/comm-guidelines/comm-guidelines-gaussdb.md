# GaussDB 通信规范

## 接口清单

| 接口名 | 协议 | 调用位置 | 业务场景 |
|---|---|---|---|
| GET /service/api/getGaussdbInfor | HTTPS | src/dao/db_init.go（getDataSourceFromDBService） | 启动时获取 GaussDB 连接信息 |
| GaussDB 数据库连接 | DB 连接（openGauss/pq） | src/dao/db_init.go（EnsureConnectGaussDB 等） | 业务数据持久化 |

## HTTP

### GET /service/api/getGaussdbInfor

- 业务场景：服务启动建连阶段，向 DB 管理微服务查询 GaussDB 的连接信息（主 IP/端口/用户名/新密码/库名）
- 接口功能：GET `https://{host}:{port}/service/api/getGaussdbInfor?serviceName={serviceName}&dbName={dbName}`，响应体经 `#`→`"` 字符替换后解析为 gaussDbInfo，拼出 DB 连接串
- 调用位置：src/dao/db_init.go（getDataSourceUrl → getDataSourceFromDBService）
- 协议信息：
  - 协议：HTTPS GET；host/port 来自 CSE 发现的 DB 服务实例（dbServiceName 取环境变量 `DB_SERVICE_NAME` 或配置 `gaussdb::servicename`）
  - 封装方式：统一封装层 src/common/https/builder.go，客户端为 `https.InnerInstance()`
  - 超时重试：未调用 WithRetry，未设置显式重试；超时为封装层框架默认（总超时 240s，ResponseHeader/TLS 握手 120s，src/common/https/client.go）
  - 错误码处理：非 2xx 或出错返回 error；当 dbServiceName=="GaussDB" 时降级为本地配置 `gaussdb::*` 拼接连接串（getDataSourceByConfig）

## 外部存储

### GaussDB 数据库连接

- 业务场景：全部业务数据的持久化（用户绑定、配置、插件包、流量统计等）；LOCAL_MODE 下由嵌入式 SQLite 替代，不影响生产链路
- 接口功能：通过 Beego ORM 读写 GaussDB，连接串形如 `host=... port=... user=... password=... dbname=...`
- 调用位置：src/dao/db_init.go（EnsureConnectGaussDB / getDataSourceUrl）；各 DAO 经 BaseInterface 使用
- 协议信息：
  - 协议：openGauss/pq 数据库协议（Beego ORM 注册驱动）
  - 封装方式：Beego ORM 统一封装（src/dao）
  - 超时重试：未识别到显式连接超时设置（框架默认）
  - 错误码处理：建连失败重试逻辑在 EnsureConnectGaussDB；查询错误由各 DAO 上抛
