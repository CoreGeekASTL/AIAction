# 服务发现

> 用途：service-discovery　实例数：2　返回 [README.md](README.md)

## 1. 核心作用

承载 CSE 服务发现相关状态——浏览器网关实例表以并发安全的封装 map 实现、CSE 回调写入，链路端点列表动态 append 后运行期只读，是服务发现链路的数据中枢。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| cse.browserGWInstances | 浏览器网关实例表（CSE 服务发现回调写入） | common/cse/cse.go |
| chainEndpoints | CSE 链路端点列表（动态 append） | common/cse/cse.go |

## 3. 实例详解

- **cse.browserGWInstances**
  - 结构：`type cse struct { register api.Registry; appid string; browserGWInstances sync.Map; chainEndpoints []string }`（common/cse/cse.go:31）
  - 关键字段：browserGWInstances（sync.Map，key=实例标识，value=browsergateway.ServiceInstance）
  - 典型操作：WatchServiceCallBack 回调写入；GetAllBrowserGateWayInstances Range 遍历收集；GetBrowserGateWayInstanceByInnerEndpoint Range 命中即止
  - 使用点：common/cse/cse.go:41（Range 查询）、:75（Init 初始化 sync.Map{}）、:96（回调写入）
  - 并发模型：sync.Map，CSE 回调线程写 + 业务线程并发读，原生并发安全
- **chainEndpoints**
  - 结构：`chainEndpoints []string`（common/cse/cse.go:35），cse struct 字段
  - 关键操作：Init 时 `make([]string, 0, 2)` 预分配；AddChainEndpoint 动态 append 端点
  - 使用点：common/cse/cse.go:35（字段定义）、:76（Init 预分配）、AddChainEndpoint 方法 append
  - 并发模型：单 goroutine 初始化，运行期只读；无锁

## 4. 使用模式与约定

- CSE 服务发现实例表用 sync.Map + Range 查询，不转普通 map 加锁
- 动态增长列表用 `[]string` + `make([]T, 0, cap)` 预分配 + append（chainEndpoints 范式）
