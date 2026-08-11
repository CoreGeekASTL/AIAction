# 过载限流 业务规则

> 生成时间：2026-08-11
> 覆盖入口：全局 BeforeRouter 过滤器 controllers.OverLoadFilter（src/routers/beego_router.go 对内外部全部路由 "*" 注册）

## 概述

过载限流功能域在所有业务请求进入路由前按"路径+方法"维度做准入控制，保护服务不过载。规则提取自全局过滤器链路，覆盖条件分支与错误码返回两类规则点。

## 规则表

| 规则名 | 条件 | 动作 | 依据 | 来源入口 |
| --- | --- | --- | --- | --- |
| over-load-filter-dim | 请求进入路由前，按 dimNameValues[APIService] = URL.Path + "/" + Method 维度做限流判定 | overloadcontroller.Process 判定放行或拒绝 | src/controllers/filter.go | 全部路由（BeforeRouter） |
| over-load-filter-reject | Process 返回 isGranted=false | 返回 HTTP 429，响应头 Retry-After=3（retryAfter），响应体 "Too many requests, please retry later."，请求不再进入业务 handler | src/controllers/filter.go | 全部路由（BeforeRouter） |
| over-load-filter-error | overloadcontroller.Process 自身执行出错 | 仅记录错误日志，按 isGranted 结果处理（出错不放行时同样走 429 分支） | src/controllers/filter.go | 全部路由（BeforeRouter） |

## 补充说明

限流为全局前置短路规则，先于一切业务参数校验与业务分支执行，被 429 拒绝的请求不产生任何业务副作用。限流维度为"路径+方法"（FilterConfKey="APIService"），具体阈值配置由 greatwall-sdk-go 外部组件承载，代码未体现，待确认。
