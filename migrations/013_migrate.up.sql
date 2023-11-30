CREATE TABLE production_order_device_config
(
    id              VARCHAR(50)  NOT NULL,
    production_order_id VARCHAR(50)  NOT NULL,
    device_id       VARCHAR(50)  NOT NULL,
    device_config   JSONB,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT production_order_device_config_pk PRIMARY KEY (id)
);