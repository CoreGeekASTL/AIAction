# 测试设计要求

> 本文件为 story-detail-design skill 的测试设计参考。SKILL.md 中通过引用链接加载本文件。

## Mock方案

| 外部依赖 | Mock方案 | 工具 |
| --- | --- | --- |
| **数据库** | sqlmock模拟SQL执行 | `github.com/DATA-DOG/go-sqlmock` |
| **HTTP服务** | Mock HTTP Server | `httptest.Server` |
| **SDK接口** | Mock接口结构体 | 自定义Mock |

## UT测试（单元测试）

- **DAO层**：测试每个方法，使用 sqlmock 模拟数据库
- **Service层**：测试核心逻辑，Mock DAO依赖
- 测试文件命名：`*_test.go`
- 使用 `stretchr/testify/assert` 断言

## DT测试（集成测试）

- 使用 Mock 模拟外部依赖
- 测试完整业务流程
- 测试文件命名：`*_dt_test.go`

## 测试覆盖率

| 模块 | 覆盖率要求 |
| --- | --- |
| DAO层 | >= 80% |
| Service层 | >= 85% |
