# 查询全部插件包（GetPluginPackages） 交互模型

> 生成时间：2026-08-11
> 流程入口：`POST /plugin/v1/getAll` → controllers.PluginController.GetPluginPackages（内部 HTTP 服务）

## 概述

内部调用方查询全部插件包列表，controllers 经 service 到 dao 全量读取 t_plugin_package 并返回。

## 主链路时序图

```mermaid
sequenceDiagram
    participant A as "内部调用方"
    participant C as "controllers"
    participant S as "service"
    participant D as "dao"
    participant DB as "DB（SQLite/GaussDB）"
    A->>C: "POST /plugin/v1/getAll"
    C->>S: "GetPluginPackages()"
    S->>D: "PluginPackageDao.List(ppl)"
    D->>DB: "查询 t_plugin_package"
    DB-->>D: "插件包列表"
    D-->>S: "列表"
    S-->>C: "[]PluginPackage"
    C-->>A: "200 + 插件包列表"
```

## 参与方说明

| 参与方 | 类型 | 代码位置 | 本流程中的职责 |
| --- | --- | --- | --- |
| 内部调用方 | 外部触发者 | -（仓外） | 查询全部插件包 |
| controllers | 模块 | src/controllers/plugin_controller.go | 入口与响应回写 |
| service | 模块 | src/service/plugin_service.go | 全量列表查询 |
| dao | 模块 | src/dao/plugin.go | t_plugin_package 读取 |
| DB（SQLite/GaussDB） | 中间件 | src/dao/db_local_sqlite.go、src/dao/db_init.go | 插件包持久化 |

## 补充说明

本流程只读，无实体状态变更。

分支与异常逻辑：归业务规则资产承载（docs/biz/rules/，待补）。
