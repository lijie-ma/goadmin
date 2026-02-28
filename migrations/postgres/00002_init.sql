-- +goose Up
-- SQL in this section is executed when the migration is applied.

--  admin 123456
INSERT INTO users (id, username, password, email, status, role_code) VALUES
(1, 'admin', '$2y$10$oPdNH9rQsliggddN1lifkeIhHFGVSUjKnfT.EcyPrPMZiTRWnBdHm', '', 1, 'sup_admin');

INSERT INTO roles (code, name, description, status, system_flag) VALUES
('sup_admin', '超级管理员', '超级管理员', 1, 1);

INSERT INTO role_permissions (role_code, permission_code) VALUES
 ('sup_admin', 'captcha_get'),
 ('sup_admin', 'captcha_set'),
 ('sup_admin', 'server_set'),
 ('sup_admin', 'server_get'),
 ('sup_admin', 'server_encrypted'),
 ('sup_admin', 'server_decrypted'),
 ('sup_admin', 'server_batch_get'),
 ('sup_admin', 'operate_log'),
 ('sup_admin', 'tenant_list'),
 ('sup_admin', 'tenant_info'),
 ('sup_admin', 'tenant_create'),
 ('sup_admin', 'tenant_update'),
 ('sup_admin', 'tenant_delete'),
 ('sup_admin', 'user_resetPwd'),
 ('sup_admin', 'user_list'),
 ('sup_admin', 'user_create'),
 ('sup_admin', 'user_update'),
 ('sup_admin', 'user_delete'),
 ('sup_admin', 'role_active'),
 ('sup_admin', 'role_list'),
 ('sup_admin', 'role_all'),
 ('sup_admin', 'role_info'),
 ('sup_admin', 'role_create'),
 ('sup_admin', 'role_update'),
 ('sup_admin', 'role_delete'),
 ('sup_admin', 'role_perm_set'),
 ('sup_admin', 'role_perm_info');

INSERT INTO role_permissions (role_code, permission_code) VALUES
('sup_admin', 'position_list'),
('sup_admin', 'position_info'),
('sup_admin', 'position_create'),
('sup_admin', 'position_update'),
('sup_admin', 'position_delete');

INSERT INTO permissions (code, name, description, path, module, global_flag) VALUES
('captcha_get',    '查询验证码',   '', 'admin/v1/setting/get_captcha_switch',   'server', 0),
('captcha_set',    '设置验证码',   '', 'admin/v1/setting/set_captcha_switch',   'server', 0),
('server_set',    '系统设置',   '', 'admin/v1/setting/set',   'server', 0),
('server_get',    '系统设置查询',   '', 'admin/v1/setting/get',   'server', 0),
('server_encrypted',    '系统设置(加密)',    '', 'admin/v1/setting/encrypted',   'server', 0),
('server_decrypted',    '系统设置(解密)查询', '', 'admin/v1/setting/decrypted',   'server', 0),
('server_batch_get','系统设置批量查询',   '', 'admin/v1/setting/batch',   'server', 0),
('upload_file',     '上传文件',          '', 'admin/v1/upload/file',     'upload', 1),
('operate_log',     '操作日志',          '', 'admin/v1/operate_log/list', 'operate_log', 0),
('tenant_list',     '租户列表',          '', 'admin/v1/tenant/list',   'tenant', 0),
('tenant_info',     '租户详情',          '', 'admin/v1/tenant/get',    'tenant', 0),
('tenant_create',   '租户创建',          '', 'admin/v1/tenant/create', 'tenant', 0),
('tenant_update',   '租户更新',          '', 'admin/v1/tenant/update', 'tenant', 0),
('tenant_delete',   '租户删除',          '', 'admin/v1/tenant/delete', 'tenant', 0),
('user_logout',    '退出',         '', 'admin/v1/user/logout',              'user',   1),
('user_changePwd', '修改密码',     '', 'admin/v1/user/change_pwd',           'user',   1),
('user_info',      '当前用户详情',   '', 'admin/v1/user/info',               'user',   1),
('user_list',      '管理员列表',     '', 'admin/v1/user/list',           'user',   0),
('user_resetPwd',  '重置密码',     '', 'admin/v1/user/reset_pwd',           'user',   0),
('user_create', '管理员创建',     '', 'admin/v1/user/create',           'user',   0),
('user_update', '管理员编辑',     '', 'admin/v1/user/update',           'user',   0),
('user_delete', '管理员删除',     '', 'admin/v1/user/delete',           'user',   0),
('role_active',    '活跃角色',     '', 'admin/v1/role/active',              'role',   0),
('role_list',      '角色列表',     '', 'admin/v1/role/list',                'role',   0),
('role_all',      '全部角色',     '', 'admin/v1/role/all',                'role',   0),
('role_info',      '角色详情',     '', 'admin/v1/role/get',                 'role',   0),
('role_create',    '角色创建',     '', 'admin/v1/role/create',              'role',   0),
('role_update',    '角色更新',     '', 'admin/v1/role/update',              'role',   0),
('role_delete',    '角色删除',     '', 'admin/v1/role/delete',              'role',   0),
('role_perm_set',  '角色权限设置', '', 'admin/v1/role/permissions/assign',  'role',   0),
('role_perm_info', '角色权限查看', '', 'admin/v1/role/permissions/get',     'role',   0);

INSERT INTO permissions (code, name, description, path, module, global_flag) VALUES
('position_list',      '位置列表',     '', 'admin/v1/position/list',                'position',   0),
('position_info',      '位置详情',     '', 'admin/v1/position/get',                 'position',   0),
('position_create',    '位置创建',     '', 'admin/v1/position/create',              'position',   0),
('position_update',    '位置更新',     '', 'admin/v1/position/update',              'position',   0),
('position_delete',    '位置删除',     '', 'admin/v1/position/delete',              'position',   0);

INSERT INTO server_setting (id, name, value) VALUES
(1, 'captcha_switch', '{"admin":1, "web":1}'),
(2, 'system_config', '{"system_name":"admin", "logo":"", "language":"zh-CN"}');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
delete from roles where code = 'sup_admin';
delete from users where username = 'admin';
delete from permissions where id>0;
delete from role_permissions where id>0;
delete from server_setting where name in ('captcha_switch');
