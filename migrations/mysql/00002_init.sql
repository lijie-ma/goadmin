-- +goose Up
-- SQL in this section is executed when the migration is applied.

--  admin 123456
INSERT INTO users (id, username, password, email, status, role_code) VALUES
('1', 'admin', '$2y$10$oPdNH9rQsliggddN1lifkeIhHFGVSUjKnfT.EcyPrPMZiTRWnBdHm', '', '1', 'sup_admin');

INSERT INTO `roles` (`code`, `name`, `description`, `status`, `system_flag`) VALUES
('sup_admin', '超级管理员', '超级管理员', 1, 1);

INSERT INTO `role_permissions` (`role_code`, `permission_code`) VALUES
 ('sup_admin', 'captcha_get'),
 ('sup_admin', 'captcha_set'),
 ('sup_admin', 'server_set'),
 ('sup_admin', 'server_get'),
 ('sup_admin', 'server_batch_get'),
 ('sup_admin', 'role_active'),
 ('sup_admin', 'role_list'),
 ('sup_admin', 'role_info'),
 ('sup_admin', 'role_create'),
 ('sup_admin', 'role_update'),
 ('sup_admin', 'role_delete'),
 ('sup_admin', 'role_perm_set'),
 ('sup_admin', 'role_perm_info');

INSERT INTO `permissions` (`code`, `name`, `description`, `path`, `module`, `global_flag`) VALUES
('user_logout',    '退出',         '', 'admin/v1/user/logout',              'user',   1),
('user_changePwd', '修改密码',     '', 'admin/v1/user/changePwd',           'user',   1),
('captcha_get',    '查询验证码',   '', 'admin/v1/setting/get_captcha_switch',   'server', 0),
('captcha_set',    '设置验证码',   '', 'admin/v1/setting/set_captcha_switch',   'server', 0),
('server_set',    '系统设置',   '', 'admin/v1/setting/set',   'server', 0),
('server_get',    '系统设置查询',   '', 'admin/v1/setting/get',   'server', 0),
('server_batch_get','系统设置批量查询',   '', 'admin/v1/setting/batch',   'server', 0),
('role_active',    '活跃角色',     '', 'admin/v1/role/active',              'role',   0),
('role_list',      '角色列表',     '', 'admin/v1/role',                     'role',   0),
('role_info',      '角色详情',     '', 'admin/v1/role/get',                 'role',   0),
('role_create',    '角色创建',     '', 'admin/v1/role/create',              'role',   0),
('role_update',    '角色更新',     '', 'admin/v1/role/update',              'role',   0),
('role_delete',    '角色删除',     '', 'admin/v1/role/delete',              'role',   0),
('role_perm_set',  '角色权限设置', '', 'admin/v1/role/permissions/assign',  'role',   0),
('role_perm_info', '角色权限查看', '', 'admin/v1/role/permissions/get',     'role',   0);


INSERT INTO server_setting (id, name, value) VALUES
('1', 'captcha_switch', '{"admin":1, "web":1}');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
delete from roles where code = 'sup_admin';
delete from users where username = 'admin';
delete from permissions where id>0;
delete from role_permissions where id>0;
delete from server_setting where name in ('captcha_switch');