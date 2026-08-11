# service 模块结构文档

> 生成时间：2026-08-11
> 所属仓：AIAction（GIDS）
> 模块路径：src/service

## 模块职责

业务 Service 层：按「接口 + 小写实现 + 包级变量 + sync.Once 单例」模式组织，承载浏览器实例交付（browser）、插件（plugin）、文件（file）、用户（user）、事件（event）、缓存（cache）、监控（monitor）、告警（alarm）、配置中心（config_center）、流量统计（traffic_stats）、远端调用（remote）、终端鉴权（auth/auth_cache/whitelist_manage）等业务逻辑；traffic_stats_service_mock.go 提供对应 mock。

## 子模块关系图

本模块为扁平包结构，无子模块。

## 子模块说明

| 关键文件 | 说明 |
| --- | --- |
| browser_service.go | 云浏览器实例交付核心业务（引用 cse/conf/dao/models 等） |
| plugin_service.go | 插件业务逻辑 |
| file_service.go | 文件业务逻辑 |
| user_service.go | 用户业务逻辑 |
| event_service.go | 事件业务逻辑（引用 common/event） |
| cache_service.go | 缓存业务逻辑 |
| monitor_service.go | 监控业务逻辑（引用 utils/monitorutil） |
| alarm_service.go | 告警业务逻辑 |
| config_center_service.go | 配置中心业务逻辑 |
| traffic_stats_service.go / traffic_stats_service_mock.go | 流量统计业务及 mock |
| remote_service.go | 远端服务调用封装 |
| auth_service.go / auth_cache.go | 终端白名单联合鉴权与进程内鉴权缓存（引用 dao、models/db） |
| whitelist_manage_service.go | 白名单 CSV 导入/导出管理 |
