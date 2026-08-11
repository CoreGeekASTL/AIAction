# encoding/json + yaml.v2 使用现状（序列化/编解码）

## 用途定位
JSON 序列化全仓统一用标准库 `encoding/json`（43 处调用点，无第三方 JSON 库）：Controller 请求/响应体（`src/controllers/controller.go` 的 `RequestBodyUnmarshalTo`、`writeHeaderAndJSON`）、https builder 的 body 构建与 `ResponseToStruct`（`src/common/https/builder.go`）、DB text 字段内容解析（`src/controllers/management_controller.go` 的 config content）、监控配置 monitor.json 解析（`src/service/monitor_service.go`）、审计事件序列化（`src/common/logger/auditlog.go`）。模型结构体同时带 `json` 与 `orm` tag（`src/models/db/*.go`）。

YAML 仅一处：`gopkg.in/yaml.v2` 解析监控指标 SQL 模板 `sql.yaml`（`src/service/traffic_stats_service.go` 的 `yaml.Unmarshal`）。

## 使用模式

请求解析骨架（读 body → Unmarshal → Validate）：

```go
// 来源：src/controllers/controller.go
func (c *BaseController) RequestBodyUnmarshalTo(param req.IRequest) error {
	inputBody, err := io.ReadAll(c.Body())
	// ...
	err = json.Unmarshal(inputBody, param)
	// ...
	err = param.Validate()
	// ...
}
```

响应写出骨架：

```go
// 来源：src/controllers/controller.go
c.AddHeader("Content-Type", contentType)
c.Ctx.ResponseWriter.WriteHeader(status)
return json.NewEncoder(c.Ctx.ResponseWriter).Encode(v)
```
