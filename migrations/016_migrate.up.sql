CREATE TABLE device_working_history
(
    id                               VARCHAR(50) NOT NULL,
    production_order_stage_device_id VARCHAR(50) NOT NULL,
    device_id                        VARCHAR(50) NOT NULL,
    working_date                     DATE        NOT NULL,
    qty                              INT         DEFAULT 0,
    updated_at                       TIMESTAMPTZ NULL,
    created_at                       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT device_working_history_pk PRIMARY KEY (id)
);
