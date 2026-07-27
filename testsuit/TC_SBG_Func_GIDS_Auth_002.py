#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# TC_SBG_Func_GIDS_Auth_002.py
# 测试用例: IMEI+IMSI联合鉴权正向/反向/异常/逃生态
#
# 看护点(逐条对齐testsuitcase.md):
#   S1 (002_001) IMEI+IMSI同时命中 → code=200
#   S2 (002_002) IMEI命中IMSI未命中 → code=401
#   S3 (002_003) IMSI命中IMEI未命中 → code=401
#   S4 (002_004) IMEI与IMSI均不在白名单 → code=401
#   S5 (002_005) IMEI缺失 → code=401, msg含"format invalid"
#   S6 (002_006) IMSI缺失 → code=401, msg含"format invalid"
#   S7 (002_007) IMEI非15位纯数字 → code=401, msg含"format invalid"
#   S8 (002_008) IMSI非15位纯数字 → code=401, msg含"format invalid"
#   S9 (002_009) 逃生态-白名单表为空时放行 → code=200
#
# 前置：GIDS服务(127.0.0.1:9090)已启动
#        S1~S8需白名单已导入; S9需白名单表为空(全新gids.db或导出后清空)
# 依赖：requests
# 用法：python TC_SBG_Func_GIDS_Auth_002.py

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
INVALID_IMEI_ALPHA = '6258412454025ABC'
INVALID_IMSI_ALPHA = '46001ABC4567890'

AUTH_CODE_SUCCESS = 200
AUTH_CODE_AUTH_FAILED = 401


def auth_imei_imsi(imei, imsi):
    """POST /auth/v1/authIMEI 校验IMEI+IMSI联合鉴权"""
    url = f"{GIDS_ADDR}/auth/v1/authIMEI"
    body = {'imei': imei, 'imsi': imsi}
    resp = requests.post(url, json=body, timeout=TIMEOUT, verify=False)
    return resp


def check_response(resp, expected_code, expected_msg_fragment=None, step_name=""):
    """通用响应校验"""
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
            print("[INFO] 提示: 白名单可能为空(逃生态), 所有请求均放行")
            print("[INFO] 请先执行 TC_SBG_Func_GIDS_Auth_001 导入白名单后再测试")
        return False
    if expected_msg_fragment and expected_msg_fragment.lower() not in actual_msg.lower():
        print(f"[FAIL] {step_name} 预期msg含'{expected_msg_fragment}', 实际msg={actual_msg}")
        return False
    print(f"[PASS] {step_name} code={actual_code}, msg={actual_msg}")
    return True


