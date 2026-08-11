# db 模块结构文档

> 生成时间：2026-08-11
> 所属仓：AIAction（GIDS）
> 模块路径：src/db

## 模块职责

数据库驱动适配层：唯子包 driver 代理 openGauss-connector-go-pq 的 driver，提取 SQL 表名等以适配 Beego ORM 的驱动注册要求，由 dao/db_init.go 匿名导入生效。

## 子模块关系图

本模块仅 1 个子包 driver，无子包间依赖。

## 子模块说明

| 子模块 | 路径 | 职责 | 主要依赖 | 被依赖 |
| --- | --- | --- | --- | --- |
| driver | src/db/driver | 代理 gauss db driver，适配 beego orm。证据：src/db/driver/driver.go | - | - |
