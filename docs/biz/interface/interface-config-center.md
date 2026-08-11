# 配置中心

> 功能域：配置中心　接口数：2　所属 server：内部
> 子文档 of [README.md](README.md)

## 1. 定位

以 key-value 形式持久化动态配置（如 Muen 云服务地址覆盖项 `moon::configEndpoint` 等），供仓内各服务读取。仅注册在内部 server。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| InsertOrUpdate | 写入或更新一条配置 | controllers/config_center_controller.go | POST /configCenter/v1/ |
| GetFromDB | 按 key 查询配置 | controllers/config_center_controller.go | POST /configCenter/v1/get |

## 3. 数据结构说明

- **InsertOrUpdate / GetFromDB**
  - 请求 `db.ConfigCenter`（models/db/config_center.go）：key（必填，空则返回 client 错误）、value、describe、enable
  - 响应：GetFromDB 返回命中的 `db.ConfigCenter` 记录；InsertOrUpdate 返回 `resp.BaseResponse`
