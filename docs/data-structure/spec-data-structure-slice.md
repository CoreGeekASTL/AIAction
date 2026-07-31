# slice（切片）

> 类型：slice　实例数：2　返回 [README.md](README.md)

## 1. 定位

slice 在本仓作为"端点/枚举列表"的轻量承载——CSE 链路端点动态 append、全局告警 ID 枚举字面量，是配置型列表的主要形式。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| chainEndpoints | CSE 链路端点列表（动态 append） | common/cse/cse.go |
| AlarmList | 全局告警 ID 枚举列表 | service/alarm_service.go |

## 3. 实例详解

- **chainEndpoints**
  - 结构：`chainEndpoints []string`（common/cse/cse.go:35），cse struct 字段
  - 关键操作：Init 时 `make([]string, 0, 2)` 预分配；AddChainEndpoint 动态 append 端点
  - 使用点：common/cse/cse.go:35（字段定义）、:76（Init 预分配）、AddChainEndpoint 方法 append
  - 并发模型：单 goroutine 初始化，运行期只读；无锁
- **AlarmList**
  - 结构：`var AlarmList = []string{AlarmId300010}`（service/alarm_service.go:189），包级字面量
  - 关键操作：定义需订阅的告警 ID 集合，告警初始化时遍历注册
  - 使用点：service/alarm_service.go:189（字面量定义）、告警注册处遍历
  - 并发模型：字面量初始化后只读，无锁

## 4. 使用模式与约定

- 动态增长列表用 `[]string` + `make([]T, 0, cap)` 预分配 + append（chainEndpoints 范式）
- 全局枚举列表用包级 `var x = []T{...}` 字面量，运行期只读（AlarmList 范式）

## 5. AI 编码指南

1. 动态增长切片用 `make([]T, 0, 预估容量)` 预分配，禁止裸 `[]T{}` 频繁扩容（依据：common/cse/cse.go:76，cap=2 预分配）
2. 全局枚举集合用包级 `var` 字面量切片，运行期禁止改写（依据：service/alarm_service.go:189）
