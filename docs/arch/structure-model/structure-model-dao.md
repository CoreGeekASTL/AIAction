# dao 模块结构文档

> 生成时间：2026-08-11
> 所属仓：AIAction（GIDS）
> 模块路径：src/dao

## 模块职责

数据访问层：基于 Beego ORM 适配 GaussDB（base_dao.go 提供 BaseInterface 与批量能力，db_init.go 初始化 ORM 并匿名导入 GIDS/db/driver），db_local_sqlite.go 提供 LOCAL_MODE 嵌入式 SQLite 初始化；其余文件按实体提供 browser_config、config_center、file、plugin、traffic_stats、user、white_list 等 CRUD，donothing_base_dao.go 提供空实现。测试辅助上，dao 的 _test.go 引用已被排除的测试辅助目录 src/test/util（testutil），不影响生产依赖。

## 子模块关系图

本模块为扁平包结构，无子模块。

## 子模块说明

| 关键文件 | 说明 |
| --- | --- |
| base_dao.go | 高斯数据库适配基类（BaseInterface、批量操作） |
| db_init.go | ORM/数据库初始化，匿名导入 GIDS/db/driver |
| db_local_sqlite.go | LOCAL_MODE SQLite 初始化 |
| donothing_base_dao.go | BaseInterface 空实现（降级/测试用） |
| browser_config.go / config_center.go / file.go / plugin.go / traffic_stats_dao.go / user.go / white_list.go | 各实体的 DAO 实现（white_list.go 含 ClearAndInsert 事务） |
