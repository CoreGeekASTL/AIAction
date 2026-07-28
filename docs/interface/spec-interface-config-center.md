# 配置中心

> 功能域：配置中心　接口数：2　所属 server：内部(HTTP)
> 子文档 of [README.md](README.md)

## 1. 定位

键值配置中心，按 key 增删改查，innerServer（HTTP）暴露，ConfigCenterController 承载。底层走 service.ConfigCenterService → dao（t_config_center 表）。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| InsertOrUpdate | 插入或更新配置 | controllers/config_center_controller.go | POST /configCenter/v1/ |
| GetFromDB | 按 key 查询配置 | controllers/config_center_controller.go | POST /configCenter/v1/get |

## 3. 数据结构说明

- **InsertOrUpdate / GetFromDB**
  - 请求 `db.ConfigCenter`（models/db）：Key（必填，非空校验）；Value；Describe；Enable（默认 1）
  - 响应：GetFromDB 返回 `db.ConfigCenter`；InsertOrUpdate 返回 retcode 标准结构

## 4. 风险与注意点

- **无 value 长度限制**：controllers/config_center_controller.go:48（仅校验 Key 非空，Value 长度未约束，可能写入超大配置）
