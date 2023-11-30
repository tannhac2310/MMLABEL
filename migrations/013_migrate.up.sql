CREATE TABLE production_order_device_config
(
    id                  VARCHAR(50) NOT NULL,
    production_order_id VARCHAR(50) NOT NULL,
    device_id           VARCHAR(50),
    color               VARCHAR(50),
    description         TEXT,
    search              TEXT,
    device_config       JSONB,
    created_by         VARCHAR(50) NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_by         VARCHAR(50) NOT NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    CONSTRAINT production_order_device_config_pk PRIMARY KEY (id)
);