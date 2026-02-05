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
  UNIQUE (username)
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

CREATE TABLE operate_log (
    id SERIAL PRIMARY KEY,
    content VARCHAR(512),
    username VARCHAR(64) NOT NULL DEFAULT '',
    ip VARCHAR(45) NOT NULL DEFAULT '',
    mtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE operate_log IS '操作日志';
COMMENT ON COLUMN operate_log.content IS '详情内容';
COMMENT ON COLUMN operate_log.username IS '操作用户';
COMMENT ON COLUMN operate_log.ip IS '操作人ip';
COMMENT ON COLUMN operate_log.mtime IS '修改时间';
COMMENT ON COLUMN operate_log.ctime IS '记录创建时间';

CREATE TABLE position (
    id SERIAL PRIMARY KEY,                             -- 主键ID，自增
    city VARCHAR(64) NOT NULL DEFAULT '' ,              -- 城市名称
    location VARCHAR(128) NOT NULL DEFAULT '',          -- 详细位置（如街道/建筑）
    longitude NUMERIC(10,6) NOT NULL,                   -- 经度
    latitude NUMERIC(10,6) NOT NULL,                    -- 纬度
    custom_name VARCHAR(128) DEFAULT NULL,              -- 自定义名称
    creator_id INTEGER DEFAULT 0,                       -- 创建人ID
    creator VARCHAR(64) NOT NULL DEFAULT '',            -- 创建人
    ctime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    mtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
    CONSTRAINT chk_coordinates CHECK (longitude BETWEEN -180 AND 180 AND latitude BETWEEN -90 AND 90)
);

-- 索引
CREATE INDEX idx_city ON position(city);
CREATE INDEX idx_location ON position(location);

-- 表注释
COMMENT ON TABLE position IS '位置信息表';

-- 字段注释
COMMENT ON COLUMN position.city IS '城市名称';
COMMENT ON COLUMN position.location IS '详细位置（如街道/建筑）';
COMMENT ON COLUMN position.longitude IS '经度';
COMMENT ON COLUMN position.latitude IS '纬度';
COMMENT ON COLUMN position.custom_name IS '自定义名称';
COMMENT ON COLUMN position.creator_id IS '创建人ID';
COMMENT ON COLUMN position.creator IS '创建人';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS server_setting;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS operate_log;
DROP TABLE IF EXISTS position;