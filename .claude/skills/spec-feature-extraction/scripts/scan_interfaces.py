#!/usr/bin/env python3
"""scan_interfaces.py — 存量代码仓接口定义扫描脚本（interface-feature-analyzer 第 2 步用）

扫描四类接口定义并输出带 文件:行号 证据的清单：
  1. IDL/契约类：proto service/rpc、thrift service、OpenAPI paths、GraphQL Query/Mutation
  2. 框架路由类：Go(gin/echo/grpc/net-http)、Java(Spring/Dubbo)、Python(FastAPI/Flask/Django)、C/C++ 注册点
  3. 语言级接口：Go interface、Java interface、C++ 纯虚抽象类、Python ABC/Protocol、C ops 结构体
  4. 消息/事件/定时入口：MQ 消费订阅、事件注册、定时任务

自动排除：生成代码（DO NOT EDIT 头、gen/proto 生成目录）、第三方目录（vendor/node_modules 等）。

用法：
  python3 scan_interfaces.py <repo_path> [-o scan_result.md] [--lang go|java|cpp|python] [--format json]
"""

import argparse
import json
import os
import re
import sys
from collections import defaultdict

# ---- 排除规则 ----

EXCLUDE_DIRS = {
    "vendor", "node_modules", "third_party", "thirdparty", "3rd", "3rdparty",
    ".git", ".idea", ".vscode", "dist", "build", "out", "target", "bin",
    "__pycache__", ".tox", "venv", ".venv", "test", "tests", "testdata",
}
EXCLUDE_DIR_PARTS = ("/gen/", "/generated/", "/gen-go/", "/proto_gen/")

GENERATED_HEADER = re.compile(r"(DO NOT EDIT|Code generated|auto-?generated)", re.I)
GENERATED_FILE_SUFFIX = (
    ".pb.go", ".pb.cc", ".pb.h", "_grpc.pb.go", ".pb2.py",
    ".connect.go", "_ttypes.py",
)

# ---- 模式库：{类别: [(子类型, 语言, 正则, 说明)]} ----
# 语言字段用于 --lang 过滤：any 表示跨语言（如 IDL 文件）

# IDL 模式按扩展名区分文件，避免 thrift 模式误中 proto 等同名关键字
IDL_PATTERN_EXT = {
    "proto-service": (".proto",),
    "proto-rpc": (".proto",),
    "thrift-service": (".thrift",),
    "openapi-path": (".yaml", ".yml", ".json"),
    "graphql-op": (".graphql", ".gql"),
}

PATTERNS = {
    "idl_contract": [
        ("proto-service", "any", re.compile(r"^\s*service\s+(\w+)\s*\{"), "proto service 定义"),
        ("proto-rpc", "any", re.compile(r"^\s*rpc\s+(\w+)\s*\("), "proto rpc 方法"),
        ("thrift-service", "any", re.compile(r"^\s*service\s+(\w+)"), "thrift service 定义"),
        ("openapi-path", "any", re.compile(r'^\s{2,}(/[\w./{}-]+)\s*:\s*$'), "OpenAPI path"),
        ("graphql-op", "any", re.compile(r"^\s*type\s+(Query|Mutation|Subscription)\b"), "GraphQL 入口类型"),
    ],
    "framework_route": [
        ("go-gin", "go", re.compile(r'\b(\w+(?:\.\w+)*)\.(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS|Any|Handle)\s*\(\s*"([^"]*)"'), "Go gin/echo/chi 路由"),
        ("go-grpc-register", "go", re.compile(r"\bRegister(\w+)Server\s*\("), "Go gRPC 服务注册"),
        ("go-nethttp", "go", re.compile(r'\b(?:http\.)?Handle(?:Func)?\s*\(\s*"([^"]*)"'), "Go net/http 路由"),
        ("java-spring", "java", re.compile(r"@(GetMapping|PostMapping|PutMapping|DeleteMapping|PatchMapping|RequestMapping)\s*\(?\s*(?:value\s*=\s*)?[\{\"]?([^)\}\"]*)"), "Spring MVC 路由注解"),
        ("java-dubbo", "java", re.compile(r"@(DubboService|Service)\b"), "Dubbo 服务暴露"),
        ("py-fastapi", "python", re.compile(r'@(\w+)\.(get|post|put|delete|patch|api_route|route)\s*\(\s*["\']([^"\']*)'), "FastAPI/Flask 路由"),
        ("py-django-url", "python", re.compile(r'\bpath\s*\(\s*["\']([^"\']*)["\']'), "Django urlpatterns"),
        ("cpp-register", "cpp", re.compile(r"\b(Reg(?:ister)?(?:Handler|Service|Msg|Cmd|Api))\w*\s*\("), "C/C++ handler/服务注册点"),
    ],
    "language_interface": [
        ("go-interface", "go", re.compile(r"^\s*type\s+(\w+)\s+interface\s*\{"), "Go interface"),
        ("java-interface", "java", re.compile(r"^\s*(?:public\s+)?(?:abstract\s+)?interface\s+(\w+)"), "Java interface"),
        ("cpp-pure-virtual", "cpp", re.compile(r"virtual\s+[\w:<>,\s*&]+\s+(\w+)\s*\([^)]*\)\s*(?:const\s*)?=\s*0\s*;"), "C++ 纯虚函数"),
        ("c-ops-struct", "cpp", re.compile(r"struct\s+(\w+_(?:ops|operations|vtable|fns))\s*\{"), "C 函数指针表结构"),
        ("py-abc", "python", re.compile(r"class\s+(\w+)\s*\([^)]*(?:ABC|Protocol)[^)]*\)"), "Python ABC/Protocol"),
    ],
    "async_entry": [
        ("mq-consume", "any", re.compile(r"\b(Subscribe|Consume|AddConsumer|RegisterListener)\w*\s*\(|@KafkaListener|@RocketMQMessageListener"), "MQ 消费订阅"),
        ("event-register", "any", re.compile(r"\b(eventbus\.|EventBus\.|\.On|\.AddObserver|RegisterCallback)\s*\("), "事件/回调注册"),
        ("cron-entry", "any", re.compile(r"@Scheduled\b|AddFunc\s*\(\s*[\"']@|cron\.\w+\s*\("), "定时任务入口"),
    ],
}

