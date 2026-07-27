#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# TC_SBG_Func_GIDS_Auth_004.py
# 测试用例: 事件上报链路IMEI+IMSI联合鉴权
#
# 看护点(逐条对齐testsuitcase.md):
#   S1 (004_001) sendClientEvent IMEI+IMSI均命中 → code=200
#   S2 (004_002) sendClientEvent IMSI不匹配 → code=401, msg含"auth rejected"
#   S3 (004_003) sendAppUseTimesEvent IMEI+IMSI均命中 → code=200
#   S4 (004_004) sendAppUseTimesEvent IMEI不匹配 → code=401, msg含"auth rejected"
#
# 前置：GIDS服务(127.0.0.1:9090)已启动, 白名单已导入
# 依赖：requests
# 用法：python TC_SBG_Func_GIDS_Auth_004.py

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
AUTH_CODE_AUTH_FAILED = 401


def send_client_event(imei, imsi):
    """POST sendClientEvent 发送客户端事件"""
    url = f"{GIDS_ADDR}/app-api/center/public/client/sendClientEvent"
    body = {
        "hsman": "Huawei",
        "hstype": "Mate60",
        "appType": "1",
        "imei": imei,
        "imsi": imsi,
        "type": "1"
    }
    resp = requests.post(url, json=body, timeout=TIMEOUT, verify=False)
    return resp


def send_app_use_times_event(imei, imsi):
    """POST sendAppUseTimesEvent 发送使用时长事件"""
    url = f"{GIDS_ADDR}/app-api/center/public/client/sendAppUseTimesEvent"
    body = {
        "useTimes": "100000",
        "hsman": "Huawei",
        "hstype": "Mate60",
        "appType": "1",
        "appId": "1",
        "scheight": "320",
        "scwidth": "240",
        "exttype": "default",
        "imei": imei,
        "imsi": imsi,
        "playMode": "1"
    }
    resp = requests.post(url, json=body, timeout=TIMEOUT, verify=False)
    return resp


def check_event_response(resp, expected_code, expected_msg_fragment=None, step_name=""):
    """通用事件上报响应校验"""
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
            print("[INFO] 请先执行 TC_SBG_Func_GIDS_Auth_001 导入白名单后再测试")
        return False
    if expected_msg_fragment and expected_msg_fragment.lower() not in actual_msg.lower():
        print(f"[FAIL] {step_name} 预期msg含'{expected_msg_fragment}', 实际msg={actual_msg}")
        return False
    print(f"[PASS] {step_name} code={actual_code}, msg={actual_msg}")
    return True


def step1_client_event_both_hit():
    """S1 (004_001): sendClientEvent IMEI+IMSI均命中 → code=200"""
    print("\n[INFO] ========== S1 (004_001) sendClientEvent联合鉴权通过 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    try:
        resp = send_client_event(WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_event_response(resp, AUTH_CODE_SUCCESS, step_name="S1 sendClientEvent通过")


def step2_client_event_imsi_mismatch():
    """S2 (004_002): sendClientEvent IMSI不匹配 -> code=401, msg含 auth rejected"""
    print("\n[INFO] ========== S2 (004_002) sendClientEvent联合鉴权拒绝(IMSI不匹配) ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}(命中), IMSI={NON_WHITELIST_IMSI}(不匹配)")
    try:
        resp = send_client_event(WHITELIST_IMEI, NON_WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_event_response(resp, AUTH_CODE_AUTH_FAILED, "auth rejected", "S2 sendClientEvent拒绝")


def step3_app_use_times_both_hit():
    """S3 (004_003): sendAppUseTimesEvent IMEI+IMSI均命中 → code=200"""
    print("\n[INFO] ========== S3 (004_003) sendAppUseTimesEvent联合鉴权通过 ==========")
    print(f"[INFO] IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    try:
        resp = send_app_use_times_event(WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_event_response(resp, AUTH_CODE_SUCCESS, step_name="S3 sendAppUseTimesEvent通过")


def step4_app_use_times_imei_mismatch():
    """S4 (004_004): sendAppUseTimesEvent IMEI不匹配 -> code=401, msg含 auth rejected"""
    print("\n[INFO] ========== S4 (004_004) sendAppUseTimesEvent联合鉴权拒绝(IMEI不匹配) ==========")
    print(f"[INFO] IMEI={NON_WHITELIST_IMEI}(不匹配), IMSI={WHITELIST_IMSI}(命中)")
    try:
        resp = send_app_use_times_event(NON_WHITELIST_IMEI, WHITELIST_IMSI)
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败")
        return False
    return check_event_response(resp, AUTH_CODE_AUTH_FAILED, "auth rejected", "S4 sendAppUseTimesEvent拒绝")


def main():
    print("=" * 66)
    print("  事件上报链路IMEI+IMSI联合鉴权  TC_SBG_Func_GIDS_Auth_004")
    print("=" * 66)
    print(f"[INFO] GIDS_ADDR = {GIDS_ADDR}")
    print(f"[INFO] 白名单内: IMEI={WHITELIST_IMEI}, IMSI={WHITELIST_IMSI}")
    print(f"[INFO] 白名单外: IMEI={NON_WHITELIST_IMEI}, IMSI={NON_WHITELIST_IMSI}")

    results = []
    results.append(("S1 (004_001) sendClientEvent联合鉴权通过", step1_client_event_both_hit()))
    results.append(("S2 (004_002) sendClientEvent IMSI不匹配拒绝", step2_client_event_imsi_mismatch()))
    results.append(("S3 (004_003) sendAppUseTimesEvent联合鉴权通过", step3_app_use_times_both_hit()))
    results.append(("S4 (004_004) sendAppUseTimesEvent IMEI不匹配拒绝", step4_app_use_times_imei_mismatch()))

    print("\n[INFO] ========== 测试结果汇总 ==========")
    for name, ok in results:
        print(f"{'[PASS]' if ok else '[FAIL]'} {name}")

    all_passed = all(ok for _, ok in results)
    print(f"\n{'[SUCCESS]' if all_passed else '[FAILED]'} ========== 测试完成 ==========")
    return 0 if all_passed else 1


if __name__ == "__main__":
    sys.exit(main())
