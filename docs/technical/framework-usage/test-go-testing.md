# Go 测试框架使用指导（测试框架）

## 用途定位
Go 单元测试（DT）。主流打法是 testing + testify/assert 表驱动；goconvey 仅 3 个早期文件；gomockit 用于 mock（经 stub）。E2E 集成测试为 Python testsuit/（另述，不在本仓 Go 框架范围）。


## 使用模式

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
