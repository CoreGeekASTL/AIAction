# 鉴权缓存

> 用途：auth-cache　实例数：1　返回 [README.md](README.md)

## 1. 核心作用

承载鉴权结果的内存缓存，以"带锁/带淘汰策略/并发安全"的封装 map 实现，是业务流程的数据中枢。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| authCache | 鉴权结果内存缓存（TTL 30min + 容量惰性清理） | service/auth_cache.go |

## 3. 实例详解

- **authCache**
  - 结构：`type authCache struct { sync.RWMutex; items map[string]cacheEntry }`（service/auth_cache.go:27）
  - 关键字段：items（key=IMEI+IMSI 组合，value=cacheEntry{result bool, expireAt time.Time}）；内嵌 RWMutex
  - 典型操作：get 走 RLock 查+过期判断；set 走 Lock 写入后若 `len>authCacheCapacity(1000)` 触发 cleanLocked 按 expireAt 升序删最旧 authCacheCleanCount(500) 条
  - 使用点：service/auth_service.go:64（get 命中跳过 DB）、:68（set 回源后回填）、:73（clear 白名单变更后清空）
  - 并发模型：读写锁，读多写少；惰性清理在写锁内完成，无独立 goroutine

## 4. 使用模式与约定

- 鉴权/告警类缓存统一"struct 内嵌 RWMutex + map 字段"封装，get/set/clear 各自加对应锁（authCache 范式）
- 容量型缓存（authCache）走惰性清理：写后判超限、按过期时间升序删最旧一批，无独立清理 goroutine
