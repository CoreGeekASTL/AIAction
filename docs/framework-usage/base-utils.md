# 自研公共工具库使用指导（基础库）

> 版本：— ｜ 涉及包：src/utils/*、src/common/constants/* ｜ 基线：main (6c93561)

## 用途定位
自研公共能力是"事实上的框架"，新代码优先复用而非新造。

| 包 | 能力 | 关键文件 |
| --- | --- | --- |
| `common/constants/retcode` | HTTP 返回码常量（Success=200/InternalFailed=-1/ClientFailed=-2/AuthFailed=401） | `src/common/constants/retcode/retcode.go` |
| `common/constants` | 全局常量（ServiceName、CleanupMonths、EnableHTTP 等） | `src/common/constants/base.go` |
| `utils/flagutil` | 结构体 flag tag 递归注册命令行参数 | `src/utils/flagutil/flags.go:13` |
| `utils/response` | 响应辅助 | `src/utils/response/response_util.go` |
| `utils/fileutil` | 文件/zip 操作（事件日志滚动删除依赖） | `src/utils/fileutil/fileutil.go`、`zip_util.go` |
| `utils/monitorutil` | 话统时间窗口计算（GetLastFiveMinuteWindow） | `src/utils/monitorutil/time_util.go` |
| `test/util` | 测试辅助 It()（goconvey 串行流程） | `src/test/util/utils.go` |
| `github.com/google/uuid` v1.6.0 | token 生成（`uuid.New().String()`） | `src/service/browser_service.go:204` |

## 核心使用模式

命令行参数（来源：`src/main.go:52`、`src/common/conf/config.go:52-57`）：

```go
type Config struct {
	Logger LoggerConfig `flag:"log"`   // → --log-file / --log-level / --log-event
	Redis  RedisConfig  `flag:"redis"` // → --redis-endpoint / --redis-db
}
c := conf.Instance()
flagutil.Parse(c) // 反射注册并 flag.Parse()
```

UUID token（来源：`src/service/browser_service.go:203-206`）：

```go
uid := uuid.New()
u.Token = uid.String()
```

## 约定与规范
- 返回码一律引用 retcode 常量，**禁止**散落数字字面量（12 个文件已遵循）。
- 新常量放 `common/constants`；新工具函数放 `utils/<领域>util/` 包（命名领域+util）。
- UUID 生成只用 `github.com/google/uuid.New()`（AGENTS.md 质量基线）。

## 已知问题与反模式
- `utils/response` 与 `BaseController.OK/Failed` 能力部分重叠，新响应优先用 BaseController 方法。

## AI 编码指南
- 新增返回码/常量：先查 `retcode` 与 `constants` 是否已有，没有再新增；**禁止**接口返回硬编码数字（依据：`src/common/constants/retcode/retcode.go`）。
- 新增工具函数：放 `utils/<领域>util/`，配套 `_test.go`（依据：存量 utils 包均有测试）。
- ID/token 生成：`uuid.New().String()`（依据：`src/service/browser_service.go:204`）。
