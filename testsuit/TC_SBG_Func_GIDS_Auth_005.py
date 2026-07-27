#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# TC_SBG_Func_GIDS_Auth_005.py
# 测试用例: 边界与缓存覆盖(IMEI+IMSI联合鉴权)
#
# 看护点(逐条对齐testsuitcase.md):
#   S1 (005_001) IMEI空字符串 → code=401, msg含"format invalid"
#   S2 (005_002) IMSI空字符串 → code=401, msg含"format invalid"
#   S3 (005_003) deviceLoginAuth IMEI+IMSI均命中 → code=200
#   S4 (005_004) deviceLoginAuth IMSI不匹配 → code=-2, msg含"auth rejected"
#   S5 (005_005) 缓存命中后重复鉴权 → code=200(缓存直接返回)
#   S6 (005_006) 缓存未命中后重复鉴权 → code=401(缓存false直接返回)
#   S7 (005_007) 逃生态失效-白名单导入后不在白名单被拒绝 → code=401
#
# 前置：GIDS服务(127.0.0.1:9090)已启动, 白名单已导入
# 依赖：requests
# 用法：python TC_SBG_Func_GIDS_Auth_005.py

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
AUTH_CODE_LOGIN_REJECTED = -2
AUTH_CODE_AUTH_FAILED = 401


def make_login_request(imei, imsi):
    """构造登录请求体(19字段, AppType=1)"""
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


def auth_imei_imsi(imei, imsi):
    """POST /auth/v1/authIMEI 联合鉴权"""
    url = f"{GIDS_ADDR}/auth/v1/authIMEI"
    body = {'imei': imei, 'imsi': imsi}
    resp = requests.post(url, json=body, timeout=TIMEOUT, verify=False)
    return resp


def device_login_auth(imei, imsi):
    """POST deviceLoginAuth 设备登录鉴权"""
    url = f"{GIDS_ADDR}/app-api/devicetcp/app/login/v1/deviceLoginAuth"
    body = make_login_request(imei, imsi)
    resp = requests.post(url, json=body, timeout=TIMEOUT, verify=False)
    return resp


def check_auth_response(resp, expected_code, expected_msg_fragment=None, step_name=""):
    """通用鉴权响应校验"""
    print(f"[INFO] 响应状态码: {resp.status_code}")
    print(f"[INFO] 响应数据: {resp.text}")
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
        return False
    if expected_msg_fragment and expected_msg_fragment.lower() not in actual_msg.lower():
        print(f"[FAIL] {step_name} 预期msg含'{expected_msg_fragment}', 实际msg={actual_msg}")
        return False
    print(f"[PASS] {step_name} code={actual_code}, msg={actual_msg}")
    return True


