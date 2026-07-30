# 连通性测试

> 功能域：test-echo　接口数：1（仅测试使用）　所属 server：外部(HTTPS)
> 子文档 of [README.md](README.md)

## 1. 定位

外部链路连通性验证接口，仅测试/冒烟使用，无业务逻辑。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| GetData | 连通性测试（仅测试使用），固定返回成功 | src/controllers/test_controller.go | GET /test/v1/get |

## 3. 数据结构说明

- **GetData**
  - 请求：无参数
  - 响应 `resp.DataResponse`（src/models/resp/response_entity.go）：`BaseResponse{code:0,msg:"test success"}` + `data: true`

## 4. 风险与注意点

- **生产暴露**：src/routers/beego_router.go:24（测试接口注册在外部 HTTPS server 的生产路由表中，未按环境隔离）
