# 自研单例/组件模式使用指导（依赖注入/组件管理）

## 用途定位
项目无 Spring/Wire 类 DI 容器，组件装配靠 Go 包级变量 + 构造函数约定实现。这是事实上的组件管理框架，新代码必须遵守。


## 使用模式
Service 层标准结构（AGENTS.md 代码风格基线 + 存量证据）：

```go
// 接口 + 小写/大写实现类 + NewXxxService 构造函数
type XxxService interface { ... }
type xxxServiceImpl struct { dep1 *dao.XxxDao; dep2 cse.Cse }
var _ XxxService = &xxxServiceImpl{}   // 编译期接口断言
func NewXxxService() XxxService { ... }
```

DAO 装配：`XxxDao{ BaseInterface }`，构造时塞 `&BaseDao{EntityType: &db.Xxx{}}`（`src/dao/user.go:7-17`）。
