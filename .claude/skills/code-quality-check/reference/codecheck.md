# CodeCheck 规则参考

CodeCheck代码质量检查规则参考，按语言分类。

## Go规则

详见 [codecheck-go.md](codecheck-go.md)，涵盖：
- 函数复杂度（行数/深度）
- 冗余代码检查
- 注释规范（版权声明、包注释、注释位置）
- 命名规范（chan方向、错误命名、文件命名、命名风格、包名一致性）
- 魔法数字
- 安全规范（资源泄漏、错误返回值检查、文件权限、switch default）

## Java规则

详见 [codecheck-java.md](codecheck-java.md)，涵盖：
- 格式规范（NeedBrace、LineLength、UpperEll）
- 类型安全（Locale、AssignCharset）
- 异常处理（EmptyCatch、ThrowRawException）
- 日志规范（LogModifier、LogWithoutChinese、DoNotLogWithSystemPrint）
- 注释规范（TodoComment）
- 资源管理（ResourceRelease）
- 代码清洁（CodeInComment、UnusedImport）

## Python规则（pylint）

详见 [codecheck-python.md](codecheck-python.md)，涵盖：
- 异常处理（H0719: try块内只包含可能抛异常的代码）
- 格式规范（C0301: 行宽不超过120字符）
- 导入排序（C0411: 标准库→第三方库→本地模块）
- 类方法排序（H3606: properties在regular methods之前）
