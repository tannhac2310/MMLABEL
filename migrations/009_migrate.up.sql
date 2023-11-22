ALTER TABLE production_order_stage_devices
    ADD estimated_complete_at TIMESTAMPTZ;
ALTER TABLE production_order_stage_devices
    ADD assigned_quantity INT DEFAULT 0;
CREATE TABLE event_logs
(

    id           SERIAL PRIMARY KEY,
    device_id    VARCHAR(50) NOT NULL,
    stage_id     VARCHAR(50),
    stage_status SMALLINT             DEFAULT 1,
    quantity     FLOAT                DEFAULT 0,
    msg          TEXT,
    date         VARCHAR(50), -- YYYY-MM-DD
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ
);

CREATE TABLE history
(
    id          SERIAL PRIMARY KEY,
    table_name  VARCHAR(255) NOT NULL,
    row_id      VARCHAR(50)  NOT NULL,
    column_name VARCHAR(255),
    old_value   JSONB,
    new_value   JSONB,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
