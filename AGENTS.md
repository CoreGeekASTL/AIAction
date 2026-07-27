# GIDS (GlobalInstanceDeliverService)

云浏览器全局实例交付服务。Go + Beego v2，源码在 `src/`，go module 为 `GIDS`。

---

## 项目架构

```
GlobalInstanceDeliverService/
├── src/                    # Go 源码（go module: GIDS）
│   ├── main.go             # 启动入口
│   ├── models/db/          # 数据实体（orm标签+TableName+init注册）
│   ├── dao/                # DAO层（继承BaseInterface，EntityType设置）
│   ├── service/            # Service层（接口+小写实现类+包级变量+sync.Once）
│   ├── controllers/        # Beego Controller层
│   ├── routers/            # Beego路由注册
│   ├── common/             # 公共工具（retcode等）
│   ├── conf/               # 配置文件
│   ├── scheduler/          # 定时任务
│   ├── utils/              # 辅助函数
│   ├── stubs/              # 日志桩（空实现）
│   ├── test/               # Go单元测试
│   ├── data/               # LOCAL_MODE SQLite数据（gids.db）
│   └── db/                 # DDL初始化脚本
├── testsuit/               # E2E集成测试脚本（Python）
├── docs/                   # 设计文档（按需求编号组织）
│   └── 27.0/终端鉴权/      # SR→架构→实现设计→Story详设
├── build/                  # 构建脚本
└── .opencode/skills/       # 项目级skill定义
```

---

## 构建与运行

在 `src/` 目录下执行：

```powershell
go build -o gids.exe .
$env:LOCAL_MODE="true"; .\gids.exe          # PowerShell
# cmd: set LOCAL_MODE=true && gids.exe
# 或:  go run .   (go run 会每次重新编译)
```

- `LOCAL_MODE=true` 触发嵌入式 SQLite（纯 Go `modernc.org/sqlite`，无 CGO、无独立进程），
  数据落 `src/data/gids.db`（首次自建）。
- 内部 HTTP 监听 `127.0.0.1:9090`。
- 不设 `LOCAL_MODE` 时走原 GaussDB 链路（需 CSE 服务发现 + 真实 GaussDB），生产不受影响。

> 日志桩（`stubs/Go-chassis-extend/.../lager`）为空实现，启动时控制台无输出属正常；
> 用 `netstat -ano | findstr :9090` 确认监听。

---

## 开发流程（训战五步法）

本项目采用 superspec 全链路流程，从需求到交付五步闭环：

```mermaid
flowchart LR
    SR["① SEHarness<br/>需求分析+SR文档"] --> ARCH["② 系统架构设计"]
    ARCH --> SVC["③ 服务实现设计"]
    SVC --> STORY["④ Story详设"]
    STORY --> CODE["⑤ 代码生成+测试验证"]
```

每一步的产出物在 `docs/{需求编号}/{功能名}/` 下：

| 步骤 | Skill | 产出物 | 评审要求 |
| --- | --- | --- | --- |
| ① 需求分析 | se-harness | SR文档（`{编号}{功能名}SR文档.md`） | 用户确认需求完整性 |
| ② 系统架构设计 | sw-architecture-design | 系统实现架构设计文档 | 用户确认技术选型 |
| ③ 服务实现设计 | sw-service-architecture-design | 服务实现设计文档（主设计文档） | 用户确认接口+数据流 |
| ④ Story详设 | story-detail-design | Story-1~N详设文档 | 用户确认每个Story可独立开发 |
| ⑤ 代码生成+验证 | code-generation-quality-loop | 代码+测试+总报告 | 所有TC SUCCESS |

### ⑤ 代码生成内部流程

```
代码生成 → 质量检查(CodeCheck+臆造+风格) → DT测试 → 集成测试(testsuit) → 总报告 → 提交
```

详见 `.opencode/skills/code-generation-quality-loop/SKILL.md`。

---

## Skill体系

项目级skill在 `.opencode/skills/` 下（35个），关键skill：

| 分类 | Skill | 何时使用 |
| --- | --- | --- |
| **全链路** | superspec | 端到端开发：SR→架构→实现→Story→代码 |
| **需求** | se-harness | 需求分析+SR文档生成（任何创造性工作之前） |
| **设计** | sw-architecture-design | 系统架构设计 |
| | sw-service-architecture-design | 服务实现设计 |
| | story-detail-design | Story详设文档生成 |
| **代码** | code-generation-quality-loop | 代码生成+质量检查+DT测试+集成测试+总报告 |
| | code-quality-check | 代码质量检查（Go 17条/Java 13条规则） |
| **测试** | test-code-generation-loop | pytest测试用例生成+执行验证 |
| | test-driven-development | TDD：先写测试再写实现 |
| **验证** | verification-before-completion | 宣称完成前必须运行验证命令 |
| **调试** | systematic-debugging | bug/测试失败时的系统化调试 |
| **评审** | requesting-code-review | 完成功能后请求代码评审 |
| | receiving-code-review | 收到评审反馈后的技术严谨执行 |
| | chinese-code-review | 中文评审沟通（仅 /chinese-code-review 触发） |
| **规划** | writing-plans | 多步骤任务动手前先写计划 |
| | executing-plans | 在新会话中执行书面实现计划 |
| **并行** | dispatching-parallel-agents | 2+独立任务并行派发子代理 |
| | subagent-driven-development | 在当前会话中派子代理执行独立任务 |
| **Git** | finishing-a-development-branch | 实现完成后的合并/PR/清理决策 |
| | using-git-worktrees | 需要隔离工作区时使用git worktree |
| **文档** | document-selector | 分析需求描述选择最相关模块文档 |
| | doc-header | 为markdown文档添加YAML frontmatter头 |

