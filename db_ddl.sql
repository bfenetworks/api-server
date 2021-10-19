DROP DATABASE IF EXISTS `open_bfe`;
CREATE DATABASE open_bfe;

USE open_bfe;

-- create bfe_clusters
DROP TABLE IF EXISTS `bfe_clusters`;
CREATE TABLE `bfe_clusters` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `pool_name` varchar(255) NOT NULL DEFAULT '',
  `capacity` bigint(20) NOT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `gtc_enabled` tinyint(1) NOT NULL DEFAULT '1',
  `gtc_manual_enabled` tinyint(1) NOT NULL DEFAULT '1',
  `exempt_traffic_check` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_uni` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- create products
DROP TABLE IF EXISTS `products`;
CREATE TABLE `products` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL ,
  `mail_list` varchar(4096) NOT NULL ,
  `contact_person` varchar(4096) NOT NULL ,
  `sms_list` varchar(4096) NOT NULL DEFAULT "no sms" ,
  `gtc_priority` int NOT NULL DEFAULT 1,
  `description` varchar(1024) NOT NULL DEFAULT "no desc" ,
  `created_at` datetime NOT NULL ,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP ,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_uni` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- create domains
DROP TABLE IF EXISTS `domains`;
CREATE TABLE `domains` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `product_id` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `using_advanced_redirect` tinyint(1) NOT NULL DEFAULT 0,
  `using_advanced_hsts` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_uni` (`name`),
  INDEX `product_id` (`product_id`),
  INDEX `type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- create clusters
