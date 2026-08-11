# 配置中心 业务规则

> 生成时间：2026-08-11
> 覆盖入口：POST /configCenter/v1/ → controllers.ConfigCenterController.InsertOrUpdate；POST /configCenter/v1/get → controllers.ConfigCenterController.GetFromDB

## 概述

配置中心功能域提供键值配置的写入与查询，配置落库并由内存缓存（5 分钟定时刷新）供运行期读取。规则提取自写入与查询两条入口链路，覆盖参数校验与条件分支两类规则点。

## 规则表

| 规则名 | 条件 | 动作 | 依据 | 来源入口 |
| --- | --- | --- | --- | --- |
| config-center-key-empty | 请求 Key 为空 | 返回 retcode.ClientFailed(-2)，Message "key cannot be empty"，拒绝读写 | src/controllers/config_center_controller.go | POST /configCenter/v1/；POST /configCenter/v1/get |
| insert-or-update-config | 按 Key 查询配置不存在（orm.ErrNoRows） | 事务内插入新配置（UpdatedAt 置当前时间）；已存在则沿用原 ID 更新并刷新 UpdatedAt | src/service/config_center_service.go | POST /configCenter/v1/ |
| get-config-cache | 运行期读配置（GetConfig） | 只读内存缓存 configs，不查库 | src/service/config_center_service.go | POST /configCenter/v1/；POST /configCenter/v1/get |
| refresh-interval | 配置缓存刷新定时任务 | 每 RefreshInterval(5 分钟) 全量重载配置表到内存缓存 | src/service/config_center_service.go | POST /configCenter/v1/；POST /configCenter/v1/get |
| get-from-db-miss | GetFromDB 按 Key 查库失败 | 返回空 ConfigCenter 与 false（接口直接返回空对象，不报错） | src/service/config_center_service.go；src/controllers/config_center_controller.go | POST /configCenter/v1/get |

## 补充说明

写入路径走事务 + 查库，读取分两层：GetFromDB 直接查库（管理面），GetConfig 读内存缓存（运行面，供配置同步等内部逻辑消费），两层间存在最长 5 分钟的缓存延迟窗口——写入后运行面读取不保证立即可见，代码未体现写后主动刷新，待确认是否为有意设计。
