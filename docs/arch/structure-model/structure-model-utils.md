# utils 模块结构文档

> 生成时间：2026-08-11
> 所属仓：AIAction（GIDS）
> 模块路径：src/utils

## 模块职责

辅助工具层：fileutil 提供文件操作工具，flagutil 提供命令行 flag 解析（被 main.go 使用），monitorutil 提供监控工具，response 提供统一响应组装（Success/失败响应，基于 retcode 与 models/resp）。

## 子模块关系图

子包间无 import 依赖关系。

## 子模块说明

| 子模块 | 路径 | 职责 | 主要依赖 | 被依赖 |
| --- | --- | --- | --- | --- |
| fileutil | src/utils/fileutil | 文件操作工具。证据：src/utils/fileutil/fileutil.go | - | - |
| flagutil | src/utils/flagutil | 命令行 flag 解析。证据：src/main.go 引用 | - | - |
| monitorutil | src/utils/monitorutil | 监控工具。证据：src/service/monitor_service.go 引用 | - | - |
| response | src/utils/response | 统一响应组装。证据：src/utils/response/response_util.go | - | - |
