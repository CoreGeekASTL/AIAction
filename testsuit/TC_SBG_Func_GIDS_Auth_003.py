#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# TC_SBG_Func_GIDS_Auth_003.py
# 测试用例: 登录链路IMEI+IMSI联合鉴权
#
# 看护点(逐条对齐testsuitcase.md):
#   S1 (003_001) gridLoginAuth IMEI+IMSI均命中 → code=200
#   S2 (003_002) gridLoginAuth IMSI不匹配 → code=-2, msg含"auth rejected"
#   S3 (003_003) gridLoginAuthOpenBrowser IMEI+IMSI均命中 → code=200
#   S4 (003_004) gridLoginAuthOpenBrowser IMEI不匹配 → code=-2, msg含"auth rejected"
#
# 前置：GIDS服务(127.0.0.1:9090)已启动, 白名单已导入
# 依赖：requests
# 用法：python TC_SBG_Func_GIDS_Auth_003.py

import sys
import os
import urllib3

import requests

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

GIDS_ADDR = os.getenv('GIDS_ADDR', 'http://127.0.0.1:9090').rstrip('/')
TIMEOUT = 30

WHITELIST_IMEI = '625841245402541'
WHITELIST_IMSI = '460011234567890'
NON_WHITELIST_IMEI = '625841245402599'
NON_WHITELIST_IMSI = '460019999999999'

AUTH_CODE_SUCCESS = 200
AUTH_CODE_AUTH_FAILED = -2


def make_login_request(imei, imsi):
    """构造登录请求体(19字段)"""
    return {
        "imsi": imsi,
        "imei": imei,
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
        "clientLanguage": "en",
        "deviceType": "1000"
    }


def grid_login_auth(imei, imsi):
    """POST gridLoginAuth 宫格登录认证"""
    url = f"{GIDS_ADDR}/app-api/devicetcp/app/login/v1/gridLoginAuth"
    body = make_login_request(imei, imsi)
    resp = requests.post(url, json=body, timeout=TIMEOUT, verify=False)
    return resp


def grid_login_auth_open_browser(imei, imsi):
    """POST gridLoginAuthOpenBrowser 宫格登录打开浏览器"""
    url = f"{GIDS_ADDR}/app-api/devicetcp/app/login/v1/gridLoginAuthOpenBrowser"
    body = make_login_request(imei, imsi)
    resp = requests.post(url, json=body, timeout=TIMEOUT, verify=False)
    return resp


def check_login_response(resp, expected_code, expected_msg_fragment=None, step_name=""):
    """通用登录响应校验"""
    print(f"[INFO] 响应状态码: {resp.status_code}")
    print(f"[INFO] 响应数据: {resp.text[:500]}")
    try:
        body = resp.json()
    except Exception:
        print(f"[FAIL] 响应非JSON: {resp.text}")
        return False
    actual_code = body.get('code')
    actual_msg = body.get('msg', '')
    if actual_code != expected_code:
        print(f"[FAIL] {step_name} 预期code={expected_code}, 实际code={actual_code}, msg={actual_msg}")
        if actual_code == AUTH_CODE_SUCCESS and expected_code != AUTH_CODE_SUCCESS:
            print("[INFO] 提示: 白名单可能为空(逃生态)")
            print("[INFO] 请先执行 TC_SBG_Func_GIDS_Auth_001 导入白名单后再测试")
        return False
    if expected_msg_fragment and expected_msg_fragment.lower() not in actual_msg.lower():
        print(f"[FAIL] {step_name} 预期msg含'{expected_msg_fragment}', 实际msg={actual_msg}")
        return False
    print(f"[PASS] {step_name} code={actual_code}, msg={actual_msg}")
    return True


def step1_grid_login_both_hit():
    """S1 (003_001): gridLoginAuth IMEI+IMSI均命中 → code=200"""
    print("\n[INFO] ========== S1 (003_001) gridLoginAuth联合鉴权通过 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    try:
        resp = grid_login_auth(WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_login_response(resp, AUTH_CODE_SUCCESS, step_name="S1 gridLoginAuth通过")


def step2_grid_login_imsi_mismatch():
    """S2 (003_002): gridLoginAuth IMSI不匹配 -> code=-2, msg含 auth rejected"""
    print("\n[INFO] ========== S2 (003_002) gridLoginAuth联合鉴权拒绝(IMSI不匹配) ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}(命中), IMSI={NON_WHITELIST_IMSI}(不匹配)")
    try:
        resp = grid_login_auth(WHITELIST_IMEI, NON_WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_login_response(resp, AUTH_CODE_AUTH_FAILED, "auth rejected", "S2 gridLoginAuth拒绝")


def step3_open_browser_both_hit():
    """S3 (003_003): gridLoginAuthOpenBrowser IMEI+IMSI均命中 → code=200"""
    print("\n[INFO] ========== S3 (003_003) gridLoginAuthOpenBrowser联合鉴权通过 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    try:
        resp = grid_login_auth_open_browser(WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_login_response(resp, AUTH_CODE_SUCCESS, step_name="S3 gridLoginAuthOpenBrowser通过")


def step4_open_browser_imei_mismatch():
    """S4 (003_004): gridLoginAuthOpenBrowser IMEI不匹配 -> code=-2, msg含 auth rejected"""
    print("\n[INFO] ========== S4 (003_004) gridLoginAuthOpenBrowser联合鉴权拒绝(IMEI不匹配) ==========")
    print(f"[INFO] IMEI={NON_WHITELIST_IMEI}(不匹配), IMSI={WHITELIST_IMSI}(命中)")
    try:
        resp = grid_login_auth_open_browser(NON_WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_login_response(resp, AUTH_CODE_AUTH_FAILED, "auth rejected", "S4 gridLoginAuthOpenBrowser拒绝")


def main():
    print("=" * 66)
    print("  登录链路IMEI+IMSI联合鉴权  TC_SBG_Func_GIDS_Auth_003")
    print("=" * 66)
    print(f"[INFO] GIDS_ADDR = {GIDS_ADDR}")
    print(f"[INFO] 白名单内: IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    print(f"[INFO] 白名单外: IMEI={NON_WHITELIST_IMEI}, IMSI={NON_WHITELIST_IMSI}")

    results = []
    results.append(("S1 (003_001) gridLoginAuth联合鉴权通过", step1_grid_login_both_hit()))
    results.append(("S2 (003_002) gridLoginAuth IMSI不匹配拒绝", step2_grid_login_imsi_mismatch()))
    results.append(("S3 (003_003) gridLoginAuthOpenBrowser联合鉴权通过", step3_open_browser_both_hit()))
    results.append(("S4 (003_004) gridLoginAuthOpenBrowser IMEI不匹配拒绝", step4_open_browser_imei_mismatch()))

    print("\n[INFO] ========== 测试结果汇总 ==========")
    for name, ok in results:
        print(f"{'[PASS]' if ok else '[FAIL]'} {name}")

    all_passed = all(ok for _, ok in results)
    print(f"\n{'[SUCCESS]' if all_passed else '[FAILED]'} ========== 测试完成 ==========")
    return 0 if all_passed else 1


if __name__ == "__main__":
    sys.exit(main())
