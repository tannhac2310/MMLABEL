CREATE TABLE casbin_rules
(
    p_type VARCHAR(256) NOT NULL DEFAULT '':::STRING,
    v0     VARCHAR(256) NOT NULL DEFAULT '':::STRING,
    v1     VARCHAR(256) NOT NULL DEFAULT '':::STRING,
    v2     VARCHAR(256) NOT NULL DEFAULT '':::STRING,
    v3     VARCHAR(256) NOT NULL DEFAULT '':::STRING,
    v4     VARCHAR(256) NOT NULL DEFAULT '':::STRING,
    v5     VARCHAR(256) NOT NULL DEFAULT '':::STRING,
    INDEX  idx_casbin_rules_p_type(p_type ASC),
    INDEX  idx_casbin_rules_v0(v0 ASC),
    INDEX  idx_casbin_rules_v1(v1 ASC),
    INDEX  idx_casbin_rules_v2(v2 ASC),
    INDEX  idx_casbin_rules_v3(v3 ASC),
    INDEX  idx_casbin_rules_v4(v4 ASC),
    INDEX  idx_casbin_rules_v5(v5 ASC),
    FAMILY "primary"(p_type, v0, v1, v2, v3, v4, v5, rowid)
);


INSERT INTO casbin_rules (p_type, v0, v1)
VALUES ('p', 'root', 'root'),
       ('p', 'user', '/hydra/user/update-profile'),
       ('p', 'user', '/hydra/user/update-fcm-token'),
       ('p', 'user', '/hydra/user/register-login-account'),
       ('p', 'user', '/hydra/user/get-current-profile'),
       ('p', 'account_manager', '/hydra/user/find-users'),
       ('p', 'account_manager', '/hydra/user/find-user-by-id'),
       ('p', 'account_manager', '/hydra/user/edit-user'),
       ('p', 'account_manager', '/hydra/user/create-user'),
       ('p', 'user', '/hydra/user/check-user-name'),
       ('p', 'account_manager', '/hydra/user/change-user-password'),
       ('p', 'user', '/hydra/user/change-password'),
       ('p', 'user', '/hydra/upload/file'),
       ('p', 'account_manager', '/hydra/role/remove-roles-for-user'),
       ('p', 'account_manager', '/hydra/role/find-role-by-id'),
       ('p', 'account_manager', '/hydra/role/find-role'),
       ('p', 'account_manager', '/hydra/role/edit-role'),
       ('p', 'account_manager', '/hydra/role/create-role'),
       ('p', 'account_manager', '/hydra/role/add-roles-for-user'),
       ('p', 'system_admin', '/hydra/notification/push-notification'),
       ('p', 'user', '/hydra/notification/mark-seen-notifications'),
       ('p', 'user', '/hydra/notification/mark-read-notifications'),
       ('p', 'user', '/hydra/notification/find-notifications'),
       ('p', 'user', '/hydra/notification/count-notification'),
       ('p', 'account_manager', '/hydra/group/remove-groups-for-user'),
       ('p', 'account_manager', '/hydra/group/find-group-by-id'),
       ('p', 'account_manager', '/hydra/group/find-group'),
       ('p', 'account_manager', '/hydra/group/edit-group'),
       ('p', 'account_manager', '/hydra/group/create-group'),
       ('p', 'account_manager', '/hydra/group/add-groups-for-user'),
       ('p', 'user', '/hydra/auth/reset-password'),
       ('p', 'user', '/hydra/auth/request-otp'),
       ('p', 'user', '/hydra/auth/refresh-token'),
       ('p', 'user', '/hydra/auth/login-username-password'),
       ('p', 'user', '/hydra/auth/login-otp'),
       ('p', 'user', '/hydra/auth/login-firebase'),
       ('p', 'user', '/gezu/config/app-config'),
       ('p', 'user', '/gezu/ws')
;
INSERT INTO "users"
(id, name, avatar, phone_number, email, linked, "status","type", created_at, updated_at, deleted_at)
VALUES ('01F3A27E8JFY7WZTV8FDDBWMJ3', 'Admin', '', '0799590102', 'hoanggiangco94@gmail.com', '', 1,3, '2021-04-15 13:06:54.205',
        '2021-04-15 13:06:54.205', NULL);

INSERT INTO username_passwords
(id, user_id, phone_number, email, "password", created_at, updated_at, deleted_at)
VALUES ('01F3A27E8WVA4R21R016QBKKTQ', '01F3A27E8JFY7WZTV8FDDBWMJ3', NULL, 'admin@gmail.com',
        '$2a$04$Y0EgK8do8/jCFWkuO.HTuunEDPwo1nP0aARmsbAQhIF/aVjxH4DD.', '2021-04-15 13:06:54.205',
        '2021-04-15 13:06:54.205', NULL);

INSERT INTO "roles"
(id, "name", "priority", created_at, updated_at, deleted_at)
VALUES ('root', 'Root', 0, '2021-04-15 13:06:54.205', '2021-04-15 13:06:54.205', NULL),
       ('user', 'User', 0, '2021-04-15 13:06:54.205', '2021-04-15 13:06:54.205', NULL);

INSERT INTO user_role
(id, user_id, role_id, created_at, updated_at, deleted_at)
VALUES ('01F3A27E8X2XZXNA3WGSAKYT1T', '01F3A27E8JFY7WZTV8FDDBWMJ3', 'root', '2021-04-15 13:06:54.205',
        '2021-04-15 13:06:54.205', NULL);

