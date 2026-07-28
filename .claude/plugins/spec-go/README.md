# Spec-go

Spec-go 是一套面向存量代码仓的规格化分析 skill 体系，由一组可组合的 spec skill 和一段确保 agent 主动使用它们的 bootstrap 指令构成。

## 快速开始

给你的 coding agent 装上 Spec-go：[OpenCode](#opencode)、[Claude Code](#claude-code)。

## 工作原理

从你打开一个存量代码仓那一刻起，Spec-go 就开始工作。一旦它发现你在做"梳理结构、盘点接口、归纳功能、梳理框架、审需求、写 story"这类分析任务，它不会上来就乱翻代码或凭印象下结论，而是先加载对应的 spec skill，按统一格式产出可消费的文档资产——mermaid 结构图、接口清单、feature 文档、框架使用指导、逻辑审核 HTML、story 设计——沉淀到 `docs/` 下，供人审、供 AI 后续编码消费。

每个 skill 都自动触发，你无需做任何特殊操作。你的 coding agent 只是"拥有 spec 能力"。

## 安装

各 harness 安装方式不同。若同时使用多个，分别安装。

### OpenCode

OpenCode 用自己的插件机制；即使你在别的 harness 装过，OpenCode 仍需单独安装一次。

- 在 `~/.config/opencode/opencode.json` 的 `plugin` 数组加入本插件目录：

```json
{
  "plugin": ["~/.config/opencode/plugins/spec-go"]
}
```

- 重启 OpenCode。配置项须是**包目录**（不是 `.js` 文件路径）——OpenCode 会读 `package.json` 的 `main`（`opencode.js`）作为入口，`config` hook 注册 skills 目录，`messages.transform` hook 注入 bootstrap。

### Claude Code

- 注册本地 marketplace：

```
/plugin marketplace add ~/.config/opencode/plugins/spec-go
```

- 从该 marketplace 安装插件：

```
/plugin install spec-go@spec-go
```

插件安装后，Claude Code 读 `hooks/hooks.json` 在 SessionStart 跑 `session-start` 脚本注入 bootstrap，并发现 `skills/` 目录下的 6 个 skill。

## 基本工作流

1. **spec-structure-analyze** - 代码仓结构摸底：第一层目录与包间依赖，产出 mermaid 依赖图 + 模块说明表。
2. **spec-interface-analyze** - 对外接口盘点：HTTP 路由 / RPC service / 消息订阅 handler / IDL 契约，主文档 README + 功能域子文档，归档 `docs/interface/`。
3. **spec-feature-analyze** - 接口归纳为功能域：同一业务的多个接口归为一个 feature，产出与 `docs/story/` 同构的 feature 文档。
4. **spec-framework-usage-analyze** - 框架使用模式：识别基础框架及调用点分布，每框架一篇使用指导，归档 `docs/framework-usage/`。
5. **spec-logic-audit** - 需求逻辑审核：把 Spec 翻译成多彩建模图，逻辑断裂点自然暴露，ask-human 补齐，产出 HTML 供编码人员审核。
6. **spec-story-design** - 需求到 story：接收 SR/特性设计，结合存量架构资产，产出与 `docs/story/` 同构的新功能 story 设计文档。

**agent 在任何分析任务前都会检查是否有相关 skill。** 这是必须遵循的流程，不是建议。

## 内含什么

### 技能库

**结构**
- **spec-structure-analyze** - 代码仓结构文档（mermaid 依赖图 + 模块说明表）

**接口**
- **spec-interface-analyze** - 对外接口盘点（HTTP/RPC/消息订阅/IDL，主文档+子文档）
- **spec-feature-analyze** - 接口按业务功能归纳为 feature 文档

**框架**
- **spec-framework-usage-analyze** - 基础框架使用模式分析，每框架一篇使用指导

**需求**
- **spec-logic-audit** - 多彩建模 + 业务逻辑完备性审核
- **spec-story-design** - 按存量模板产出新功能 story 设计

## 理念

- **设计先行于编码** - 先沉淀 spec/设计文档，再动代码
- **多彩建模** - 把文档翻译成事件-角色-实体的因果图，让逻辑断点无处遁形
- **证据先行** - 基于代码仓实际接口/调用点产出，不臆造
- **统一产出** - 所有 skill 按既定格式归档到 `docs/`，人与 AI 共同消费

## 更新

Spec-go 的 skill 内容在 `skills/*/SKILL.md`。改了任一 skill 后，重新生成 bootstrap：

```bash
cd ~/.config/opencode/plugins/spec-go
bun scripts/generate-bootstrap.mjs
```

然后重启 coding agent 使新 bootstrap 生效。

## 许可证

MIT License。
