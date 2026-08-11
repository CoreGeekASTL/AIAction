# 测试联通

> 功能域：测试联通　接口数：1　所属 server：外部
> 子文档 of [README.md](README.md)

## 1. 定位

连通性自测接口，仅测试使用，验证 externalServer 路由与响应链路可用。

## 2. 接口清单

| 接口名 | 作用 | 所在文件 | 方法/路径 |
|---|---|---|---|
| GetData | 返回固定成功响应（仅测试使用） | controllers/test_controller.go | GET /test/v1/get |

## 3. 数据结构说明

- **GetData**
  - 请求：无参数
  - 响应 `resp.DataResponse`（models/resp/response_entity.go）：code=成功、msg="test success"、data=true
