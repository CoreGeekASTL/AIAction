# 序列化框架使用指导（序列化/编解码）

## 用途定位
- `encoding/json`：全部 HTTP 请求/响应体、DB text 字段、配置内容、监控模板（monitor.json）的编解码。
- `gopkg.in/yaml.v2`：仅用于话统 SQL 配置文件 sql.yaml（`src/service/traffic_stats_service.go:16`）。
- 无 Protobuf/MsgPack；Redis 存取走 `encoding.BinaryMarshaler` 接口由业务对象自实现（见 [storage-redis.md](storage-redis.md)）。


## 使用模式

HTTP 请求体解析（来源：`src/controllers/controller.go:71-90`）：

```go
inputBody, _ := io.ReadAll(c.Body())
err = json.Unmarshal(inputBody, param) // param 实现 req.IRequest
err = param.Validate()
```

HTTP 响应写出（来源：`src/controllers/controller.go:137-146`）：

```go
c.AddHeader("Content-Type", "application/json")
c.Ctx.ResponseWriter.WriteHeader(status)
json.NewEncoder(c.Ctx.ResponseWriter).Encode(v)
```

DB text 字段存取（来源：`src/controllers/management_controller.go:103,164`）：

```go
json.Unmarshal([]byte(cfg.Content), bc) // 读
contentB, _ := json.Marshal(bc)         // 写
```

yaml 配置加载（来源：`src/service/traffic_stats_service.go:80-95`）：

```go
type SQLConfig struct {
	Queries map[string]struct {
		SQL    string   `yaml:"sql"`
		Params []string `yaml:"params"`
	} `yaml:"queries"`
}
data, _ := os.ReadFile(configPath)
yaml.Unmarshal(data, &cfg)
```
