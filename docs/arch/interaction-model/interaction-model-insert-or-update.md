# 配置中心写入（InsertOrUpdate） 交互模型

> 生成时间：2026-08-11
> 流程入口：`POST /configCenter/v1/` → controllers.ConfigCenterController.InsertOrUpdate（内部 HTTP 服务）

## 概述

内部调用方写入配置项，service 在事务中按键判断存在性，不存在则插入、存在则更新配置中心表。

## 主链路时序图

```mermaid
sequenceDiagram
    participant A as "内部调用方"
    participant C as "controllers"
    participant S as "service"
    participant D as "dao"
    participant DB as "DB（SQLite/GaussDB）"
    A->>C: "POST /configCenter/v1/"
    C->>S: "InsertOrUpdateConfig(config)"
    S->>D: "ConfigCenterDao.Get(Key)"
    D->>DB: "查询 t_config_center"
    DB-->>D: "记录（或不存在）"
    D-->>S: "结果"
    S->>D: "ConfigCenterDao.Insert / Update（事务内）"
    D->>DB: "写入 t_config_center"
    DB-->>D: "成功"
    D-->>S: "成功"
    S-->>C: "成功"
    C-->>A: "200"
```

## 参与方说明

| 参与方 | 类型 | 代码位置 | 本流程中的职责 |
| --- | --- | --- | --- |
| 内部调用方 | 外部触发者 | -（仓外） | 写入配置项 |
| controllers | 模块 | src/controllers/config_center_controller.go | 解析并校验 key 非空，回写响应 |
| service | 模块 | src/service/config_center_service.go | 事务内 upsert 配置项 |
| dao | 模块 | src/dao/config_center.go | t_config_center 读写 |
| DB（SQLite/GaussDB） | 中间件 | src/dao/db_local_sqlite.go、src/dao/db_init.go | 配置项持久化 |

## 补充说明

关键实体状态变更：t_config_center 对应 Key 记录被插入或更新（UpdatedAt 刷新）。写入后内存缓存 configs 由定时刷新任务（见 interaction-model-refresh-config-task）周期同步，非实时失效。

分支与异常逻辑：归业务规则资产承载（docs/biz/rules/，待补）。
