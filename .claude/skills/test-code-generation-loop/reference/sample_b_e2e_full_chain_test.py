#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
样例B: E2E完整链路测试 (使用项目lib库 + TestContext断言模式)

测试目标: 功能机按键 - 宫格登录认证验证 (完整链路)
涉及服务: GIDS(9090) + BGW TCP(30001) + Browser Proxy(8000)

测试场景:
  - 宫格登录认证: GIDS三步登录鉴权获取token+tcpAddr
  - TCP控制通道: 连接BGW控制端口, TLV LOGIN获取mediaAddr
  - TCP媒体通道: 连接媒体端口, 验证视频帧到达
  - 浏览器渲染: proxy查找browser+context, 截图相似度验证九宫格界面
  - 清理释放: 关闭通道+删除浏览器实例

测试步骤 (对齐mobile BrowserContext.deviceLogin完整流程):
  阶段1: GIDS三步登录鉴权 → token + tcpAddr(控制通道地址)
  阶段2: TCP连接控制端口 → TLV LOGIN → ACK + RETURN_MEDIA(含mediaAddr)
  阶段3: TCP连接媒体端口(mediaAddr) → TLV LOGIN → ACK → 视频流开始推送
  阶段4: proxy查找BGW创建的browser+context → 截图相似度验证视频流 + 九宫格界面验证
  阶段5: 清理 + 报告

断言条件:
  阶段1: GIDS三步登录成功(login_ok=True) [fatal]
  阶段2: TCP连接成功(connect_ok=True) [fatal]
         控制通道LOGIN成功(tlv_ok=True) [fatal]
         RETURN_MEDIA包含mediaAddr(media_addr非空) [fatal]
  阶段3: 媒体通道LOGIN成功(media_ok=True) [fatal]
         视频帧到达(video_ok=True)
  阶段4: proxy找到browser+context(find_ok=True) [fatal]
         九宫格截图与基线相似度>=85%(is_similar_ss=True)
         完整链路: 截图相似+视频帧到达(chain_ok=True)
  阶段5: 浏览器实例删除成功(deleted=True)

此样例展示:
  1. 多服务协同的E2E测试流程
  2. TestContext.add_assertion 断言模式 (每条断言有名称 + 实际验证条件)
  3. 非伪断言: 每条断言验证实际行为而非恒真
  4. fatal=True 标记关键断言 (失败则终止后续步骤)
  5. 阶段化流程 (Phase 1→2→3→4→5)
  6. 清理 + 报告生成
  7. 测试函数docstring包含: 测试场景 + 测试步骤 + 断言条件
