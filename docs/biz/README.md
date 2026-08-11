# AIAction 业务要素（biz）资产索引

> 生成时间：2026-08-11（终端鉴权变更刷新）
> 生成工具：all-index（自动聚合产物，同名覆盖更新；资产变更后重跑 all-index 刷新，请勿手改）

## interface/（接口）

服务对外接口清单，按功能域聚类

| 文件 | 说明 |
| --- | --- |
| [README.md](interface/README.md) | 对外接口总览 |
| [interface-auth.md](interface/interface-auth.md) | 终端鉴权与白名单管理 |
| [interface-cache.md](interface/interface-cache.md) | 缓存管理 |
| [interface-config-center.md](interface/interface-config-center.md) | 配置中心 |
| [interface-event.md](interface/interface-event.md) | 客户端事件上报 |
| [interface-file.md](interface/interface-file.md) | 文件管理 |
| [interface-login.md](interface/interface-login.md) | 登录鉴权与用户绑定 |
| [interface-management.md](interface/interface-management.md) | 浏览器配置同步 |
| [interface-plugin.md](interface/interface-plugin.md) | 插件管理 |
| [interface-stats.md](interface/interface-stats.md) | 流量统计 |
| [interface-test.md](interface/interface-test.md) | 测试联通 |

## rules/（业务规则）

按需求类整理"条件 → 动作 + 依据"规则条目

| 文件 | 说明 |
| --- | --- |
| [rules-auth.md](rules/rules-auth.md) | 终端鉴权与白名单管理 业务规则 |
| [rules-cache.md](rules/rules-cache.md) | 用户缓存删除 业务规则 |
| [rules-config-center.md](rules/rules-config-center.md) | 配置中心 业务规则 |
| [rules-event.md](rules/rules-event.md) | 客户端事件上报 业务规则 |
| [rules-file.md](rules/rules-file.md) | 文件上传下载 业务规则 |
| [rules-login.md](rules/rules-login.md) | 设备登录鉴权 业务规则 |
| [rules-management.md](rules/rules-management.md) | 浏览器配置同步 业务规则 |
| [rules-overload.md](rules/rules-overload.md) | 过载限流 业务规则 |
| [rules-plugin.md](rules/rules-plugin.md) | 插件包管理 业务规则 |
| [rules-traffic-stats.md](rules/rules-traffic-stats.md) | 流量统计 业务规则 |

## object-model/（对象模型）

实体、值对象、聚合、领域服务、领域事件

| 文件 | 说明 |
| --- | --- |
| [object-model-config-center.md](object-model/object-model-config-center.md) | 配置中心 对象模型 |
| [object-model-config.md](object-model/object-model-config.md) | 浏览器配置 对象模型 |
| [object-model-file.md](object-model/object-model-file.md) | 文件 对象模型 |
| [object-model-info.md](object-model/object-model-info.md) | 事件 对象模型 |
| [object-model-monitor-config.md](object-model/object-model-monitor-config.md) | 监控 对象模型 |
| [object-model-plugin-package.md](object-model/object-model-plugin-package.md) | 插件包 对象模型 |
| [object-model-schedule-election.md](object-model/object-model-schedule-election.md) | 定时任务选主 对象模型 |
| [object-model-service-instance.md](object-model/object-model-service-instance.md) | 浏览器网关实例 对象模型 |
| [object-model-session-stats.md](object-model/object-model-session-stats.md) | 流量统计 对象模型 |
| [object-model-user-bind.md](object-model/object-model-user-bind.md) | 用户实例绑定 对象模型 |
| [object-model-user.md](object-model/object-model-user.md) | 用户 对象模型 |
| [object-model-white-list.md](object-model/object-model-white-list.md) | 白名单 对象模型 |

## data-model/（数据模型）

持久态表结构、缓存数据结构、字段关系与数据生命周期

| 文件 | 说明 |
| --- | --- |
| [data-model-config-center.md](data-model/data-model-config-center.md) | 配置中心（ConfigCenter）数据模型 |
| [data-model-config.md](data-model/data-model-config.md) | 浏览器配置（Config）数据模型 |
| [data-model-control-traffic-stats.md](data-model/data-model-control-traffic-stats.md) | 控制流量统计（ControlTrafficStats）数据模型 |
| [data-model-file.md](data-model/data-model-file.md) | 文件（File）数据模型 |
| [data-model-media-traffic-stats.md](data-model/data-model-media-traffic-stats.md) | 媒体流量统计（MediaTrafficStats）数据模型 |
| [data-model-plugin-package.md](data-model/data-model-plugin-package.md) | 插件包（PluginPackage）数据模型 |
| [data-model-session-stats.md](data-model/data-model-session-stats.md) | 会话统计（SessionStats）数据模型 |
| [data-model-user-bind.md](data-model/data-model-user-bind.md) | 用户绑定（UserBind）数据模型 |
| [data-model-user.md](data-model/data-model-user.md) | 用户（User）数据模型 |
| [data-model-white-list.md](data-model/data-model-white-list.md) | 白名单（WhiteList）数据模型 |

## lexicon/（领域词典）

业务与代码共用的受控词汇集

| 文件 | 说明 |
| --- | --- |
| [lexicon.md](lexicon/lexicon.md) | GIDS 领域词典 |
| [lexicon-auth.md](lexicon/lexicon-auth.md) | 终端鉴权与白名单管理 领域词典 |
| [lexicon-cache.md](lexicon/lexicon-cache.md) | 缓存管理 领域词典 |
| [lexicon-config-center.md](lexicon/lexicon-config-center.md) | 配置中心 领域词典 |
| [lexicon-event.md](lexicon/lexicon-event.md) | 客户端事件上报 领域词典 |
| [lexicon-file.md](lexicon/lexicon-file.md) | 文件管理 领域词典 |
| [lexicon-login.md](lexicon/lexicon-login.md) | 登录鉴权与用户绑定 领域词典 |
| [lexicon-management.md](lexicon/lexicon-management.md) | 浏览器配置同步 领域词典 |
| [lexicon-plugin.md](lexicon/lexicon-plugin.md) | 插件管理 领域词典 |
| [lexicon-stats.md](lexicon/lexicon-stats.md) | 流量统计 领域词典 |
| [lexicon-test.md](lexicon/lexicon-test.md) | 测试联通 领域词典 |
