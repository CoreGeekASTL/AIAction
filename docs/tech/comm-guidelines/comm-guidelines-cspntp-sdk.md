# cspntp-sdk 通信规范

## 平台 SDK

### ntp.Init（NTP 校时初始化）

- 业务场景：服务启动阶段初始化 NTP 校时 SDK，保证本节点时间与 NTP 服务同步
- 接口功能：`ntp.Init()`
- 调用位置：src/main.go
- 协议信息：
  - 协议：CSPNTP_SDK_GO（src/stubs/CSPNTP_SDK_GO；SDK 内部 NTP 通道，仓内为日志桩空实现）
  - 封装方式：平台 SDK 直连（`CSPNTP_SDK_GO/api`）
  - 超时重试：未设置（SDK 内部处理）
  - 错误码处理：未识别（SDK 内部处理）
