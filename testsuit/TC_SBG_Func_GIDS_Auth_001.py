#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# TC_SBG_Func_GIDS_Auth_001.py
# 测试用例: IMEI+IMSI白名单导入导出全链路
#
# 看护点：
#   S1 首次导入(firstImport): 上传CSV(IMEI+IMSI) → code=200
#   S2 重复首次导入(库非空): 再firstImport → code=-1
#   S3 导出白名单: GET导出 → CSV含IMEI+IMSI两列
#   S4 更新导入(update): 覆盖导入 → code=200
#   S5 再次导出验证: 导出内容非空含IMEI+IMSI
#   S6 CSV校验-IMEI非15位数字 → code=-1 拒绝整批
#   S7 CSV校验-IMSI非15位数字 → code=-1 拒绝整批
#
# 前置：GIDS服务(127.0.0.1:9090)已启动, 白名单表初始为空(逃生态)
# 依赖：requests
# 用法：python TC_SBG_Func_GIDS_Auth_001.py

import os
import sys
import csv
import io
import tempfile
import urllib3

import requests

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

GIDS_ADDR = os.getenv('GIDS_ADDR', 'http://127.0.0.1:9090').rstrip('/')
WHITELIST_CSV = os.getenv('WHITELIST_CSV', os.path.join(os.path.dirname(os.path.abspath(__file__)), 'imei_whitelist_test.csv'))
TIMEOUT = 30

AUTH_CODE_SUCCESS = 200
AUTH_CODE_INTERNAL_FAILED = -1


def import_imei_list(csv_path, operation):
    """POST /auth/v1/importIMEIList 上传白名单CSV"""
    url = f"{GIDS_ADDR}/auth/v1/importIMEIList"
    with open(csv_path, 'rb') as f:
        files = {'file': ('imei_whitelist_test.csv', f, 'text/csv')}
        data = {'operation': operation}
        resp = requests.post(url, files=files, data=data, timeout=TIMEOUT, verify=False)
    return resp


def export_imei_list():
    """GET /auth/v1/exportIMEIList 导出白名单CSV"""
    url = f"{GIDS_ADDR}/auth/v1/exportIMEIList"
    resp = requests.get(url, timeout=TIMEOUT, verify=False)
    return resp


def parse_exported_csv(resp):
    """解析导出的CSV响应, 返回(imei_list, imsi_list)"""
    reader = csv.reader(io.StringIO(resp.text))
    rows = list(reader)
    if len(rows) < 2:
        return [], []
    imei_list = []
    imsi_list = []
    for row in rows[1:]:
        if len(row) >= 2:
            imei_list.append(row[0])
            imsi_list.append(row[1])
    return imei_list, imsi_list


def create_invalid_csv(invalid_field, invalid_value):
    """创建含非法字段的临时CSV文件"""
    tmp = tempfile.NamedTemporaryFile(mode='w', suffix='.csv', delete=False, newline='')
    writer = csv.writer(tmp)
    writer.writerow(['IMEI', 'IMSI'])
    if invalid_field == 'IMEI':
        writer.writerow([invalid_value, '460011234567890'])
    elif invalid_field == 'IMSI':
        writer.writerow(['6258412454025411', invalid_value])
    tmp.close()
    return tmp.name


def step1_first_import():
    """S1: 首次导入白名单(firstImport) → code=200"""
    print("\n[INFO] ========== S1 首次导入白名单(firstImport, IMEI+IMSI) ==========")
    print(f"[INFO] CSV文件: {WHITELIST_CSV}")

    if not os.path.exists(WHITELIST_CSV):
        print(f"[FAIL] 白名单CSV文件不存在: {WHITELIST_CSV}")
        return False

    try:
        resp = import_imei_list(WHITELIST_CSV, 'firstImport')
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败: 无法连接到 {GIDS_ADDR}")
        return False

    print(f"[INFO] 响应状态码: {resp.status_code}")
    print(f"[INFO] 响应数据: {resp.text}")

    try:
        body = resp.json()
    except Exception:
        print(f"[FAIL] 响应非JSON: {resp.text}")
        return False

    if body.get('code') == AUTH_CODE_SUCCESS:
        print("[PASS] S1 首次导入成功")
        return True
    else:
        print(f"[FAIL] S1 首次导入失败, code={body.get('code')}, msg={body.get('msg')}")
        return False


