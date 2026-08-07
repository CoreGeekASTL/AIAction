# 关键类清单

> 分析基准：ready/27.0-终端鉴权 分支（2026-08-07）　类总数：17

| 类名 | 类的职责 |
| --- | --- |
| BaseController | 全部控制器基类，统一请求体解析校验与 JSON 响应封装 |
| LoginController | 登录鉴权主链路入口，编排建档、实例路由与事件上报 |
| BrowserService | 登录核心：为用户分配/校验浏览器网关实例并签发 token |
| UserService | 用户档案与 UserBind 绑定关系的查询、更新与过期处理 |
| TrafficStatsService | 会话/流量统计的批量入库、配置化查询、CSV 导出与清理 |
| PluginService | 插件包上传与加载编排，向 BrowserGW 分发并跟踪进度 |
| ConfigCenterService | 配置中心单例缓存，承载动态配置读取与 5 分钟定时刷新 |
| EventService | 登录/业务事件上报统一入口，经存储工厂落本地审计文件 |
| MonitorService | CSP 话统注册与 5 分钟粒度指标定时上报调度 |
| AlarmService | 告警发送/清除单例，channel 异步消峰并抑制重复告警 |
| BaseDao | 全部 DAO 基座，封装 Beego ORM 的 CRUD、原生 SQL 与事务 |
| dbConnection | GaussDB 连接管理：主库发现、健康检查与故障切换 |
| Cse | 服务发现单例，维护 BrowserGW 实例注册表，路由分配数据源 |
| request（https.Builder） | 出站 HTTP 请求 builder，带重试，全仓外呼统一入口 |
| UserBind | 用户-浏览器实例绑定核心模型，承载 token 与端点状态 |
| ServiceInstance | BrowserGW 实例模型，含容量/健康字段与负载排序规则 |
| DataCleanupScheduler | 统计数据每日 2 点清理调度器，带重试与优雅停止 |
