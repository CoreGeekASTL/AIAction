# custom-container（自定义容器）

> 类型：custom-container　实例数：4　返回 [README.md](README.md)

## 1. 定位

GIDS 的核心状态承载全部走自定义容器——鉴权结果缓存、配置缓存、告警抑制计数、网关实例表各以"带锁/带淘汰策略/并发安全"的封装 map 实现，是业务流程的数据中枢。

## 2. 关键实例清单

| 实例名 | 作用 | 定义位置 |
|---|---|---|
| authCache | 鉴权结果内存缓存（TTL 30min + 容量惰性清理） | service/auth_cache.go |
| configCenter.configs | 配置中心缓存（5min 定时整体刷新） | service/config_center_service.go |
| alarmService.alarms | 告警抑制计数（AlarmID → 上次发送时间戳） | service/alarm_service.go |
| cse.browserGWInstances | 浏览器网关实例表（CSE 服务发现回调写入） | common/cse/cse.go |

## 3. 实例详解

- **authCache**
  - 结构：`type authCache struct { sync.RWMutex; items map[string]cacheEntry }`（service/auth_cache.go:27）
  - 关键字段：items（key=IMEI+IMSI 组合，value=cacheEntry{result bool, expireAt time.Time}）；内嵌 RWMutex
  - 典型操作：get 走 RLock 查+过期判断；set 走 Lock 写入后若 `len>authCacheCapacity(1000)` 触发 cleanLocked 按 expireAt 升序删最旧 authCacheCleanCount(500) 条
  - 使用点：service/auth_service.go:64（get 命中跳过 DB）、:68（set 回源后回填）、:73（clear 白名单变更后清空）
  - 并发模型：读写锁，读多写少；惰性清理在写锁内完成，无独立 goroutine
- **configCenter.configs**
  - 结构：`type configCenterServiceImpl struct { configs map[string]string; dao *dao.ConfigCenterDao; stopChan chan struct{} }`（service/config_center_service.go:30）
  - 关键字段：configs（key=配置 Key，value=配置 Value）；stopChan（控制刷新 goroutine 退出）
  - 典型操作：Refresh() 重新 List 全表构造新 map 后整体赋值 `c.configs = configMap`；GetConfig() 裸读 `configCenter.configs[key]`
  - 使用点：service/config_center_service.go:60（Refresh 整体替换）、:65（GetConfig 裸读）、:92-95（init 单例+stopChan）
  - 并发模型：**无锁**，独立 goroutine 定时 Refresh 整体替换，主线程并发裸读——存在并发读写隐患（见 README 全局风险）
- **alarmService.alarms**
  - 结构：`type alarmServiceImpl struct { alarms map[string]int64; alarmManager base.CSPAlarmManager }`（service/alarm_service.go:48）
  - 关键字段：alarms（key=AlarmID，value=上次发送 UnixMilli 时间戳）
  - 典型操作：sendAlarm 读 lastSendTime 判断 10min 内抑制；发送成功更新时间戳
  - 使用点：service/alarm_service.go:91（读抑制判断）、:94（抑制比较）、:69（init 初始化 map+channel+启消费 goroutine）
  - 并发模型：由 `handleEvent` 单 goroutine 消费 alarmEventChanel 时调用 sendAlarm/clearAlarm，alarms 访问限该 goroutine 内，当前无锁安全；跨 goroutine 直调须加锁
- **cse.browserGWInstances**
  - 结构：`type cse struct { register api.Registry; appid string; browserGWInstances sync.Map; chainEndpoints []string }`（common/cse/cse.go:31）
  - 关键字段：browserGWInstances（sync.Map，key=实例标识，value=browsergateway.ServiceInstance）
  - 典型操作：WatchServiceCallBack 回调写入；GetAllBrowserGateWayInstances Range 遍历收集；GetBrowserGateWayInstanceByInnerEndpoint Range 命中即止
  - 使用点：common/cse/cse.go:41（Range 查询）、:75（Init 初始化 sync.Map{}）、:96（回调写入）
  - 并发模型：sync.Map，CSE 回调线程写 + 业务线程并发读，原生并发安全

## 4. 使用模式与约定

- 鉴权/告警类缓存统一"struct 内嵌 RWMutex + map 字段"封装，get/set/clear 各自加对应锁（authCache 范式）
- CSE 服务发现实例表用 sync.Map + Range 查询，不转普通 map 加锁
- 容量型缓存（authCache）走惰性清理：写后判超限、按过期时间升序删最旧一批，无独立清理 goroutine

## 5. AI 编码指南

1. 新增"key→业务对象"缓存时，复用 authCache 同款"struct 内嵌 RWMutex + map"封装 + 容量惰性清理，禁止裸 map + 外部加锁（依据：service/auth_cache.go:27，仓内统一范式）
2. 新增定时刷新缓存须加读写锁或改 sync.Map，**禁止**沿用 configCenter 的无锁裸 map 整体替换——会触发 `concurrent map read and map write`（依据：service/config_center_service.go:31/60，现存隐患）
3. CSE/服务发现类并发实例表用 sync.Map + Range，禁止加锁转普通 map（依据：common/cse/cse.go:34，回调线程与业务线程并发）

## 6. 风险与注意点

- **configCenter.configs 并发读写隐患**：service/config_center_service.go:31（无锁裸 map，Refresh 整体替换与 GetConfig 并发读会 fatal，待修复）
- **authCache 无容量上限兜底之外无主动过期**：service/auth_cache.go:52（仅写时惰性清理，长期未写的过期条目不主动清，但 TTL 命中判断兜底，影响有限）
