# 测试框架使用指导（测试框架）

> 版本：testify v1.8.4 + goconvey v1.8.1 + gomockit v1.1.0（stub）+ Python testsuit ｜ 调用点：17 个 Go 测试文件 + 5 个 E2E 脚本 ｜ 涉及文件：src 各模块 `_test.go` + `testsuit/` ｜ 基线：main @ 5e78a48

## 用途定位

双层测试体系：
- **Go 单元测试**（`src/**/*_test.go`）：testify assert 为主、goconvey 为辅，DAO 层经 build tag 隔离真实 DB
- **Python E2E 集成测试**（`testsuit/TC_*.py`）：对运行中的服务（LOCAL_MODE, 127.0.0.1:9090）发真实 HTTP 验证全链路

## 初始化与配置

- 运行 Go UT：`cd src && go test -v ./service/... ./dao/... ./controllers/...`
- 静态检查：`go vet ./...`
- 运行 E2E：先 `LOCAL_MODE=true` 起服务，再 `python testsuit/TC_SBG_Func_GIDS_Auth_00X.py`（AGENTS.md「常用命令」）
- DB 隔离：`dao/db_init.go:1` 有 `//go:build !test` tag——test 构建时整个 GaussDB 初始化被剔除，由 `donothing_base_dao.go` 与测试内 stub 接管
- 测试替身：`service/traffic_stats_service_mock.go`（手写 mock）、`service/master_election_service_stub_test.go`（stub）、`gomockit`（内部 mock 工具，stub）

## 核心使用模式

```go
// testify 风格（主流，来源：src/utils/monitorutil/time_util_test.go 等 17 文件）
import "github.com/stretchr/testify/assert"
assert.Equal(t, expected, actual)
```

```go
// goconvey 风格（少数，来源：src/db/driver/driver_test.go:40-73、src/dao/base_dao_test.go:245）
convey.Convey("test driver", t, func() {
	convey.Convey("test select", func() { ... })
})
```

```go
// 串行流程补充工具（src/test/util/utils.go:6）——goconvey 场景内串行步骤标记
```

```python
# E2E 风格（testsuit/TC_SBG_Func_GIDS_Auth_001.py）：requests 调 9090 端口，断言 code/数据
```

## 约定与规范

- 新 UT 默认用 **testify**（17:2 的存量比例）；goconvey 仅出现在 db/driver 与 base_dao 两处历史文件。
- 测试文件放被测包同目录（`xxx_test.go`）。
- 涉及 DB 的 DAO/Service 测试依赖 `!test` tag 隔离 + 手写 mock，不连真实库。
- AGENTS.md 流程要求：DT 测试通过后**必须主动运行 testsuit**（已踩坑 #7）。

## 已知问题与反模式

- 双断言库并存（testify/goconvey），新代码统一 testify。
- `gomockit` 为内部工具且本地 stub，实际 mock 多靠手写（`*_mock.go`/`*_stub_test.go`）。
- testsuit 依赖服务手工启动，无自动化拉起/回收。

## AI 编码指南

- 新单元测试：同目录 `xxx_test.go` + testify `assert/require`；**禁止**新增 goconvey 用法。依据：上文「约定与规范」。
- 测 DAO/Service：利用 `//go:build !test` 隔离的 dao 层 + 手写 mock 结构体实现对应接口（参考 `traffic_stats_service_mock.go`、`master_election_service_stub_test.go`）。依据：上文「初始化与配置」。
- 功能交付前：`go build` + `go vet ./...` + `go test ./...` + 起 LOCAL_MODE 服务跑对应 `testsuit/TC_*.py`，全部通过才可宣称完成。依据：AGENTS.md 开发流程与已踩坑 #7。
