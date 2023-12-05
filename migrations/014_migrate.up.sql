CREATE TABLE device_progress_status_history
(
    id                               VARCHAR(50) NOT NULL,
    production_order_stage_device_id VARCHAR(50) NOT NULL,
    device_id                        VARCHAR(50) NOT NULL,
    process_status                   INT         NOT NULL,
    created_at                       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
