ALTER TABLE device_progress_status_history
    ADD COLUMN is_resolved SMALLINT DEFAULT 0,
    ADD COLUMN updated_at TIMESTAMPTZ,
    ADD COLUMN updated_by VARCHAR(50) NOT NULL,
    ADD COLUMN error_code VARCHAR(50),
    ADD COLUMN description VARCHAR(255);