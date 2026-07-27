# 基础库使用指导（基础库）

> 版本：google/uuid v1.6.0 + 自研 utils ｜ 调用点：~10 ｜ 涉及文件：5 ｜ 基线：main @ 5e78a48

## 用途定位

- **google/uuid**：全局唯一 ID 生成（会话/记录标识）
- **`utils/fileutil`**：文件与 zip 操作（事件日志转储用）：`CopyFile/ZipFile` + 权限常量（`PermissionForLogDir/PermissionForLogFile/PermissionForZipFile`）
- **`utils/monitorutil`**：话统时间窗口工具 `GetLastFiveMinuteWindow`（`monitor_service.go:190`）
- **`utils/response`**：响应辅助 `utils/response/response_util.go`
- **`utils/flagutil`**：命令行解析（见 config md）

## 核心使用模式

```go
// 唯一 ID（来源：src/service/browser_service.go:204）
uid := uuid.New()

// 文件转储（来源：src/common/event/local_storage.go:168-192）
util.CopyFile(logTempFile, storage.logFile)
util.ZipFile(logTempFile, backFile)
os.Chmod(backFile, util.PermissionForZipFile)
os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, util.PermissionForLogFile)
```

## 约定与规范

- ID 生成统一 `uuid.New()`（AGENTS.md 基线，禁止自造随机串）。
- 新增文件落盘代码：权限一律用 `fileutil` 的常量，不要手写 `0o755/0o644` 字面量（例外：`db_local_sqlite.go:131` 用了 `0o755`，属不一致个例）。
- 本机 IP 获取统一 `https.GetLocalIP(ethEnv, defaultEth)`（AGENTS.md 基线，`common/https/https_server.go:70-85`），失败回退 127.0.0.1。
- 错误码统一 `common/constants/retcode`，常量统一 `common/constants/base.go`。

## 已知问题与反模式

- `common/event/local_storage.go:9,226` 使用已废弃的 `io/ioutil`。
- `db_local_sqlite.go:131` 硬编码 `0o755` 未复用 fileutil 常量。

## AI 编码指南

- 生成 ID：`uuid.New()`；获取本机 IP：`https.GetLocalIP(...)`；文件权限：`fileutil.PermissionForXxx`。三者均有 AGENTS.md 基线约束，**禁止**自造实现。
- 新工具函数按领域放 `utils/<领域>util/` 目录（fileutil/monitorutil 同款命名），不放 `common/`。依据：现有目录结构。
