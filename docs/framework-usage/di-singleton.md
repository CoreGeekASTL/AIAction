# 自研单例/组件模式使用指导（依赖注入/组件管理）

> 版本：无 DI 框架，纯约定 ｜ 调用点：全部 service/dao/common 包 ｜ 涉及文件：12+ ｜ 基线：main (6c93561)

## 用途定位
项目无 Spring/Wire 类 DI 容器，组件装配靠 Go 包级变量 + 构造函数约定实现。这是事实上的组件管理框架，新代码必须遵守。

## 初始化与配置
三种单例形态并存：

形态一：包级 var + init()（来源：`src/service/config_center_service.go:89-96`）：

```go
var configCenter *configCenterServiceImpl
func init() {
	configCenter = &configCenterServiceImpl{dao: dao.NewConfigCenterDao()}
	configCenter.stopChan = make(chan struct{})
}
func NewConfigCenterService() ConfigCenterService { return configCenter } // 返回单例
```

形态二：sync.Once 延迟初始化（来源：`src/service/event_service.go:19-46`）：

```go
var once = sync.Once{}
func NewEventService() *EventServiceImpl {
	once.Do(initEventStorageFactory) // 首次调用时初始化工厂
	...
}
```

形态三：每次 new（无共享状态，来源：`src/service/browser_service.go:47-53`）：

```go
func NewBrowserService() BrowserService {
	return &BrowserServiceImpl{
		ubd:        dao.NewUserBindDao(),
		cse:        cse.NewCse(),        // 内部是包级单例
		httpClient: https.Instance(),    // 进程级 client
	}
}
```

## 核心使用模式
Service 层标准结构（AGENTS.md 代码风格基线 + 存量证据）：

```go
// 接口 + 小写/大写实现类 + NewXxxService 构造函数
type XxxService interface { ... }
type xxxServiceImpl struct { dep1 *dao.XxxDao; dep2 cse.Cse }
var _ XxxService = &xxxServiceImpl{}   // 编译期接口断言
func NewXxxService() XxxService { ... }
```

DAO 装配：`XxxDao{ BaseInterface }`，构造时塞 `&BaseDao{EntityType: &db.Xxx{}}`（`src/dao/user.go:7-17`）。

## 封装层与扩展点
- 可替换性来自接口：Controller 只持有接口字段（`src/controllers/login_controller.go:19-24`），单测注入 fake（`src/service/browser_service_test.go` 注入 fake `cse.Cse`）。
- 扩展点：`var _ Iface = &Impl{}` 断言保证接口一致。

## 并发与线程模型
- 包级单例内共享状态必须自保护（configCenter 的 stopChan、alarmService 的 channel 串行化）。
- Controller 在 `Prepare()` 中调用 `NewXxxService()`，每请求重新装配，因此 service 实现必须是无状态或内部单例（证据：`src/controllers/login_controller.go:41-45`）。

## 错误处理与容错
init() 中初始化失败只能打日志（如 alarm SDK init，`src/service/alarm_service.go:63-72`），无法上抛。

## 约定与规范
- 命名：`XxxService` 接口、`xxxServiceImpl` 实现、`NewXxxService()` 构造（存量混有 `BrowserServiceImpl` 大写，新增一律小写）。
- 单例初始化：优先 `sync.Once`（AGENTS.md 质量基线）；包 init() 仅用于无依赖顺序问题的场景。

## 已知问题与反模式
- `NewXxxService` 命名与返回单例语义不符（configCenter/alarm 每次调用返回同一实例），调用方勿假设是新对象。
- `BrowserServiceImpl`、`TrafficStatsServiceImpl` 等实现类大写导出（历史遗留）。

## AI 编码指南
- 新增 Service：`XxxService` 接口 + `xxxServiceImpl` + `var _ XxxService = &xxxServiceImpl{}` + `NewXxxService()`；有共享状态用包级单例（sync.Once 保护），无状态每次 new（依据：上文三种形态）。
- **禁止**在 service 间直接引用对方实现类，依赖走接口字段注入（依据：`src/service/browser_service.go:55-59`）。
- Controller 依赖只能在 `Prepare()` 里通过 `NewXxxService()` 装配（依据：`src/controllers/login_controller.go:41`）。