CATEGORY_TITLE = {
    "idl_contract": "IDL/契约类接口",
    "framework_route": "框架路由类接口（HTTP/RPC 入口）",
    "language_interface": "语言级接口（模块间契约）",
    "async_entry": "消息/事件/定时入口",
}

LANG_EXT = {
    "go": (".go",),
    "java": (".java",),
    "cpp": (".c", ".cc", ".cpp", ".cxx", ".h", ".hpp"),
    "python": (".py",),
}
IDL_EXT = (".proto", ".thrift", ".graphql", ".gql", ".fbs", ".capnp", ".yaml", ".yml")

MAX_MATCHES_PER_PATTERN = 500  # 单模式命中上限，防止宏/表驱动注册刷屏


def is_generated(path):
    if path.endswith(GENERATED_FILE_SUFFIX):
        return True
    norm = path.replace(os.sep, "/")
    return any(part in f"/{norm}" for part in EXCLUDE_DIR_PARTS)


def file_lang(path):
    for lang, exts in LANG_EXT.items():
        if path.endswith(exts):
            return lang
    return "any" if path.endswith(IDL_EXT) else None


def scan(repo, lang_filter=None):
    """返回 {类别: [(子类型, 相对路径, 行号, 行内容)]}"""
    hits = defaultdict(list)
    pattern_count = defaultdict(int)
    for root, dirs, files in os.walk(repo):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS and not d.startswith(".")]
        for fname in files:
            fpath = os.path.join(root, fname)
            rel = os.path.relpath(fpath, repo)
            if is_generated(rel):
                continue
            f_lang = file_lang(fname)
            if f_lang is None:
                continue
            if lang_filter and f_lang not in (lang_filter, "any"):
                continue
            try:
                with open(fpath, encoding="utf-8", errors="ignore") as f:
                    head = f.read(2048)
                    if GENERATED_HEADER.search(head):
                        continue
                    lines = head.splitlines() + f.readlines()
            except OSError:
                continue
            for category, plist in PATTERNS.items():
                for subtype, plang, regex, _desc in plist:
                    key = (category, subtype)
                    if pattern_count[key] >= MAX_MATCHES_PER_PATTERN:
                        continue
                    if lang_filter and plang not in ("any", lang_filter):
                        continue
                    if plang != "any" and f_lang != plang:
                        continue
                    # IDL 模式只匹配 IDL 文件；any 类 async 模式匹配所有语言文件
                    if category == "idl_contract" and f_lang != "any":
                        continue
                    if category != "idl_contract" and f_lang == "any":
                        continue
                    if category == "idl_contract":
                        ext = IDL_PATTERN_EXT.get(subtype)
                        if ext and not fname.endswith(ext):
                            continue
                    for lineno, line in enumerate(lines, 1):
                        if regex.search(line):
                            hits[category].append((subtype, rel, lineno, line.strip()[:160]))
                            pattern_count[key] += 1
    return hits


def to_markdown(hits, repo):
    out = [f"# 接口扫描结果\n", f"> 代码仓：`{os.path.abspath(repo)}`\n"]
    total = 0
    for category in PATTERNS:
        rows = hits.get(category, [])
        total += len(rows)
        out.append(f"\n## {CATEGORY_TITLE[category]}（{len(rows)} 处）\n")
        if not rows:
            out.append("（未命中——需人工确认是确实没有还是扫描遗漏）\n")
            continue
        by_subtype = defaultdict(list)
        for subtype, rel, lineno, line in rows:
            by_subtype[subtype].append((rel, lineno, line))
        for subtype, items in sorted(by_subtype.items()):
            out.append(f"\n### {subtype}（{len(items)}）\n")
            for rel, lineno, line in items[:100]:
                out.append(f"- `{rel}:{lineno}` — `{line}`")
            if len(items) > 100:
                out.append(f"- ……（其余 {len(items) - 100} 条略，用 grep 补全）")
    out.insert(2, f"\n共命中 **{total}** 处接口定义/注册点。本清单为线索入口，聚类与精读需结合 grep 补全。\n")
    return "\n".join(out) + "\n"


def main():
    ap = argparse.ArgumentParser(description="存量代码仓接口定义扫描")
    ap.add_argument("repo", help="代码仓路径")
    ap.add_argument("-o", "--output", help="输出文件（默认 stdout）")
    ap.add_argument("--lang", choices=["go", "java", "cpp", "python"], help="限定语言")
    ap.add_argument("--format", choices=["md", "json"], default="md")
    args = ap.parse_args()

    if not os.path.isdir(args.repo):
        sys.exit(f"错误：{args.repo} 不是目录")

    hits = scan(args.repo, args.lang)
    if args.format == "json":
        data = {
            cat: [
                {"subtype": st, "file": rel, "line": ln, "text": text}
                for st, rel, ln, text in rows
            ]
            for cat, rows in hits.items()
        }
        result = json.dumps(data, ensure_ascii=False, indent=2)
    else:
        result = to_markdown(hits, args.repo)

    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            f.write(result)
        print(f"已写入 {args.output}（{sum(len(v) for v in hits.values())} 处命中）")
    else:
        print(result)


if __name__ == "__main__":
    main()