DROP TABLE IF EXISTS `clusters`;
CREATE TABLE `clusters` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `description` varchar(1024) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT "no desc",
  `product_id` bigint(20) NOT NULL,
  `max_idle_conn_per_host` smallint(6) NOT NULL DEFAULT '2',
  `timeout_conn_serv` int(11) NOT NULL DEFAULT '50000',
  `timeout_response_header` int(11) NOT NULL DEFAULT '50000',
  `timeout_readbody_client` int(11) NOT NULL DEFAULT '30000',
  `timeout_read_client_again` int(11) NOT NULL DEFAULT '30000',
  `timeout_write_client` int(11) NOT NULL DEFAULT '60000',
  `healthcheck_schem` varchar(16) NOT NULL DEFAULT 'http',
  `healthcheck_interval` int(11) NOT NULL DEFAULT '1000',
  `healthcheck_failnum` int(11) NOT NULL DEFAULT '10',
  `healthcheck_host` varchar(255) NOT NULL,
  `healthcheck_uri` varchar(255) NOT NULL,
  `healthcheck_statuscode` int(11) NOT NULL DEFAULT '200',
  `clientip_carry` tinyint(4) NOT NULL DEFAULT '0',
  `port_carry` tinyint(1) NOT NULL DEFAULT '0',
  `max_retry_in_cluster` tinyint(4) NOT NULL DEFAULT '3',
  `max_retry_cross_cluster` tinyint(4) NOT NULL DEFAULT '0',
  `ready` tinyint(1) NOT NULL DEFAULT '1',
  `hash_strategy` int NOT NULL DEFAULT '0',
  `cookie_key` varchar(255) NOT NULL DEFAULT 'BAIDUID',
  `hash_header` varchar(255) NOT NULL DEFAULT 'Cookie:BAIDUID',
  `session_sticky` tinyint(1) NOT NULL DEFAULT '0',
  `req_write_buffer_size` int(11) NOT NULL DEFAULT '512',
  `req_flush_interval` int(11) NOT NULL DEFAULT '0',
  `res_flush_interval` int(11) NOT NULL DEFAULT '20',
  `cancel_on_client_close` tinyint(1) NOT NULL DEFAULT '0',
  `failure_status` tinyint(1) NOT NULL DEFAULT '0',
  `max_conns_per_host` int(11) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_index` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- create lb_matrices
DROP TABLE IF EXISTS `lb_matrices`;
CREATE TABLE `lb_matrices` (
  `cluster_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `lb_matrix` varchar(8192) NOT NULL,
  `product_id` bigint(20) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`cluster_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- create sub_clusters
DROP TABLE IF EXISTS `sub_clusters`;
CREATE TABLE `sub_clusters` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `cluster_id` bigint(20) NOT NULL,
  `product_id` bigint(20) NOT NULL,
  `description` varchar(1024) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT "no desc",
  `bns_name_id` bigint(20) NOT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_product_index` (`name`, `product_id`),
  INDEX `cluster_id` (`cluster_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- create pools
DROP TABLE IF EXISTS `pools`;
CREATE TABLE `pools` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `product_id` bigint(20) NOT NULL DEFAULT 0,
  `ready`  boolean NOT NULL DEFAULT 1,
  `instance_detail` mediumtext,
  `type` tinyint(4) NOT NULL DEFAULT 1,
  `tag` tinyint(4) NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_uni` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- create route_basic_rules
DROP TABLE IF EXISTS `route_basic_rules`;
CREATE TABLE `route_basic_rules` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `description` varchar(1024) NOT NULL DEFAULT '',
  `product_id` bigint(20) NOT NULL,
  `host_names` text NOT NULL,
  `paths` text NOT NULL,
  `cluster_id` bigint(20) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- create route_advance_rules
DROP TABLE IF EXISTS `route_advance_rules`;
CREATE TABLE `route_advance_rules` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(1024) NOT NULL DEFAULT '',
  `product_id` bigint(20) NOT NULL,
  `expression` varchar(4096) binary NOT NULL,
  `cluster_id` bigint(20) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- create forward_cases
DROP TABLE IF EXISTS `route_cases`;
CREATE TABLE `route_cases` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `description` varchar(1024) NOT NULL DEFAULT '',
  `product_id` bigint(20) NOT NULL,
  `url` varchar(4096) NOT NULL,
  `method` varchar(255) NOT NULL DEFAULT "",
  `protocol` varchar(255) NOT NULL DEFAULT "",
  `header` varchar(4096) NOT NULL,
  `body` varchar(4096) NOT NULL DEFAULT "",
  `expect_cluster` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- create users
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL DEFAULT '',
  `session_key` varchar(20) NOT NULL DEFAULT '',
  `session_key_created_at` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `roles` varchar(2048) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_uni` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- create user_products
DROP TABLE IF EXISTS `user_products`;
CREATE TABLE `user_products` (
  `user_id` bigint(20) NOT NULL,
  `product_id` bigint(20) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`, `product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- create certificates
DROP TABLE IF EXISTS `certificates`;
CREATE TABLE `certificates` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `cert_name` varchar(255) NOT NULL,
  `description` varchar(1024) NOT NULL DEFAULT 'no desc',
  `is_default` tinyint(1) NOT NULL DEFAULT '0',
  `expired_date` varchar(255) NOT NULL,
  `cert_file_name` varchar(255) NOT NULL,
  `cert_file_path` varchar(255) NOT NULL,
  `key_file_name` varchar(255) NOT NULL,
  `key_file_path` varchar(255) NOT NULL,

  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`cert_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- create extra_files
DROP TABLE IF EXISTS `extra_files`;
CREATE TABLE `extra_files` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `product_id` bigint(20) NOT NULL DEFAULT 0,
  `description` varchar(1024) NOT NULL DEFAULT '',
  `md5` varchar(64) NOT NULL,
  `content` mediumtext,
  `created_at` datetime NOT NULL ,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP ,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_product` (`name`, `product_id`),
  INDEX `product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `config_versions`;
CREATE TABLE `config_versions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `data_sign` varchar(255) NOT NULL,
  `version` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


insert into users (id, name, roles, password, created_at) values(1, 'admin', 'admin', 'admin', now());
insert into products (id, name, `description`,                      mail_list,      contact_person, created_at) values
                     (1, 'BFE', 'Build-in Product, User by Admin', 'bfe@baidu.com', 'bfe',          now());