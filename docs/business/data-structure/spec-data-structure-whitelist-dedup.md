# 白名单去重

> 用途：whitelist-dedup　实例数：1　返回 [README.md](README.md)

## 1. 核心作用

承载白名单导入去重的"存在性查询"——以 `map[K]struct{}` 模拟集合语义，跳过重复的 IMEI+IMSI 行。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| seen | 白名单导入去重集合（跳过重复 IMEI+IMSI 行） | service/auth_manage_service.go |

## 3. 实例详解

- **seen**
  - 结构：`seen := make(map[string]struct{}, len(rows))`（service/auth_manage_service.go:132），局部 set
  - 关键操作：遍历白名单 CSV 行，以 IMEI+IMSI 拼 key 写入 seen，命中即跳过重复行
  - 使用点：service/auth_manage_service.go:132（创建）、后续行写入/查询去重
  - 并发模型：单 goroutine 内局部变量，无并发

## 4. 使用模式与约定

- 去重场景一律 `map[K]struct{}`（零内存），key 即集合元素，存在性靠 `_, ok := m[k]`
