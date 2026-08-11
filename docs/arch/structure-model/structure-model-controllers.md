# controllers 模块结构文档

> 生成时间：2026-08-11
> 所属仓：AIAction（GIDS）
> 模块路径：src/controllers

## 模块职责

Beego Controller 层：接收并校验 HTTP 请求，调用 Service 层完成业务，按统一响应结构返回。覆盖登录鉴权（login/exlogin）、终端鉴权与白名单管理（auth）、事件上报（event）、插件（plugin）、文件（file/exfile）、缓存（cache）、配置中心（config_center）、管理（management）、流量统计（traffic_stats）等接口；filter.go 提供过滤器，controller.go 提供公共基类能力。

## 子模块关系图

本模块为扁平包结构，无子模块。

## 子模块说明

| 关键文件 | 说明 |
| --- | --- |
| controller.go | Controller 公共基类与通用处理（引用 dao、models/req/resp） |
| filter.go | 请求过滤器（鉴权等前置处理） |
| login_controller.go / exlogin_controller.go | 终端登录链路接口 |
| event_controller.go | 事件上报链路接口 |
| plugin_controller.go | 插件相关接口 |
| file_controller.go / exfile_controller.go | 文件导入导出接口 |
| cache_controller.go | 缓存管理接口 |
| config_center_controller.go | 配置中心接口 |
| management_controller.go | 管理面接口 |
| traffic_stats_controller.go | 流量统计接口（引用 utils） |
| test_controller.go | 测试用途接口 |
| auth_controller.go | 终端鉴权与白名单管理接口（仅内网注册，引用 service/auth、retcode） |
