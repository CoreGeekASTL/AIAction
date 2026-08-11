# AIAction 分支与变更规范（branch-guidelines）

| 元信息 | 值 |
|--------|-----|
| 适用仓库 | AIAction（/home/lele/project/work/csp-ysj/AIAction，origin: github.com/CoreGeekASTL/AIAction） |
| 默认分支 | main（证据：远端仅存 remotes/origin/main 主线；origin/HEAD symbolic-ref 未设置） |
| 更新日期 | 2026-08-11 |
| Skill | qual-branch-guidelines-analyze |
| 运行模式 | 起草模式 |
| 采样口径 | 分支：全量 6 条本地 + 2 条远端唯一（去重后 6 条）；commit：main 分支全量 34 条（不足 100 条，取全部） |

> 本规范为 guidelines 形态资产：指导性规范（"应该"遵守），违反出报告提示改进，不做 CI 拦截。
> 本文档为活文档，同名覆盖更新；各章「约定」小节均为**建议，待团队确认**，与「现状」事实分节呈现，不得混淆。

## 一、分支模型

### 现状

**分支命名形态分布**（证据：`git branch -a` 全量，去重后 6 条）

| 命名形态 | 示例分支名（真实） | 样本数 | 占比 |
| --- | --- | --- | --- |
| personal/{用户}/{简述}，英文小写 | personal/houle/test、personal/houle/test2、personal/houle/test3 | 3 | 50% |
| 无前缀，snake_case 英文 | new_skill_test | 1 | 17% |
| ready/{版本号}-{中文功能名} | ready/27.0-终端鉴权 | 1 | 17% |
| 默认分支 | main | 1 | 17% |

**长期分支与生命周期**：无 develop / release / 集成分支等常驻分支，仅默认分支 main；已合入 main 未删除分支 0 条（证据：`git branch --merged main` 仅 main 自身）；personal/houle/test3 与 ready/27.0-终端鉴权 已推送到远端且未合入 main。

**tag 与发布形态**：无 tag（证据：`git tag` 输出为空），无发布版本管理痕迹。

### 约定（建议，待团队确认）

- 分支命名统一 `{type}/{简述}`，type ∈ feature / fix / hotfix / docs / chore；个人试验分支沿用 `personal/{用户}/{简述}` 前缀并限本地或短期保留。
- 需求交付分支沿用 `ready/{版本号}-{功能名}` 形态（如 ready/27.0-终端鉴权），作为待发分支语义。
- 分支合入 main 后即删除（含远端）；长期分支仅保留 main。
- 发布节点打 tag，建议语义化版本 v{X}.{Y}.{Z}。
- 现状分歧说明：当前无前缀分支（new_skill_test）与前缀分支并存，建议收敛为全部带 type 前缀。

## 二、commit message 规范

### 现状

**类型前缀分布**（证据：main 全量 34 条 commit 采样）

| 类型前缀 | 样本数 | 占比 | 证据（commit 短 hash） |
| --- | --- | --- | --- |
| conventional commits（feat:/docs:/refactor:/chore:/revert:） | 19 | 55.9% | ba05b1c（feat）、750f376（docs）、f29a814（refactor）、a93135c（chore）、b5fa785（revert） |
| 无前缀（中文直述） | 15 | 44.1% | c239285、656732c、9925bc7 |

**语言分布**：中文 33 条（97.1%）/ 英文 1 条（2.9%，证据：66c80dd "bugfix for repo-structure-doc"）/ 混合 0 条

**subject 长度分布**：≤50 字符 17 条 / 51~72 字符 10 条 / >72 字符 7 条（证据：ba05b1c "feat: spec skill 体系增强（issues #1/#2/#3/#5/#6/#7）+ 迁移至 spec-go 插件"）

**引用形态**：工单/issue 引用出现率 2.9%（1/34，证据：ba05b1c 含 "issues #1/#2/#3/#5/#6/#7"）；MR/PR 编号后缀（"(#N)"）出现率 0%

### 约定（建议，待团队确认）

- 采用 conventional commits 格式 `{type}: {subject}`，type ∈ feat / fix / docs / style / refactor / perf / test / build / ci / chore / revert；现状已有 55.9% 命中，建议收敛为 100%。
- subject 使用中文、≤72 字符、一句话说明变更意图；禁止无上下文 message（如"修改""修复"类无信息量表述，现状反例：656732c "修改插件名字"、9925bc7 "优化结构和接口分析skill"）。
- 涉 issue/工单的变更在 subject 或 trailer 中带引用编号。

## 三、merge 策略

### 现状

main 全量 34 条历史中 **0 条 merge commit**（证据：`git log main --merges` 输出为空），历史完全线性；同时无 "(#N)" squash 痕迹（出现率 0%）。判定为**线性历史，rebase 与 fast-forward 不可区分**；结合 personal/houle/test3、ready/27.0-终端鉴权 均未合入 main，现状实质为**单人/少数成员直推 main 为主**，尚无分支合入主线的实践样本。

### 约定（建议，待团队确认）

- 统一经分支开发、合入 main 保持线性历史，建议采用 squash merge 或 rebase + fast-forward，禁止 merge commit 与 squash 混用。
- 禁止直推 main 的约定是否启用，待团队规模确认后定稿（现状单人开发直推为主）。

## 四、MR 与评审要求

### 现状

**git 历史内评审痕迹**：merge commit message 中 MR/PR 编号出现率 0%（无 merge commit）；Reviewed-by / Approved-by / Signed-off-by trailer 出现率 0%（证据：main 全量 34 条 commit 的 trailers 均为空）。

**平台侧评审事实**：未识别（原因：环境中无 gh / glab CLI，无平台访问权限）。

**仓内成文约定**：无。仓内未发现 CODEOWNERS、CONTRIBUTING.md、.gitmessage、commitlint.config.*、.husky/、.github/pull_request_template.*、.gitlab/merge_request_templates/ 等协作配置文件。

### 约定（建议，待团队确认）

- 一切变更经 MR 合入 main，禁止直推默认分支（团队 >1 人活跃开发时启用）。
- 至少 1 人评审通过方可合入；评审结论以平台 approval 记录为准，不强制 commit trailer。
- 建议补充 CODEOWNERS 与 MR 模板，将评审要求成文固化。

## 附注

- 探测受限说明：环境中无 gh / glab CLI，平台侧评审数据未采集；origin/HEAD symbolic-ref 未设置，默认分支依据远端分支结构判定为 main。
- 采样口径说明：main 全量仅 34 条 commit（不足 100 条采样窗口），结论基于全量历史；分支为全量 6 条，未触发 50 条采样上限。