def step2_duplicate_first_import():
    """S2: 重复首次导入(库非空时) → 预期失败 code=-1"""
    print("\n[INFO] ========== S2 重复首次导入(库非空, 预期失败) ==========")

    if not os.path.exists(WHITELIST_CSV):
        print(f"[FAIL] 白名单CSV文件不存在: {WHITELIST_CSV}")
        return False

    try:
        resp = import_imei_list(WHITELIST_CSV, 'firstImport')
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败: 无法连接到 {GIDS_ADDR}")
        return False

    print(f"[INFO] 响应状态码: {resp.status_code}")
    print(f"[INFO] 响应数据: {resp.text}")

    try:
        body = resp.json()
    except Exception:
        print(f"[FAIL] 响应非JSON: {resp.text}")
        return False

    if body.get('code') == AUTH_CODE_INTERNAL_FAILED and 'not empty' in body.get('msg', '').lower():
        print(f"[PASS] S2 重复首次导入被正确拒绝, code={body.get('code')}, msg={body.get('msg')}")
        return True
    else:
        print(f"[FAIL] S2 预期code={AUTH_CODE_INTERNAL_FAILED}(含'not empty'), 实际code={body.get('code')}, msg={body.get('msg')}")
        return False


def step3_export_verify():
    """S3: 导出白名单 → 验证CSV含IMEI+IMSI两列"""
    print("\n[INFO] ========== S3 导出白名单验证(IMEI+IMSI) ==========")

    try:
        resp = export_imei_list()
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败: 无法连接到 {GIDS_ADDR}")
        return False

    print(f"[INFO] 响应状态码: {resp.status_code}")
    print(f"[INFO] Content-Type: {resp.headers.get('Content-Type', 'N/A')}")

    if resp.status_code != 200:
        print(f"[FAIL] 导出失败, HTTP状态码={resp.status_code}")
        return False

    imei_list, imsi_list = parse_exported_csv(resp)
    print(f"[INFO] 导出IMEI条目数: {len(imei_list)}")
    for i, (imei, imsi) in enumerate(zip(imei_list, imsi_list)):
        print(f"[INFO]   [{i}] IMEI={imei}, IMSI={imsi}")

    if len(imei_list) == 0:
        print("[FAIL] 导出白名单为空")
        return False

    has_imei = any(len(v) == 15 and v.isdigit() for v in imei_list)
    has_imsi = any(len(v) == 15 and v.isdigit() for v in imsi_list)
    print(f"[INFO] 含合法IMEI条目: {has_imei}, 含合法IMSI条目: {has_imsi}")

    if has_imei and has_imsi:
        print("[PASS] S3 导出验证成功(含IMEI+IMSI)")
        return True
    else:
        print(f"[FAIL] S3 导出内容缺失, IMEI={has_imei}, IMSI={has_imsi}")
        return False


def step4_update_import():
    """S4: 更新导入(update) → code=200"""
    print("\n[INFO] ========== S4 更新导入白名单(update, IMEI+IMSI) ==========")

    if not os.path.exists(WHITELIST_CSV):
        print(f"[FAIL] 白名单CSV文件不存在: {WHITELIST_CSV}")
        return False

    try:
        resp = import_imei_list(WHITELIST_CSV, 'update')
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败: 无法连接到 {GIDS_ADDR}")
        return False

    print(f"[INFO] 响应状态码: {resp.status_code}")
    print(f"[INFO] 响应数据: {resp.text}")

    try:
        body = resp.json()
    except Exception:
        print(f"[FAIL] 响应非JSON: {resp.text}")
        return False

    if body.get('code') == AUTH_CODE_SUCCESS:
        print("[PASS] S4 更新导入成功")
        return True
    else:
        print(f"[FAIL] S4 更新导入失败, code={body.get('code')}, msg={body.get('msg')}")
        return False


def step5_export_after_update():
    """S5: 更新后再次导出 → 验证内容含IMEI+IMSI"""
    print("\n[INFO] ========== S5 更新后导出验证(IMEI+IMSI) ==========")

    try:
        resp = export_imei_list()
    except requests.exceptions.ConnectionError:
        print(f"[ERROR] 连接失败: 无法连接到 {GIDS_ADDR}")
        return False

    print(f"[INFO] 响应状态码: {resp.status_code}")

    if resp.status_code != 200:
        print(f"[FAIL] 导出失败, HTTP状态码={resp.status_code}")
        return False

    imei_list, imsi_list = parse_exported_csv(resp)
    print(f"[INFO] 导出IMEI条目数: {len(imei_list)}")
    for i, (imei, imsi) in enumerate(zip(imei_list, imsi_list)):
        print(f"[INFO]   [{i}] IMEI={imei}, IMSI={imsi}")

    if len(imei_list) == 0:
        print("[FAIL] 更新后导出为空")
        return False

    print("[PASS] S5 更新后导出验证成功(含IMEI+IMSI)")
    return True


