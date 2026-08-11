# routers 模块结构文档

> 生成时间：2026-08-11
> 所属仓：AIAction（GIDS）
> 模块路径：src/routers

## 模块职责

Beego 路由注册层：在 beego_router.go 中集中定义 URL 路径与 Controller 的映射，由 main.go 导入生效。

## 子模块关系图

本模块为扁平包结构，无子模块。

## 子模块说明

| 关键文件 | 说明 |
| --- | --- |
| beego_router.go | 定义 Beego 路由配置，注册全部 Controller 路由 |