def step1_both_hit():
    """S1 (002_001): IMEI+IMSI同时命中白名单 → code=200"""
    print("\n[INFO] ========== S1 (002_001) IMEI+IMSI同时命中鉴权通过 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    try:
        resp = auth_imei_imsi(WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败: 无法连接到 {GIDS_ADDR}")
        return False
    return check_response(resp, AUTH_CODE_SUCCESS, step_name="S1 IMEI+IMSI同时命中")


def step2_imei_hit_imsi_miss():
    """S2 (002_002): IMEI命中但IMSI未命中 → code=401"""
    print("\n[INFO] ========== S2 (002_002) IMEI命中IMSI未命中鉴权拒绝 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}(命中), IMSI={NON_WHITELIST_IMSI}(未命中)")
    try:
        resp = auth_imei_imsi(WHITELIST_IMEI, NON_WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_response(resp, AUTH_CODE_AUTH_FAILED, step_name="S2 IMEI命中IMSI未命中")


def step3_imsi_hit_imei_miss():
    """S3 (002_003): IMSI命中但IMEI未命中 → code=401"""
    print("\n[INFO] ========== S3 (002_003) IMSI命中IMEI未命中鉴权拒绝 ==========")
    print(f"[INFO] IMEI={NON_WHITELIST_IMEI}(未命中), IMSI={WHITELIST_IMSI}(命中)")
    try:
        resp = auth_imei_imsi(NON_WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_response(resp, AUTH_CODE_AUTH_FAILED, step_name="S3 IMSI命中IMEI未命中")


def step4_both_miss():
    """S4 (002_004): IMEI与IMSI均不在白名单 → code=401"""
    print("\n[INFO] ========== S4 (002_004) IMEI与IMSI均不在白名单鉴权拒绝 ==========")
    print(f"[INFO] IMEI={NON_WHITELIST_IMEI}, IMSI={NON_WHITELIST_IMSI}")
    try:
        resp = auth_imei_imsi(NON_WHITELIST_IMEI, NON_WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_response(resp, AUTH_CODE_AUTH_FAILED, step_name="S4 均不在白名单")


def step5_imei_missing():
    """S5 (002_005): IMEI缺失 -> code=401, msg含 format invalid"""
    print("\n[INFO] ========== S5 (002_005) IMEI缺失参数错误 ==========")
    print(f"[INFO] body: imsi={WHITELIST_IMSI} (no imei)")

    try:
        url = f"{GIDS_ADDR}/auth/v1/authIMEI"
        resp = requests.post(url, json={'imsi': WHITELIST_IMSI}, timeout=TIMEOUT, verify=False)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_response(resp, AUTH_CODE_AUTH_FAILED, "format invalid", "S5 IMEI缺失")


def step6_imsi_missing():
    """S6 (002_006): IMSI缺失 -> code=401, msg含 format invalid"""
    print("\n[INFO] ========== S6 (002_006) IMSI缺失参数错误 ==========")
    print(f"[INFO] body: imei={WHITELIST_IMEI} (no imsi)")

    try:
        url = f"{GIDS_ADDR}/auth/v1/authIMEI"
        resp = requests.post(url, json={'imei': WHITELIST_IMEI}, timeout=TIMEOUT, verify=False)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_response(resp, AUTH_CODE_AUTH_FAILED, "format invalid", "S6 IMSI缺失")


def step7_invalid_imei():
    """S7 (002_007): IMEI非15位纯数字 -> code=401, msg含 format invalid"""
    print("\n[INFO] ========== S7 (002_007) IMEI非15位纯数字参数错误 ==========")
    print(f"[INFO] IMEI={INVALID_IMEI_ALPHA}, IMSI={WHITELIST_IMSI}")
    try:
        resp = auth_imei_imsi(INVALID_IMEI_ALPHA, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_response(resp, AUTH_CODE_AUTH_FAILED, "format invalid", "S7 IMEI非纯数字")


def step8_invalid_imsi():
    """S8 (002_008): IMSI非15位纯数字 -> code=401, msg含 format invalid"""
    print("\n[INFO] ========== S8 (002_008) IMSI非15位纯数字参数错误 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}, IMSI={INVALID_IMSI_ALPHA}")
    try:
        resp = auth_imei_imsi(WHITELIST_IMEI, INVALID_IMSI_ALPHA)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_response(resp, AUTH_CODE_AUTH_FAILED, "format invalid", "S8 IMSI非纯数字")


def step9_escape_mode():
    """S9 (002_009): 逃生态-白名单表为空时放行 → code=200"""
    print("\n[INFO] ========== S9 (002_009) 逃生态-白名单表为空时放行 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    print("[INFO] 前置: 白名单表为空(全新gids.db或未导入)")
    try:
        resp = auth_imei_imsi(WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_response(resp, AUTH_CODE_SUCCESS, step_name="S9 逃生态放行")


def main():
    print("=" * 66)
    print("  IMEI+IMSI联合鉴权正向/反向/异常/逃生态  TC_SBG_Func_GIDS_Auth_002")
    print("=" * 66)
    print(f"[INFO] GIDS_ADDR = {GIDS_ADDR}")
    print(f"[INFO] 白名单内: IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    print(f"[INFO] 白名单外: IMEI={NON_WHITELIST_IMEI}, IMSI={NON_WHITELIST_IMSI}")
    print("[INFO] S1~S8需白名单已导入; S9需白名单表为空")

    results = []
    results.append(("S1 (002_001) IMEI+IMSI同时命中", step1_both_hit()))
    results.append(("S2 (002_002) IMEI命中IMSI未命中", step2_imei_hit_imsi_miss()))
    results.append(("S3 (002_003) IMSI命中IMEI未命中", step3_imsi_hit_imei_miss()))
    results.append(("S4 (002_004) 均不在白名单", step4_both_miss()))
    results.append(("S5 (002_005) IMEI缺失", step5_imei_missing()))
    results.append(("S6 (002_006) IMSI缺失", step6_imsi_missing()))
    results.append(("S7 (002_007) IMEI非15位纯数字", step7_invalid_imei()))
    results.append(("S8 (002_008) IMSI非15位纯数字", step8_invalid_imsi()))
    results.append(("S9 (002_009) 逃生态放行", step9_escape_mode()))

    print("\n[INFO] ========== 测试结果汇总 ==========")
    for name, ok in results:
        print(f"{'[PASS]' if ok else '[FAIL]'} {name}")

    all_passed = all(ok for _, ok in results)
    print(f"\n{'[SUCCESS]' if all_passed else '[FAILED]'} ========== 测试完成 ==========")
    return 0 if all_passed else 1


if __name__ == "__main__":
    sys.exit(main())
