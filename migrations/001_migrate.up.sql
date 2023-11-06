CREATE TABLE organizations
(
    id           VARCHAR(50)  NOT NULL,
    name         VARCHAR(255) NOT NULL,
    avatar       VARCHAR(512),
    phone_number VARCHAR(50),
    email        VARCHAR(255),
    status       INT2                  DEFAULT 1,
    type         INT2                  DEFAULT 1,
    address      VARCHAR(512),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT pk_organizations PRIMARY KEY (id ASC)
);
INSERT INTO organizations
(id, name, avatar, phone_number, email, "status", "type", created_at, updated_at, deleted_at)
VALUES ('1', 'Tech 4 edu', '', '0799590102', 'hoanggiangco94@gmail.com', 1, 3, '2021-04-15 13:06:54.205',
        '2021-04-15 13:06:54.205', NULL);

CREATE TABLE users
(
    id           VARCHAR(50)  NOT NULL,
    name         VARCHAR(255) NOT NULL,
    code         VARCHAR(255) NOT NULL,
    avatar       VARCHAR(512) NOT NULL DEFAULT '':::STRING,
    phone_number VARCHAR(50)  NOT NULL DEFAULT '':::STRING,
    email        VARCHAR(255) NOT NULL DEFAULT '':::STRING,
    linked       VARCHAR(255),
    departments  text, -- `BGD,HCNX,IC`
    status       INT2                  DEFAULT 1,
    type         INT2                  DEFAULT 1,
    language_id  INT2                  default 1,
    blocked      boolean               default false,
    address      VARCHAR(512) NOT NULL DEFAULT '':::STRING,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT pk_users PRIMARY KEY (id ASC),
    UNIQUE INDEX idx_users__code (code ASC)
);

COMMENT ON COLUMN users.status IS '1: active, 2: ban';

CREATE TABLE groups
(
    id         VARCHAR(50)  NOT NULL,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_groups PRIMARY KEY (id ASC)
);

CREATE TABLE roles
(
    id         VARCHAR(50)  NOT NULL,
    name       VARCHAR(255) NOT NULL,
    priority   INT                   DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_roles PRIMARY KEY (id ASC)
);

CREATE TABLE user_group
(
    id         VARCHAR(50) NOT NULL,
    user_id    VARCHAR(50) NOT NULL,
    group_id   VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_user_group PRIMARY KEY (id ASC),
    CONSTRAINT fk__user_group__users FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk__user_group__groups FOREIGN KEY (group_id) REFERENCES groups (id)
);

CREATE TABLE user_role
(
    id         VARCHAR(50) NOT NULL,
    user_id    VARCHAR(50) NOT NULL,
    role_id    VARCHAR(50) NOT NULL,
    created_by VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_user_role PRIMARY KEY (id ASC),
    CONSTRAINT fk__user_role__users FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk__user_role__roles FOREIGN KEY (role_id) REFERENCES roles (id)
);

CREATE TABLE user_firebases
(
    id         VARCHAR(50) NOT NULL,
    user_id    VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_user_firebases PRIMARY KEY (id ASC),
    CONSTRAINT fk__user_firebases__users FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE username_passwords
(
    id           VARCHAR(50) NOT NULL,
    user_id      VARCHAR(50) NOT NULL,
    username     VARCHAR(50),
    phone_number VARCHAR(50),
    email        VARCHAR(255),
    password     VARCHAR(512),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT pk_username_passwords PRIMARY KEY (id ASC),
    CONSTRAINT fk__username_password__user_id FOREIGN KEY (user_id) REFERENCES users (id),
    UNIQUE INDEX idx_username_passwords__username (username ASC),
    UNIQUE INDEX idx_username_passwords__email (email ASC),
    UNIQUE INDEX idx_username_passwords__phone_number (phone_number ASC)
);

CREATE TABLE user_fcm_tokens
(
    id         VARCHAR(50) NOT NULL,
    user_id    VARCHAR(50) NOT NULL,
    device_id  VARCHAR(255),
    token      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_user_fcm_tokens PRIMARY KEY (id ASC),
    CONSTRAINT fk__user_fcm_tokens__user_id FOREIGN KEY (user_id) REFERENCES users (id),
    UNIQUE INDEX idx__user_fcm_tokens__user_id__device_id (user_id, device_id)
);
create table departments
(
    id         varchar(50)                            NOT NULL,
    name       varchar(255)                           not null,
    short_name varchar(255)                           not null,
    code       varchar(255)                           not null,
    stage_code varchar(255)                           not null,
    priority   bigint                   default 100,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null,
    deleted_at timestamp with time zone,
    CONSTRAINT pk_departments PRIMARY KEY (id ASC),
    UNIQUE INDEX idx_departments__code (code ASC)
);