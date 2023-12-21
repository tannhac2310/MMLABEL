CREATE TABLE role_permissions
(
    id          VARCHAR(50) NOT NULL,
    role_id     VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id   VARCHAR(50) NOT NULL,
    created_by   VARCHAR(50) NOT NULL,
    updated_by   VARCHAR(50) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ
);

alter table user_role drop constraint fk__user_role__users;
alter table user_role drop constraint fk__user_role__roles;