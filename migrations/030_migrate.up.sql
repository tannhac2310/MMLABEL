-- create table production_order_stage_responsible

CREATE TABLE production_order_stage_responsible
(
    id                               VARCHAR(50) NOT NULL,
    production_order_stage_device_id varchar(50) NOT NULL,
    user_id                          varchar(50) NOT NULL,
    created_at                       TIMESTAMPTZ NOT NULL DEFAULT now()::TIMESTAMPTZ,
    updated_at                       TIMESTAMPTZ NOT NULL DEFAULT now()::TIMESTAMPTZ,
    deleted_at                       TIMESTAMPTZ,
    CONSTRAINT pk_production_order_stage_responsible PRIMARY KEY (id ASC)
);
create index idx_posp_pos_id on production_order_stage_responsible(production_order_stage_device_id);
create index idx_posp_user_id on production_order_stage_responsible(user_id);


ALTER TABLE production_order_stage_devices
    ADD COLUMN estimated_start_at TIMESTAMP;

ALTER TABLE production_order_stage_devices
    ADD COLUMN color VARCHAR(250);

ALTER TABLE production_order_stage_devices
    ADD COLUMN start_at TIMESTAMP;

ALTER TABLE production_order_stage_devices
    ADD COLUMN complete_at TIMESTAMP;