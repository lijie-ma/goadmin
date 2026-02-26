-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE tenants (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  code VARCHAR(64) UNIQUE NOT NULL COMMENT '租户唯一编码',
  contact_email VARCHAR(128),
  contact_phone VARCHAR(32),
  status TINYINT DEFAULT 1 COMMENT '1启用 2停用',
  config JSON NULL COMMENT '扩展配置，如logo、域名、自定义参数',
  ctime DATETIME DEFAULT CURRENT_TIMESTAMP,
  mtime DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建permissions表
CREATE TABLE `permissions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `code` varchar(32)  NOT NULL DEFAULT '',
  `name` varchar(50)  NOT NULL DEFAULT '',
  `description` varchar(200)  DEFAULT '',
  `path` varchar(200)  NOT NULL DEFAULT '',
  `module` varchar(50)  NOT NULL DEFAULT '',
  `global_flag` TINYINT DEFAULT 2 COMMENT '1:全局权限,2:局部权限',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建roles表
CREATE TABLE `roles` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `code` varchar(32) NOT NULL DEFAULT '',
  `name` varchar(50) NOT NULL DEFAULT '',
  `description` varchar(200) DEFAULT '',
  `status` int DEFAULT '1' COMMENT '1:active,2:inactive',
  `system_flag` TINYINT DEFAULT 2 COMMENT '1:系统内置,2:用户自定义',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 创建users表
CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `username` varchar(50) NOT NULL DEFAULT '',
  `password` varchar(100) NOT NULL DEFAULT '',
  `email` varchar(100) DEFAULT '',
  `role_code` VARCHAR(32) NOT NULL DEFAULT '',
  `status` int DEFAULT '1' COMMENT '0:inactive,1:active,2:locked,3:deleted',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 创建role_permissions表
CREATE TABLE `role_permissions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `role_code` varchar(32) NOT NULL DEFAULT '',
  `permission_code` varchar(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_permission` (`role_code`,`permission_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建permissions表 10000 以内为系统内置权限
 CREATE TABLE `server_setting` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `value` varchar(1024) not null default '',
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COMMENT='服务端配置';

CREATE TABLE `operate_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(512) DEFAULT NULL COMMENT '详情内容',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '操作用户',
  `ip` varchar(45) NOT NULL DEFAULT '' COMMENT '操作人ip',
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='操作日志';


CREATE TABLE `position` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `city` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '城市名称',
  `location` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '详细位置（如街道/建筑）',
  `longitude` DECIMAL(10,6) NOT NULL COMMENT '经度',
  `latitude` DECIMAL(10,6) NOT NULL COMMENT '纬度',
  `custom_name` VARCHAR(128) DEFAULT NULL COMMENT '自定义名称',
  `creator_id` int DEFAULT 0 COMMENT '创建人',
  `creator` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '创建人',
  `ctime` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_city` (`city`),
  KEY `idx_location` (`location`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='位置信息表';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS server_setting;
DROP TABLE IF EXISTS operate_log;
DROP TABLE IF EXISTS position;