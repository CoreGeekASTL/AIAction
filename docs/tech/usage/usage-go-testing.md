# Go testing + testify + goconvey 使用现状（测试框架）

## 用途定位
单测用标准 `testing` 包 + `github.com/stretchr/testify`（assert，17 个测试文件，集中在 controllers/utils/dao/common 各包 `*_test.go`）；`github.com/smartystreets/goconvey` 用于老测试（`src/dao/base_dao_test.go`、`src/db/driver/driver_test.go`），新测试以 testify 为主。`src/test/util/utils.go` 提供 goconvey 串行流程辅助 `It`。HTTP 层测试通过 `beego.NewHttpServerWithCfg` + httptest 风格构造；DAO 测试用 SQLite driver 与 `orm.DoNothingOrm`（`src/dao/donothing_base_dao.go`）做隔离。E2E 集成测试不在 Go 侧，在 `testsuit/`（Python 脚本）。

## 使用模式

```go
// 来源：src/controllers/management_controller_test.go（ testify 风格）
import "github.com/stretchr/testify/assert"
```

```go
// 来源：src/dao/base_dao_test.go（goconvey 风格）
import convey "github.com/smartystreets/goconvey/convey"
```
