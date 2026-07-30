# 序列化框架使用指导（序列化/编解码）

> 版本：encoding/json（stdlib）+ gopkg.in/yaml.v2 v2.4.0 ｜ 调用点：json 43 / yaml 1 ｜ 涉及文件：19 ｜ 基线：main (6c93561)

## 用途定位
- `encoding/json`：全部 HTTP 请求/响应体、DB text 字段、配置内容、监控模板（monitor.json）的编解码。
- `gopkg.in/yaml.v2`：仅用于话统 SQL 配置文件 sql.yaml（`src/service/traffic_stats_service.go:16`）。
- 无 Protobuf/MsgPack；Redis 存取走 `encoding.BinaryMarshaler` 接口由业务对象自实现（见 [storage-redis.md](storage-redis.md)）。

## 初始化与配置
无初始化。

## 核心使用模式

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

## 封装层与扩展点
无独立封装层；约定通过结构体 tag 表达：
- 请求/响应结构体放 `src/models/req|resp/`，字段带 `json:"camelCase"` tag。
- DB 实体（`src/models/db/`）同时带 `orm:"column(xxx)"` 与 `json` tag；不希望暴露给前端的字段用 `json:"-"`（`src/models/db/user.go:22` 的 CreatedAt）。

## 并发与线程模型
json 包并发安全；无共享状态。

## 错误处理与容错
- Unmarshal 失败统一转 `errors.New(文案)` 并 Errorf 日志（`src/controllers/controller.go:78-88`）。
- `json.RawMessage` 用于异构批量数据延迟解析（`src/service/traffic_stats_service.go` BatchInsertStats）。

## 约定与规范
- 对外 JSON 字段一律 camelCase（证据：`src/models/db/user.go:12-21`）。
- yaml 仅允许配置文件场景，接口报文禁止 yaml。

## 已知问题与反模式
- DB Service 返回体需预处理（`#`→`"`、去首尾引号）再 Unmarshal（`src/dao/db_init.go:366-370`），属对端协议妥协，非通用模式。
- yaml.v2 与 indirect yaml.v3 并存于 go.sum，代码实际只用 v2。

## AI 编码指南
- 新增请求结构体：放 `src/models/req/`，camelCase json tag，实现 `Validate()`；响应对应放 `src/models/resp/`（依据：`src/controllers/controller.go:71-90`）。
- DB 实体新增字段必须同时标注 `orm:"column(snake_case)"` 与 `json:"camelCase"`；内部字段用 `json:"-"`（依据：`src/models/db/user.go:11-23`）。
- 接口报文只用 encoding/json；**禁止**引入新序列化框架（依据：go.mod 依赖清单与存量一致性）。
