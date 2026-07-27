# Python CodeCheck 规则参考 (pylint)

## G.ERR.01 — try块内只包含可能抛异常的代码 (H0719)

**规则**：try代码块内只包含可能抛出异常的代码段，不会抛异常的初始化/赋值语句应移到try之前。

**常见场景**：
- 列表初始化 `data = []` 放在try内 → 应移到try之前
- 变量赋值 `result = None` 放在try内 → 应移到try之前

**修正示例**：
```python
# 问题代码
def get_items():
    try:
        result = []
        for item in items:
            result.append(item.json())
        return result
    except Exception:
        ...

# 修正后：初始化语句移到try之前
def get_items():
    result = []
    try:
        for item in items:
            result.append(item.json())
        return result
    except Exception:
        ...
```

---

## G.FMT.02 — 行宽不超过120个字符 (C0301)

**规则**：每行不超过120个窄字符（CJK字符算2个窄字符）。

**修正方式**：
1. 用括号隐式续行（推荐）
2. 用反斜杠显式续行

**修正示例**：
```python
# 问题代码（131 chars）
allowlisted_extension_id = f"{allowlisted_extension_id},{target_ext_id}" if allowlisted_extension_id else target_ext_id

# 修正后：括号隐式续行
allowlisted_extension_id = (
    f"{allowlisted_extension_id},{target_ext_id}"
    if allowlisted_extension_id
    else target_ext_id
)

# 问题代码：函数签名过长
def __init__(self, browser: Browser, browser_type: BrowserType, browser_id: str, userdata: str, playwright: Playwright):

# 修正后：参数换行
def __init__(
    self, browser: Browser, browser_type: BrowserType,
    browser_id: str, userdata: str, playwright: Playwright
):

# 问题代码：字符串值过长（215 chars）
"value": "{%22optional%22:true%2C%22ga%22:true%2C%22af%22:true...}"

# 修正后：字符串拼接续行
"value": (
    "{%22optional%22:true%2C%22ga%22:true"
    "%2C%22af%22:true%2C%22fbp%22:true"
    "%2C%22lip%22:true%2C%22bing%22:true"
    "%2C%22ttads%22:true%2C%22reddit%22:true"
    "%2C%22hubspot%22:true%2C%22version%22:%22v10%22}"
),
```

---

## G.FMT.07 — 导入排序：标准库→第三方库→本地模块 (C0411)

**规则**：import部分按照标准库、第三方库、应用程序自定义模块的顺序排列，每组之间空一行。

**排序标准**：
1. **标准库**：`import json`, `import os`, `import logging`, `from pathlib import Path`, `from typing import ...`
2. **第三方库**：`from fastapi import ...`, `from playwright.async_api import ...`
3. **本地模块**：`from browser_proxy.api.common import ...`, `from browser_proxy.api.logger_config import ...`

**修正示例**：
```python
# 问题代码：标准库import混在第三方库之后
import json
import os
from fastapi import APIRouter, HTTPException, Request
from playwright.async_api import async_playwright
import logging          # 标准库应在前
from pathlib import Path  # 标准库应在前
from browser_proxy.api.common import BrowserType, BrowserWrapper, browser_list, browser_get
from browser_proxy.api.logger_config import log

# 修正后：标准库→第三方库→本地模块，每组空一行
import json
import logging
import os
from pathlib import Path

from fastapi import APIRouter, HTTPException, Request
from playwright.async_api import async_playwright

from browser_proxy.api.common import BrowserType, BrowserWrapper, browser_list, browser_get
from browser_proxy.api.logger_config import log
```

---

## G.CLS.06 — 类方法排序：properties在regular methods之前 (H3606)

**规则**：类的方法建议统一排列，properties应放在regular methods之前。

**排序标准**：
1. `__init__`
2. `@property` 方法（属性访问器）
3. 普通 `def` 方法（实例方法）
4. `async def` 方法（异步方法）

**修正示例**：
```python
# 问题代码：set_record_data和append_page在current属性之前
class ContextWrapper:
    def __init__(self, ...):
        ...
    @property
    def id(self) -> str: ...
    @property
    def context(self) -> BrowserContext: ...
    @property
    def pages(self) -> List[PageWrapper]: ...
    @property
    def record_data(self): ...
    def set_record_data(self, data): ...    # 方法应在属性之后
    def append_page(self, page: PageWrapper): ...  # 方法应在属性之后
    @property
    def current(self) -> PageWrapper: ...   # 属性应在前

# 修正后：所有properties在前，methods在后
class ContextWrapper:
    def __init__(self, ...):
        ...
    @property
    def id(self) -> str: ...
    @property
    def context(self) -> BrowserContext: ...
    @property
    def pages(self) -> List[PageWrapper]: ...
    @property
    def record_data(self): ...
    @property
    def current(self) -> PageWrapper: ...
    def set_record_data(self, data): ...
    def append_page(self, page: PageWrapper): ...
```

```python
# 问题代码：close/as_json/append_context在contexts属性之前
class BrowserWrapper:
    def __init__(self, ...): ...
    @property
    def used(self) -> int: ...
    @property
    def id(self) -> str: ...
    @property
    def browser_type(self) -> BrowserType: ...
    @property
    def browser(self) -> Browser: ...
    @property
    def userdata(self) -> str: ...
    async def close(self) -> None: ...      # 方法应在属性之后
    def as_json(self): ...                   # 方法应在属性之后
    def append_context(self, context): ...   # 方法应在属性之后
    @property
    def contexts(self) -> List[ContextWrapper]: ...  # 属性应在前

# 修正后：所有properties在前，methods在后
class BrowserWrapper:
    def __init__(self, ...): ...
    @property
    def used(self) -> int: ...
    @property
    def id(self) -> str: ...
    @property
    def browser_type(self) -> BrowserType: ...
    @property
    def browser(self) -> Browser: ...
    @property
    def userdata(self) -> str: ...
    @property
    def contexts(self) -> List[ContextWrapper]: ...
    async def close(self) -> None: ...
    def as_json(self): ...
    def append_context(self, context): ...
```
