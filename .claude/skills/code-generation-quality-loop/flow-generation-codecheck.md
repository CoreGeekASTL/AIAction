# 生成即CodeCheck流程图

```mermaid
flowchart TB
    subgraph Phase1["阶段一：代码生成"]
        A1[读取Story详设文档] --> A2[分析新增文件清单]
        A2 --> A2c[1.2 语言识别: Go / Java / Python]

        A2c --> A2b{"1.3 规则加载<br/>按语言只读对应文件"}

        A2b -->|Go| R1["读取 codecheck-go.md<br/>17条全量规则"]
        A2b -->|Java| R2["读取 codecheck-java.md<br/>13条全量规则"]
        A2b -->|Python| R3["读取 codecheck-python.md<br/>4条全量规则"]

        R1 --> A3
        R2 --> A3
        R3 --> A3

        A3["1.4 代码生成<br/>grep验证导入/方法/配置<br/>复用现有builder,不臆造"] --> A4[生成代码文件]

        A4 --> A5{"1.3 生成自检<br/>对照全量规则逐条自检<br/>第2道防线"}
        A5 -->|违规| A4["修正后重新生成"]
        A5 -->|通过| A6[生成测试文件]
    end

    A6 --> B0

    subgraph Phase2["阶段二：质量检查（兜底门禁）"]
        B0["CodeCheck全量扫描<br/>调用code-quality-check skill<br/>grep全量规则<br/>第3道防线"] --> B1[AI臆造验证]
        B1 --> B2[语法与并发检查]
        B2 --> B3{发现缺陷?}
        B3 -->|有缺陷| B4[修正代码]
        B4 --> B1
        B3 -->|无缺陷| C1
    end

    B0 -.违规.-> A4

    subgraph Phase3["阶段三：DT测试验证"]
        C1[运行UT/DT测试] --> C2{测试通过?}
        C2 -->|失败| C3[分析失败原因]
        C3 --> C4[修复问题]
        C4 --> C1
    end

    C2 -->|通过| E1

    subgraph Phase3_5["阶段3.5：集成测试验证"]
        E1["运行testsuit<br/>TC_001→TC_002→…→TC_N"] --> E2{全部SUCCESS?}
        E2 -->|失败| E3["修正脚本或代码"]
        E3 --> E1
    end

    E2 -->|通过| F1

    subgraph Phase3_6["阶段3.6：代码生成总报告"]
        F1["读取report模板<br/>聊天区输出7章节汇总"] --> D1
    end

    subgraph Phase4["阶段四：提交"]
        D1[提交代码变更]
    end
```

## 三道防线

| 防线 | 位置 | 机制 | 规则范围 |
| --- | --- | --- | --- |
| 第1道：规则前置 | 1.3 规则加载 | 生成前读对应语言全量规则，作为生成硬约束 | 全量 |
| 第2道：生成自检 | 1.3 生成自检 | 写入前对照全量规则逐条自检，违规即修 | 全量 |
| 第3道：兜底扫描 | 阶段二 | 调用 code-quality-check skill，grep 全量规则扫描 | 全量 |

## 核心特点

- **语言隔离**：生成 Go 只加载 `codecheck-go.md`，Java/Python 规则不进 context
- **全量覆盖**：三道防线均用全量规则，无子集、无红线
- **闭环修复**：任何一道发现违规 → 修正 → 回到自检重来，不放过
