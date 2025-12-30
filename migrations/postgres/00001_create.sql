-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- 创建 permissions 表
CREATE TABLE permissions (
  id SERIAL PRIMARY KEY,
  ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  mtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  code VARCHAR(32) NOT NULL DEFAULT '',
  name VARCHAR(50) NOT NULL DEFAULT '',
  description VARCHAR(200),
  path VARCHAR(200) NOT NULL DEFAULT '',
  module VARCHAR(50) NOT NULL DEFAULT '',
  global_flag SMALLINT DEFAULT 2, -- 1:全局权限,2:局部权限
  UNIQUE (code),
  UNIQUE (name)
);

-- 创建 roles 表
CREATE TABLE roles (
  id SERIAL PRIMARY KEY,
  ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  mtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  code VARCHAR(32) NOT NULL DEFAULT '',
  name VARCHAR(50) NOT NULL DEFAULT '',
  description VARCHAR(200),
  status INT DEFAULT 1, -- 1:active,2:inactive
  system_flag SMALLINT DEFAULT 2, -- 1:系统内置,2:用户自定义
  UNIQUE (code),
  UNIQUE (name)
);

-- 创建 users 表
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  mtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  username VARCHAR(50) NOT NULL DEFAULT '',
  password VARCHAR(100) NOT NULL DEFAULT '',
  email VARCHAR(100),
  role_code VARCHAR(32) NOT NULL DEFAULT '',
  status INT DEFAULT 1, -- 0:inactive,1:active,2:locked,3:deleted
  UNIQUE (username),
  UNIQUE (email)
);

-- 创建 role_permissions 表
CREATE TABLE role_permissions (
  id SERIAL PRIMARY KEY,
  ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  mtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  role_code VARCHAR(32) NOT NULL DEFAULT '',
  permission_code VARCHAR(32) NOT NULL DEFAULT '',
  UNIQUE (role_code, permission_code)
);

-- 创建 server_setting 表 (10000以内为系统内置)
CREATE TABLE server_setting (
  id SERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL DEFAULT '',
  value VARCHAR(1024) NOT NULL DEFAULT '',
  mtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (name)
);

COMMENT ON TABLE server_setting IS '服务端配置';
COMMENT ON COLUMN server_setting.id IS '10000以内为系统内置配置';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS server_setting;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;