"""

import sys
import os
import time

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from lib.gids_client import GIDSClient
from lib.proxy_client import ProxyClient
from lib.tlv_client import BGWTlvClient
from lib.assertions import (
    TestContext, url_contains, screenshot_similarity, video_frame_arrived
)

GIDS_ADDR = os.environ.get("GIDS_ADDR", "http://127.0.0.1:9090")
PROXY_ADDR = os.environ.get("PROXY_ADDR", "http://127.0.0.1:8000")
BGW_TCP_ADDR = os.environ.get("BGW_TCP_ADDR", "127.0.0.1:30001")

DEVICE_IMEI = os.environ.get("DEVICE_IMEI", "6258412454025411")
DEVICE_IMSI = os.environ.get("DEVICE_IMSI", "68510155565211")


def test_keyboard_grid_login_auth():
    """
    测试场景: 功能机按键 - 宫格登录认证验证(完整链路)
    涉及服务: GIDS(9090) + BGW TCP(30001) + Browser Proxy(8000)

    测试步骤 (对齐mobile BrowserContext.deviceLogin完整流程):
      阶段1: GIDS三步登录鉴权 → token + tcpAddr(控制通道地址)
      阶段2: TCP连接控制端口 → TLV LOGIN → ACK + RETURN_MEDIA(含mediaAddr)
      阶段3: TCP连接媒体端口(mediaAddr) → TLV LOGIN → ACK → 视频流开始推送
      阶段4: proxy查找BGW创建的browser+context → 截图相似度验证视频流 + 九宫格界面验证
      阶段5: 清理 + 报告

    断言条件:
      阶段1: GIDS三步登录成功(login_ok=True) [fatal]
      阶段2: TCP连接成功(connect_ok=True) [fatal]
             控制通道LOGIN成功(tlv_ok=True) [fatal]
             RETURN_MEDIA包含mediaAddr(media_addr非空) [fatal]
      阶段3: 媒体通道LOGIN成功(media_ok=True) [fatal]
             视频帧到达(video_ok=True)
      阶段4: proxy找到browser+context(find_ok=True) [fatal]
             九宫格截图与基线相似度>=85%(is_similar_ss=True)
             完整链路: 截图相似+视频帧到达(chain_ok=True)
      阶段5: 浏览器实例删除成功(deleted=True)
    """
    ctx = TestContext(
        test_id="TC_SBG_Func_E2E_keyboard_grid_login",
        test_name="功能机按键 - 宫格登录认证验证(完整链路)",
    )
    gids = GIDSClient(GIDS_ADDR, DEVICE_IMEI, DEVICE_IMSI)
    proxy = ProxyClient(PROXY_ADDR)
    tlv = BGWTlvClient(BGW_TCP_ADDR)

    try:
        # 阶段1: GIDS三步登录鉴权 → token + tcpAddr
        ctx.next_step("阶段1: GIDS三步登录鉴权")
        login_ok = gids.three_step_login(width=240, height=320, device_type=2, app_type=1)
        ctx.add_assertion("GIDS三步登录成功", login_ok, f"login_ok={login_ok}", fatal=True)
        if not login_ok:
            return 1

        # 阶段2: TCP连接控制端口 → TLV LOGIN → ACK + RETURN_MEDIA(含mediaAddr)
        ctx.next_step("阶段2: TCP连接BGW控制端口")
        connect_ok = tlv.connect()
        ctx.add_assertion("TCP连接BGW控制端口成功", connect_ok, fatal=True)
        if not connect_ok:
            return 1

        ctx.next_step("阶段2: 控制通道发送TLV LOGIN")
        login_result = tlv.send_login(
            imei=DEVICE_IMEI, imsi=DEVICE_IMSI,
            token=gids.response.token,
            lcd_width=240, lcd_height=320,
            device_type=2, app_type=1,
            client_language="en_US",
        )
        tlv_ok = login_result.get("success", False)
        ctx.add_assertion(
            "控制通道LOGIN成功(ACK+RETURN_MEDIA)",
            tlv_ok,
            f"login_result={login_result}",
            fatal=True,
        )
        if not tlv_ok:
            return 1

        media_addr = login_result.get("media_addr", "")
        ctx.add_assertion(
            "RETURN_MEDIA包含mediaAddr",
            media_addr != "",
            f"media_addr={media_addr}",
            fatal=True,
        )
        if not media_addr:
            return 1

        # 阶段3: TCP连接媒体端口(mediaAddr) → TLV LOGIN → ACK → 视频流开始推送
        ctx.next_step("阶段3: TCP连接媒体端口并发送LOGIN")
        media_result = tlv.media_login(
            media_addr=media_addr,
            imei=DEVICE_IMEI, imsi=DEVICE_IMSI,
            token=gids.response.token,
            lcd_width=240, lcd_height=320,
            device_type=2, app_type=1,
            client_language="en_US",
        )
        media_ok = media_result.get("success", False)
        ctx.add_assertion(
            "媒体通道LOGIN成功(ACK code=200)",
            media_ok,
            f"media_result={media_result}",
            fatal=True,
        )
        if not media_ok:
            return 1

        ctx.next_step("阶段3: 等待视频帧数据到达")
        video_ok = video_frame_arrived(tlv, min_frames=1, timeout=10.0)
        ctx.add_assertion("模拟mobile客户端收到视频帧", video_ok)
        if not video_ok:
            return 1

        # 阶段4: proxy查找BGW创建的browser+context → 截图相似度验证
        ctx.next_step("阶段4: 查找BGW创建的browser+context")
        find_ok = proxy.find_ids(browser_type="KEYS")
        ctx.add_assertion(
            "proxy找到browser+context",
            find_ok,
            f"browser_id={proxy.browser_id}, context_id={proxy.context_id}",
            fatal=True,
        )
        if not find_ok:
            return 1

        ctx.next_step("阶段4: 验证九宫格界面加载")
        baseline_path = os.path.join(ctx.baselines_dir, f"{ctx.test_id}_grid_initial_state.png")
        is_similar_ss, sim_ratio_ss = screenshot_similarity(proxy, baseline_path, threshold=0.85)
        ctx.add_assertion(
            "九宫格截图与基线一致(相似度>=85%)",
            is_similar_ss,
            f"similarity={sim_ratio_ss:.4f}",
        )

        ctx.next_step("阶段4: 完整链路验证")
        chain_ok = is_similar_ss and video_ok
        ctx.add_assertion(
            "完整链路验证(截图相似>=85%+视频帧到达)",
            chain_ok,
            f"截图相似度={sim_ratio_ss:.4f}",
        )

        # 阶段5: 清理 + 报告
        ctx.next_step("阶段5: 清理控制+媒体通道及浏览器实例")
        tlv.close()
        if proxy.browser_id:
            deleted = proxy.delete_browser()
            ctx.add_assertion("浏览器实例删除成功", deleted)

    except AssertionError as e:
        print(f"\n[ERROR] 断言失败: {str(e)}")
    except Exception as e:
        print(f"\n[ERROR] 测试异常: {str(e)}")

    _, _, all_passed = ctx.summary()
    ctx.generate_html_report()
    return 0 if all_passed else 1


if __name__ == "__main__":
    sys.exit(test_keyboard_grid_login_auth())
