#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
样例A: 纯HTTP接口测试 (pytest + requests 风格)

测试目标: GIDS 宫格登录接口 (gridLoginAuthOpenBrowser)
涉及服务: GIDS Mock (端口 9090)

测试场景:
  - 正向验证: 正确参数登录成功, 返回token+tcpAddr
  - 异常验证: IMSI格式错误/为空/过短/缺少imei字段, 应被拒绝
  - 边界验证: IMSI长度过短(5位), 应被拒绝
  - 状态验证: GIDS服务健康检查通过

此样例展示:
  1. BASE_URL 从环境变量读取 (非硬编码)
  2. 正向场景必须断言状态码 + 响应体关键字段值 (非伪断言)
  3. 异常场景断言错误码 + 错误信息字段
  4. 边界场景断言具体拒绝行为
  5. 每条断言验证实际值而非恒真/空断言
  6. 测试函数docstring包含: 测试场景 + 测试步骤 + 断言条件
"""

import os
import pytest
import requests
import urllib3

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

BASE_URL = os.environ.get("GIDS_ADDR", "http://127.0.0.1:9090")
GIDS_LOGIN_PATH = "/app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser"

VALID_IMEI = os.environ.get("DEVICE_IMEI", "6258412454025411")
VALID_IMSI = os.environ.get("DEVICE_IMSI", "685101555652111")


def _make_request(imei: str, imsi: str, **overrides) -> requests.Response:
    body = {
        "imei": imei,
        "imsi": imsi,
        "manufacturer": "Huawei",
        "model": "Mate60",
        "appType": "1",
        "extendModel": "default",
        "country": "CN",
        "platform": "1",
        "width": "240",
        "height": "320",
        "mcc": "460",
        "mnc": "00",
        "lac": "100",
        "ci": "5.21",
        "rxlev": "-72",
        "totalKb": "1424122",
        "freeKb": "1424122",
        "clientLanguage": "en_US",
        "deviceType": "1000",
    }
    body.update(overrides)
    return requests.post(
        f"{BASE_URL}{GIDS_LOGIN_PATH}",
        json=body,
        timeout=30,
        verify=False,
    )


class TestGridLoginAuth:
    """GIDS宫格登录接口测试"""

    def test_valid_params_success(self):
        """
        测试场景: 正向验证 — 正确参数登录成功
        测试步骤:
          1. 使用正确IMEI+IMSI发送POST请求到gridLoginAuthOpenBrowser
          2. 解析响应JSON
        断言条件:
          - resp.status_code == 200 (HTTP状态码)
          - data.code == 200 (业务状态码)
          - data.data.token 非空 (返回有效token)
          - data.data.tcpAddr 非空 (返回控制通道地址)
        """
        resp = _make_request(VALID_IMEI, VALID_IMSI)

        assert resp.status_code == 200
        data = resp.json()
        assert data.get("code") == 200
        assert "data" in data
        assert data["data"].get("token") not in ("", None)
        assert data["data"].get("tcpAddr") not in ("", None)

    def test_imsi_with_letters_rejected(self):
        """
        测试场景: 异常验证 — IMSI包含字母字符
        测试步骤:
          1. 使用IMSI="68510155565ABC"(包含字母)发送POST请求
          2. 解析响应JSON
        断言条件:
          - resp.status_code == 200 (HTTP层面正常)
          - data.code != 200 (业务层面拒绝)
          - data中包含message或msg字段 (返回错误信息)
        """
        resp = _make_request(VALID_IMEI, "68510155565ABC")

        assert resp.status_code == 200
        data = resp.json()
        assert data.get("code") != 200
        assert "message" in data or "msg" in data

    def test_imsi_empty_rejected(self):
        """
        测试场景: 异常验证 — IMSI为空字符串
        测试步骤:
          1. 使用IMSI=""(空字符串)发送POST请求
          2. 解析响应JSON
        断言条件:
          - resp.status_code == 200 (HTTP层面正常)
          - data.code != 200 (业务层面拒绝)
          - data中包含message或msg字段 (返回错误信息)
        """
        resp = _make_request(VALID_IMEI, "")

        assert resp.status_code == 200
        data = resp.json()
        assert data.get("code") != 200
        assert "message" in data or "msg" in data

    def test_imsi_too_short_rejected(self):
        """
        测试场景: 边界验证 — IMSI长度过短(5位, 正常应为15位)
        测试步骤:
          1. 使用IMSI="12345"(仅5位)发送POST请求
          2. 解析响应JSON
        断言条件:
          - resp.status_code == 200 (HTTP层面正常)
          - data.code != 200 (业务层面拒绝, IMSI长度不合规)
          - data中包含message或msg字段 (返回错误信息)
        """
        resp = _make_request(VALID_IMEI, "12345")

        assert resp.status_code == 200
        data = resp.json()
        assert data.get("code") != 200
        assert "message" in data or "msg" in data

    def test_missing_imei_field_rejected(self):
        """
        测试场景: 异常验证 — 请求体缺少imei字段
        测试步骤:
          1. 发送POST请求, 请求体仅包含imsi, 不包含imei
          2. 解析响应JSON
        断言条件:
          - resp.status_code in (200, 400) (HTTP可能返回400)
          - data.code != 200 (业务层面拒绝)
          - data中包含message或msg字段 (返回错误信息)
        """
        resp = requests.post(
            f"{BASE_URL}{GIDS_LOGIN_PATH}",
            json={"imsi": VALID_IMSI},
            timeout=30,
            verify=False,
        )

        assert resp.status_code in (200, 400)
        data = resp.json()
        assert data.get("code") != 200
        assert "message" in data or "msg" in data

    def test_health_endpoint(self):
        """
        测试场景: 状态验证 — GIDS服务健康检查
        测试步骤:
          1. 发送GET请求到 /health 端点
        断言条件:
          - resp.status_code == 200 (服务正常运行)
        """
        resp = requests.get(f"{BASE_URL}/health", timeout=10, verify=False)

        assert resp.status_code == 200