**使用方式**：Skill是被动加载的——当任务描述匹配skill时，AI会自动调用；也可手动 `/skill-name` 触发。

---

## 代码风格规范

### Go代码质量基线

| 要求 | 标准 | 检查方法 |
| --- | --- | --- |
| 单例初始化 | `sync.Once` 保护 | 搜索 `once.Do` |
| context参数 | `context.TODO()` 不传nil | 搜索 `ContextDo` |
| 错误处理 | 所有error必须检查处理 | `go vet ./...` |
| 注释风格 | 中文：`// 函数名 功能说明` | 对比 `src/service/*.go` |
| 版权声明 | 文件首行 `// Copyright (c) Huawei Technologies Co., Ltd. 2026.` | 搜索 `Copyright` |
| 接口定义 | `type XXXService interface` + `xxxServiceImpl` | 对比存量Service |
| DAO继承 | `type XxxDao struct { BaseInterface }` | 对比存量DAO |
| 实体注册 | `orm:"pk;column(xxx)"` + `TableName()` + `init(){ orm.RegisterModel() }` | 对比存量Model |
| HTTP请求 | 优先 `https.NewRequest().WithRetry()` builder | 搜索 `NewRequest` |
| UUID生成 | `github.com/google/uuid.New()` | 搜索 `uuid.New` |
| 本地IP | `https.GetLocalIP(ethEnv, defaultEth)` | 搜索 `GetLocalIP` |

完整规则17条见 `.opencode/skills/code-quality-check/reference/codecheck-go.md`。

---

## 已踩坑记录

> 这些是实际开发中犯过的错误，新人务必注意。

| # | 坑 | 正确做法 | 来源 |
| --- | --- | --- | --- |
| 1 | IMEI写成16位数字 | IMEI/IMSI必须是15位纯数字，正则 `^[0-9]{15}$` | 终端鉴权TC_001 |
| 2 | login拒绝返回code=401 | login路径拒绝返回 `retcode.ClientFailed(-2)`，event路径拒绝返回 `retcode.AuthFailed(401)`，必须区分 | 终端鉴权TC_003 |
| 3 | Python docstring嵌套中文引号 | `"""msg含"xxx"""`" 语法错误，改为去掉嵌套引号 | 终端鉴权TC_002 |
| 4 | f-string中嵌套dict字面量 | `f"{'key': '{val}'}"` 格式错误，改为简单字符串拼接 | 终端鉴权TC_002 |
| 5 | CSV带header行 | 设计要求纯数据CSV（无header），parseCSV所有行当数据读 | 终端鉴权Story-2 |
| 6 | Beego ORM给IMSI加pk tag | Beego不支持复合pk，IMEI设pk，IMSI不加pk，用DDL UNIQUE INDEX兜底 | 终端鉴权Story-1 |
| 7 | 集成测试不主动运行 | DT测试通过后必须主动运行testsuit（若存在），不可等用户提醒 | 终端鉴权整个流程 |

---

## 常用命令

```powershell
# 构建
cd src; go build -o gids.exe .

# 本地运行（SQLite模式）
$env:LOCAL_MODE="true"; .\gids.exe

# 验证监听
netstat -ano | findstr :9090

# Go单元测试
cd src; go test -v ./service/... ./dao/... ./controllers/...

# Go vet（检查Unhandled error等）
cd src; go vet ./...

# E2E集成测试（testsuit存在时必做）
python testsuit\TC_SBG_Func_GIDS_Auth_001.py   # 白名单导入导出
python testsuit\TC_SBG_Func_GIDS_Auth_002.py   # 鉴权正常/异常/逃生态
python testsuit\TC_SBG_Func_GIDS_Auth_003.py   # 登录链路鉴权
python testsuit\TC_SBG_Func_GIDS_Auth_004.py   # 事件链路鉴权
python testsuit\TC_SBG_Func_GIDS_Auth_005.py   # 边界与缓存覆盖

# Git分支管理
git checkout feature/27.0-终端鉴权
```
