# MinIO 数据访问规范

| 元信息 | 值 |
|--------|-----|
| 分支 | 27.0 分支 (2026-08-11) |
| 更新日期 | 2026-08-11 |
| Skill | tech-data-access-guidelines-analyze |
| 运行模式 | 起草模式 |

## 用途定位

MinIO（`github.com/minio/minio-go/v7`）封装为对象存储客户端（src/common/storage/oss/minio.go），规划上承载插件包/文件等对象存取；连接配置取 app.conf `oss::endpoint`，AccessKey/SecretKey 当前硬编码默认值 `minioadmin`（src/common/conf/config.go）。**当前仓内未发现任何业务调用方**：文件内容实际落 GaussDB t_file 表（src/dao/file.go），MinIO 属于预留能力，详见附注。

## 访问点分布

| 访问点 | 所在文件 | 所在函数 | 访问方式 | 业务场景 |
| --- | --- | --- | --- | --- |
| 客户端初始化与健康检查 | src/common/storage/oss/minio.go | Init | 封装层自身（minio.New + IsOnline） | 服务启动初始化（当前仓内未发现 Init 调用方） |
| 对象读写接口 | src/common/storage/oss/minio.go | ossClient.PutObject/GetObject/DeleteObject/EnsureBucket | 封装层：src/common/storage/oss/minio.go | 预留，无业务调用方 |

## 连接与客户端管理

### 现状

- 包级单例 `var client Client`，经 `Init(config conf.OSSConfig)` 创建（`minio.New`，StaticV4 凭证，`Secure:false` 即明文 HTTP），`Instance()` 获取（src/common/storage/oss/minio.go）。
- 配置来源：`oss::endpoint`（app.conf）；AccessKey/SecretKey 硬编码默认值 `minioadmin/minioadmin`，Token 空（src/common/conf/config.go）。
- 初始化后 `IsOnline()` 健康检查，失败返回 error（src/common/storage/oss/minio.go）。
- 无连接池概念（HTTP client 由 minio-go 内部管理，参数未显式配置）；无显式关闭点。
- 未发现 Init 的启动调用方（src/main.go 等启动路径未调用），当前进程内 client 保持 nil。

### 约定（起草待评审）

1. 启用前必须在启动路径统一调用 `oss.Init` 并处理失败，业务代码一律经 `oss.Instance()` 获取单例。
2. 凭证禁止硬编码默认值上线，AccessKey/SecretKey 必须从配置/密钥管理注入；生产环境应启用 TLS（`Secure:true`）。
3. `Instance()` 使用前须有 nil 防护，或 Init 失败即拒绝启动。

## 事务使用

### 现状

本中间件不涉及（对象存储无事务语义；仓内亦无"对象存储+DB"跨中间件一致性代码）。

### 约定（起草待评审）

1. 对象与元数据分离存储时（如插件包：对象在 MinIO、元数据在 DB），须定义一致性与补偿路径（先写对象成功再落元数据；失败清理孤儿对象），启用前在 service 层落地。

## 分页与批量

### 现状

未提供 ListObjects/批量删除接口；`GetObject` 返回 `io.ReadCloser` 为流式读取（src/common/storage/oss/minio.go）。

### 约定（起草待评审）

1. 大对象读取保持流式（io.ReadCloser），禁止一次性读入内存；调用方必须关闭返回的 ReadCloser。
2. 需要列举对象时在封装层补 ListObjects 分页接口，业务层不得直连 minio-go。

## SQL 拼接与注入防护

### 现状

本中间件不涉及。

### 约定（起草待评审）

本中间件不涉及。

## 缓存读写模式

### 现状

本中间件不涉及。

### 约定（起草待评审）

本中间件不涉及。

## 错误处理

### 现状

- `EnsureBucket`：`BucketExists` 失败上抛，不存在则 `MakeBucket` 创建，失败上抛（src/common/storage/oss/minio.go）。
- `PutObject`/`DeleteObject`/`GetObject` 错误原样上抛并记日志；`DeleteObject` 使用 `ForceDelete:true`（src/common/storage/oss/minio.go）。
- 无重试、无降级逻辑；minio-go 内部默认重试行为未显式配置。
- 未区分"对象不存在（NoSuchKey）"与系统错误，统一按 error 上抛。

### 约定（起草待评审）

1. 调用方须显式区分 NoSuchKey（业务 miss）与连接/超时故障（系统错误），封装层宜提供类似 `ErrNotExist` 的统一语义（对齐 redis 封装）。
2. 写失败不重试；读失败可按业务降级（如回源 DB 内嵌内容）。
3. 启用前评估并显式配置 minio-go 的重试/超时参数。

## 附注

- **预留死代码**：src/common/storage/oss/minio.go 整套 Client 封装在当前仓内无任何业务调用方，`oss.Init` 亦无启动调用点；文件与插件包内容当前经 src/dao/file.go 落 GaussDB t_file 表（bytea），不经过对象存储。
- **配置声明 vs 实际访问差异**：src/common/conf/config.go 声明了 `oss::endpoint` 与 `minioadmin/minioadmin` 默认凭证，但代码无实际访问路径，启用前需先接启动初始化并迁移文件存储路径。

<!-- 规范未建，本次为现状盘点与约定起草。 -->
