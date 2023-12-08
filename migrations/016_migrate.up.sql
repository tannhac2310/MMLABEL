CREATE TABLE device_working_history
(
    id                               VARCHAR(50) NOT NULL,
    production_order_stage_device_id VARCHAR(50) NOT NULL,
    device_id                        VARCHAR(50) NOT NULL,
    date                             VARCHAR(50) NOT NULL, -- YYYY-MM-DD
    quantity                         INT                  DEFAULT 0, -- sl trong ngay
    working_time                     INT                  DEFAULT 0, -- thoi gian lam viec trong ngay(s)
    updated_at                       TIMESTAMPTZ NULL,
    created_at                       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT device_working_history_pk PRIMARY KEY (id)
);
