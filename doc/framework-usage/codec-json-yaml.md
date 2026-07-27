# 序列化使用指导（序列化）

> 版本：encoding/json（标准库）+ gopkg.in/yaml.v2 v2.4.0 + encoding/csv ｜ 调用点：json 43 / yaml 1 / csv 若干 ｜ 涉及文件：18 ｜ 基线：main @ 5e78a48

## 用途定位

- **encoding/json**：唯一的数据交换格式——HTTP 请求/响应体、审计日志、CSE 实例属性、Redis 对象、DB text 字段
- **yaml.v2**：仅用于外置 SQL 配置 `conf/sql.yaml`（`traffic_stats_service.go:16,87-100`）
- **encoding/csv**：流量统计导出（`traffic_stats_service.go:8`）、IMEI 白名单 CSV 导入（testsuit 配套）
- 无 Protobuf/MessagePack 等二进制协议（go.mod 中 protobuf 为间接依赖）

## 核心使用模式

```go
// 请求体解析 + 校验（controllers/controller.go:71-90）
err = json.Unmarshal(inputBody, param)
err = param.Validate()

// 响应输出（controller.go:137-146）
return json.NewEncoder(c.Ctx.ResponseWriter).Encode(v)

// 结构体 tag 规范（models/db/user.go:11-23）——json 与 orm tag 并存
Manufacturer string `json:"manufacturer" orm:"column(manufacturer)"`
CreatedAt    string `json:"-" orm:"column(created_at)"`   // json:"-" 不对外暴露
```

```go
// yaml 加载 SQL 配置（traffic_stats_service.go:80-100）
type SQLConfig struct {
	Queries map[string]struct {
		SQL    string   `yaml:"sql"`
		Params []string `yaml:"params"`
	} `yaml:"queries"`
}
err = yaml.Unmarshal(data, &config)
```

## 封装层与扩展点

- HTTP 层的 json 编解码已封装在 `BaseController.RequestBodyUnmarshalTo / writeHeaderAndJSON`（见 rpc-beego-web.md）。
- HTTP 客户端侧封装在 `Response.ResponseToStruct`（`common/https/builder.go:423-430`）。
- 双重序列化特例：审计日志 body 需 `json.Marshal` 两次（对端协议要求，`auditlog.go:97-107`）。

## 约定与规范

- 请求结构体放 `models/req/`、响应放 `models/resp/`，字段用 camelCase json tag；DB 实体放 `models/db/`，内部字段（created_at 等）打 `json:"-"`。
- DB 实体的时间字段统一 string + `time.DateTime` 格式，不用 time.Time（`user_service.go:106`）。
- 响应统一内嵌 `resp.BaseResponse{Code, Message}`。

## 已知问题与反模式

- `db_init.go:366-367` 对 DB 服务返回体做 `strings.Replace(#→")` 脏修复后再 Unmarshal——对端协议畸形的兼容代码，勿扩散此模式。
- `filter.go` 等少数文件 response 直接 `Write([]byte(...))` 明文而非 JSON（`filter.go:44`），仅限 429 特例。
- CSV 文件无 header 行，解析时所有行都当数据（AGENTS.md 已踩坑 #5）。

## AI 编码指南

- HTTP 出入参：走 `BaseController` / `https.Builder` 的既有封装，**禁止**在业务代码里裸 `json.NewEncoder(w)`。依据：上文「封装层与扩展点」。
- 新结构体同时服务 HTTP 与 DB 时：json tag camelCase + orm tag snake_case + 内部字段 `json:"-"`，对照 `models/db/user.go:10-24`。
- 需要外置 SQL/规则的新场景：沿用 yaml.v2 + `LoadXxxConfig(path)` 模式（模板 `traffic_stats_service.go:87-100`）；**禁止**为配置引入新的序列化格式（toml/protobuf 等）。
