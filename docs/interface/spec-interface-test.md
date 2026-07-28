# 测试接口

> 功能域：测试接口　接口数：1（仅测试）　所属 server：外部(HTTPS)
> 子文档 of [README.md](README.md)

## 1. 定位

externalServer（HTTPS）暴露的连通性测试接口，仅用于健康检查，非业务功能。TestController 承载。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| GetData | 测试连通性 | controllers/test_controller.go | GET /test/v1/get |

## 3. 数据结构说明

- **GetData**
  - 请求：无参数
  - 响应 `resp.DataResponse{BaseResponse: {Code: retcode.Success, Message: "test success"}, Data: true}`

## 4. 风险与注意点

- **测试接口暴露在生产 externalServer**：controllers/test_controller.go:20（`/test/v1/get` 注册在 RegisterExternalRouter，生产 HTTPS 可访问，应下线或加白名单）
