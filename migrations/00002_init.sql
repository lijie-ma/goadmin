-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO roles (id, ctime, mtime, code, name, description, status)
VALUES ('1', '2025-09-19 10:43:03', '2025-09-19 10:43:03', 'sup_admin', '超级管理员', '超级管理员', '1');


INSERT INTO users (id, ctime, mtime, username, password, email, status, role_code) VALUES
('1', '2025-09-19 10:38:10', '2025-09-19 10:38:10', 'admin', '$2y$10$VZ6HMgqqnPb3cyxMbTbjl.PxYwlPwNTHJX2PCdDZ4aRehGi34AbkW', '', '1', 'sup_admin');


INSERT INTO settings (id, name, value, mtime, ctime) VALUES
('1', 'session_salt', 'Qwert', '2025-09-20 01:33:40', '2025-09-20 01:33:40'),
('2', 'captcha_switch', '{"server":0, "web":1}', '2025-09-20 01:33:40', '2025-09-20 01:33:40');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
delete from roles where code = 'sup_admin';
delete from users where username = 'admin';
delete from settings where name in ('session_salt', 'captcha_switch');
