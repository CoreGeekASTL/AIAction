# 配置中心 领域词典

> 子文档 of [lexicon.md](lexicon.md)；词汇口径、来源说明与全仓待确认清单见主文档。子域锚点 `config-center`，功能域口径与 docs/biz/interface/ 一致。

## 实体与业务概念

| 术语 | 释义 | 语境边界 | 代码命名映射 |
| --- | --- | --- | --- |
| ConfigCenter | 配置中心表（t_config_center）：config_key/config_value/config_describe/enable/updated_at | - | `db.ConfigCenter`，`src/models/db/config_center.go` |
| Key | 配置键（config_key 列） | - | `db.ConfigCenter.Key`，`src/models/db/config_center.go` |
| Value | 配置值（config_value 列） | - | `db.ConfigCenter.Value`，`src/models/db/config_center.go` |
| Describe | 配置描述（config_describe 列） | - | `db.ConfigCenter.Describe`，`src/models/db/config_center.go` |
| Enable | 配置启用开关（bool） | - | `db.ConfigCenter.Enable`，`src/models/db/config_center.go` |