def step1_empty_imei():
    """S1 (005_001): IMEI空字符串 -> code=401, msg含 format invalid"""
    print("\n[INFO] ========== S1 (005_001) IMEI空字符串参数错误 ==========")
    print(f"[INFO] IMEI='', IMSI={WHITELIST_IMSI}")
    try:
        resp = auth_imei_imsi('', WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_auth_response(resp, AUTH_CODE_AUTH_FAILED, "format invalid", "S1 IMEI空字符串")


def step2_empty_imsi():
    """S2 (005_002): IMSI空字符串 -> code=401, msg含 format invalid"""
    print("\n[INFO] ========== S2 (005_002) IMSI空字符串参数错误 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}, IMSI=''")
    try:
        resp = auth_imei_imsi(WHITELIST_IMEI, '')
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_auth_response(resp, AUTH_CODE_AUTH_FAILED, "format invalid", "S2 IMSI空字符串")


def step3_device_login_both_hit():
    """S3 (005_003): deviceLoginAuth IMEI+IMSI均命中 → code=200"""
    print("\n[INFO] ========== S3 (005_003) deviceLoginAuth联合鉴权通过 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    try:
        resp = device_login_auth(WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    print(f"[INFO] 响应状态码: {resp.status_code}")
    print(f"[INFO] 响应数据: {resp.text[:500]}")
    try:
        body = resp.json()
    except Exception:
        print(f"[FAIL] 响应非JSON: {resp.text}")
        return False
    if body.get('code') == AUTH_CODE_SUCCESS:
        print(f"[PASS] S3 deviceLoginAuth通过, code={body.get('code')}")
        return True
    else:
        print(f"[FAIL] S3 预期code={AUTH_CODE_SUCCESS}, 实际code={body.get('code')}, msg={body.get('msg')}")
        return False


def step4_device_login_imsi_mismatch():
    """S4 (005_004): deviceLoginAuth IMSI不匹配 -> code=-2, msg含 auth rejected"""
    print("\n[INFO] ========== S4 (005_004) deviceLoginAuth联合鉴权拒绝(IMSI不匹配) ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}(命中), IMSI={NON_WHITELIST_IMSI}(不匹配)")
    try:
        resp = device_login_auth(WHITELIST_IMEI, NON_WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    print(f"[INFO] 响应状态码: {resp.status_code}")
    print(f"[INFO] 响应数据: {resp.text}")
    try:
        body = resp.json()
    except Exception:
        print(f"[FAIL] 响应非JSON: {resp.text}")
        return False
    actual_code = body.get('code')
    actual_msg = body.get('msg', '')
    if actual_code == AUTH_CODE_LOGIN_REJECTED and 'auth rejected' in actual_msg.lower():
        print(f"[PASS] S4 deviceLoginAuth拒绝, code={actual_code}, msg={actual_msg}")
        return True
    elif actual_code == AUTH_CODE_SUCCESS:
        print(f"[FAIL] S4 应被拒绝但登录成功, code={actual_code}")
        print("[INFO] 提示: 白名单可能为空(逃生态)")
        return False
    else:
        print(f"[FAIL] S4 预期code={AUTH_CODE_LOGIN_REJECTED}(auth rejected), 实际code={actual_code}, msg={actual_msg}")
        return False


def step5_cache_hit_pass():
    """S5 (005_005): 缓存命中后重复鉴权 → code=200"""
    print("\n[INFO] ========== S5 (005_005) 缓存命中后重复鉴权直接返回 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    print("[INFO] 步骤: 先鉴权一次(命中), 再用相同参数鉴权(应命中缓存)")
    try:
        resp1 = auth_imei_imsi(WHITELIST_IMEI, WHITELIST_IMSI)
        print(f"[INFO] 第一次鉴权: code={resp1.json().get('code')}")
        resp2 = auth_imei_imsi(WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_auth_response(resp2, AUTH_CODE_SUCCESS, step_name="S5 缓存命中重复鉴权")


def step6_cache_miss_reject():
    """S6 (005_006): 缓存未命中后重复鉴权 → code=401"""
    print("\n[INFO] ========== S6 (005_006) 缓存未命中后重复鉴权仍返回拒绝 ==========")
    print(f"[INFO] IMEI={NON_WHITELIST_IMEI}, IMSI={NON_WHITELIST_IMSI}")
    print("[INFO] 步骤: 先鉴权一次(未命中), 再用相同参数鉴权(应命中false缓存)")
    try:
        resp1 = auth_imei_imsi(NON_WHITELIST_IMEI, NON_WHITELIST_IMSI)
        print(f"[INFO] 第一次鉴权: code={resp1.json().get('code')}")
        resp2 = auth_imei_imsi(NON_WHITELIST_IMEI, NON_WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_auth_response(resp2, AUTH_CODE_AUTH_FAILED, step_name="S6 缓存false重复鉴权")


def step7_escape_disabled():
    """S7 (005_007): 逃生态失效-白名单导入后不在白名单被拒绝 → code=401"""
    print("\n[INFO] ========== S7 (005_007) 逃生态失效-白名单导入后不在白名单被拒绝 ==========")
    print(f"[INFO] IMEI={NON_WHITELIST_IMEI}, IMSI={NON_WHITELIST_IMSI}")
    print("[INFO] 前置: 白名单已导入(逃生态已关闭), 此IMEI+IMSI不在白名单")
    try:
        resp = auth_imei_imsi(NON_WHITELIST_IMEI, NON_WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_auth_response(resp, AUTH_CODE_AUTH_FAILED, step_name="S7 逃生态失效正常拒绝")


def main():
    print("=" * 66)
    print("  边界与缓存覆盖(IMEI+IMSI联合鉴权)  TC_SBG_Func_GIDS_Auth_005")
    print("=" * 66)
    print(f"[INFO] GIDS_ADDR = {GIDS_ADDR}")
    print(f"[INFO] 白名单内: IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    print(f"[INFO] 白名单外: IMEI={NON_WHITELIST_IMEI}, IMSI={NON_WHITELIST_IMSI}")

    results = []
    results.append(("S1 (005_001) IMEI空字符串参数错误", step1_empty_imei()))
    results.append(("S2 (005_002) IMSI空字符串参数错误", step2_empty_imsi()))
    results.append(("S3 (005_003) deviceLoginAuth联合鉴权通过", step3_device_login_both_hit()))
    results.append(("S4 (005_004) deviceLoginAuth IMSI不匹配拒绝", step4_device_login_imsi_mismatch()))
    results.append(("S5 (005_005) 缓存命中后重复鉴权通过", step5_cache_hit_pass()))
    results.append(("S6 (005_006) 缓存未命中后重复鉴权拒绝", step6_cache_miss_reject()))
    results.append(("S7 (005_007) 逃生态失效正常拒绝", step7_escape_disabled()))

    print("\n[INFO] ========== 测试结果汇总 ==========")
    for name, ok in results:
        print(f"{'[PASS]' if ok else '[FAIL]'} {name}")

    all_passed = all(ok for _, ok in results)
    print(f"\n{'[SUCCESS]' if all_passed else '[FAILED]'} ========== 测试完成 ==========")
    return 0 if all_passed else 1


if __name__ == "__main__":
    sys.exit(main())