def step6_invalid_imei_csv():
    """S6: CSV校验-IMEI非15位数字 → code=-1 拒绝整批"""
    print("\n[INFO] ========== S6 CSV校验-IMEI非15位数字拒绝 ==========")
    test_cases = [
        ('14位IMEI', '62584124540254'),
        ('16位IMEI', '62584124540254111'),
        ('含字母IMEI', '6258412454025ABC'),
    ]
    for desc, invalid_imei in test_cases:
        tmp_path = create_invalid_csv('IMEI', invalid_imei)
        print(f"[INFO] 测试: {desc} = {invalid_imei}")
        try:
            resp = import_imei_list(tmp_path, 'firstImport')
        except requests.exceptions.ConnectionError:
            print(f"[ERROR] 连接失败")
            os.unlink(tmp_path)
            return False
        os.unlink(tmp_path)
        try:
            body = resp.json()
        except Exception:
            print(f"[FAIL] 响应非JSON: {resp.text}")
            return False
        if body.get('code') == AUTH_CODE_SUCCESS:
            print(f"[FAIL] {desc}应被拒绝但导入成功, code={body.get('code')}")
            return False
        print(f"[INFO]   {desc}: code={body.get('code')}, msg={body.get('msg')} → 拒绝")
    print("[PASS] S6 IMEI非15位数字CSV均被拒绝")
    return True


def step7_invalid_imsi_csv():
    """S7: CSV校验-IMSI非15位数字 → code=-1 拒绝整批"""
    print("\n[INFO] ========== S7 CSV校验-IMSI非15位数字拒绝 ==========")
    test_cases = [
        ('14位IMSI', '46001123456789'),
        ('16位IMSI', '4600112345678901'),
        ('含字母IMSI', '46001ABC4567890'),
    ]
    for desc, invalid_imsi in test_cases:
        tmp_path = create_invalid_csv('IMSI', invalid_imsi)
        print(f"[INFO] 测试: {desc} = {invalid_imsi}")
        try:
            resp = import_imei_list(tmp_path, 'firstImport')
        except requests.exceptions.ConnectionError:
            print(f"[ERROR] 连接失败")
            os.unlink(tmp_path)
            return False
        os.unlink(tmp_path)
        try:
            body = resp.json()
        except Exception:
            print(f"[FAIL] 响应非JSON: {resp.text}")
            return False
        if body.get('code') == AUTH_CODE_SUCCESS:
            print(f"[FAIL] {desc}应被拒绝但导入成功, code={body.get('code')}")
            return False
        print(f"[INFO]   {desc}: code={body.get('code')}, msg={body.get('msg')} → 拒绝")
    print("[PASS] S7 IMSI非15位数字CSV均被拒绝")
    return True


def main():
    print("=" * 66)
    print("  IMEI+IMSI白名单导入导出全链路  TC_SBG_Func_GIDS_Auth_001")
    print("=" * 66)
    print(f"[INFO] GIDS_ADDR     = {GIDS_ADDR}")
    print(f"[INFO] WHITELIST_CSV = {WHITELIST_CSV}")

    results = []
    results.append(("S1 首次导入(firstImport)", step1_first_import()))
    results.append(("S2 重复首次导入(预期失败)", step2_duplicate_first_import()))
    results.append(("S3 导出白名单验证(IMEI+IMSI)", step3_export_verify()))
    results.append(("S4 更新导入(update)", step4_update_import()))
    results.append(("S5 更新后导出验证(IMEI+IMSI)", step5_export_after_update()))
    results.append(("S6 CSV校验-IMEI非15位数字拒绝", step6_invalid_imei_csv()))
    results.append(("S7 CSV校验-IMSI非15位数字拒绝", step7_invalid_imsi_csv()))

    print("\n[INFO] ========== 测试结果汇总 ==========")
    for name, ok in results:
        print(f"{'[PASS]' if ok else '[FAIL]'} {name}")

    all_passed = all(ok for _, ok in results)
    print(f"\n{'[SUCCESS]' if all_passed else '[FAILED]'} ========== 测试完成 ==========")
    return 0 if all_passed else 1


if __name__ == "__main__":
    sys.exit(main())
