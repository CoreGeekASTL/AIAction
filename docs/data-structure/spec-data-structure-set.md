# set（集合）

> 类型：set　实例数：2　返回 [README.md](README.md)

## 1. 定位

GIDS 无原生 set，统一以 `map[K]struct{}` 或 `map[K]bool` 模拟集合语义，承担白名单导入去重与 TLS 弱算法判定的"存在性查询"职责。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| seen | 白名单导入去重集合（跳过重复 IMEI+IMSI 行） | service/auth_manage_service.go |
| weakAlgorithms | TLS 弱签名算法集合（证书校验判定） | common/https/tls.go |

## 3. 实例详解

- **seen**
  - 结构：`seen := make(map[string]struct{}, len(rows))`（service/auth_manage_service.go:132），局部 set
  - 关键操作：遍历白名单 CSV 行，以 IMEI+IMSI 拼 key 写入 seen，命中即跳过重复行
  - 使用点：service/auth_manage_service.go:132（创建）、后续行写入/查询去重
  - 并发模型：单 goroutine 内局部变量，无并发
- **weakAlgorithms**
  - 结构：`weakAlgorithms := map[x509.SignatureAlgorithm]bool{...}`（common/https/tls.go:205），包内字面量集合
  - 关键操作：证书校验时 `if weakAlgorithms[cert.SignatureAlgorithm]` 判定弱算法并告警
  - 使用点：common/https/tls.go:205（字面量定义）、后续判定查询
  - 并发模型：字面量初始化后只读，无并发写

## 4. 使用模式与约定

- 去重场景一律 `map[K]struct{}`（零内存），key 即集合元素，存在性靠 `_, ok := m[k]`
- 包级常量集合可用 `map[K]bool` 字面量（weakAlgorithms 历史写法），新代码建议 `struct{}` value

## 5. AI 编码指南

1. 新增进程内去重集合用 `map[K]struct{}`，禁止 `map[K]bool`（依据：service/auth_manage_service.go:132，仓内去重范式；weakAlgorithms 的 bool 为历史遗留）
2. 存在性查询直接 `_, ok := m[k]`，禁止遍历集合判定（依据：set 语义）
