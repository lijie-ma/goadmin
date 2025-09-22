-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO roles (id, code, name, description, status)
VALUES ('1', 'sup_admin', '超级管理员', '超级管理员', '1');


--  admin 123456
INSERT INTO users (id, username, password, email, status, role_code) VALUES
('1', 'admin', '$2y$10$oPdNH9rQsliggddN1lifkeIhHFGVSUjKnfT.EcyPrPMZiTRWnBdHm', '', '1', 'sup_admin');

INSERT INTO settings (id, name, value) VALUES
('1', 'captcha_switch', '{"admin":0, "web":1}');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
delete from roles where code = 'sup_admin';
delete from users where username = 'admin';
delete from settings where name in ('captcha_switch');
