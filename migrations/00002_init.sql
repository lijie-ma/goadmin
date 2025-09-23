-- +goose Up
-- SQL in this section is executed when the migration is applied.

--  admin 123456
INSERT INTO users (id, username, password, email, status, role_code) VALUES
('1', 'admin', '$2y$10$oPdNH9rQsliggddN1lifkeIhHFGVSUjKnfT.EcyPrPMZiTRWnBdHm', '', '1', 'sup_admin');

INSERT INTO roles (id, code, name, description, status)
VALUES ('1', 'sup_admin', '超级管理员', '超级管理员', '1');

INSERT INTO `role_permissions` (`role_code`, `permission_code`) VALUES
 ('sup_admin', 'captcha_get'),
 ('sup_admin', 'captcha_set');

INSERT INTO `permissions` (`code`, `name`, `description`, `path`, `module`, `global_flag`) VALUES
('captcha_get', '查询验证码', '', 'admin/v1/setting/captcha-switch', 'server', 0),
('captcha_set', '设置验证码', '', 'admin/v1/setting/captcha-switch', 'server', 0),
('user_logout', '退出', '', 'admin/v1/user/logout', 'user', 1),
('user_changePwd', '修改密码', '', 'admin/v1/user/changePwd', 'user', 1);

INSERT INTO server_setting (id, name, value) VALUES
('1', 'captcha_switch', '{"admin":0, "web":1}');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
delete from roles where code = 'sup_admin';
delete from users where username = 'admin';
delete from permissions where id>0;
delete from server_setting where name in ('captcha_switch');