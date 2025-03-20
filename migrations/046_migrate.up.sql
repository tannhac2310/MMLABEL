CREATE TABLE device_broken_history (
                                       id                               VARCHAR(50)  NOT NULL,
                                       production_order_stage_device_id VARCHAR(50)  NOT NULL,
                                       device_id                        VARCHAR(50)  NOT NULL,
                                       process_status                   INT          NOT NULL,
                                       is_resolved                      SMALLINT     DEFAULT 0,
                                       updated_at                        TIMESTAMPTZ,
                                       updated_by                        VARCHAR(50),
                                       created_by                        VARCHAR(50),
                                       error_code                        VARCHAR(50),
                                       error_reason                      VARCHAR(255),
                                       description                       VARCHAR(255),
                                       created_at                        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);