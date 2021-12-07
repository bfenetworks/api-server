# 升级指南

本文档描述如何从一个已经部署的较早版本进行升级。

## v0.0.2

### 升级路径

可以从如下版本升级至v0.0.2:

- v0.0.1

### mysql 数据库表更新

```
ALTER TABLE users ADD COLUMN `type` tinyint(1) NOT NULL DEFAULT '0' AFTER name;
ALTER TABLE users ADD COLUMN `scopes` varchar(2048) NOT NULL DEFAULT '' AFTER `type`;

UPDATE users SET type = 0, scopes = 'System' WHERE roles = 'admin';
UPDATE users SET type = 1, scopes = 'Support' WHERE roles = 'inner';

ALTER TABLE users CHANGE COLUMN  `session_key`   `ticket` varchar(20) NOT NULL DEFAULT '';
ALTER TABLE users CHANGE COLUMN  `session_key_created_at`  `ticket_created_at` datetime NOT NULL DEFAULT '0000-01-01 00:00:00';

ALTER TABLE users DROP COLUMN `roles`;

ALTER TABLE users DROP INDEX  name_uni;
ALTER TABLE users ADD   UNIQUE KEY `name_uni` (`name`, `type`);
```

### Dashboard 版本升级
请升级 Dashboard 到 v0.02 或更新的版本。

### Conf-Agent 版本升级
需要 v0.0.1 或更新版本的 Conf-Agent 。

注意: 如果你使用 v0.0.1 版本的 Conf-Agent , 请按如下方式编辑 `conf/conf-agent.toml`:

```
# old:
{"Authorization" = "Session {Token}"}

# now:
{"Authorization" = "Token {Token}"}
```