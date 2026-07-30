# 配置中心读写

> 功能域：config-center　接口数：2　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

提供 KV 配置项的内部读写入口，供集群内其他组件写入/查询配置（t_config_center 表）。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| InsertOrUpdate | 写入或更新配置项 | src/controllers/config_center_controller.go | POST /configCenter/v1/ |
| GetFromDB | 按 key 查询配置项 | src/controllers/config_center_controller.go | POST /configCenter/v1/get |

## 3. 数据结构说明

- **InsertOrUpdate**
  - 请求 `db.ConfigCenter`（src/models/db/config_center.go）：Key（必填，controller 校验非空）、Value、Describe、Enable（bool）
  - 响应 `resp.BaseResponse`；失败也返回 ClientFailed + HTTP 400（src/controllers/config_center_controller.go:60）
- **GetFromDB**
  - 请求 `db.ConfigCenter`：仅需 Key（必填非空）
  - 响应 `db.ConfigCenter`（查不到时返回零值结构，`GetFromDB` 的 bool 结果被忽略，src/controllers/config_center_controller.go:42-43）

## 4. 风险与注意点

- **查询缺失语义**：src/controllers/config_center_controller.go:42（key 不存在时返回 200 + 零值 ConfigCenter，调用方无法区分"不存在"与"空值"）
- **GetFromDB 用 POST**：查询接口以 POST + body 传 key，路径尾带 `/get`，与 REST 惯例不一致
