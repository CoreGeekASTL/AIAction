# Go 测试框架使用指导（测试框架）

> 版本：testing（stdlib）+ testify v1.8.4（replace 到 v1.8.2）+ goconvey v1.8.1 + gomockit v1.1.0（stub）｜ 调用点：17 测试文件 ｜ 涉及文件：17 ｜ 基线：main (6c93561)

## 用途定位
Go 单元测试（DT）。主流打法是 testing + testify/assert 表驱动；goconvey 仅 3 个早期文件；gomockit 用于 mock（经 stub）。E2E 集成测试为 Python testsuit/（另述，不在本仓 Go 框架范围）。

## 初始化与配置
- 运行：`cd src; go test -v ./service/... ./dao/... ./controllers/...`（AGENTS.md）。
- 构建标签：`dao/db_init.go` 带 `//go:build !test`，测试编译期可用 `test` tag 替换 DB 初始化（`src/dao/db_init.go:1`）。
- 测试替身：`dao.DoNothingBase` 替换 DAO（`src/dao/donothing_base_dao.go`）；`https.HTTPDoer`、`cse.Cse`、`service.TrafficStatsService`（`src/service/traffic_stats_service_mock.go`）接口注入 fake。

## 核心使用模式

表驱动 + testify（来源：`src/service/browser_service_test.go:42-60`）：

```go
func TestBrowserServiceImplGetAllReadyServiceInstances(t *testing.T) {
	tests := []struct {
		name string
		cse  cse.Cse      // 接口注入 fake
		want []browsergateway.ServiceInstance
	}{
		{name: "no instances", cse: &noInstanceCse{}, want: nil},
		{name: "one ready", cse: &oneReadyInstanceCse{}, want: ...},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			...
			assert.Equal(t, tt.want, got)
		})
	}
}
```

HTTP 过滤器测试（来源：`src/controllers/filter_test.go:17-40`）：`httptest` + 构造 `beecontext.Context` + `overloadcontroller` 策略 JSON 初始化。

goconvey 串行辅助（来源：`src/test/util/utils.go`）：

```go
util.It("step desc", func() { ... }) // 补充 goconvey 缺失的串行流程表达
```

测试文件头固定中文模板：测试用例描述/预置条件/操作步骤/预期结果/修改历史（`src/service/browser_service_test.go:22-37`）。

## 封装层与扩展点
- 可测性来自接口化：service/dao/cse/https client 全部接口注入（见 [di-singleton.md](di-singleton.md)）。
- `redis.InitForTest` 提供测试清理钩子（`src/common/storage/redis/redis.go:44`）。

## 约定与规范
- 测试文件与源码同包同目录，命名 `xxx_test.go`。
- 断言用 `testify/assert`；新测试**不再引入 goconvey**（存量 3 个文件为历史）。
- 用例头部必须写中文描述模板（描述/预置/步骤/预期/修改历史）。

## 已知问题与反模式
- `DoNothingBase` 多个方法 `panic("implement me")`（`src/dao/donothing_base_dao.go:15-22`），用到未实现方法会直接 panic。
- testify 版本被 replace 固定到 v1.8.2（`src/go.mod:79`），升级需注意。

## AI 编码指南
- 新增单测：testing + testify/assert 表驱动；依赖经接口注入 fake；用例头写中文模板（依据：`src/service/browser_service_test.go`）。
- 涉及 DB 的测试：注入 `dao.DoNothingBase` 或自定义 BaseInterface fake，**禁止**在单测中连真实数据库（依据：`src/dao/donothing_base_dao.go`）。
- DT 通过后必须主动运行 testsuit/ 下 Python E2E（若存在），不可等提醒（AGENTS.md 坑记录 #7）。
